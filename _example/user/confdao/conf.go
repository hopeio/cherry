package confdao

import (
	"database/sql"
	"fmt"
	"github.com/hopeio/cherry/initialize/conf_dao/gormdb/postgres"
	"github.com/hopeio/cherry/initialize/conf_dao/log"
	"github.com/hopeio/cherry/initialize/conf_dao/mail"
	"github.com/hopeio/cherry/initialize/conf_dao/pebble"
	"github.com/hopeio/cherry/initialize/conf_dao/redis"
	"github.com/hopeio/cherry/initialize/conf_dao/ristretto"
	"github.com/hopeio/cherry/initialize/conf_dao/server"
	"github.com/hopeio/cherry/utils/io/fs"
	timei "github.com/hopeio/cherry/utils/time"
	"runtime"
	"time"
)

var (
	Conf      = &config{}
	Dao  *dao = &dao{}
)

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    server.Config
	Log       log.Config
}

type serverConfig struct {
	Volume fs.Dir

	PassSalt string
	// 天数
	TokenMaxAge time.Duration
	TokenSecret string
	PageSize    int8

	LuosimaoSuperPW   string
	LuosimaoVerifyURL string
	LuosimaoAPIKey    string

	QrCodeSaveDir fs.Dir //二维码保存路径
	PrefixUrl     string
	FontSaveDir   fs.Dir //字体保存路径

}

func (c *config) InitBeforeInject() {
	c.Customize.TokenMaxAge = timei.Day
}

func (c *config) InitAfterInject() {
	if runtime.GOOS == "windows" {
	}

	c.Customize.TokenMaxAge = timei.StdDuration(c.Customize.TokenMaxAge, time.Hour)
}

// dao dao.
type dao struct {
	// GORMDB 数据库连接
	GORMDB   postgres.DB
	StdDB    *sql.DB
	PebbleDB pebble.DB
	// RedisPool Redis连接池
	Redis redis.Redis
	Cache ristretto.Cache
	//elastic
	Mail mail.Mail `init:"config:Mail"`
}

func (d *dao) InitBeforeInject() {
	d.GORMDB.Conf.GormConfig.NowFunc = time.Now
}

func (d *dao) InitAfterInjectConfig() {
	fmt.Println("这里后执行")
}

func (d *dao) InitAfterInject() {
	db := d.GORMDB
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")

	d.StdDB, _ = db.DB.DB()
}
