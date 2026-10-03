package facades

import (
	"database/sql"
	"errors"
	"fmt"

	"hospitaladmission/internal/equipmentreserver"
	"hospitaladmission/internal/insuranceverifier"
	"hospitaladmission/internal/pharmacydispatcher"
)

type IntakeFacade struct {
	Gateway   Gateway
	Insurance insuranceverifier.Verifier
	Equipment equipmentreserver.Reserver
	Pharmacy  pharmacydispatcher.Dispatcher
}

func (f IntakeFacade) Start(a Admission) error {
	if a.ID == "" {
		return errors.New("admission ID is required")
	}
	if err := f.Insurance.Rule.Check(a.IIN, a.Coverage); err != nil {
		return err
	}
	if err := f.Equipment.Rule.Check(a.Equipment); err != nil {
		return err
	}
	if err := f.Pharmacy.Rule.Check(a.Medicine); err != nil {
		return err
	}
	if err := f.Gateway.Create(a); err != nil {
		return fmt.Errorf("create admission: %w", err)
	}
	return f.ensureA(a)
}
func (f IntakeFacade) Admit(a Admission) error {
	if err := f.Start(a); err != nil {
		return err
	}
	return f.Resume(a.ID)
}
func (f IntakeFacade) ensureA(a Admission) error {
	v, err := f.Insurance.Repository.Get(a.ID)
	if err == sql.ErrNoRows {
		if err = f.Insurance.Verify(a.ID, a.IIN, a.Coverage); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if v.Status != insuranceverifier.Verified {
		return fmt.Errorf("insurance state %s cannot continue", v.Status)
	}
	return f.Gateway.Update(a.ID, PendingB, "")
}
func (f IntakeFacade) Resume(id string) error {
	a, err := f.Gateway.Get(id)
	if err != nil {
		return err
	}
	switch a.Status {
	case Completed:
		return nil
	case Compensated:
		return errors.New("admission was compensated")
	case Compensating:
		return f.compensate(a.ID, a.Reason)
	case PendingA, PendingB, PendingC:
	default:
		return fmt.Errorf("unknown admission status %s", a.Status)
	}
	if a.Status == PendingA {
		if err = f.ensureA(a); err != nil {
			return err
		}
	}
	v, err := f.Equipment.Repository.Get(id)
	if err == sql.ErrNoRows {
		if err = f.Equipment.Reserve(id, a.Equipment); err != nil {
			return f.fail(id, err)
		}
	} else if err != nil {
		return err
	} else if v.Status != equipmentreserver.Reserved {
		return fmt.Errorf("equipment state %s cannot continue", v.Status)
	}
	if err = f.Gateway.Update(id, PendingC, ""); err != nil {
		return err
	}
	p, err := f.Pharmacy.Repository.Get(id)
	if err == sql.ErrNoRows {
		if err = f.Pharmacy.Dispatch(id, a.Medicine); err != nil {
			return f.fail(id, err)
		}
	} else if err != nil {
		return err
	} else if p.Status != pharmacydispatcher.Dispatched {
		return fmt.Errorf("pharmacy state %s cannot continue", p.Status)
	}
	return f.Gateway.Update(id, Completed, "")
}
func (f IntakeFacade) fail(id string, cause error) error {
	if err := f.compensate(id, cause.Error()); err != nil {
		return fmt.Errorf("%v; compensation failed: %w", cause, err)
	}
	return cause
}
func (f IntakeFacade) compensate(id, reason string) error {
	if err := f.Gateway.Update(id, Compensating, reason); err != nil {
		return err
	}
	b, err := f.Equipment.Repository.Get(id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil && b.Status == equipmentreserver.Reserved {
		if err = f.Equipment.Release(id); err != nil {
			return err
		}
	}
	a, err := f.Insurance.Repository.Get(id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil && a.Status == insuranceverifier.Verified {
		if err = f.Insurance.Reverse(id); err != nil {
			return err
		}
	}
	return f.Gateway.Update(id, Compensated, reason)
}
