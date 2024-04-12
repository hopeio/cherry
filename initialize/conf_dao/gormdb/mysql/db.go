package mysql

import (
	"fmt"
	pkdb "github.com/hopeio/cherry/initialize/conf_dao/gormdb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config pkdb.Config

func (c *Config) InitBeforeInject() {
	(*pkdb.Config)(c).InitBeforeInject()
}
func (c *Config) InitAfterInject() {
	(*pkdb.Config)(c).InitAfterInject()
}

func (c *Config) Build() *gorm.DB {
	c.InitAfterInject()
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host,
		c.Port, c.Database, c.Charset)
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
