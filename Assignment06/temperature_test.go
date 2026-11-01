package main

import (
	"encoding/json"
	"testing"
)

func TestTemperatureContributesToNumericReports(t *testing.T) {
	patient := NewPatient("Patient", NewECG("ECG", 100, 1, 2), NewTemperature("Temperature", 40, 3, 8))
	if got := Apply[float64](patient, NewQualityVisitor()); got != 70 {
		t.Fatalf("quality = %v, want (100+40)/2 = 70", got)
	}
	if got := Apply[float64](patient, NewLatencyVisitor()); got != 5 {
		t.Fatalf("latency = %v, want (2+8)/2 = 5", got)
	}
	alarms := Apply[map[string]int](patient, NewAlarmVisitor())
	if alarms["ecg"] != 1 || alarms["temperature"] != 3 {
		t.Fatalf("alarms = %v, want ecg=1 and temperature=3", alarms)
	}
}

func TestTemperatureCanBeFilteredAndExported(t *testing.T) {
	temperature := NewTemperature("Temperature", 40, 3, 8)
	patient := NewPatient("Patient", NewECG("ECG", 100, 1, 2), temperature)
	if matches := Apply[[]Element](patient, NewFilterVisitor(60)); len(matches) != 1 || matches[0] != temperature {
		t.Fatalf("matches = %v, want temperature only", matches)
	}
	var output struct {
		Children []struct {
			Type      string  `json:"type"`
			Name      string  `json:"name"`
			Quality   float64 `json:"quality"`
			Alarms    int     `json:"alarms"`
			LatencyMS float64 `json:"latency_ms"`
		} `json:"children"`
	}
	if err := json.Unmarshal([]byte(Apply[string](patient, NewJSONVisitor())), &output); err != nil {
		t.Fatal(err)
	}
	if len(output.Children) != 2 {
		t.Fatalf("exported %d children, want 2", len(output.Children))
	}
	got := output.Children[1]
	if got.Type != "temperature" || got.Name != "Temperature" || got.Quality != 40 || got.Alarms != 3 || got.LatencyMS != 8 {
		t.Fatalf("temperature JSON = %+v", got)
	}
}
