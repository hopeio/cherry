package client

import (
	"github.com/hopeio/cherry/utils/log"
	grpci "github.com/hopeio/cherry/utils/net/http/grpc"
	"google.golang.org/grpc"
)

type Config struct {
	Addr    string
	TLS     bool
	Options []grpc.DialOption
}

func (c *Config) InitBeforeInject() {
}

func (c *Config) Init() {

}

func (c *Config) Build() *grpc.ClientConn {
	c.Init()
	var conn *grpc.ClientConn
	var err error
	if c.TLS {
		conn, err = grpci.NewTLSClient(c.Addr, c.Options...)
	} else {
		conn, err = grpci.NewClient(c.Addr, c.Options...)
	}
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

type Client struct {
	Conn *grpc.ClientConn
	Conf Config
}

func (c *Client) Config() any {
	return &c.Conf
}

func (c *Client) Set() {
	c.Conn = c.Conf.Build()
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
