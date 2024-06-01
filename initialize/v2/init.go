package v2

import (
	"github.com/hopeio/cherry/initialize"
	"github.com/hopeio/cherry/initialize/initconf"
	"github.com/spf13/viper"
	"sync"
)

type globalConfig[C initialize.Config, D initialize.Dao] struct {
	InitConfig initconf.InitConfig `mapstructure:",squash"`
	initialize.BuiltinConfig

	conf C
	dao  D

	*viper.Viper

	/*
		cacheConf      any*/
	editTimes   uint32
	defers      []func()
	initialized bool
	lock        sync.RWMutex
}
