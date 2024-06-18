package influxdb

import influxdb2 "github.com/influxdata/influxdb-client-go/v2"

type Config struct {
	ServerURL string
	AuthToken string
	options   *influxdb2.Options
}

func (c *Config) Build() influxdb2.Client {
	client := influxdb2.NewClientWithOptions(c.ServerURL, c.AuthToken, c.options)
	return client
}

type Client struct {
	Client influxdb2.Client
	Conf   Config
}

func (c *Client) Config() any {
	return c.Conf
}

func (c *Client) Init() error {
	c.Client = c.Conf.Build()
	return nil
}

func (c *Client) Close() error {
	c.Client.Close()
	return nil
}
