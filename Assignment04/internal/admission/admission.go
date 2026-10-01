package admission

import (
	"database/sql"
	"errors"
	"fmt"

	"hospitaladmission/internal/equipment"
	"hospitaladmission/internal/insurance"
	"hospitaladmission/internal/pharmacy"
)

const (
	PendingA     = "PENDING_A"
	PendingB     = "PENDING_B"
	PendingC     = "PENDING_C"
	Completed    = "COMPLETED"
	Compensating = "COMPENSATING"
	Compensated  = "COMPENSATED"
)

type Admission struct{ ID, IIN, Policy, Equipment, Medicine, Status, Reason string }

type Gateway struct{ DB *sql.DB }

func (g Gateway) Create(a Admission) error {
	_, err := g.DB.Exec(`INSERT INTO admissions(id,iin,policy,equipment,medicine,status,reason) VALUES(?,?,?,?,?,?,?)`, a.ID, a.IIN, a.Policy, a.Equipment, a.Medicine, PendingA, "")
	return err
}
func (g Gateway) Get(id string) (Admission, error) {
	var a Admission
	err := g.DB.QueryRow(`SELECT id,iin,policy,equipment,medicine,status,reason FROM admissions WHERE id=?`, id).Scan(&a.ID, &a.IIN, &a.Policy, &a.Equipment, &a.Medicine, &a.Status, &a.Reason)
	return a, err
}
func (g Gateway) Update(id, status, reason string) error {
	_, err := g.DB.Exec(`UPDATE admissions SET status=?,reason=? WHERE id=?`, status, reason, id)
	return err
}

// IntakeFacade is the only object a normal client needs for the chain.
type IntakeFacade struct {
	Gateway   Gateway
	Insurance insurance.Verifier
	Equipment equipment.Reserver
	Pharmacy  pharmacy.Dispatcher
}

func (f IntakeFacade) Start(a Admission) error {
	if a.ID == "" {
		return errors.New("admission ID is required")
	}
	if err := f.Insurance.Rule.Check(a.IIN, a.Policy); err != nil {
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
		if err = f.Insurance.Verify(a.ID, a.IIN, a.Policy); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if v.Status != insurance.Verified {
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
	} else if v.Status != equipment.Reserved {
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
	} else if p.Status != pharmacy.Dispatched {
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
	// Read the committed rows. A restart or a status update failure may leave the
	// progress marker behind the actual subsystem state.
	b, err := f.Equipment.Repository.Get(id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil && b.Status == equipment.Reserved {
		if err = f.Equipment.Release(id); err != nil {
			return err
		}
	}
	a, err := f.Insurance.Repository.Get(id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil && a.Status == insurance.Verified {
		if err = f.Insurance.Reverse(id); err != nil {
			return err
		}
	}
	return f.Gateway.Update(id, Compensated, reason)
}

type View struct {
	Admission                      Admission
	Insurance, Equipment, Pharmacy string
}

// AdminFacade gives operators inspection and manual compensation only.
type AdminFacade struct{ Intake IntakeFacade }

func (f AdminFacade) Inspect(id string) (View, error) {
	a, err := f.Intake.Gateway.Get(id)
	if err != nil {
		return View{}, err
	}
	v := View{Admission: a, Insurance: "NONE", Equipment: "NONE", Pharmacy: "NONE"}
	i, e := f.Intake.Insurance.Repository.Get(id)
	if e != nil && e != sql.ErrNoRows {
		return View{}, e
	}
	if e == nil {
		v.Insurance = i.Status
	}
	b, e := f.Intake.Equipment.Repository.Get(id)
	if e != nil && e != sql.ErrNoRows {
		return View{}, e
	}
	if e == nil {
		v.Equipment = b.Status
	}
	p, e := f.Intake.Pharmacy.Repository.Get(id)
	if e != nil && e != sql.ErrNoRows {
		return View{}, e
	}
	if e == nil {
		v.Pharmacy = p.Status
	}
	return v, nil
}
func (f AdminFacade) Compensate(id string) error {
	v, err := f.Inspect(id)
	if err != nil {
		return err
	}
	if v.Admission.Status == Completed {
		return errors.New("completed admission cannot be manually compensated")
	}
	if v.Pharmacy == pharmacy.Dispatched {
		return errors.New("dispatched medicine cannot be manually compensated")
	}
	if v.Admission.Status == Compensated {
		return nil
	}
	return f.Intake.compensate(id, "manual compensation")
}
