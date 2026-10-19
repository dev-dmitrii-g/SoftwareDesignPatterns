package flow

import "hospitaladmission/internal/admission"

type Processor struct{}

func (Processor) Begin(mode string, a admission.Admission) (admission.Result, error) {
	selected, err := (Selector{}).Manual(mode)
	if err != nil {
		return admission.Result{}, err
	}
	return NewContext(selected).Begin(a)
}
