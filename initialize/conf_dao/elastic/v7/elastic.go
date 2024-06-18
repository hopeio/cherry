package elastic

import (
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/config"
)

type Config config.Config

func (c *Config) InitBeforeInject() {
}

func (c *Config) Init() {

}

func (c *Config) Build() (*elastic.Client, error) {
	c.Init()
	return elastic.NewClientFromConfig((*config.Config)(c))
}

type Client struct {
	*elastic.Client
	Conf Config
}

func (es *Client) Config() any {
	return &es.Conf
}

func (es *Client) Init() error {
	var err error
	es.Client, err = es.Conf.Build()
	return err
}

func (es *Client) Close() error {
	return nil
}
