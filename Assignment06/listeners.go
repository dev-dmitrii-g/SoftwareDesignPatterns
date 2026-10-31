package main

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type Metrics struct {
	store        *Store
	mu           sync.Mutex
	count        int
	totalLatency time.Duration
}

func NewMetrics(store *Store) *Metrics {
	return &Metrics{store: store}
}

func (m *Metrics) OnEvent(e Event) {
	elapsed := time.Since(e.created)
	m.store.RecordLatency(e.channel, float64(elapsed)/float64(time.Millisecond))
	m.mu.Lock()
	m.count++
	m.totalLatency += elapsed
	m.mu.Unlock()
}

func (m *Metrics) Snapshot() (int, float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.count == 0 {
		return 0, 0
	}
	return m.count, float64(m.totalLatency) / float64(time.Millisecond) / float64(m.count)
}

type Preview struct {
	mu     sync.Mutex
	output io.Writer
	limit  int
	count  int
	done   func()
}

func NewPreview(output io.Writer, limit int, done func()) *Preview {
	return &Preview{output: output, limit: limit, done: done}
}

func (p *Preview) OnEvent(e Event) {
	p.mu.Lock()
	if p.count >= p.limit {
		p.mu.Unlock()
		return
	}
	fmt.Fprintf(p.output, "%s value=%.3f quality=%.1f%%\n", e.channel, e.value, e.quality)
	p.count++
	finished := p.count == p.limit
	p.mu.Unlock()
	if finished {
		p.done()
	}
}

type Summary struct {
	mu          sync.Mutex
	filtered    int
	inspections int
	latest      Event
}

func NewSummary() *Summary {
	return &Summary{}
}

func (s *Summary) OnEvent(e Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch e.Type() {
	case Filtered:
		s.filtered++
	case Inspection:
		s.inspections++
		s.latest = e
	}
}

func (s *Summary) Snapshot() (int, int, Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.filtered, s.inspections, s.latest
}
