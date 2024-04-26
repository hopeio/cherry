package sqlite

import (
	pkdb "github.com/hopeio/cherry/initialize/conf_dao/gormdb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config pkdb.Config

func (c *Config) InitBeforeInject() {
	(*pkdb.Config)(c).InitBeforeInject()
}

func (c *Config) Build() *gorm.DB {
	(*pkdb.Config)(c).Init()
	return (*pkdb.Config)(c).Build(sqlite.Open(c.Sqlite.Path))
}

type DB pkdb.DB

func (db *DB) Config() any {
	return (*Config)(&db.Conf)
}

func (db *DB) SetEntity(entity interface{}) {
	db.DB = (*Config)(&db.Conf).Build()
}

func (db *DB) Close() error {
	return nil
}
