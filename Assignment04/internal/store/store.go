package store

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Open creates the small, durable database used by the three independent gateways.
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err = db.Exec(`PRAGMA busy_timeout = 5000;
	CREATE TABLE IF NOT EXISTS insurance (admission_id TEXT PRIMARY KEY, iin TEXT NOT NULL, policy TEXT NOT NULL, status TEXT NOT NULL);
	CREATE TABLE IF NOT EXISTS equipment (admission_id TEXT PRIMARY KEY, item TEXT NOT NULL, status TEXT NOT NULL);
	CREATE TABLE IF NOT EXISTS pharmacy (admission_id TEXT PRIMARY KEY, medicine TEXT NOT NULL, status TEXT NOT NULL);
	CREATE TABLE IF NOT EXISTS stock (kind TEXT NOT NULL, name TEXT NOT NULL, quantity INTEGER NOT NULL CHECK(quantity >= 0), PRIMARY KEY(kind, name));
	CREATE TABLE IF NOT EXISTS admissions (id TEXT PRIMARY KEY, iin TEXT NOT NULL, policy TEXT NOT NULL, equipment TEXT NOT NULL, medicine TEXT NOT NULL, status TEXT NOT NULL, reason TEXT NOT NULL DEFAULT '');`); err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize database: %w", err)
	}
	return db, nil
}
