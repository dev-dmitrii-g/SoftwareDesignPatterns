package equipmentreserver

import (
	"database/sql"
	"errors"
	"fmt"
)

const (
	New      = "NEW"
	Reserved = "RESERVED"
	Released = "RELEASED"
)

type Reservation struct{ AdmissionID, Item, Status string }

func (v *Reservation) Transition(to string) error {
	if (v.Status == New && to == Reserved) || (v.Status == Reserved && to == Released) {
		v.Status = to
		return nil
	}
	return fmt.Errorf("equipment: illegal transition %s -> %s", v.Status, to)
}

type Rule struct{}

func (Rule) Check(item string) error {
	if item == "" {
		return errors.New("equipment name is required")
	}
	return nil
}

type Repository struct{ DB *sql.DB }

func (r Repository) Get(id string) (Reservation, error) {
	var v Reservation
	err := r.DB.QueryRow(`SELECT admission_id,item,status FROM equipment WHERE admission_id=?`, id).Scan(&v.AdmissionID, &v.Item, &v.Status)
	return v, err
}
func (r Repository) SetStock(item string, n int) error {
	if n < 0 {
		return errors.New("stock cannot be negative")
	}
	_, err := r.DB.Exec(`INSERT INTO stock(kind,name,quantity) VALUES('equipment',?,?) ON CONFLICT(kind,name) DO UPDATE SET quantity=excluded.quantity`, item, n)
	return err
}
func (r Repository) Reserve(v Reservation) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := tx.Exec(`UPDATE stock SET quantity=quantity-1 WHERE kind='equipment' AND name=? AND quantity>0`, v.Item)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("equipment %q unavailable", v.Item)
	}
	if _, err = tx.Exec(`INSERT INTO equipment(admission_id,item,status) VALUES(?,?,?)`, v.AdmissionID, v.Item, v.Status); err != nil {
		return err
	}
	return tx.Commit()
}
func (r Repository) Release(v Reservation) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := tx.Exec(`UPDATE equipment SET status=? WHERE admission_id=? AND status=?`, Released, v.AdmissionID, Reserved)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("reservation is not active")
	}
	if _, err = tx.Exec(`UPDATE stock SET quantity=quantity+1 WHERE kind='equipment' AND name=?`, v.Item); err != nil {
		return err
	}
	return tx.Commit()
}

type Reserver struct {
	Rule       Rule
	Repository Repository
}

func (s Reserver) Reserve(id, item string) error {
	if err := s.Rule.Check(item); err != nil {
		return err
	}
	v := Reservation{id, item, New}
	if err := v.Transition(Reserved); err != nil {
		return err
	}
	return s.Repository.Reserve(v)
}
func (s Reserver) Release(id string) error {
	v, err := s.Repository.Get(id)
	if err != nil {
		return err
	}
	if err = v.Transition(Released); err != nil {
		return err
	}
	return s.Repository.Release(v)
}
