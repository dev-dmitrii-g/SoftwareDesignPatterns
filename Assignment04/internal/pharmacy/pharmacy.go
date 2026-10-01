package pharmacy

import (
	"database/sql"
	"errors"
	"fmt"
)

const (
	New        = "NEW"
	Dispatched = "DISPATCHED"
	Cancelled  = "CANCELLED"
)

type Dispatch struct{ AdmissionID, Medicine, Status string }

func (v *Dispatch) Transition(to string) error {
	if (v.Status == New && to == Dispatched) || (v.Status == Dispatched && to == Cancelled) {
		v.Status = to
		return nil
	}
	return fmt.Errorf("pharmacy: illegal transition %s -> %s", v.Status, to)
}

type Rule struct{}

func (Rule) Check(medicine string) error {
	if medicine == "" {
		return errors.New("medicine name is required")
	}
	return nil
}

type Repository struct{ DB *sql.DB }

func (r Repository) Get(id string) (Dispatch, error) {
	var v Dispatch
	err := r.DB.QueryRow(`SELECT admission_id,medicine,status FROM pharmacy WHERE admission_id=?`, id).Scan(&v.AdmissionID, &v.Medicine, &v.Status)
	return v, err
}
func (r Repository) SetStock(medicine string, n int) error {
	if n < 0 {
		return errors.New("stock cannot be negative")
	}
	_, err := r.DB.Exec(`INSERT INTO stock(kind,name,quantity) VALUES('pharmacy',?,?) ON CONFLICT(kind,name) DO UPDATE SET quantity=excluded.quantity`, medicine, n)
	return err
}
func (r Repository) Dispatch(v Dispatch) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := tx.Exec(`UPDATE stock SET quantity=quantity-1 WHERE kind='pharmacy' AND name=? AND quantity>0`, v.Medicine)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("medicine %q unavailable", v.Medicine)
	}
	if _, err = tx.Exec(`INSERT INTO pharmacy(admission_id,medicine,status) VALUES(?,?,?)`, v.AdmissionID, v.Medicine, v.Status); err != nil {
		return err
	}
	return tx.Commit()
}

type Dispatcher struct {
	Rule       Rule
	Repository Repository
}

func (s Dispatcher) Dispatch(id, medicine string) error {
	if err := s.Rule.Check(medicine); err != nil {
		return err
	}
	v := Dispatch{id, medicine, New}
	if err := v.Transition(Dispatched); err != nil {
		return err
	}
	return s.Repository.Dispatch(v)
}
