package coverage

import (
	"errors"
	"strings"
	"testing"
	"time"
)

type fakeSource struct {
	record Record
	err    error
}

func (s fakeSource) Lookup(string) (Record, error) {
	return s.record, s.err
}

var testTime = time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)

func TestEachCoverageStrategy(t *testing.T) {
	cases := []struct {
		name     string
		strategy Strategy
		service  string
		approved bool
	}{
		{"GOBMP public service", GOBMP{fakeSource{record: Record{GOBMPEligible: true}}}, "emergency", true},
		{"GOBMP other service", GOBMP{fakeSource{record: Record{GOBMPEligible: true}}}, "consultation", false},
		{"OSMS active", OSMS{fakeSource{record: Record{OSMSActive: true}}}, "consultation", true},
		{"OSMS inactive", OSMS{fakeSource{record: Record{OSMSActive: false}}}, "consultation", false},
		{"OSMS rehabilitation referral", OSMS{fakeSource{record: Record{OSMSActive: true, HasReferral: true}}}, "rehabilitation", true},
		{"OSMS rehabilitation without referral", OSMS{fakeSource{record: Record{OSMSActive: true}}}, "rehabilitation", false},
		{"DMS covered", DMS{fakeSource{record: Record{DMSExpiresAt: testTime.Add(time.Hour), DMSServices: map[string]bool{"consultation": true}}}}, "consultation", true},
		{"DMS expired", DMS{fakeSource{record: Record{DMSExpiresAt: testTime.Add(-time.Hour), DMSServices: map[string]bool{"consultation": true}}}}, "consultation", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.strategy.Check(Request{IIN: "123456789123", Service: tc.service, At: testTime})
			if err != nil {
				t.Fatal(err)
			}
			if got.Approved != tc.approved || got.Program == "" {
				t.Fatalf("got %+v", got)
			}
		})
	}
}

func TestAllStrategiesReportProgram(t *testing.T) {
	source := fakeSource{record: Record{GOBMPEligible: true, OSMSActive: true, DMSExpiresAt: testTime.Add(time.Hour), DMSServices: map[string]bool{"consultation": true}}}
	strategies := map[string]Strategy{"GOBMP": GOBMP{source}, "OSMS": OSMS{source}, "DMS": DMS{source}}
	for name, candidate := range strategies {
		t.Run(name, func(t *testing.T) {
			got, err := candidate.Check(Request{IIN: "123456789123", Service: "consultation", At: testTime})
			if err != nil || got.Program != name || got.Reason == "" {
				t.Fatalf("decision=%+v error=%v", got, err)
			}
		})
	}
}

func TestStrategiesPropagateSourceFailure(t *testing.T) {
	source := fakeSource{err: errors.New("source offline")}
	for name, candidate := range map[string]Strategy{"GOBMP": GOBMP{source}, "OSMS": OSMS{source}, "DMS": DMS{source}} {
		t.Run(name, func(t *testing.T) {
			if _, err := candidate.Check(Request{IIN: "123456789123"}); err == nil || !strings.Contains(err.Error(), "offline") {
				t.Fatalf("error=%v", err)
			}
		})
	}
}

func TestSelectorAutomaticRules(t *testing.T) {
	cases := []struct {
		name, service, want string
		record              Record
		wantError           bool
	}{
		{name: "DMS first", service: "consultation", want: "DMS", record: Record{OSMSActive: true, DMSExpiresAt: testTime.Add(time.Hour), DMSServices: map[string]bool{"consultation": true}}},
		{name: "OSMS fallback", service: "consultation", want: "OSMS", record: Record{OSMSActive: true, DMSExpiresAt: testTime.Add(-time.Hour), DMSServices: map[string]bool{"consultation": true}}},
		{name: "GOBMP fallback", service: "emergency", want: "GOBMP", record: Record{GOBMPEligible: true}},
		{name: "none", service: "consultation", record: Record{}, wantError: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			selector := Selector{Source: fakeSource{record: tc.record}}
			selected, err := selector.Automatic(Request{IIN: "123456789123", Service: tc.service, At: testTime})
			if tc.wantError {
				if err == nil {
					t.Fatal("expected no coverage error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			decision, err := selected.Check(Request{IIN: "123456789123", Service: tc.service, At: testTime})
			if err != nil || decision.Program != tc.want || !decision.Approved {
				t.Fatalf("decision=%+v error=%v", decision, err)
			}
		})
	}
}

func TestSelectorManualRules(t *testing.T) {
	selector := Selector{Source: fakeSource{record: Record{}}}
	for _, name := range []string{"GOBMP", "OSMS", "DMS"} {
		t.Run(name, func(t *testing.T) {
			selected, err := selector.Manual(strings.ToLower(name))
			if err != nil {
				t.Fatal(err)
			}
			decision, err := selected.Check(Request{IIN: "123456789123", At: testTime})
			if err != nil || decision.Program != name {
				t.Fatalf("decision=%+v error=%v", decision, err)
			}
		})
	}
	if _, err := selector.Manual("other"); err == nil {
		t.Fatal("unknown program accepted")
	}
}

func TestContextSwitchesCoverageStrategy(t *testing.T) {
	source := fakeSource{record: Record{OSMSActive: true}}
	context := NewContext(GOBMP{Source: source})
	request := Request{IIN: "123456789123", Service: "consultation", At: testTime}
	before, err := context.Check(request)
	if err != nil || before.Approved {
		t.Fatalf("before=%+v error=%v", before, err)
	}
	context.SetStrategy(OSMS{Source: source})
	after, err := context.Check(request)
	if err != nil || !after.Approved {
		t.Fatalf("after=%+v error=%v", after, err)
	}
}
