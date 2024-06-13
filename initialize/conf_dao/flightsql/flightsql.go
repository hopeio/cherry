package flightsql

import (
	"database/sql"
	_ "github.com/apache/arrow-adbc/go/adbc/sqldriver/flightsql"
	"github.com/hopeio/cherry/utils/log"
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

func (c *Config) Build() *sql.DB {
	db, err := sql.Open("flightsql", c.DNS)
	if err != nil {
		log.Fatal(err)
	}
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
