package facades

import (
	"database/sql"
	"errors"

	"hospitaladmission/internal/subsystemc"
)

type View struct {
	Admission                      Admission
	Insurance, Equipment, Pharmacy string
}

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
	if v.Pharmacy == subsystemc.Dispatched {
		return errors.New("dispatched medicine cannot be manually compensated")
	}
	if v.Admission.Status == Compensated {
		return nil
	}
	return f.Intake.compensate(id, "manual compensation")
}
