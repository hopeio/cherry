package initialize

import (
	"github.com/hopeio/cherry/utils/encoding"
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
		struct2Map(reflect.ValueOf(&gc.InitConfig.BasicConfig).Elem(), confMap)
		struct2Map(reflect.ValueOf(&gc.InitConfig.EnvConfig).Elem(), confMap)
		delete(confMap, fixedFieldNameConfigCenter)
	}
	struct2Map(reflect.ValueOf(gc.conf).Elem(), confMap)
	if gc.dao != nil {
		daoConfigTemplateMap(format, reflect.ValueOf(gc.dao).Elem(), confMap)
	}

	encoderRegistry := reflect.ValueOf(gc.Viper).Elem().FieldByName(fixedFieldNameEncoderRegistry).Elem()
	fieldValue := reflect.NewAt(encoderRegistry.Type(), unsafe.Pointer(encoderRegistry.UnsafeAddr()))
	data, err := fieldValue.Interface().(Encoder).Encode(string(format), confMap)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(dir+filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func daoConfigTemplateMap(format encoding.Format, value reflect.Value, confMap map[string]any) {
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Addr().Type().Implements(DaoFieldType) {
			newconfMap := make(map[string]any)
			fieldType := typ.Field(i)
			name := fieldType.Tag.Get(string(format))
			if name == "" {
				name = fieldType.Name
				tagSettings := ParseInitTagSettings(fieldType.Tag.Get(initTagName))
				if tagSettings.ConfigName != "" {
					name = stringsi.UpperCaseFirst(tagSettings.ConfigName)
				}
			}

			confMap[name] = newconfMap
			struct2Map(reflect.ValueOf(field.Addr().Interface().(DaoField).Config()).Elem(), newconfMap)
		}
	}
}
