package store

import (
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
