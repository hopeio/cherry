package v7

import (
	"github.com/hopeio/cherry/utils/log"
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/config"
)

type Config config.Config

func (c *Config) InitBeforeInject() {
}

func (c *Config) InitAfterInject() {

}

func (c *Config) Build() *elastic.Client {
	c.InitAfterInject()
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

func (es *Client) SetEntity() {
	es.Client = es.Conf.Build()
}

func (es *Client) Close() error {
	return nil
}
