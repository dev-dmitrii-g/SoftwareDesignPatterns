package app

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunDistributesSamples(t *testing.T) {
	for _, workers := range []int{1, 2, 4} {
		t.Run(fmt.Sprint(workers), func(t *testing.T) {
			var counts [3]atomic.Int64
			var sources []*Source
			for index, kind := range []string{"ECG", "PPG", "EEG"} {
				bus := NewBus[Event]()
				bus.Subscribe(Sample, ObserverFunc[Event](func(Event) { counts[index].Add(1) }))
				sources = append(sources, NewSource(Channel{kind, "patient-1", kind}, bus))
			}
			accepted, err := Run(context.Background(), sources, workers, 1200)
			if err != nil || accepted != 1200 {
				t.Fatalf("accepted = %d, error = %v", accepted, err)
			}
			for index := range counts {
				if got := counts[index].Load(); got != 400 {
					t.Fatalf("source %d delivered %d, want 400", index, got)
				}
			}
		})
	}
}

func TestRunCancellationDrainsAcceptedSamples(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus := NewBus[Event]()
	var published atomic.Int64
	bus.Subscribe(Sample, ObserverFunc[Event](func(Event) {
		if published.Add(1) == 10 {
			cancel()
		}
	}))
	sources := []*Source{NewSource(Channel{"ecg", "patient-1", "ECG"}, bus)}
	type result struct {
		accepted int
		err      error
	}
	finished := make(chan result, 1)
	go func() {
		accepted, err := Run(ctx, sources, 4, 100000)
		finished <- result{accepted, err}
	}()
	select {
	case got := <-finished:
		if got.err != nil || got.accepted < 10 || got.accepted >= 100000 {
			t.Fatalf("cancelled run = %+v", got)
		}
		if delivered := published.Load(); delivered != int64(got.accepted) {
			t.Fatalf("accepted %d, delivered %d", got.accepted, delivered)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("run did not finish after cancellation")
	}
}

func TestRunRejectsInvalidInputs(t *testing.T) {
	sources := []*Source{NewSource(Channel{"ecg", "patient-1", "ECG"}, NewBus[Event]())}
	for _, test := range []struct {
		workers int
		samples int
		sources []*Source
	}{
		{0, 1, sources},
		{-1, 1, sources},
		{9, 1, sources},
		{1, 0, sources},
		{1, -1, sources},
		{1, 1, nil},
	} {
		if _, err := Run(context.Background(), test.sources, test.workers, test.samples); err == nil {
			t.Fatalf("accepted invalid input: %+v", test)
		}
	}
}
