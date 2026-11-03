package app

import (
	"encoding/json"
	"reflect"
	"testing"
)

func visitorTree() *Hospital {
	return NewHospital("Hospital", NewWard("Ward",
		NewPatient("Patient 1", NewECG("ECG", 90, 2, 1), NewPPG("PPG", 60, 1, 2), NewDevice("Monitor", 0, 3, 3)),
		NewPatient("Patient 2", NewEEG("EEG", 30, 4, 6)),
	))
}

func TestQualityMeanExcludesDevice(t *testing.T) {
	if got := Apply[float64](visitorTree(), NewQualityVisitor()); got != 60 {
		t.Fatalf("quality = %v, want (90+60+30)/3 = 60", got)
	}
}

func TestAlarmsAccumulateByType(t *testing.T) {
	got := Apply[map[string]int](visitorTree(), NewAlarmVisitor())
	want := map[string]int{"ecg": 2, "ppg": 1, "eeg": 4, "device": 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("alarms = %v, want %v", got, want)
	}
}

func TestFilterUsesStrictThreshold(t *testing.T) {
	got := Apply[[]Element](visitorTree(), NewFilterVisitor(60))
	if len(got) != 1 || got[0].Name() != "EEG" {
		t.Fatalf("matches = %v, want only EEG below 60", got)
	}
}

func TestVisitorOnlyVisitsRequestedSubtree(t *testing.T) {
	patient := visitorTree().Children()[0].(*Ward).Children()[0]
	if got := Apply[float64](patient, NewQualityVisitor()); got != 75 {
		t.Fatalf("patient quality = %v, want (90+60)/2 = 75", got)
	}
	matches := Apply[[]Element](patient, NewFilterVisitor(70))
	if len(matches) != 1 || matches[0].Name() != "PPG" {
		t.Fatalf("patient matches = %v, want PPG only", matches)
	}
}

func TestJSONPreservesNestingAndEscapesStrings(t *testing.T) {
	name := "Hospital \"A\"\n\\\t\b\f\x00"
	root := NewHospital(name, visitorTree().Children()...)
	type jsonNode struct {
		Type      string     `json:"type"`
		Name      string     `json:"name"`
		Children  []jsonNode `json:"children"`
		Quality   float64    `json:"quality"`
		Alarms    int        `json:"alarms"`
		LatencyMS float64    `json:"latency_ms"`
	}
	var got jsonNode
	output := Apply[string](root, NewJSONVisitor())
	if err := json.Unmarshal([]byte(output), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Type != "hospital" || got.Name != name || len(got.Children) != 1 {
		t.Fatalf("wrong hospital: %+v", got)
	}
	ward := got.Children[0]
	if ward.Type != "ward" || len(ward.Children) != 2 {
		t.Fatalf("wrong ward: %+v", ward)
	}
	first, second := ward.Children[0], ward.Children[1]
	if first.Type != "patient" || second.Type != "patient" || len(first.Children) != 3 || len(second.Children) != 1 {
		t.Fatalf("wrong patients: %+v", ward.Children)
	}
	ec := first.Children[0]
	if ec.Type != "ecg" || ec.Quality != 90 || ec.Alarms != 2 || ec.LatencyMS != 1 {
		t.Fatalf("wrong ECG: %+v", ec)
	}
	if first.Children[1].Type != "ppg" || first.Children[2].Type != "device" || second.Children[0].Type != "eeg" {
		t.Fatalf("missing leaf types: %+v", ward.Children)
	}
}

func TestEmptyVisitorResults(t *testing.T) {
	root := NewHospital("Empty")
	if Apply[float64](root, NewQualityVisitor()) != 0 {
		t.Fatal("empty quality must be zero")
	}
	if len(Apply[map[string]int](root, NewAlarmVisitor())) != 0 || len(Apply[[]Element](root, NewFilterVisitor(60))) != 0 {
		t.Fatal("empty tree produced alarms or matches")
	}
	if got := Apply[string](root, NewJSONVisitor()); got != `{"type":"hospital","name":"Empty","children":[]}` {
		t.Fatalf("empty JSON = %s", got)
	}
}

func TestContainerCopiesChildren(t *testing.T) {
	children := []Element{NewECG("ECG", 90, 0, 1)}
	patient := NewPatient("Patient", children...)
	children[0] = NewEEG("EEG", 10, 0, 1)
	patient.Children()[0] = children[0]
	if got := Apply[float64](patient, NewQualityVisitor()); got != 90 {
		t.Fatalf("caller changed snapshot quality to %v", got)
	}
}

func TestAlarmResultDoesNotExposeVisitorState(t *testing.T) {
	visitor := NewAlarmVisitor()
	result := Apply[map[string]int](visitorTree(), visitor)
	result["ecg"] = 100
	if visitor.Result()["ecg"] != 2 {
		t.Fatal("result exposes the visitor's map")
	}
}
