package main

import "sync"

type state struct {
	quality float64
	alarms int
	latencyMS float64
}

type Store struct {
	mu sync.RWMutex
	channels []Channel
	readings map[string]state
}

func NewStore(channels []Channel) *Store {
	readings := make(map[string]state)
	for _, channel := range channels { readings[channel.name] = state{quality: 100} }
	return &Store{channels: append([]Channel(nil), channels...), readings: readings}
}

func (s *Store) OnEvent(e Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	reading := s.readings[e.channel]
	reading.quality = e.quality
	if e.quality < 60 { reading.alarms++ }
	s.readings[e.channel] = reading
}

func (s *Store) RecordLatency(channel string, milliseconds float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	reading := s.readings[channel]
	reading.latencyMS = milliseconds
	s.readings[channel] = reading
}

func (s *Store) Tree(target string) Element {
	s.mu.RLock()
	defer s.mu.RUnlock()
	patient := func(name string) Element {
		children := []Element{}
		for _, channel := range s.channels {
			if channel.patient != name { continue }
			r := s.readings[channel.name]
			switch channel.kind {
			case "ECG": children = append(children, NewECG(channel.name, r.quality, r.alarms, r.latencyMS))
			case "PPG": children = append(children, NewPPG(channel.name, r.quality, r.alarms, r.latencyMS))
			case "EEG": children = append(children, NewEEG(channel.name, r.quality, r.alarms, r.latencyMS))
			}
		}
		children = append(children, NewDevice(name+" monitor", 0, 0, 0))
		return NewPatient(name, children...)
	}
	if target != "hospital" { return patient(target) }
	return NewHospital("hospital", NewWard("ward", patient("patient-1"), patient("patient-2")))
}
