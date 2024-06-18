package nacos

import (
	"github.com/hopeio/cherry/initialize/conf_dao/nacos"
	"github.com/hopeio/cherry/utils/log"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/cache"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/file"
	"github.com/nacos-group/nacos-sdk-go/v2/util"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"os"
)

var ConfigCenter = &Nacos{}

type Nacos struct {
	Conf   Config
	Client config_client.IConfigClient
}

type Config struct {
	nacos.Config
	vo.ConfigParam
}

func (cc *Nacos) Type() string {
	return "nacos"
}

func (cc *Nacos) Config() any {
	return &cc.Conf
}

func (cc *Nacos) Handle(handle func([]byte)) error {
	if cc.Client == nil {
		var err error
		cc.Client, err = cc.Conf.Config.Build()
		if err != nil {
			return err
		}
	}

	config, err := cc.Client.GetConfig(cc.Conf.ConfigParam)
	if err != nil {
		log.Fatal(err)
	}
	// nacos-go-sdk的问题，首次拉取的配置缓存在cache目录，listen拉取的缓存在cache/config，listen是异步的，如果要先同步获取配置且不在未更改配置的情况下触发listen的Onchange，就要把配置写进listen的目录，来回读取写入，浪费性能
	cacheDir := file.GetCurrentPath() + string(os.PathSeparator) + "cache/config"
	cacheKey := util.GetConfigCacheKey(cc.Conf.DataId, cc.Conf.Group, cc.Conf.ClientConfig.NamespaceId)
	cache.WriteConfigToFile(cacheKey, cacheDir, config)
	handle([]byte(config))
	cc.Conf.OnChange = func(namespace, group, dataId, data string) {
		handle([]byte(data))
	}

	return cc.Client.ListenConfig(cc.Conf.ConfigParam)
}

func (cc *Nacos) Close() error {
	cc.Client.CloseClient()
	return nil
}
