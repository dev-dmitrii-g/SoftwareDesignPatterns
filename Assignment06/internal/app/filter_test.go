package app

import (
	"testing"
	"time"
)

func TestFilterAdaptsToLoadAndLatency(t *testing.T) {
	filter := NewAdaptiveFilter(NewBus[Event](), []Channel{{name: "ecg"}})
	for _, test := range []struct {
		load int
		age  time.Duration
		size int
	}{
		{0, -time.Second, 4}, {8, -time.Second, 16}, {0, time.Second, 16}, {0, -time.Second, 4},
	} {
		filter.OnEvent(Event{channel: "ecg", created: time.Now().Add(-test.age), load: test.load})
		if filter.Size() != test.size {
			t.Fatalf("window = %d, want %d", filter.Size(), test.size)
		}
	}
}

func TestFilterAveragesAndFlushesOnce(t *testing.T) {
	bus := NewBus[Event]()
	filter := NewAdaptiveFilter(bus, []Channel{{name: "ecg"}})
	values := []float64{}
	bus.Subscribe(Filtered, ObserverFunc[Event](func(e Event) { values = append(values, e.value) }))
	for _, value := range []float64{2, 4, 6, 8, 10, 12} {
		filter.OnEvent(Event{channel: "ecg", value: value, created: time.Now().Add(time.Second)})
	}
	filter.Flush()
	filter.Flush()
	if len(values) != 2 || values[0] != 5 || values[1] != 9 {
		t.Fatalf("averages = %v, want [5 9]", values)
	}
}

func TestFilterKeepsChannelsSeparate(t *testing.T) {
	bus := NewBus[Event]()
	filter := NewAdaptiveFilter(bus, []Channel{{name: "ecg"}, {name: "ppg"}})
	values := map[string]float64{}
	bus.Subscribe(Filtered, ObserverFunc[Event](func(e Event) { values[e.channel] = e.value }))
	for _, event := range []Event{{channel: "ecg", value: 2}, {channel: "ppg", value: 10}, {channel: "ecg", value: 6}} {
		event.created = time.Now().Add(time.Second)
		filter.OnEvent(event)
	}
	filter.Flush()
	if values["ecg"] != 4 || values["ppg"] != 10 {
		t.Fatalf("averages = %v", values)
	}
}
