package initialize

import (
	"github.com/hopeio/cherry/initialize/conf_center"
	"github.com/hopeio/cherry/initialize/initconf"
	"github.com/hopeio/cherry/utils/log"
	"github.com/hopeio/cherry/utils/reflect/mtos"
	"github.com/spf13/viper"
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
	EnvConfig *initconf.EnvConfig // field name can be dev,test,prod ... and anything you like
}

const (
	fixedFieldNameEnvConfig       = "EnvConfig"
	fixedFieldNameBasicConfig     = "BasicConfig"
	fixedFieldNameConfigCenter    = "ConfigCenter"
	fixedFieldNameEnv             = "Env"
	fixedFieldNameEncoderRegistry = "encoderRegistry"
	prefixConfigTemplate          = "config.template."
	prefixLocalTemplate           = "local.template."
	skipTypeTlsConfig             = "tls.Config"
)

func (gc *globalConfig) setEnvConfig() {
	if gc.InitConfig.Env == "" {
		log.Warn("lack of env configuration, try single config file")
		return
	}
	format := gc.InitConfig.ConfigCenter.Format

	defer func() {
		if gc.InitConfig.EnvConfig.ConfigTemplateDir != "" {
			// template
			confMap := make(map[string]any)
			struct2Map(&gc.InitConfig.BasicConfig, confMap)
			envMap := make(map[string]any)
			struct2Map(&gc.InitConfig.EnvConfig, envMap)
			confMap[gc.InitConfig.Env] = envMap
			ccMap := make(map[string]any)
			struct2Map(&gc.InitConfig.EnvConfig.ConfigCenter, ccMap)
			envMap[fixedFieldNameConfigCenter] = ccMap
			for name, v := range conf_center.GetRegisteredConfigCenter() {
				cc := make(map[string]any)
				struct2Map(v, cc)
				ccMap[name] = cc
			}
			// unsafe
			encoderRegistry := reflect.ValueOf(gc.Viper).Elem().FieldByName(fixedFieldNameEncoderRegistry).Elem()
			fieldValue := reflect.NewAt(encoderRegistry.Type(), unsafe.Pointer(encoderRegistry.UnsafeAddr()))
			data, err := fieldValue.Interface().(Encoder).Encode(string(format), confMap)

			dir := gc.InitConfig.EnvConfig.ConfigTemplateDir
			if dir[len(dir)-1] != '/' {
				dir += "/"
			}
			err = os.WriteFile(dir+prefixConfigTemplate+string(format), data, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	envConfig, ok := gc.Viper.Get(gc.InitConfig.Env).(map[string]any)
	if !ok {
		log.Warn("lack of environment configuration, try single config file")
		return
	}
	err := mtos.Unmarshal(&gc.InitConfig.EnvConfig, envConfig)
	if err != nil {
		log.Fatal(err)
	}
	applyFlagConfig(nil, &gc.InitConfig.EnvConfig)
	if gc.InitConfig.EnvConfig.ConfigCenter.Format == "" {
		log.Fatalf("lack of configCenter format, support format:%v", viper.SupportedExts)
	}
	gc.Viper.SetConfigType(gc.InitConfig.EnvConfig.ConfigCenter.Format)
	if gc.InitConfig.EnvConfig.ConfigCenter.Type == "" {
		log.Warn("lack of configCenter type, try single config file")
		return
	}

	configCenter, ok := conf_center.GetRegisteredConfigCenter()[strings.ToLower(gc.InitConfig.EnvConfig.ConfigCenter.Type)]
	if !ok {
		log.Warn("lack of registered configCenter, try single config file")
		return
	}

	applyFlagConfig(gc.Viper, configCenter)
	gc.InitConfig.EnvConfig.ConfigCenter.ConfigCenter = configCenter

	configCenterConfig, ok := gc.Viper.Get(gc.InitConfig.Env + ".configCenter." + gc.InitConfig.EnvConfig.ConfigCenter.Type).(map[string]any)
	if !ok {
		log.Warn("lack of configCenter config, try single config file")
		return
	}
	err = mtos.Unmarshal(configCenter, configCenterConfig)
	if err != nil {
		log.Fatal(err)
	}
}
