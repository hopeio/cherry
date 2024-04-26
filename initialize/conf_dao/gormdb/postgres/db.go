package postgres

import (
	"fmt"
	pkdb "github.com/hopeio/cherry/initialize/conf_dao/gormdb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Config pkdb.Config

func (c *Config) InitBeforeInject() {
	(*pkdb.Config)(c).InitBeforeInject()
}

func (c *Config) Build() *gorm.DB {
	(*pkdb.Config)(c).Init()
	url := fmt.Sprintf("host=%s user=%s dbname=%s port=%d sslmode=%s password=%s TimeZone=%s",
		c.Host, c.User, c.Database, c.Port, c.Postgres.SSLMode, c.Password, c.TimeZone)
	return (*pkdb.Config)(c).Build(postgres.Open(url))
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

func (db *DB) Table(name string) *gorm.DB {
	name = db.TableName(name)
	gdb := db.DB.Clauses()
	gdb.Statement.TableExpr = &clause.Expr{SQL: gdb.Statement.Quote(name)}
	gdb.Statement.Table = name
	return gdb
}

func (db *DB) TableName(name string) string {
	if db.Conf.Postgres.Schema != "" {
		return db.Conf.Postgres.Schema + "." + name
	}
	return name
}

// TODO:
func (db *DB) Inject(configName string) {

}
