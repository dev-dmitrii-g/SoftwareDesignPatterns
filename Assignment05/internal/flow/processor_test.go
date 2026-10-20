package flow

import (
	"testing"
	"time"

	"hospitaladmission/internal/admission"
)

func sampleAdmission() admission.Admission {
	return admission.Admission{
		IIN: "123456789123", Service: "consultation",
		RequestedMedicine: "Activated Charcoal", RequestedEquipment: "IV Stand",
		ScheduledAt: time.Date(2030, 10, 25, 10, 0, 0, 0, time.UTC),
	}
}

func TestProcessorModes(t *testing.T) {
	cases := []struct {
		mode, status string
		steps        int
	}{
		{"automatic", admission.Completed, 3},
		{"manual", admission.PendingApproval, 1},
		{"scheduled", admission.Scheduled, 1},
	}
	for _, tc := range cases {
		t.Run(tc.mode, func(t *testing.T) {
			result, err := (Processor{}).Begin(tc.mode, sampleAdmission())
			if err != nil {
				t.Fatal(err)
			}
			if result.Status != tc.status || len(result.Steps) != tc.steps {
				t.Fatalf("got %+v", result)
			}
		})
	}
}

func TestProcessorRejectsInvalidInput(t *testing.T) {
	a := sampleAdmission()
	a.IIN = "invalid"
	if _, err := (Processor{}).Begin("automatic", a); err == nil {
		t.Fatal("invalid IIN accepted")
	}
	a = sampleAdmission()
	a.ScheduledAt = time.Time{}
	if _, err := (Processor{}).Begin("scheduled", a); err == nil {
		t.Fatal("missing scheduled time accepted")
	}
	if _, err := (Processor{}).Begin("missing", sampleAdmission()); err == nil {
		t.Fatal("unknown mode accepted")
	}
}

func TestAllModesPreserveIIN(t *testing.T) {
	for _, mode := range []string{"automatic", "manual", "scheduled"} {
		t.Run(mode, func(t *testing.T) {
			a := sampleAdmission()
			result, err := (Processor{}).Begin(mode, a)
			if err != nil {
				t.Fatal(err)
			}
			if result.IIN != a.IIN {
				t.Fatalf("IIN changed: %q", result.IIN)
			}
		})
	}
}
