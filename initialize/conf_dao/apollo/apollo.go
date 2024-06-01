package apollo

import (
	"github.com/hopeio/cherry/utils/configor/apollo"
	"github.com/hopeio/cherry/utils/log"
)

type Config apollo.Config

func (c *Config) InitBeforeInject() {
}

func (c *Config) Init() {

}

func (c *Config) Build() *apollo.Client {
	c.Init()
	//初始化更新配置，这里不需要，开启实时更新时初始化会更新一次
	client, err := (*apollo.Config)(c).NewClient()
	if err != nil {
		log.Fatal(err)
	}
	return client
}

type Client struct {
	*apollo.Client
	Conf Config
}

func (c *Client) Config() any {
	return &c.Conf
}

func (c *Client) Set() {
	c.Client = c.Conf.Build()
}

func (c *Client) Close() error {
	return c.Client.Close()
}
