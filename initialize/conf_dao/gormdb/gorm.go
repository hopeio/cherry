package gormdb

import (
	initialize2 "github.com/hopeio/cherry/initialize"
	gormi "github.com/hopeio/cherry/utils/dao/db/gorm"
	loggeri "github.com/hopeio/cherry/utils/dao/db/gorm/logger"
	"github.com/hopeio/cherry/utils/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
	stdlog "log"
	"os"
)

type Config gormi.Config

func (c *Config) InitBeforeInject() {
	c.EnableStdLogger = initialize2.GlobalConfig().Debug
}

func (c *Config) InitAfterInject() {
	(*gormi.Config)(c).Init()
}

func (c *Config) Build(dialector gorm.Dialector) *gorm.DB {

	dbConfig := &c.Gorm
	dbConfig.NamingStrategy = c.NamingStrategy

	// 日志
	if c.EnableStdLogger {
		// 默认日志
		logger.Default = logger.New(stdlog.New(os.Stdout, "\r", stdlog.LstdFlags), c.Logger)
	} else {
		logger.Default = &loggeri.Logger{Logger: log.GetCallerSkipLogger(2).Logger, Config: &c.Logger}
	}

	db, err := gorm.Open(dialector, dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	if c.EnablePrometheus {
		if c.Type == gormi.MYSQL {
			for _, pc := range c.PrometheusConfigs {
				c.Prometheus.MetricsCollector = append(c.Prometheus.MetricsCollector, &prometheus.MySQL{
					Prefix:        pc.Prefix,
					Interval:      pc.Interval,
					VariableNames: pc.VariableNames,
				})
			}

		}
		if c.Type == gormi.POSTGRES {
			for _, pc := range c.PrometheusConfigs {
				c.Prometheus.MetricsCollector = append(c.Prometheus.MetricsCollector, &prometheus.Postgres{
					Prefix:        pc.Prefix,
					Interval:      pc.Interval,
					VariableNames: pc.VariableNames,
				})
			}
		}
		err = db.Use(prometheus.New(c.Prometheus))
		if err != nil {
			log.Fatal(err)
		}
	}

	rawDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	rawDB.SetMaxIdleConns(c.MaxIdleConns)
	rawDB.SetMaxOpenConns(c.MaxOpenConns)
	rawDB.SetConnMaxLifetime(c.ConnMaxLifetime)
	rawDB.SetConnMaxIdleTime(c.ConnMaxIdleTime)
	return db
}

type DB struct {
	*gorm.DB
	Conf Config
}

func (db *DB) Table(name string) *gorm.DB {
	gdb := db.DB.Clauses()
	gdb.Statement.TableExpr = &clause.Expr{SQL: gdb.Statement.Quote(name)}
	gdb.Statement.Table = name
	return gdb
}
