package initialize

import (
	"github.com/hopeio/cherry/utils/log"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"os"
	"reflect"
	"unsafe"
)

func (gc *globalConfig) genConfigTemplate(singleTemplateFileConfig bool) {
	dir := gc.InitConfig.ConfigTemplateDir
	if dir == "" {
		return
	}
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}

	format := gc.InitConfig.ConfigCenter.Format
	filename := prefixLocalTemplate + string(format)

	confMap := make(map[string]any)
	if singleTemplateFileConfig {
		filename = prefixConfigTemplate + string(format)
		struct2Map(&gc.InitConfig.BasicConfig, confMap)
		delete(confMap, fixedFieldNameEnv)
		struct2Map(&gc.InitConfig.EnvConfig, confMap)
		delete(confMap, fixedFieldNameConfigCenter)
	}
	struct2Map(gc.conf, confMap)
	if gc.dao != nil {
		daoConfig2Map(reflect.ValueOf(gc.dao).Elem(), confMap)
	}

	encoderRegistry := reflect.ValueOf(gc.Viper).Elem().FieldByName(fixedFieldNameEncoderRegistry).Elem()
	fieldValue := reflect.NewAt(encoderRegistry.Type(), unsafe.Pointer(encoderRegistry.UnsafeAddr()))
	data, err := fieldValue.Interface().(Encoder).Encode(format, confMap)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(dir+filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func daoConfig2Map(value reflect.Value, confMap map[string]any) {
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Addr().Type().Implements(DaoFieldType) {
			newconfMap := make(map[string]any)
			fieldType := typ.Field(i)
			name := fieldType.Name
			tagSettings := ParseInitTagSettings(fieldType.Tag.Get(initTagName))
			if tagSettings.ConfigName != "" {
				name = stringsi.UpperCaseFirst(tagSettings.ConfigName)
			}

			confMap[name] = newconfMap
			struct2Map(field.Addr().Interface().(DaoField).Config(), newconfMap)
		}
	}
}
