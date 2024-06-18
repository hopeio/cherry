package log

import (
	"github.com/hopeio/cherry/initialize/initconf"
	"github.com/hopeio/cherry/utils/log"
)

// 全局变量,只一个实例,只提供config
type Config log.Config

func (c *Config) InitBeforeInjectWithInitConfig(conf *initconf.InitConfig) {
	c.Development = conf.Debug
}

func (c *Config) InitAfterInject() {
	log.SetDefaultLogger((*log.Config)(c))
}

func (c *Config) Build() *log.Logger {
	return (*log.Config)(c).NewLogger()
}

type Logger struct {
	*log.Logger
	Conf Config
}

func (l *Logger) Config() any {
	return &l.Conf
}

func (l *Logger) Init() error {
	l.Logger = l.Conf.Build()
	return nil
}

func (l *Logger) Close() error {
	if l.Logger == nil {
		return nil
	}
	return l.Logger.Sync()
}
