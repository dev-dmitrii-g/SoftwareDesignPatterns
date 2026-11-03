package app

import (
	"bytes"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPreviewUnsubscribesAtLimit(t *testing.T) {
	bus := NewBus[Event]()
	var output bytes.Buffer
	var unsubscribe func()
	callbacks := 0
	preview := NewPreview(&output, 2, func() {
		callbacks++
		unsubscribe()
	})
	unsubscribe = bus.Subscribe(Sample, preview)
	for range 4 {
		bus.Publish(Event{kind: Sample, channel: "ecg-1", value: 1.25, quality: 98})
	}
	want := strings.Repeat("ecg-1 value=1.250 quality=98.0%\n", 2)
	if output.String() != want || callbacks != 1 {
		t.Fatalf("output = %q, callbacks = %d", output.String(), callbacks)
	}
}

func TestPreviewIgnoresInFlightDeliveryAfterLimit(t *testing.T) {
	bus := NewBus[Event]()
	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	bus.Subscribe(Sample, ObserverFunc[Event](func(Event) {
		entered <- struct{}{}
		<-release
	}))
	var output bytes.Buffer
	var callbacks atomic.Int32
	var unsubscribe func()
	preview := NewPreview(&output, 1, func() {
		callbacks.Add(1)
		unsubscribe()
	})
	unsubscribe = bus.Subscribe(Sample, preview)
	var workers sync.WaitGroup
	for range 2 {
		workers.Add(1)
		go func() {
			defer workers.Done()
			bus.Publish(Event{kind: Sample, channel: "ppg-1", value: 0.5, quality: 90})
		}()
	}
	for range 2 {
		select {
		case <-entered:
		case <-time.After(time.Second):
			close(release)
			t.Fatal("publish did not snapshot both listeners")
		}
	}
	close(release)
	workers.Wait()
	bus.Publish(Event{kind: Sample, channel: "ppg-1"})
	if output.String() != "ppg-1 value=0.500 quality=90.0%\n" || callbacks.Load() != 1 {
		t.Fatalf("output = %q, callbacks = %d", output.String(), callbacks.Load())
	}
}

func TestSummaryCountsTypesAndKeepsLastInspection(t *testing.T) {
	summary := NewSummary()
	filtered, inspections, latest := summary.Snapshot()
	if filtered != 0 || inspections != 0 || latest != (Event{}) {
		t.Fatal("new summary is not empty")
	}
	last := Event{kind: Inspection, patient: "patient-2", quality: 55}
	for _, event := range []Event{
		{kind: Sample},
		{kind: Filtered},
		{kind: Inspection, patient: "patient-1", quality: 45},
		last,
		{kind: Filtered},
	} {
		summary.OnEvent(event)
	}
	filtered, inspections, latest = summary.Snapshot()
	if filtered != 2 || inspections != 2 || latest != last {
		t.Fatalf("filtered = %d, inspections = %d, latest = %+v", filtered, inspections, latest)
	}
}
