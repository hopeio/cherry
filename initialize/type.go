package initialize

import (
	"github.com/hopeio/cherry/initialize/initconf"
	"io"
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

var DaoFieldType = reflect.TypeOf((*DaoField)(nil)).Elem()

type DaoField interface {
	Config() any
	SetEntity()
	io.Closer
}

// TODO
type DaoFieldCloserWithError = io.Closer
type DaoFieldCloser interface {
	Close()
}

type CloseFunc func() error

type DaoConfig[D any] interface {
	Build() (*D, CloseFunc)
}

type DaoEntity[C DaoConfig[D], D any] struct {
	Conf   C
	Client *D
	close  CloseFunc
}

func (d *DaoEntity[C, D]) Config() any {
	return d.Conf
}

func (d *DaoEntity[C, D]) SetEntity() {
	d.Client, d.close = d.Conf.Build()
}

func (d *DaoEntity[C, D]) Close() error {
	if d.close != nil {
		return d.close()
	}
	return nil
}

type Marshal = func(any) ([]byte, error)
type Unmarshal = func([]byte, any) error
