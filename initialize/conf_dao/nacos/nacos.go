package nacos

import (
	"github.com/hopeio/cherry/utils/log"
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

func (c *Config) Build() config_client.IConfigClient {
	client, err := clients.NewConfigClient(c.NacosClientParam)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

type ConfigClient struct {
	Client config_client.IConfigClient
	Conf   Config
}

func (m *ConfigClient) Config() any {
	return &m.Conf
}

func (m *ConfigClient) Set() {
	m.Client = m.Conf.Build()
}

func (m *ConfigClient) Close() error {
	m.Client.CloseClient()
	return nil
}
