package mysql

import (
	"fmt"
	pkdb "github.com/hopeio/cherry/initialize/conf_dao/gormdb"
	"github.com/hopeio/cherry/initialize/initconf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config pkdb.Config

func (c *Config) InitBeforeInjectWithInitConfig(conf *initconf.InitConfig) {
	(*pkdb.Config)(c).InitBeforeInjectWithInitConfig(conf)
}

func (c *Config) Build() *gorm.DB {
	(*pkdb.Config)(c).Init()
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		c.User, c.Password, c.Host,
		c.Port, c.Database, c.Charset, c.Mysql.ParseTime, c.Mysql.Loc)
	return (*pkdb.Config)(c).Build(mysql.Open(url))
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
