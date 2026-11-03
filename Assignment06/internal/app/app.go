package app

import (
	"context"
	"io"
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
