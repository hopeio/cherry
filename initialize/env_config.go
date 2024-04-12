package initialize

import (
	"github.com/hopeio/cherry/initialize/conf_center"
	"github.com/hopeio/cherry/utils/encoding/common"
	"github.com/hopeio/cherry/utils/log"
	strings "github.com/hopeio/cherry/utils/strings"
	"os"
	"reflect"
)

// FileConfig
// unused
// example
// 配置文件映射结构体,每个启动都有一个必要的配置文件,用于初始化基本配置及配置中心配置
type FileConfig struct {
	// 模块名
	Module    string
	Env       string
	EnvConfig *EnvConfig
}

type EnvConfig struct {
	Debug             bool   `flag:"name:debug;short:d;default:debug;usage:是否测试;env:DEBUG" json:"debug" toml:"debug"`
	ConfigTemplateDir string `flag:"name:conf_tmpl_dir;short:t;usage:是否生成配置模板;env:CONFIG_TEMPLATE_DIR" json:"config_template_dir"`
	// 代理, socks5://localhost:1080
	Proxy    string `flag:"name:proxy;short:p;default:'socks5://localhost:1080';usage:代理;env:HTTP_PROXY" json:"proxy"`
	NoInject []string
	// config字段顺序不能变,ConfigCenter 保持在最后
	ConfigCenter conf_center.Config
}

func (gc *globalConfig) setEnvConfig(data []byte) {

	fromat := gc.ConfigCenter.Format
	var structFields []reflect.StructField
	confCenterValue := reflect.ValueOf(&gc.EnvConfig.ConfigCenter).Elem()
	confCenterTyp := confCenterValue.Type()
	for i := 0; i < confCenterTyp.NumField(); i++ {
		field := confCenterTyp.Field(i)
		if field.Name == "ConfigCenter" {
			continue
		}
		structFields = append(structFields, reflect.StructField{Name: field.Name, Type: field.Type, Tag: field.Tag})
	}

	for name, v := range conf_center.GetRegisteredConfigCenter() {
		structFields = append(structFields, reflect.StructField{Name: strings.UpperCaseFirst(name), Type: reflect.TypeOf(v)})
	}
	newConfCenterTyp := reflect.StructOf(structFields)
	structFields = structFields[:0]
	envConfigValue := reflect.ValueOf(&gc.EnvConfig).Elem()
	envConfigTyp := envConfigValue.Type()
	for i := 0; i < envConfigTyp.NumField(); i++ {
		field := envConfigTyp.Field(i)
		if field.Name == "ConfigCenter" {
			structFields = append(structFields, reflect.StructField{Name: "ConfigCenter", Type: newConfCenterTyp, Tag: field.Tag})
			continue
		}
		structFields = append(structFields, reflect.StructField{Name: field.Name, Type: field.Type, Tag: field.Tag})
	}
	newEnvConfigTyp := reflect.StructOf(structFields)

	structFields = structFields[:0]
	fileConfigTyp := reflect.TypeOf(FileConfig{})
	for i := 0; i < fileConfigTyp.NumField(); i++ {
		field := fileConfigTyp.Field(i)
		if field.Name == "EnvConfig" {
			structFields = append(structFields, reflect.StructField{Name: strings.UpperCaseFirst(gc.Env), Type: newEnvConfigTyp, Tag: reflect.StructTag(genEncodingTag(gc.Env))})
			continue
		}
		structFields = append(structFields, reflect.StructField{Name: field.Name, Type: field.Type, Tag: field.Tag})
	}

	newFileConfigTyp := reflect.StructOf(structFields)
	tmpFileConfigValue := reflect.New(newFileConfigTyp)
	tmpFileConfig := tmpFileConfigValue.Interface()
	confMap := make(map[string]any)
	newConfig(gc.ConfigCenter.Format, reflect.ValueOf(tmpFileConfig).Elem(), confMap)
	ndata, err := common.Marshal(gc.ConfigCenter.Format, confMap)
	if err != nil {
		log.Fatal(err)
	}
	_, err = common.Unmarshal(gc.ConfigCenter.Format, data, tmpFileConfig)
	if err != nil {
		log.Fatal(err)
	}
	tmpEnvConfigValue := tmpFileConfigValue.Elem().FieldByName(strings.UpperCaseFirst(gc.Env))
	if !tmpEnvConfigValue.IsValid() || tmpEnvConfigValue.IsZero() {
		log.Fatal("缺少环境配置")
	}

	// config字段顺序不能变,ConfigCenter 保持在最后
	for i := 0; i < envConfigValue.NumField(); i++ {
		field := envConfigValue.Field(i)
		structField := envConfigTyp.Field(i)
		if structField.Name == "ConfigCenter" {
			tmpccField := tmpEnvConfigValue.Field(i)
			for j := 0; j < confCenterValue.NumField(); j++ {
				ccField := confCenterValue.Field(j)
				ccstructField := confCenterTyp.Field(j)
				if ccstructField.Name == "ConfigCenter" {
					ccField.Set(tmpccField.FieldByName(strings.UpperCaseFirst(gc.EnvConfig.ConfigCenter.ConfigType)))
					continue
				}
				ccField.Set(tmpccField.Field(j))
			}
			continue
		}
		field.Set(tmpEnvConfigValue.Field(i))
	}

	if gc.EnvConfig.ConfigTemplateDir != "" {
		dir := gc.EnvConfig.ConfigTemplateDir
		if dir[len(dir)-1] != '/' {
			dir += "/"
		}
		err = os.WriteFile(dir+"config.template."+string(fromat), ndata, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	if gc.ConfigCenter.Format == "" {
		gc.ConfigCenter.Format = fromat
	}
}
