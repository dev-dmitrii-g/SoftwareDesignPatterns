package insurance

import (
	"database/sql"
	"errors"
	"fmt"
)

const (
	New      = "NEW"
	Verified = "VERIFIED"
	Reversed = "REVERSED"
)

type Verification struct{ AdmissionID, IIN, Policy, Status string }

func (v *Verification) Transition(to string) error {
	if (v.Status == New && to == Verified) || (v.Status == Verified && to == Reversed) {
		v.Status = to
		return nil
	}
	return fmt.Errorf("insurance: illegal transition %s -> %s", v.Status, to)
}

type Rule struct{}

func (Rule) Check(iin, policy string) error {
	if len(iin) != 12 {
		return errors.New("IIN must contain 12 digits")
	}
	for _, r := range iin {
		if r < '0' || r > '9' {
			return errors.New("IIN must contain 12 digits")
		}
	}
	if policy == "" {
		return errors.New("policy number is required")
	}
	return nil
}

type Repository struct{ DB *sql.DB }

func (r Repository) Save(v Verification) error {
	_, err := r.DB.Exec(`INSERT INTO insurance(admission_id,iin,policy,status) VALUES(?,?,?,?)
	ON CONFLICT(admission_id) DO UPDATE SET status=excluded.status`, v.AdmissionID, v.IIN, v.Policy, v.Status)
	return err
}
func (r Repository) Get(id string) (Verification, error) {
	var v Verification
	err := r.DB.QueryRow(`SELECT admission_id,iin,policy,status FROM insurance WHERE admission_id=?`, id).Scan(&v.AdmissionID, &v.IIN, &v.Policy, &v.Status)
	return v, err
}

type Verifier struct {
	Rule       Rule
	Repository Repository
}

func (s Verifier) Verify(id, iin, policy string) error {
	if err := s.Rule.Check(iin, policy); err != nil {
		return err
	}
	v := Verification{id, iin, policy, New}
	if err := v.Transition(Verified); err != nil {
		return err
	}
	return s.Repository.Save(v)
}
func (s Verifier) Reverse(id string) error {
	v, err := s.Repository.Get(id)
	if err != nil {
		return err
	}
	if err = v.Transition(Reversed); err != nil {
		return err
	}
	return s.Repository.Save(v)
}
