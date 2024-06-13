package initialize

import (
	"github.com/hopeio/cherry/initialize/initconf"
	"reflect"
)

type InitBeforeInject interface {
	InitBeforeInject()
}

type InitBeforeInjectWithInitConfig interface {
	InitBeforeInjectWithInitConfig(*initconf.InitConfig)
}

type InitAfterInject interface {
	InitAfterInject()
}

type InitAfterInjectWithInitConfig interface {
	InitAfterInjectWithInitConfig(*initconf.InitConfig)
}

type InitAfterInjectConfig interface {
	InitAfterInjectConfig()
}

type InitAfterInjectConfigWithInitConfig interface {
	InitAfterInjectConfigWithInitConfig(*initconf.InitConfig)
}

type Config interface {
	// 注入之前设置默认值
	InitBeforeInject
	// 注入之后初始化
	InitAfterInject
}

type Dao interface {
	InitBeforeInject
	// 注入config后执行
	InitAfterInjectConfig
	// 注入dao后执行
	InitAfterInject
}

type EmbeddedPresets struct {
}

func (u EmbeddedPresets) InitBeforeInject() {
}
func (u EmbeddedPresets) InitAfterInjectConfig() {
}
func (u EmbeddedPresets) InitAfterInject() {
}

var EmbeddedPresetsType = reflect.TypeOf((*EmbeddedPresets)(nil)).Elem()
