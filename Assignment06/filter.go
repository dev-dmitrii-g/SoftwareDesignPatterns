package main

import (
	"sync"
	"time"
)

type window struct {
	values  []float64
	pending int
	last    Event
}

type AdaptiveFilter struct {
	Subject[Event]
	mu      sync.Mutex
	windows map[string]*window
	size    int
}

func NewAdaptiveFilter(subject Subject[Event], channels []Channel) *AdaptiveFilter {
	windows := make(map[string]*window)
	for _, channel := range channels {
		windows[channel.name] = &window{}
	}
	return &AdaptiveFilter{Subject: subject, windows: windows, size: 4}
}

func (f *AdaptiveFilter) OnEvent(e Event) {
	f.mu.Lock()
	f.size = 4
	if e.load >= 8 || time.Since(e.created) >= 10*time.Millisecond {
		f.size = 16
	}
	w := f.windows[e.channel]
	if w == nil {
		f.mu.Unlock()
		return
	}
	w.values = append(w.values, e.value)
	if len(w.values) > f.size {
		w.values = w.values[len(w.values)-f.size:]
	}
	w.pending++
	w.last = e
	ready := w.pending >= f.size
	var result Event
	if ready {
		result = average(w)
	}
	f.mu.Unlock()
	if ready {
		f.Publish(result)
	}
}

func average(w *window) Event {
	sum := 0.0
	for _, value := range w.values {
		sum += value
	}
	e := w.last
	e.kind, e.value = Filtered, sum/float64(len(w.values))
	w.pending = 0
	return e
}

func (f *AdaptiveFilter) Flush() {
	f.mu.Lock()
	results := []Event{}
	for _, w := range f.windows {
		if w.pending > 0 {
			results = append(results, average(w))
		}
	}
	f.mu.Unlock()
	for _, result := range results {
		f.Publish(result)
	}
}

func (f *AdaptiveFilter) Size() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.size
}
