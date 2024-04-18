package gormi

import (
	"github.com/hopeio/cherry/utils/configor"
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
	Type, Charset, Database, TimeZone string
	Host                              string `flag:"name:db_host;usage:数据库host"`
	Port                              int32
	User, Password                    string
	Postgres                          Postgres
	Mysql                             Mysql
	Sqlite                            Sqlite
	MaxIdleConns, MaxOpenConns        int
	ConnMaxLifetime, ConnMaxIdleTime  time.Duration

	GormConfig gorm.Config

	EnableStdLogger bool
	Logger          logger.Config

	NamingStrategy schema.NamingStrategy

	EnablePrometheus  bool
	Prometheus        prometheus.Config
	PrometheusConfigs []PrometheusConfig
}

type Postgres struct {
	Schema  string
	SSLMode string
}

type Mysql struct {
	ParseTime string
	Loc       string
}

type Sqlite struct {
	Path string
}

type PrometheusConfig struct {
	Prefix        string
	Interval      uint32
	VariableNames []string
}

func (c *Config) Init() {
	if c.Type == "" {
		c.Type = POSTGRES
	}
	configor.DurationNotify(c.Logger.SlowThreshold, 10*time.Millisecond)
	if c.TimeZone == "" {
		c.TimeZone = "Asia/Shanghai"
	}
	if c.Postgres.SSLMode == "" {
		c.Postgres.SSLMode = "disable"
	}
	if c.Mysql.Loc == "" {
		c.Mysql.Loc = "Local"
	}
	if c.Mysql.ParseTime == "" {
		c.Mysql.ParseTime = "True"
	}
	if c.Charset == "" {
		if c.Type == MYSQL {
			c.Charset = "utf8mb4"
		}
		if c.Type == POSTGRES {
			c.Charset = "utf8"
		}

	}

	if c.Port == 0 {
		if c.Type == MYSQL {
			c.Port = 3306
		}
		if c.Type == POSTGRES {
			c.Port = 5432
		}
	}

	if c.Sqlite.Path == "" {
		c.Sqlite.Path = "./sqlite.db"
	}
}
