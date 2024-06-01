package elastic

import (
	"github.com/hopeio/cherry/utils/log"
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/config"
)

type Config config.Config

func (c *Config) InitBeforeInject() {
}

func (c *Config) Init() {

}

func (c *Config) Build() *elastic.Client {
	c.Init()
	client, err := elastic.NewClientFromConfig((*config.Config)(c))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

type Client struct {
	*elastic.Client
	Conf Config
}

func (es *Client) Config() any {
	return &es.Conf
}

func (es *Client) Set() {
	es.Client = es.Conf.Build()
}

func (es *Client) Close() error {
	return nil
}
