package pebble

import (
	"github.com/cockroachdb/pebble"
	"github.com/hopeio/cherry/utils/log"
)

type Config struct {
	DirName string
	pebble.Options
}

func (c *Config) InitBeforeInject() {
}
func (c *Config) InitAfterInject() {
	if c.DirName == "" {
		log.Fatal("pebble config not set dirname")
	}
}

func (c *Config) Build() *pebble.DB {
	c.InitAfterInject()
	db, err := pebble.Open(c.DirName, &c.Options)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type DB struct {
	*pebble.DB
	Conf Config
}

func (p *DB) Config() any {
	return &p.Conf
}

func (p *DB) SetEntity(entity interface{}) {
	if client, ok := entity.(*pebble.DB); ok {
		p.DB = client
	}
}

func (p *DB) Close() error {
	if p.DB == nil {
		return nil
	}
	return p.DB.Close()
}
