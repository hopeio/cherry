package initialize

import (
	"github.com/hopeio/cherry/initialize/conf_center"
	"github.com/hopeio/cherry/initialize/conf_center/local"
	"github.com/hopeio/cherry/initialize/initconf"
	"github.com/hopeio/cherry/utils/errors/multierr"
	"github.com/hopeio/cherry/utils/io/fs"
	"github.com/hopeio/cherry/utils/slices"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/hopeio/cherry/utils/log"
)

// 约定大于配置
var (
	gConfig = &globalConfig{
		InitConfig: initconf.InitConfig{
			ConfUrl:   "",
			EnvConfig: initconf.EnvConfig{Debug: true},
		},

		Viper: viper.New(),
		lock:  sync.RWMutex{},
	}
	decoderConfigOptions = []viper.DecoderConfigOption{
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.TextUnmarshallerHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		)),
		func(config *mapstructure.DecoderConfig) {
			config.Squash = true
		},
	}
)

func GlobalConfig() *globalConfig {
	return gConfig
}

// globalConfig
// 全局配置
type globalConfig struct {
	InitConfig initconf.InitConfig `mapstructure:",squash"`
	BuiltinConfig

	conf Config
	dao  Dao

	*viper.Viper

	/*
		cacheConf      any*/
	editTimes   uint32
	defers      []func()
	initialized bool
	lock        sync.RWMutex
}

func Start(conf Config, dao Dao, configCenter ...conf_center.ConfigCenter) func() {
	if reflect.ValueOf(conf).IsNil() {
		log.Fatalf("初始化错误: 配置不能为空")
	}
	gConfig.initialized = false

	// 为支持自定义配置中心,并且遵循依赖最小化原则,配置中心改为可插拔的,考虑将配置序列话也照此重做
	// 注册配置中心,默认注册本地文件
	conf_center.RegisterConfigCenter(local.ConfigCenter)
	for _, cc := range configCenter {
		conf_center.RegisterConfigCenter(cc)
	}

	gConfig.setConfDao(conf, dao)
	gConfig.loadConfig()
	gConfig.initialized = true
	return func() {
		// 倒序调用defer
		for i := len(gConfig.defers) - 1; i > 0; i-- {
			gConfig.defers[i]()
		}
	}
}

func (gc *globalConfig) setConfDao(conf Config, dao Dao) {
	gc.conf = conf
	gc.dao = dao
	gc.defers = append(gc.defers, func() {
		log.Sync()
	})
	if dao != nil {
		gc.defers = append(gc.defers, func() {
			closeDao(dao)
		})
	}

}

const defaultConfigName = "config"

func (gc *globalConfig) loadConfig() {
	gc.Viper.AutomaticEnv()
	var format string
	// find config
	if gc.InitConfig.ConfUrl == "" {
		log.Debug("searching for config in .")
		for _, ext := range viper.SupportedExts {
			filePath := filepath.Join(".", defaultConfigName+"."+ext)
			if b := fs.Exist(filePath); b {
				log.Debug("found file", "file", filePath)
				gc.InitConfig.ConfUrl = filePath
				format = ext
				break
			}
		}
	}
	if gc.InitConfig.ConfUrl != "" {
		log.Infof("load config from: %s", gc.InitConfig.ConfUrl)
		if format == "" {
			format = path.Ext(gc.InitConfig.ConfUrl)
			if format != "" {
				// remove .
				format = format[1:]
				if !slices.Contains(viper.SupportedExts, format) {
					log.Fatalf("unsupport config format, support: %v", viper.SupportedExts)
				}
			} else {
				log.Fatalf("config path need format ext, support: %v", viper.SupportedExts)
			}
		}

		gc.InitConfig.ConfigCenter.Format = format
		gc.Viper.SetConfigType(format)
		gc.Viper.SetConfigFile(gc.InitConfig.ConfUrl)
		err := gc.Viper.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	gc.setBasicConfig()
	gc.setEnvConfig()
	for i := range gc.InitConfig.NoInject {
		gc.InitConfig.NoInject[i] = strings.ToUpper(gc.InitConfig.NoInject[i])
	}

	var singleTemplateFileConfig bool
	if gc.InitConfig.EnvConfig.ConfigCenter.ConfigCenter == nil {
		if gc.InitConfig.Env == "" {
			singleTemplateFileConfig = true
		}
		// 单配置文件
		gc.InitConfig.ConfigCenter.ConfigCenter = &local.Local{
			ConfigPath: gc.InitConfig.ConfUrl,
		}
		applyFlagConfig(gc.Viper, gc.InitConfig.ConfigCenter.ConfigCenter)
	}

	// hook function
	gc.conf.InitBeforeInject()
	if c, ok := gc.conf.(InitBeforeInjectWithInitConfig); ok {
		c.InitBeforeInjectWithInitConfig(&gc.InitConfig)
	}
	if !gc.initialized && gc.dao != nil {
		gc.dao.InitBeforeInject()
		if c, ok := gc.dao.(InitBeforeInjectWithInitConfig); ok {
			c.InitBeforeInjectWithInitConfig(&gc.InitConfig)
		}
	}

	gc.genConfigTemplate(singleTemplateFileConfig)

	cfgcenter := gc.InitConfig.ConfigCenter.ConfigCenter
	err := cfgcenter.HandleConfig(gc.UnmarshalAndSet)
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}

}

func (gc *globalConfig) DeferFunc(deferf ...func()) {
	gc.lock.Lock()
	defer gc.lock.Unlock()
	gc.defers = append(gc.defers, deferf...)
}

func RegisterDeferFunc(deferf ...func()) {
	gConfig.lock.Lock()
	defer gConfig.lock.Unlock()
	gConfig.defers = append(gConfig.defers, deferf...)
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
	var errs multierr.MultiError
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
				errs.Append(err)
			}

		}
	}

	if errs.HasErrors() {
		return &errs
	}
	return nil
}

func GetConfig[T any]() *T {
	gConfig.lock.RLock()
	defer gConfig.lock.RUnlock()
	if gConfig.initialized == false {
		log.Fatalf("配置未初始化")
	}
	conf := gConfig.conf
	value := reflect.ValueOf(conf).Elem()
	for i := 0; i < value.NumField(); i++ {
		if conf, ok := value.Field(i).Interface().(T); ok {
			return &conf
		}
	}
	return new(T)
}

func GetDao[T any]() *T {
	gConfig.lock.RLock()
	defer gConfig.lock.RUnlock()
	if gConfig.initialized == false {
		log.Fatalf("配置未初始化")
	}
	dao := gConfig.dao
	value := reflect.ValueOf(dao).Elem()
	for i := 0; i < value.NumField(); i++ {
		if dao, ok := value.Field(i).Interface().(T); ok {
			return &dao
		}
	}
	return new(T)
}
