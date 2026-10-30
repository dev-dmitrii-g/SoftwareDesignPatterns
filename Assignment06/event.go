package main

import "time"

const (
	Sample EventType = "sample"
	Filtered EventType = "filtered"
	Inspection EventType = "inspection"
)

type Event struct {
	kind EventType
	channel string
	patient string
	value float64
	quality float64
	created time.Time
	load int
}

func (e Event) Type() EventType { return e.kind }

type Channel struct {
	name string
	patient string
	kind string
}
