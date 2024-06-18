package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/hopeio/cherry/utils/net/http/auth"
)

type Config elasticsearch.Config

func (c *Config) InitBeforeInject() {
}

func (c *Config) Init() {

}

func (c *Config) Build() (*elasticsearch.Client, error) {
	c.Init()
	if c.Username != "" && c.Password != "" {
		auth.SetBasicAuth(c.Header, c.Username, c.Password)
	}

	return elasticsearch.NewClient((elasticsearch.Config)(*c))
}

type Client struct {
	*elasticsearch.Client
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
