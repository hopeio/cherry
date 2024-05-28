package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/hopeio/cherry/utils/log"
	"github.com/hopeio/cherry/utils/net/http/auth"
	"net/http"
)

type Config elasticsearch.Config

func (c *Config) InitBeforeInject() {
}
func (c *Config) Init() {
	if c.Header == nil {
		c.Header = http.Header{}
	}
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

func (es *Client) SetEntity(entity interface{}) {
	if client, ok := entity.(*elasticsearch.Client); ok {
		es.Client = client
	}
}

func (es *Client) Close() {
}
