package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	sources []*Source
	store   *Store
	filter  *AdaptiveFilter
	metrics *Metrics
	summary *Summary
}

func NewApp(output io.Writer) *App {
	channels := []Channel{
		{"ecg-1", "patient-1", "ECG"},
		{"ppg-1", "patient-1", "PPG"},
		{"eeg-2", "patient-2", "EEG"},
		{"temp-2", "patient-2", "TEMP"},
	}
	store := NewStore(channels)
	results := NewBus[Event]()
	filter := NewAdaptiveFilter(results, channels)
	inspector := NewInspector(store, results)
	metrics, summary := NewMetrics(store), NewSummary()
	results.Subscribe(Filtered, summary)
	results.Subscribe(Inspection, summary)
	stops := []func(){}
	preview := NewPreview(output, 5, func() {
		for _, stop := range stops {
			stop()
		}
	})
	sources := []*Source{}
	for _, channel := range channels {
		source := NewSource(channel, NewBus[Event]())
		source.Subscribe(Sample, store)
		source.Subscribe(Sample, filter)
		source.Subscribe(Sample, inspector)
		source.Subscribe(Sample, metrics)
		stops = append(stops, source.Subscribe(Sample, preview))
		sources = append(sources, source)
	}
	return &App{sources, store, filter, metrics, summary}
}

func simulate(ctx context.Context, workers, samples int, output io.Writer) (*App, time.Duration, error) {
	app := NewApp(output)
	started := time.Now()
	_, err := Run(ctx, app.sources, workers, samples)
	app.filter.Flush()
	return app, time.Since(started), err
}

func Report(tree Element, name string) (any, error) {
	switch name {
	case "quality":
		return Apply[float64](tree, NewQualityVisitor()), nil
	case "alarms":
		return Apply[map[string]int](tree, NewAlarmVisitor()), nil
	case "latency":
		return Apply[float64](tree, NewLatencyVisitor()), nil
	case "filter":
		names := []string{}
		for _, element := range Apply[[]Element](tree, NewFilterVisitor(60)) {
			names = append(names, element.Name())
		}
		return names, nil
	case "json":
		return json.RawMessage(Apply[string](tree, NewJSONVisitor())), nil
	default:
		return nil, fmt.Errorf("unknown report: %s", name)
	}
}

func runCLI(ctx context.Context, args []string, output, progress io.Writer) error {
	flags := flag.NewFlagSet("biosignals", flag.ContinueOnError)
	flags.SetOutput(progress)
	workers := flags.Int("workers", 3, "producer workers (1-8)")
	samples := flags.Int("samples", 12000, "total samples")
	report := flags.String("report", "quality", "quality, alarms, filter, json, latency")
	target := flags.String("target", "hospital", "hospital, patient-1, patient-2")
	bench := flags.Bool("bench", false, "compare 1, 2 and 4 workers")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}
	if flags.NArg() > 0 {
		return fmt.Errorf("unexpected argument: %s", flags.Arg(0))
	}
	if *workers < 1 || *workers > 8 || *samples < 1 {
		return fmt.Errorf("use 1-8 workers and a positive sample count")
	}
	if *target != "hospital" && *target != "patient-1" && *target != "patient-2" {
		return fmt.Errorf("unknown target: %s", *target)
	}
	if _, err := Report(NewHospital("validation"), *report); err != nil {
		return err
	}
	if *bench {
		fmt.Fprintln(output, "workers samples samples/s latency_ms")
		for _, count := range []int{1, 2, 4} {
			if ctx.Err() != nil {
				break
			}
			app, elapsed, err := simulate(ctx, count, *samples, io.Discard)
			if err != nil {
				return err
			}
			total, latency := app.metrics.Snapshot()
			fmt.Fprintf(output, "%d %d %.0f %.5f\n", count, total, float64(total)/elapsed.Seconds(), latency)
		}
		return nil
	}
	app, elapsed, err := simulate(ctx, *workers, *samples, progress)
	if err != nil {
		return err
	}
	total, latency := app.metrics.Snapshot()
	batches, inspections, latest := app.summary.Snapshot()
	fmt.Fprintf(progress, "samples=%d batches=%d inspections=%d window=%d\n", total, batches, inspections, app.filter.Size())
	fmt.Fprintf(progress, "throughput=%.0f samples/s latency=%.5f ms\n", float64(total)/elapsed.Seconds(), latency)
	if inspections > 0 {
		fmt.Fprintf(progress, "last inspection: %s quality=%.2f%%\n", latest.patient, latest.quality)
	}
	result, err := Report(app.store.Tree(*target), *report)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := runCLI(ctx, os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
