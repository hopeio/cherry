package conf_center

import (
	"github.com/hopeio/cherry/utils/log"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"io"
	"strings"
)

type ConfigType string

type ConfigCenter interface {
	Config() any
	io.Closer
	Handle(func([]byte)) error
	Type() string
}

type Config struct {
	// 配置格式
	Format string `flag:"name:format;default:toml;usage:配置格式"`
	// 配置类型
	Type string `flag:"name:conf_type;default:local;usage:配置类型"`
	// config字段顺序不能变,ConfigCenter 保持在最后
	ConfigCenter ConfigCenter
}

var configCenter = map[string]ConfigCenter{}

func RegisterConfigCenter(c ConfigCenter) {
	if c != nil {
		typ := strings.ToLower(c.Type())
		if !stringsi.IsASCIILetters(typ) {
			log.Fatal("config type must be letters")
		}
		if _, ok := configCenter[typ]; !ok {
			configCenter[typ] = c
		}
	}
}

func GetConfigCenter(configType string) ConfigCenter {
	return configCenter[configType]
}

func GetRegisteredConfigCenter() map[string]ConfigCenter {
	return configCenter
}
