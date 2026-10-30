package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

type observerTestEvent EventType

func (e observerTestEvent) Type() EventType { return EventType(e) }

func TestEventTypeSubscription(t *testing.T) {
	bus := NewBus[observerTestEvent]()
	calls := 0
	bus.Subscribe("sample", ObserverFunc[observerTestEvent](func(observerTestEvent) {
		calls++
	}))
	bus.Publish("report")
	bus.Publish("sample")
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}

func TestUnsubscribeInsideHandler(t *testing.T) {
	bus := NewBus[observerTestEvent]()
	calls := 0
	var unsubscribe func()
	unsubscribe = bus.Subscribe("sample", ObserverFunc[observerTestEvent](func(observerTestEvent) {
		calls++
		unsubscribe()
	}))
	bus.Publish("sample")
	bus.Publish("sample")
	unsubscribe()
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}

func TestSubscribeInsideHandler(t *testing.T) {
	bus := NewBus[observerTestEvent]()
	calls := 0
	var once sync.Once
	bus.Subscribe("sample", ObserverFunc[observerTestEvent](func(observerTestEvent) {
		once.Do(func() {
			bus.Subscribe("sample", ObserverFunc[observerTestEvent](func(observerTestEvent) {
				calls++
			}))
		})
	}))
	bus.Publish("sample")
	if calls != 0 {
		t.Fatal("new listener received the current event")
	}
	bus.Publish("sample")
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}

func TestInFlightSnapshotCanFinish(t *testing.T) {
	bus := NewBus[observerTestEvent]()
	var unsubscribe func()
	bus.Subscribe("sample", ObserverFunc[observerTestEvent](func(observerTestEvent) {
		unsubscribe()
	}))
	calls := 0
	unsubscribe = bus.Subscribe("sample", ObserverFunc[observerTestEvent](func(observerTestEvent) {
		calls++
	}))
	bus.Publish("sample")
	bus.Publish("sample")
	if calls != 1 {
		t.Fatalf("snapshot deliveries = %d, want 1", calls)
	}
}

func TestStressFourProducersThousandEvents(t *testing.T) {
	bus := NewBus[observerTestEvent]()
	var calls atomic.Int64
	bus.Subscribe("sample", ObserverFunc[observerTestEvent](func(observerTestEvent) {
		calls.Add(1)
	}))
	var producers sync.WaitGroup
	for i := 0; i < 4; i++ {
		producers.Add(1)
		go func() {
			defer producers.Done()
			for j := 0; j < 1000; j++ {
				bus.Publish("sample")
			}
		}()
	}
	producers.Wait()
	if got := calls.Load(); got != 4000 {
		t.Fatalf("deliveries = %d, want 4000", got)
	}
}

func TestConcurrentRegistryChanges(t *testing.T) {
	bus := NewBus[observerTestEvent]()
	var calls atomic.Int64
	bus.Subscribe("sample", ObserverFunc[observerTestEvent](func(observerTestEvent) {
		calls.Add(1)
	}))
	var producers sync.WaitGroup
	for i := 0; i < 4; i++ {
		producers.Add(1)
		go func() {
			defer producers.Done()
			for j := 0; j < 1000; j++ {
				unsubscribe := bus.Subscribe("sample", ObserverFunc[observerTestEvent](func(observerTestEvent) {}))
				bus.Publish("sample")
				unsubscribe()
			}
		}()
	}
	producers.Wait()
	if got := calls.Load(); got != 4000 {
		t.Fatalf("permanent listener deliveries = %d, want 4000", got)
	}
}
