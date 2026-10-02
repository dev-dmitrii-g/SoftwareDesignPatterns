package store

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestInitialStockAndReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hospital.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	var quantity int
	if err := db.QueryRow(`SELECT quantity FROM stock WHERE kind='equipment' AND name='wheelchair'`).Scan(&quantity); err != nil {
		t.Fatal(err)
	}
	if quantity != 5 {
		t.Fatalf("wheelchair stock = %d, want 5", quantity)
	}
	if _, err := db.Exec(`UPDATE stock SET quantity=0 WHERE kind='equipment' AND name='wheelchair'`); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.QueryRow(`SELECT quantity FROM stock WHERE kind='equipment' AND name='wheelchair'`).Scan(&quantity); err != nil {
		t.Fatal(err)
	}
	if quantity != 0 {
		t.Fatalf("stock reset on reopen: %d", quantity)
	}
}

func TestOldPolicyColumnIsRenamed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "old.db")
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE insurance(admission_id TEXT PRIMARY KEY, iin TEXT, policy TEXT, status TEXT);
	CREATE TABLE admissions(id TEXT PRIMARY KEY, iin TEXT, policy TEXT, equipment TEXT, medicine TEXT, status TEXT, reason TEXT);
	INSERT INTO insurance VALUES('A1','990101123456','OLD-POL','VERIFIED');
	INSERT INTO admissions VALUES('A1','990101123456','OLD-POL','wheelchair','paracetamol','COMPLETED','');`)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, table := range []string{"insurance", "admissions"} {
		var coverage string
		if err := db.QueryRow("SELECT coverage FROM " + table + " WHERE " + map[string]string{"insurance": "admission_id", "admissions": "id"}[table] + "='A1'").Scan(&coverage); err != nil {
			t.Fatal(err)
		}
		if coverage != "OLD-POL" {
			t.Fatalf("%s coverage = %q", table, coverage)
		}
	}
}
