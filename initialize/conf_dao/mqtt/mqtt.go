package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hopeio/cherry/utils/configor"
	"github.com/hopeio/cherry/utils/crypto/tls"
	"github.com/hopeio/cherry/utils/log"
	"time"
)

type Config struct {
	*mqtt.ClientOptions
	CertFile string `json:"cert_file,omitempty"`
	KeyFile  string `json:"key_file,omitempty"`
}

func (c *Config) InitBeforeInject() {
	c.ClientOptions = mqtt.NewClientOptions()
}

func (c *Config) InitAfterInject() {
	tlsConfig, err := tls.Certificate(c.CertFile, c.KeyFile)
	if err != nil {
		log.Fatal(err)
	}
	c.TLSConfig = tlsConfig

	configor.DurationNotify("PingTimeout", c.PingTimeout, time.Second)
	configor.DurationNotify("ConnectTimeout", c.ConnectTimeout, time.Second)
	configor.DurationNotify("MaxReconnectInterval", c.MaxReconnectInterval, time.Second)
	configor.DurationNotify("ConnectRetryInterval", c.ConnectRetryInterval, time.Second)
	configor.DurationNotify("WriteTimeout", c.WriteTimeout, time.Second)
}

func (c *Config) Build() mqtt.Client {
	client := mqtt.NewClient(c.ClientOptions)
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
	return &c.Conf
}

func (c *Client) SetEntity() {
	c.Client = c.Conf.Build()
}

func (c *Client) Close() error {
	c.Client.Disconnect(5)
	return nil
}
