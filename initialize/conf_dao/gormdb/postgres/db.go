package postgres

import (
	"fmt"
	pkdb "github.com/hopeio/cherry/initialize/conf_dao/gormdb"
	"github.com/hopeio/cherry/initialize/initconf"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Config pkdb.Config

func (c *Config) InitBeforeInjectWithInitConfig(conf *initconf.InitConfig) {
	(*pkdb.Config)(c).InitBeforeInjectWithInitConfig(conf)
}

func (c *Config) Build() (*gorm.DB, error) {
	(*pkdb.Config)(c).Init()
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%d sslmode=%s password=%s TimeZone=%s",
		c.Host, c.User, c.Database, c.Port, c.Postgres.SSLMode, c.Password, c.TimeZone)
	return (*pkdb.Config)(c).Build(postgres.Open(dsn))
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
