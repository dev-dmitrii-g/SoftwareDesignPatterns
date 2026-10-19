package flow

import (
	"errors"
	"hospitaladmission/internal/admission"
)

type Processor struct{}

func (Processor) Begin(mode string, a admission.Admission) (admission.Result, error) {
	if err := a.Validate(); err != nil {
		return admission.Result{}, err
	}
	result := admission.Result{IIN: a.IIN}
	switch mode {
	case "automatic":
		result.Status = admission.Completed
		result.Steps = []string{"verify coverage", "reserve " + a.RequestedEquipment, "dispatch " + a.RequestedMedicine}
	case "manual":
		result.Status = admission.PendingApproval
		result.Steps = []string{"request staff approval"}
	case "scheduled":
		if a.ScheduledAt.IsZero() {
			return admission.Result{}, errors.New("scheduled time is required")
		}
		result.Status = admission.Scheduled
		result.ScheduledAt = a.ScheduledAt
		result.Steps = []string{"queue admission"}
	default:
		return admission.Result{}, errors.New("unknown admission mode")
	}
	return result, nil
}
