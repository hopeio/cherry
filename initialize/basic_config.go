package initialize

import "github.com/hopeio/cherry/utils/encoding/common"

// BasicConfig
// zh: 基本配置，包含模块名
type BasicConfig struct {
	// 模块名
	Module string `flag:"name:mod;short:m;default:cherry-app;usage:模块名;env:MODULE" json:"module,omitempty"`
	// environment
	Env string `flag:"name:env;short:e;default:dev;usage:环境;env:ENV" json:"env,omitempty"`
}

// SingleFileConfig This is for illustrative purposes only and is not for practical use
type SingleFileConfig struct {
	BasicConfig
	EnvConfig
}

func (gc *globalConfig) setBasicConfig(data []byte) {
	basicConfig := &SingleFileConfig{}
	format, err := common.Unmarshal(gc.ConfigCenter.Format, data, basicConfig)
	if err != nil {
		return
	}
	gc.BasicConfig = basicConfig.BasicConfig
	gc.EnvConfig = basicConfig.EnvConfig
	gc.ConfigCenter.Format = format

	if gc.Module == "" {
		gc.Module = "cherry-app"
	}
	if gc.Env == "" {
		gc.Env = "dev"
	}
}
