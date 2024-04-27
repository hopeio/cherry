package initialize

import (
	"github.com/hopeio/cherry/initialize/initconf"
	"github.com/hopeio/cherry/utils/log"
)

// SingleFileConfig This is for illustrative purposes only and is not for practical use
type SingleFileConfig struct {
	initconf.BasicConfig `yaml:",inline"`
	initconf.EnvConfig   `yaml:",inline"`
}

func (gc *globalConfig) setBasicConfig() {
	format := gc.InitConfig.ConfigCenter.Format
	basicConfig := &SingleFileConfig{}
	/*format, err := common.Unmarshal(gc.ConfigCenter.Format, data, basicConfig)
	if err != nil {
		return
	}*/
	err := gc.Viper.Unmarshal(basicConfig, decoderConfigOptions...)
	if err != nil {
		log.Fatal(err)
	}
	parseFlag(gc.flag)
	gc.InitConfig.BasicConfig = basicConfig.BasicConfig
	gc.InitConfig.EnvConfig = basicConfig.EnvConfig
	gc.InitConfig.ConfigCenter.Format = format
	if gc.InitConfig.Module == "" {
		gc.InitConfig.Module = "cherry-app"
	}
	if gc.InitConfig.Env == "" {
		gc.InitConfig.Env = "dev"
	}
}
