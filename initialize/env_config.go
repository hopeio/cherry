package initialize

import (
	"github.com/hopeio/cherry/initialize/conf_center"
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
	BasicConfig
	EnvConfig *EnvConfig `init:"fixed"` // field name can be dev,test,prod ... and anything you like
}

type EnvConfig struct {
	Debug             bool   `flag:"name:debug;short:d;default:debug;usage:是否测试;env:DEBUG" json:"debug" toml:"debug"`
	ConfigTemplateDir string `flag:"name:conf_tmpl_dir;short:t;usage:是否生成配置模板;env:CONFIG_TEMPLATE_DIR" json:"config_template_dir"`
	// 代理, socks5://localhost:1080
	Proxy    string `flag:"name:proxy;short:p;default:'socks5://localhost:1080';usage:代理;env:HTTP_PROXY" json:"proxy"`
	NoInject []string
	// config字段顺序不能变,ConfigCenter 保持在最后
	ConfigCenter conf_center.Config `init:"fixed"`
}

const (
	fixedFieldNameEnvConfig    = "EnvConfig"
	fixedFieldNameBasicConfig  = "BasicConfig"
	fixedFieldNameConfigCenter = "ConfigCenter"
)

func (gc *globalConfig) setEnvConfig() {
	gc.Viper.AllSettings()
	format := gc.ConfigCenter.Format
	envConfig, ok := gc.Viper.Get(gc.Env).(map[string]any)
	if !ok {
		log.Warn("lack of environment configuration, try single config file")
		return
	}
	err := mtos.Unmarshal(&gc.EnvConfig, envConfig)
	if err != nil {
		log.Fatal(err)
	}
	parseFlag(gc.flag)
	if gc.EnvConfig.ConfigCenter.ConfigType == "" {
		log.Warn("lack of configCenter configType, try single config file")
		return
	}

	cc, ok := conf_center.GetRegisteredConfigCenter()[strings.ToLower(gc.EnvConfig.ConfigCenter.ConfigType)]
	if !ok {
		log.Warn("lack of registered configCenter, try single config file")
		return
	}

	ccConfig, ok := gc.Viper.Get(gc.Env + ".ConfigCenter." + gc.EnvConfig.ConfigCenter.ConfigType).(map[string]any)
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
	gc.EnvConfig.ConfigCenter.ConfigCenter = cc

	confMap := make(map[string]any)
	struct2MapHelper(format, reflect.ValueOf(&gc.BasicConfig).Elem(), confMap)
	envMap := make(map[string]any)
	struct2MapHelper(format, reflect.ValueOf(&gc.EnvConfig).Elem(), envMap)
	confMap[gc.Env] = envMap
	ccMap := envMap["ConfigCenter"].(map[string]any)
	for name, v := range conf_center.GetRegisteredConfigCenter() {
		cc := make(map[string]any)
		struct2MapHelper(format, reflect.ValueOf(v).Elem(), cc)
		ccMap[name] = cc
	}
	// unsafe
	encoderRegistry := reflect.ValueOf(gc.Viper).Elem().FieldByName("encoderRegistry").Elem()
	fieldValue := reflect.NewAt(encoderRegistry.Type(), unsafe.Pointer(encoderRegistry.UnsafeAddr()))
	ndata, err := fieldValue.Interface().(Encoder).Encode(string(format), confMap)

	if gc.EnvConfig.ConfigTemplateDir != "" {
		dir := gc.EnvConfig.ConfigTemplateDir
		if dir[len(dir)-1] != '/' {
			dir += "/"
		}
		err = os.WriteFile(dir+"config.template."+string(format), ndata, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
