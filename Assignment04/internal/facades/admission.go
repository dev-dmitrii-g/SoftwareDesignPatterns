package facades

import "database/sql"

const (
	PendingA     = "PENDING_A"
	PendingB     = "PENDING_B"
	PendingC     = "PENDING_C"
	Completed    = "COMPLETED"
	Compensating = "COMPENSATING"
	Compensated  = "COMPENSATED"
)

type Admission struct{ ID, IIN, Coverage, Equipment, Medicine, Status, Reason string }

type Gateway struct{ DB *sql.DB }

func (g Gateway) Create(a Admission) error {
	_, err := g.DB.Exec(`INSERT INTO admissions(id,iin,coverage,equipment,medicine,status,reason) VALUES(?,?,?,?,?,?,?)`, a.ID, a.IIN, a.Coverage, a.Equipment, a.Medicine, PendingA, "")
	return err
}
func (g Gateway) Get(id string) (Admission, error) {
	var a Admission
	err := g.DB.QueryRow(`SELECT id,iin,coverage,equipment,medicine,status,reason FROM admissions WHERE id=?`, id).Scan(&a.ID, &a.IIN, &a.Coverage, &a.Equipment, &a.Medicine, &a.Status, &a.Reason)
	return a, err
}
func (g Gateway) List() ([]Admission, error) {
	rows, err := g.DB.Query(`SELECT id,iin,coverage,equipment,medicine,status,reason FROM admissions ORDER BY rowid DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Admission{}
	for rows.Next() {
		var a Admission
		if err := rows.Scan(&a.ID, &a.IIN, &a.Coverage, &a.Equipment, &a.Medicine, &a.Status, &a.Reason); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, rows.Err()
}
func (g Gateway) Update(id, status, reason string) error {
	_, err := g.DB.Exec(`UPDATE admissions SET status=?,reason=? WHERE id=?`, status, reason, id)
	return err
}
