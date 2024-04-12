package initialize

import "github.com/hopeio/cherry/utils/encoding/common"

// BasicConfig
// zh: 基本配置，包含模块名
type BasicConfig struct {
	// 配置文件路径
	ConfUrl string `flag:"name:confdao;short:c;default:config.toml;usage:配置文件路径,默认./config.toml或./config/config.toml;env:CONFIG" json:"conf_url,omitempty"`
	// 模块名
	Module string `flag:"name:mod;short:m;default:cherry-app;usage:模块名;env:MODULE" json:"module,omitempty"`
	// environment
	Env string `flag:"name:env;short:e;default:dev;usage:环境;env:ENV" json:"env,omitempty"`
}

func (gc *globalConfig) setBasicConfig(data []byte) {
	basicConfig := &BasicConfig{}
	format, err := common.Unmarshal(gc.ConfigCenter.Format, data, gc)
	if err != nil {
		return
	}
	gc.ConfigCenter.Format = format

	if gc.Module == "" {
		gc.Module = "cherry-app"
		if basicConfig.Module != "" {
			gc.Module = basicConfig.Module
		}
	}
	if gc.Env == "" {
		gc.Env = "dev"
		if basicConfig.Env != "" {
			gc.Module = basicConfig.Env
		}
	}
}
