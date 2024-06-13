package common

import (
	"errors"
	pkdb "github.com/hopeio/cherry/initialize/conf_dao/gormdb"
	"github.com/hopeio/cherry/initialize/conf_dao/gormdb/mysql"
	"github.com/hopeio/cherry/initialize/conf_dao/gormdb/postgres"
	"github.com/hopeio/cherry/initialize/conf_dao/gormdb/sqlite"
	"github.com/hopeio/cherry/initialize/initconf"
	dbi "github.com/hopeio/cherry/utils/dao/database"
	"gorm.io/gorm"
)

// Deprecated 每个驱动分开，不然每次都要编译所有驱动
type Config pkdb.Config

func (c *Config) InitBeforeInjectWithInitConfig(conf *initconf.InitConfig) {
	(*pkdb.Config)(c).InitBeforeInjectWithInitConfig(conf)
}

func (c *Config) Build() (*gorm.DB, error) {
	(*pkdb.Config)(c).Init()
	if c.Type == dbi.Mysql {
		return (*mysql.Config)(c).Build()
	} else if c.Type == dbi.Postgres {
		return (*postgres.Config)(c).Build()
	} else if c.Type == dbi.Sqlite {
		return (*sqlite.Config)(c).Build()
	}

	return nil, errors.New("只支持" + dbi.Mysql + "," + dbi.Postgres + "." + dbi.Sqlite)
}

type DB pkdb.DB

func (db *DB) Config() any {
	return (*Config)(&db.Conf)
}

func (db *DB) Set() {
	var err error
	db.DB, err = (*Config)(&db.Conf).Build()
	return err
}

func (db *DB) Close() error {
	return nil
}
