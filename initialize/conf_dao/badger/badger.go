package badger

import (
	"github.com/dgraph-io/badger/v3"
)

type Config badger.Options

func (c *Config) InitBeforeInject() {
}
func (c *Config) Init() {

}
func (c *Config) Build() (*badger.DB, error) {
	c.Init()
	return badger.Open(badger.Options(*c))
}

type DB struct {
	*badger.DB
	Conf Config
}

func (c *DB) Config() any {
	return &c.Conf
}

func (c *DB) Init() error {
	var err error
	c.DB, err = c.Conf.Build()
	return err
}

func (c *DB) Close() error {
	return c.DB.Close()
}
