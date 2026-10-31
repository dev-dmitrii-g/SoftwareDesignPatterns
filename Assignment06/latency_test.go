package main

import "testing"

func TestMeanLatencyExcludesDevice(t *testing.T) {
	patient := NewPatient("Patient", NewECG("ECG", 90, 0, 1), NewPPG("PPG", 80, 0, 2),
		NewEEG("EEG", 70, 0, 6), NewDevice("Monitor", 0, 0, 100))
	if got := Apply[float64](patient, NewLatencyVisitor()); got != 3 {
		t.Fatalf("latency = %v, want (1+2+6)/3 = 3", got)
	}
}

func TestEmptyMeanLatency(t *testing.T) {
	if got := Apply[float64](NewHospital("Empty"), NewLatencyVisitor()); got != 0 {
		t.Fatalf("empty latency = %v, want 0", got)
	}
}
