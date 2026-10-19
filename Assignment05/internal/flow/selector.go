package flow

import (
	"errors"
	"hospitaladmission/internal/admission"
)

type Selector struct{}

func (Selector) Manual(name string) (Strategy, error) {
	switch name {
	case "automatic":
		return Automatic{}, nil
	case "manual":
		return Manual{}, nil
	case "scheduled":
		return Scheduled{}, nil
	default:
		return nil, errors.New("unknown admission mode")
	}
}

func (Selector) Automatic(a admission.Admission) Strategy {
	if a.Urgent {
		return Automatic{}
	}
	if !a.ScheduledAt.IsZero() {
		return Scheduled{}
	}
	return Manual{}
}
