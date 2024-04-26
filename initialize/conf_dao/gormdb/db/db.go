package db

import (
	pkdb "github.com/hopeio/cherry/initialize/conf_dao/gormdb"
	"github.com/hopeio/cherry/initialize/conf_dao/gormdb/mysql"
	"github.com/hopeio/cherry/initialize/conf_dao/gormdb/postgres"
	"github.com/hopeio/cherry/initialize/conf_dao/gormdb/sqlite"
	dbi "github.com/hopeio/cherry/utils/dao/db"
	"github.com/hopeio/cherry/utils/log"

	"gorm.io/gorm"
)

// Deprecated 每个驱动分开，不然每次都要编译所有驱动
type Config pkdb.Config

func (c *Config) InitBeforeInject() {
	(*pkdb.Config)(c).InitBeforeInject()
}

func (c *Config) Build() *gorm.DB {
	(*pkdb.Config)(c).Init()
	if c.Type == dbi.Mysql {
		(*mysql.Config)(c).Build()
	} else if c.Type == dbi.Postgres {
		(*postgres.Config)(c).Build()
	} else if c.Type == dbi.Sqlite {
		(*sqlite.Config)(c).Build()
	}

	log.Fatal("只支持" + dbi.Mysql + "," + dbi.Postgres + "." + dbi.Sqlite)
	return nil
}

type DB pkdb.DB

func (db *DB) Config() any {
	return (*Config)(&db.Conf)
}

func (db *DB) SetEntity() {
	db.DB = (*Config)(&db.Conf).Build()
}

func (db *DB) Close() error {
	return nil
}
