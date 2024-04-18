package redis

import (
	"github.com/hopeio/cherry/utils/configor"
	"github.com/hopeio/cherry/utils/crypto/tls"
	"github.com/hopeio/cherry/utils/log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	redis.Options
	CertFile string `json:"cert_file,omitempty"`
	KeyFile  string `json:"key_file,omitempty"`
}

func (c *Config) InitBeforeInject() {
}

func (c *Config) InitAfterInject() {
	tlsConfig, err := tls.Certificate(c.CertFile, c.KeyFile)
	if err != nil {
		log.Fatal(err)
	}
	c.TLSConfig = tlsConfig
	configor.DurationNotify("IdleTimeout", c.IdleTimeout, time.Second)
}

func (c *Config) Build() *redis.Client {
	c.InitAfterInject()
	return redis.NewClient(&c.Options)
}

type Redis struct {
	*redis.Client
	Conf Config
}

func (db *Redis) Config() any {
	return &db.Conf
}

func (db *Redis) SetEntity() {
	db.Client = db.Conf.Build()
}

func (db *Redis) Close() error {
	if db.Client == nil {
		return nil
	}
	return db.Client.Close()
}
