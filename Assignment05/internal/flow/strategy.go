package flow

import (
	"errors"
	"hospitaladmission/internal/admission"
)

type Strategy interface {
	Begin(admission.Admission) (admission.Result, error)
}

type Automatic struct{}

func (Automatic) Begin(a admission.Admission) (admission.Result, error) {
	if err := a.Validate(); err != nil {
		return admission.Result{}, err
	}
	return admission.Result{
		IIN:    a.IIN,
		Status: admission.Completed,
		Steps:  []string{"verify coverage", "reserve " + a.RequestedEquipment, "dispatch " + a.RequestedMedicine},
	}, nil
}

type Manual struct{}

func (Manual) Begin(a admission.Admission) (admission.Result, error) {
	if err := a.Validate(); err != nil {
		return admission.Result{}, err
	}
	return admission.Result{
		IIN:    a.IIN,
		Status: admission.PendingApproval,
		Steps:  []string{"request staff approval"},
	}, nil
}

type Scheduled struct{}

func (Scheduled) Begin(a admission.Admission) (admission.Result, error) {
	if err := a.Validate(); err != nil {
		return admission.Result{}, err
	}
	if a.ScheduledAt.IsZero() {
		return admission.Result{}, errors.New("scheduled time is required")
	}
	return admission.Result{
		IIN:         a.IIN,
		Status:      admission.Scheduled,
		ScheduledAt: a.ScheduledAt,
		Steps:       []string{"queue admission"},
	}, nil
}
