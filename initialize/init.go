package initialize

import (
	"github.com/hopeio/cherry/initialize/conf_center"
	"github.com/hopeio/cherry/initialize/conf_center/local"
	"github.com/hopeio/cherry/utils/encoding"
	"github.com/hopeio/cherry/utils/errors/multierr"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"

	"github.com/hopeio/cherry/utils/log"
)

// 约定大于配置
var (
	GlobalConfig = &globalConfig{
		BasicConfig: BasicConfig{ConfUrl: "./config.toml"},
		EnvConfig:   EnvConfig{Debug: true},
		lock:        sync.RWMutex{},
	}
)

// globalConfig
// 全局配置
type globalConfig struct {
	BasicConfig
	EnvConfig

	conf        Config
	dao         Dao
	deferFuncs  []func()
	initialized bool
	lock        sync.RWMutex
}

func Start(conf Config, dao Dao, configCenter ...conf_center.ConfigCenter) func() {
	if conf == nil {
		log.Fatalf("初始化错误: 配置不能为空")
	}
	GlobalConfig.initialized = false

	// 为支持自定义配置中心,并且遵循依赖最小化原则,配置中心改为可插拔的,考虑将配置序列话也照此重做
	// 注册配置中心,默认注册本地文件
	conf_center.RegisterConfigCenter(local.ConfigCenter)
	for _, cc := range configCenter {
		conf_center.RegisterConfigCenter(cc)
	}

	GlobalConfig.setConfDao(conf, dao)
	GlobalConfig.LoadConfig()
	GlobalConfig.initialized = true
	return func() {
		for _, f := range GlobalConfig.deferFuncs {
			f()
		}
	}
}

func (gc *globalConfig) setConfDao(conf Config, dao Dao) {
	gc.conf = conf
	gc.dao = dao

	gc.deferFuncs = []func(){
		func() { closeDao(dao) },
		func() { log.Sync() },
	}
}

func (gc *globalConfig) LoadConfig() {
	log.Infof("Load config from: %s\n", gc.ConfUrl)
	if _, err := os.Stat(gc.ConfUrl); os.IsNotExist(err) {
		log.Fatalf("配置路径错误: 请确保可执行文件和配置文件在同一目录下或在config目录下或指定配置文件")
	}
	data, err := os.ReadFile(gc.ConfUrl)
	if err != nil {
		log.Fatalf("读取配置错误: %v", err)
	}

	format := encoding.Format(path.Ext(gc.ConfUrl))
	if format != "" {
		// remove .
		format = format[1:]
		if format == "yml" {
			format = "yaml"
		}
	}

	gc.ConfigCenter.Format = format
	gc.setBasicConfig(data)
	gc.setEnvConfig(data)

	if gc.EnvConfig.ConfigCenter.ConfigCenter == nil {
		// 单配置文件
		gc.ConfigCenter.ConfigCenter = &local.Local{
			ReloadType: local.ReloadTypeFsNotify,
			ConfigPath: gc.ConfUrl,
		}
	}

	for i := range gc.NoInject {
		gc.NoInject[i] = strings.ToUpper(gc.NoInject[i])
	}

	gc.applyFlagConfig()

	gc.conf.InitBeforeInject()
	if !gc.initialized && gc.dao != nil {
		gc.dao.InitBeforeInject()
	}

	GenConfigTemplate(format, gc.conf, gc.dao, GlobalConfig.ConfigTemplateDir)

	cfgcenter := gc.ConfigCenter.ConfigCenter
	err = cfgcenter.HandleConfig(gc.UnmarshalAndSet)
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}

}

func (gc *globalConfig) RegisterDeferFunc(deferf ...func()) {
	gc.lock.Lock()
	defer gc.lock.Unlock()
	gc.deferFuncs = append(gc.deferFuncs, deferf...)
}

func (gc *globalConfig) Config() Config {
	return gc.conf
}

func (gc *globalConfig) closeDao() {
	if !gc.initialized || gc.dao == nil {
		return
	}
	err := closeDao(gc.dao)
	if err != nil {
		log.Error(err)
	}
}

func closeDao(dao Dao) error {
	var merr multierr.MultiError
	daoValue := reflect.ValueOf(dao).Elem()
	for i := 0; i < daoValue.NumField(); i++ {
		fieldV := daoValue.Field(i)
		if fieldV.Type().Kind() == reflect.Struct {
			fieldV = daoValue.Field(i).Addr()
		}
		if !fieldV.IsValid() || fieldV.IsNil() {
			continue
		}
		inter := fieldV.Interface()
		if daofield, ok := inter.(DaoField); ok {
			if err := daofield.Close(); err != nil {
				merr.Append(err)
			}

		}
	}

	if merr.HasErrors() {
		return &merr
	}
	return nil
}

func GetConfig[T any]() *T {
	iconf := GlobalConfig.Config()
	value := reflect.ValueOf(iconf).Elem()
	for i := 0; i < value.NumField(); i++ {
		if conf, ok := value.Field(i).Interface().(T); ok {
			return &conf
		}
	}
	return new(T)
}
