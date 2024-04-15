package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hopeio/cherry/utils/log"
)

type Config mqtt.ClientOptions

func (c *Config) InitBeforeInject() {

}

func (c *Config) InitAfterInject() {
}

func (c *Config) Build() mqtt.Client {
	client := mqtt.NewClient((*mqtt.ClientOptions)(c))
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	return nil
}

type Client struct {
	Conf Config
	mqtt.Client
}

func (c *Client) Config() any {
	c.Conf = Config(*mqtt.NewClientOptions())
	return &c.Conf
}

func (c *Client) SetEntity() {
	c.Client = c.Conf.Build()
}

func (c *Client) Close() error {
	c.Client.Disconnect(5)
	return nil
}
