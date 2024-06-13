package duckdb

import "testing"

func TestDuckDB(t *testing.T) {
	config := Config{
		DSN:         "./duck.db?access_mode=read_only&threads=4",
		Path:        "",
		AccessMode:  "",
		Threads:     0,
		BootQueries: nil,
	}
	db, err := config.Build()
	if err != nil {
		t.Error("Build err", err)
	}
	_, err = db.Exec(`CREATE TABLE people (id INTEGER, name VARCHAR)`)
	if err != nil {
		t.Error("Exec err", err)
	}
}
