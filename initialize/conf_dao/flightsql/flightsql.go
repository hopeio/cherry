package flightsql

import (
	"database/sql"
	_ "github.com/apache/arrow-adbc/go/adbc/sqldriver/flightsql"
	"github.com/tidwall/buntdb"
)

type Config struct {
	DNS string
	buntdb.Config
}

func (c *Config) InitBeforeInject() {

}

func (c *Config) Init() {
}

func (c *Config) Build() (*sql.DB, error) {
	return sql.Open("flightsql", c.DNS)
}

type DB struct {
	*sql.DB
	Conf Config
}

func (m *DB) Config() any {
	return &m.Conf
}

func (m *DB) Init() error {
	var err error
	m.DB, err = m.Conf.Build()
	return err
}

func (m *DB) Close() error {
	return m.DB.Close()
}
