package initialize

import (
	"github.com/hopeio/cherry/initialize/initconf"
	"github.com/hopeio/cherry/utils/log"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"os"
	"path/filepath"
)

// SingleFileConfig This is for illustrative purposes only and is not for practical use
type SingleFileConfig struct {
	initconf.BasicConfig
	initconf.EnvConfig
}

func (gc *globalConfig) setBasicConfig() {
	format := gc.InitConfig.ConfigCenter.Format
	basicConfig := &SingleFileConfig{}

	err := gc.Viper.Unmarshal(basicConfig, decoderConfigOptions...)
	if err != nil {
		log.Fatal(err)
	}
	applyFlagConfig(nil, basicConfig)
	gc.InitConfig.BasicConfig = basicConfig.BasicConfig
	gc.InitConfig.EnvConfig = basicConfig.EnvConfig
	if gc.InitConfig.ConfigCenter.Format == "" {
		gc.InitConfig.ConfigCenter.Format = format
	}
	if gc.InitConfig.Module == "" {
		gc.InitConfig.Module = stringsi.CutPart(filepath.Base(os.Args[0]), ".")
	}
}
