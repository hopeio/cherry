package buntdb

import (
	"github.com/hopeio/cherry/utils/log"
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

func (c *Config) Build() *buntdb.DB {
	db, err := buntdb.Open(c.Path)
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetConfig(c.Config)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type DB struct {
	*buntdb.DB
	Conf Config
}

func (m *DB) Config() any {
	return &m.Conf
}

func (m *DB) Set() {
	m.DB = m.Conf.Build()
}

func (m *DB) Close() error {
	return nil
}
