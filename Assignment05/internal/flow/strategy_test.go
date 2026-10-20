package flow

import (
	"reflect"
	"testing"
	"time"

	"hospitaladmission/internal/admission"
)

func TestFlowSelectorAutomaticRules(t *testing.T) {
	cases := []struct {
		name      string
		urgent    bool
		scheduled bool
		want      Strategy
	}{
		{name: "urgent takes priority", urgent: true, scheduled: true, want: Automatic{}},
		{name: "scheduled", scheduled: true, want: Scheduled{}},
		{name: "otherwise manual", want: Manual{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := sampleAdmission()
			a.Urgent = tc.urgent
			if !tc.scheduled {
				a.ScheduledAt = time.Time{}
			}
			got := (Selector{}).Automatic(a)
			if reflect.TypeOf(got) != reflect.TypeOf(tc.want) {
				t.Fatalf("got %T, want %T", got, tc.want)
			}
		})
	}
}

func TestFlowSelectorManualRules(t *testing.T) {
	for _, tc := range []struct {
		name string
		want Strategy
	}{
		{"automatic", Automatic{}},
		{"manual", Manual{}},
		{"scheduled", Scheduled{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := (Selector{}).Manual(tc.name)
			if err != nil || reflect.TypeOf(got) != reflect.TypeOf(tc.want) {
				t.Fatalf("got %T, error %v", got, err)
			}
		})
	}
	if _, err := (Selector{}).Manual("other"); err == nil {
		t.Fatal("unknown mode accepted")
	}
}

func TestFlowContextSwitch(t *testing.T) {
	context := NewContext(Manual{})
	a := sampleAdmission()
	before, err := context.Begin(a)
	if err != nil || before.Status != admission.PendingApproval {
		t.Fatalf("before=%+v error=%v", before, err)
	}
	context.SetStrategy(Automatic{})
	after, err := context.Begin(a)
	if err != nil || after.Status != admission.Completed {
		t.Fatalf("after=%+v error=%v", after, err)
	}
}

func TestAllFlowStrategiesRejectInvalidIIN(t *testing.T) {
	for name, selected := range map[string]Strategy{
		"automatic": Automatic{},
		"manual":    Manual{},
		"scheduled": Scheduled{},
	} {
		t.Run(name, func(t *testing.T) {
			a := sampleAdmission()
			a.IIN = "invalid"
			if _, err := selected.Begin(a); err == nil {
				t.Fatal("invalid IIN accepted")
			}
		})
	}
}

func TestScheduledFlowKeepsAppointment(t *testing.T) {
	a := sampleAdmission()
	result, err := (Scheduled{}).Begin(a)
	if err != nil || result.Status != admission.Scheduled || !result.ScheduledAt.Equal(a.ScheduledAt) {
		t.Fatalf("result=%+v error=%v", result, err)
	}
}
