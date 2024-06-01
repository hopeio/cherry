package etcd

import (
	"github.com/hopeio/cherry/utils/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Config clientv3.Config

func (c *Config) InitBeforeInject() {
}

func (c *Config) Init() {
}

func (c *Config) Build() *clientv3.Client {
	c.Init()
	client, err := clientv3.New((clientv3.Config)(*c))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

type Client struct {
	*clientv3.Client
	Conf Config
}

func (e *Client) Config() any {
	return &e.Conf
}

func (e *Client) Set() {
	e.Client = e.Conf.Build()
}

func (e *Client) Close() error {
	if e.Client == nil {
		return nil
	}
	return e.Client.Close()
}
