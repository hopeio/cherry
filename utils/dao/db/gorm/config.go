package gormi

import (
	"gorm.io/gorm/schema"
	"gorm.io/plugin/prometheus"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
	SQLite   = "sqlite3"
)

type Config struct {
	Type, Charset, Database, Schema, TimeZone string
	Host                                      string `flag:"name:db_host;usage:数据库host"`
	Port                                      int32
	User, Password                            string
	TimeFormat                                string
	MaxIdleConns, MaxOpenConns                int

	GormConfig gorm.Config

	EnableStdLogger bool
	Logger          logger.Config

	NamingStrategy schema.NamingStrategy

	EnablePrometheus  bool
	Prometheus        prometheus.Config
	PrometheusConfigs []PrometheusConfig
}

type PrometheusConfig struct {
	Prefix        string
	Interval      uint32
	VariableNames []string
}

func (c *Config) Init() {
	if c.Logger.SlowThreshold < 10*time.Millisecond {
		c.Logger.SlowThreshold *= time.Millisecond
	}
	if c.TimeFormat == "" {
		c.TimeZone = "Asia/Shanghai"
	}
	if c.TimeFormat == "" {
		c.TimeFormat = "2006-01-02 15:04:05"
	}
	if c.Charset == "" {
		c.Charset = "utf8"
	}
	if c.Type == "" {
		c.Type = POSTGRES
	}
	if c.Port == 0 {
		if c.Type == MYSQL {
			c.Port = 3306
		}
		if c.Type == POSTGRES {
			c.Port = 5432
		}
	}
}
