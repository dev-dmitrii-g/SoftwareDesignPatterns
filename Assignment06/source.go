package main

import (
	"math"
	"sync/atomic"
	"time"
)

type Source struct {
	Subject[Event]
	channel  Channel
	sequence atomic.Uint64
}

func NewSource(channel Channel, subject Subject[Event]) *Source {
	return &Source{Subject: subject, channel: channel}
}

func (s *Source) Produce(created time.Time, load int) {
	n := s.sequence.Add(1)
	x := float64(n) / 10
	noise := 0.04 * math.Sin(float64(n)*17.3)
	value := math.Sin(x) + noise
	switch s.channel.kind {
	case "ECG":
		value += 2 * math.Exp(-math.Pow(math.Mod(x, 2*math.Pi)-1, 2)/0.025)
	case "PPG":
		value = math.Pow((math.Sin(x)+1)/2, 2) + noise
	case "EEG":
		value = 0.4*math.Sin(8*x) + 0.2*math.Sin(13*x) + noise
	case "TEMP":
		value = 36.6 + 0.2*math.Sin(x/10) + noise
	}
	quality := 96 + 3*math.Sin(x)
	if n%40 < 4 {
		quality = 45
	}
	s.Publish(Event{Sample, s.channel.name, s.channel.patient, value, quality, created, load})
}
