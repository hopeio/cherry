package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type Config struct {
	vo.NacosClientParam
}

func (c *Config) InitBeforeInject() {

}

func (c *Config) Init() {
}

func (c *Config) Build() (config_client.IConfigClient, error) {
	return clients.NewConfigClient(c.NacosClientParam)
}

type ConfigClient struct {
	Client config_client.IConfigClient
	Conf   Config
}

func (m *ConfigClient) Config() any {
	return &m.Conf
}

func (m *ConfigClient) Init() error {
	var err error
	m.Client, err = m.Conf.Build()
	return err
}

func (m *ConfigClient) Close() error {
	m.Client.CloseClient()
	return nil
}
