package admission

import (
	"errors"
	"time"
)

type Admission struct {
	IIN                string
	Service            string
	RequestedMedicine  string
	RequestedEquipment string
	Urgent             bool
	ScheduledAt        time.Time
}

type Result struct {
	IIN         string
	Status      string
	Steps       []string
	ScheduledAt time.Time
}

const (
	Completed       = "COMPLETED"
	PendingApproval = "PENDING_APPROVAL"
	Scheduled       = "SCHEDULED"
)

func (a Admission) Validate() error {
	if len(a.IIN) != 12 {
		return errors.New("IIN must contain 12 digits")
	}
	for _, digit := range a.IIN {
		if digit < '0' || digit > '9' {
			return errors.New("IIN must contain 12 digits")
		}
	}
	if a.Service == "" || a.RequestedMedicine == "" || a.RequestedEquipment == "" {
		return errors.New("service, medicine, and equipment are required")
	}
	return nil
}
