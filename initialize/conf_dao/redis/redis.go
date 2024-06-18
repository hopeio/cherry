package redis

import (
	"context"
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

func (c *Config) Init() {
	tlsConfig, err := tls.NewServerTLSConfig(c.CertFile, c.KeyFile)
	if err != nil {
		log.Fatal(err)
	}
	c.TLSConfig = tlsConfig
	configor.DurationNotify("IdleTimeout", c.IdleTimeout, time.Second)
}

func (c *Config) Build() (*redis.Client, error) {
	c.Init()
	client := redis.NewClient(&c.Options)
	return client, client.Ping(context.Background()).Err()
}

type Client struct {
	*redis.Client
	Conf Config
}

func (db *Client) Config() any {
	return &db.Conf
}

func (db *Client) Init() error {
	var err error
	db.Client, err = db.Conf.Build()
	return err
}

func (db *Client) Close() error {
	if db.Client == nil {
		return nil
	}
	return db.Client.Close()
}
