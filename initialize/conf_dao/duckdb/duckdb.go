package duckdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/marcboeker/go-duckdb"
)

// https://github.com/marcboeker/go-duckdb/issues/115
// CGO_ENABLED=1 CGO_LDFLAGS="-L/path/to/duckdb.dll" go build -tags=duckdb_use_lib,go1.22 main.go
// LD_LIBRARY_PATH=/path/to/libs ./main
type Config struct {
	DSN         string
	Path        string
	AccessMode  string `json:"access_mode" comment:"read_only"`
	Threads     int    `json:"threads"`
	BootQueries []BootQuery
}
type BootQuery struct {
	Query string
	Args  []driver.NamedValue
}

func (c *Config) InitBeforeInject() {

}

func (c *Config) Init() {
}

func (c *Config) Build() (*sql.DB, error) {
	connector, err := duckdb.NewConnector(c.DSN, func(execer driver.ExecerContext) error {
		for _, query := range c.BootQueries {
			_, err := execer.ExecContext(context.Background(), query.Query, query.Args)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return sql.OpenDB(connector), err
}

type DB struct {
	*sql.DB
	Conf Config
}

func (m *DB) Config() any {
	return &m.Conf
}

func (m *DB) Set() error {
	var err error
	m.DB, err = m.Conf.Build()
	return err
}

func (m *DB) Close() error {
	return m.DB.Close()
}
