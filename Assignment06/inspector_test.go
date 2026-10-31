package main

import "testing"

func TestInspectionVisitsOnlyRelatedPatient(t *testing.T) {
	store := NewStore([]Channel{{"ecg", "patient-1", "ECG"}, {"ppg", "patient-1", "PPG"}, {"eeg", "patient-2", "EEG"}})
	store.OnEvent(Event{channel: "ecg", quality: 80})
	store.OnEvent(Event{channel: "eeg", quality: 0})
	source, results := NewBus[Event](), NewBus[Event]()
	summary := NewSummary()
	results.Subscribe(Inspection, summary)
	source.Subscribe(Sample, store)
	source.Subscribe(Sample, NewInspector(store, results))
	source.Publish(Event{kind: Sample, channel: "ppg", patient: "patient-1", quality: 40})
	_, count, report := summary.Snapshot()
	if count != 1 || report.patient != "patient-1" || report.quality != 60 {
		t.Fatalf("inspection = %+v, count=%d", report, count)
	}
	source.Publish(Event{kind: Sample, channel: "ppg", patient: "patient-1", quality: 90})
	_, count, _ = summary.Snapshot()
	if count != 1 {
		t.Fatal("healthy sample triggered inspection")
	}
}
