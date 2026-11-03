package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
)

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

func RunCLI(ctx context.Context, args []string, output, progress io.Writer) error {
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
