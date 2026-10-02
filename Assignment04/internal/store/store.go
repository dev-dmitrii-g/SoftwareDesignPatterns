package store

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err = db.Exec(`PRAGMA busy_timeout = 5000;
	CREATE TABLE IF NOT EXISTS insurance (admission_id TEXT PRIMARY KEY, iin TEXT NOT NULL, coverage TEXT NOT NULL, status TEXT NOT NULL);
	CREATE TABLE IF NOT EXISTS equipment (admission_id TEXT PRIMARY KEY, item TEXT NOT NULL, status TEXT NOT NULL);
	CREATE TABLE IF NOT EXISTS pharmacy (admission_id TEXT PRIMARY KEY, medicine TEXT NOT NULL, status TEXT NOT NULL);
	CREATE TABLE IF NOT EXISTS stock (kind TEXT NOT NULL, name TEXT NOT NULL, quantity INTEGER NOT NULL CHECK(quantity >= 0), PRIMARY KEY(kind, name));
	CREATE TABLE IF NOT EXISTS admissions (id TEXT PRIMARY KEY, iin TEXT NOT NULL, coverage TEXT NOT NULL, equipment TEXT NOT NULL, medicine TEXT NOT NULL, status TEXT NOT NULL, reason TEXT NOT NULL DEFAULT '');
	INSERT OR IGNORE INTO stock(kind,name,quantity) VALUES('equipment','wheelchair',5);
	INSERT OR IGNORE INTO stock(kind,name,quantity) VALUES('pharmacy','paracetamol',5);`); err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize database: %w", err)
	}
	for _, table := range []string{"insurance", "admissions"} {
		if err := migrateCoverage(db, table); err != nil {
			db.Close()
			return nil, fmt.Errorf("migrate %s: %w", table, err)
		}
	}
	return db, nil
}

func migrateCoverage(db *sql.DB, table string) error {
	rows, err := db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		return err
	}
	hasPolicy, hasCoverage := false, false
	for rows.Next() {
		var position, required, primaryKey int
		var name, dataType string
		var defaultValue sql.NullString
		if err := rows.Scan(&position, &name, &dataType, &required, &defaultValue, &primaryKey); err != nil {
			rows.Close()
			return err
		}
		if name == "policy" {
			hasPolicy = true
		}
		if name == "coverage" {
			hasCoverage = true
		}
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	if hasPolicy && !hasCoverage {
		_, err = db.Exec("ALTER TABLE " + table + " RENAME COLUMN policy TO coverage")
		return err
	}
	return nil
}
