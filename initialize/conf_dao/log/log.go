package log

import (
	initialize2 "github.com/hopeio/cherry/initialize"
	"github.com/hopeio/cherry/utils/log"
)

// 全局变量,只一个实例,只提供config
type Config log.Config

func (c *Config) InitBeforeInject() {
	c.Development = initialize2.GlobalConfig.Debug
	c.AppName = initialize2.GlobalConfig.Module
}

func (c *Config) InitAfterInject() {
	logConf := (*log.Config)(c)
	log.SetDefaultLogger(logConf)
}

/*func (c *Config) Build() *log.Logger {
	c.Init()
	return (*log.Config)(c).NewLogger()
}*/

/*type Logger struct {
	*log.Logger `
	Conf        Config
}

func (l *Logger) Config() any {
	return &l.Conf
}

func (l *Logger) SetEntity() {
	l.Logger = l.Conf.Build()
}

func (l *Logger) Close() error {
	if l.Logger == nil {
		return nil
	}
	return l.Logger.Sync()
}
*/
