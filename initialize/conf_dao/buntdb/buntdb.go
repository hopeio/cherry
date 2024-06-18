package buntdb

import (
	"github.com/tidwall/buntdb"
)

type Config struct {
	Path string
	buntdb.Config
}

func (c *Config) InitBeforeInject() {

}

func (c *Config) Init() {
}

func (c *Config) Build() (*buntdb.DB, error) {
	db, err := buntdb.Open(c.Path)
	if err != nil {
		return nil, err
	}
	return db, db.SetConfig(c.Config)
}

type DB struct {
	*buntdb.DB
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
