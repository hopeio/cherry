package badger

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/hopeio/cherry/utils/log"
)

type Config badger.Options

func (c *Config) InitBeforeInject() {
}
func (c *Config) Init() {

}
func (c *Config) Build() *badger.DB {
	c.Init()
	db, err := badger.Open(badger.Options(*c))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type DB struct {
	*badger.DB
	Conf Config
}

func (c *DB) Config() any {
	return &c.Conf
}

func (c *DB) Set() {
	c.DB = c.Conf.Build()
}

func (c *DB) Close() error {
	return c.DB.Close()
}
