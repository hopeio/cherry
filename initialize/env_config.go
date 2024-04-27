package initialize

import (
	"github.com/hopeio/cherry/initialize/conf_center"
	"github.com/hopeio/cherry/initialize/initconf"
	"github.com/hopeio/cherry/utils/log"
	"github.com/hopeio/cherry/utils/reflect/mtos"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

// FileConfig
// unused
// example
// 配置文件映射结构体,每个启动都有一个必要的配置文件,用于初始化基本配置及配置中心配置
type FileConfig struct {
	// 模块名
	initconf.BasicConfig
	EnvConfig *initconf.EnvConfig `init:"fixed"` // field name can be dev,test,prod ... and anything you like
}

const (
	fixedFieldNameEnvConfig       = "EnvConfig"
	fixedFieldNameBasicConfig     = "BasicConfig"
	fixedFieldNameConfigCenter    = "ConfigCenter"
	fixedFieldNameEncoderRegistry = "encoderRegistry"
	prefixConfigTemplate          = "config.template."
	prefixLocalTemplate           = "local.template."
	skipTypeTlsConfig             = "tls.Config"
)

func (gc *globalConfig) setEnvConfig() {
	format := gc.InitConfig.ConfigCenter.Format
	envConfig, ok := gc.Viper.Get(gc.InitConfig.Env).(map[string]any)
	if !ok {
		log.Warn("lack of environment configuration, try single config file")
		return
	}
	err := mtos.Unmarshal(&gc.InitConfig.EnvConfig, envConfig)
	if err != nil {
		log.Fatal(err)
	}
	parseFlag(gc.flag) // TODO: 缺少环境变量覆盖
	if gc.InitConfig.EnvConfig.ConfigCenter.ConfigType == "" {
		log.Warn("lack of configCenter configType, try single config file")
		return
	}

	cc, ok := conf_center.GetRegisteredConfigCenter()[strings.ToLower(gc.InitConfig.EnvConfig.ConfigCenter.ConfigType)]
	if !ok {
		log.Warn("lack of registered configCenter, try single config file")
		return
	}

	ccConfig, ok := gc.Viper.Get(gc.InitConfig.Env + ".ConfigCenter." + gc.InitConfig.EnvConfig.ConfigCenter.ConfigType).(map[string]any)
	if !ok {
		log.Warn("lack of configCenter config, try single config file")
		return
	}
	err = mtos.Unmarshal(cc, ccConfig)
	if err != nil {
		log.Fatal(err)
	}
	injectFlagConfig(gc.flag, reflect.ValueOf(cc).Elem())
	parseFlag(gc.flag)
	gc.InitConfig.EnvConfig.ConfigCenter.ConfigCenter = cc
	// template
	confMap := make(map[string]any)
	struct2Map(reflect.ValueOf(&gc.InitConfig.BasicConfig).Elem(), confMap)
	envMap := make(map[string]any)
	struct2Map(reflect.ValueOf(&gc.InitConfig.EnvConfig).Elem(), envMap)
	confMap[gc.InitConfig.Env] = envMap
	ccMap := envMap[fixedFieldNameConfigCenter].(map[string]any)
	for name, v := range conf_center.GetRegisteredConfigCenter() {
		cc := make(map[string]any)
		struct2Map(reflect.ValueOf(v).Elem(), cc)
		ccMap[name] = cc
	}
	// unsafe
	encoderRegistry := reflect.ValueOf(gc.Viper).Elem().FieldByName(fixedFieldNameEncoderRegistry).Elem()
	fieldValue := reflect.NewAt(encoderRegistry.Type(), unsafe.Pointer(encoderRegistry.UnsafeAddr()))
	data, err := fieldValue.Interface().(Encoder).Encode(string(format), confMap)

	if gc.InitConfig.EnvConfig.ConfigTemplateDir != "" {
		dir := gc.InitConfig.EnvConfig.ConfigTemplateDir
		if dir[len(dir)-1] != '/' {
			dir += "/"
		}
		err = os.WriteFile(dir+prefixConfigTemplate+string(format), data, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
