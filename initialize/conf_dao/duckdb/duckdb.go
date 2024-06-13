package duckdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/hopeio/cherry/utils/log"
	"github.com/marcboeker/go-duckdb"
)

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

func (c *Config) Build() *sql.DB {
	connector, err := duckdb.NewConnector(c.DSN, func(execer driver.ExecerContext) error {
		for _, query := range c.BootQueries {
			_, err := execer.ExecContext(context.Background(), query.Query, query.Args)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	db := sql.OpenDB(connector)
	return db
}

type DB struct {
	*sql.DB
	Conf Config
}

func (m *DB) Config() any {
	return &m.Conf
}

func (m *DB) Set() {
	m.DB = m.Conf.Build()
}

func (m *DB) Close() error {
	return m.DB.Close()
}
