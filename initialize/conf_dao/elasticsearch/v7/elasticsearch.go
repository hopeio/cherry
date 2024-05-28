package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/hopeio/cherry/utils/log"
	"github.com/hopeio/cherry/utils/net/http/auth"
)

type Config elasticsearch.Config

func (c *Config) InitBeforeInject() {
}

func (c *Config) Init() {

}

func (c *Config) Build() *elasticsearch.Client {
	c.Init()
	if c.Username != "" && c.Password != "" {
		auth.SetBasicAuth(c.Header, c.Username, c.Password)
	}
	client, err := elasticsearch.NewClient((elasticsearch.Config)(*c))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

type Client struct {
	*elasticsearch.Client
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
