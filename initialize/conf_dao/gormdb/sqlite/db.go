package sqlite

import (
	pkdb "github.com/hopeio/cherry/initialize/conf_dao/gormdb"
	"github.com/hopeio/cherry/initialize/initconf"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config pkdb.Config

func (c *Config) InitBeforeInjectWithInitConfig(conf *initconf.InitConfig) {
	(*pkdb.Config)(c).InitBeforeInjectWithInitConfig(conf)
}

func (c *Config) Build() (*gorm.DB, error) {
	(*pkdb.Config)(c).Init()
	return (*pkdb.Config)(c).Build(sqlite.Open(c.Sqlite.DSN))
}

type DB pkdb.DB

func (db *DB) Config() any {
	return (*Config)(&db.Conf)
}

func (db *DB) Set() error {
	var err error
	db.DB, err = (*Config)(&db.Conf).Build()
	return err
}

func (db *DB) Close() error {
	return nil
}
