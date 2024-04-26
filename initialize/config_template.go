package initialize

import (
	"github.com/hopeio/cherry/utils/encoding"
	"github.com/hopeio/cherry/utils/log"
	"github.com/hopeio/cherry/utils/slices"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"os"
	"reflect"
	"unsafe"
)

func (gc *globalConfig) genConfigTemplate(singleFileConfig bool) {
	dir := gc.ConfigTemplateDir
	if dir == "" {
		return
	}
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}

	format := gc.ConfigCenter.Format
	filename := "local.template." + string(format)

	confMap := make(map[string]any)
	if singleFileConfig {
		filename = "config.template." + string(format)
		struct2MapHelper(format, reflect.ValueOf(&gc.BasicConfig).Elem(), confMap)
		struct2MapHelper(format, reflect.ValueOf(&gc.EnvConfig).Elem(), confMap)
		delete(confMap, "ConfigCenter")
	}
	struct2MapHelper(format, reflect.ValueOf(gc.conf).Elem(), confMap)
	if gc.dao != nil {
		daoConfigTemplateMap(format, reflect.ValueOf(gc.dao).Elem(), confMap)
	}

	encoderRegistry := reflect.ValueOf(gc.Viper).Elem().FieldByName("encoderRegistry").Elem()
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

func struct2Map(format encoding.Format, value reflect.Value) map[string]any {
	m := make(map[string]any)
	struct2MapHelper(format, value, m)
	return m
}

// 递归的根据反射将对象中的指针变量赋值
func struct2MapHelper(format encoding.Format, value reflect.Value, confMap map[string]any) {
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fileKind := field.Kind()
		fieldType := typ.Field(i)
		// 判断field是否大写
		if !fieldType.IsExported() {
			continue
		}
		switch fileKind {
		case reflect.Func, reflect.Chan, reflect.Interface:
			continue
		case reflect.Slice, reflect.Map, reflect.Array:
			if slices.Contains([]reflect.Kind{reflect.Func, reflect.Chan, reflect.Interface}, fieldType.Type.Elem().Kind()) {
				continue
			}
		case reflect.Ptr, reflect.Struct:
			if field.CanSet() {
				// 如果是tls.Config 类型，则不处理,这里可能会干扰其他相同的定义
				typName := fieldType.Type.String()
				if fileKind == reflect.Ptr {
					typName = field.Type().Elem().String()
				}
				if typName == "tls.Config" {
					continue
				}
				newValue := field
				if fileKind == reflect.Ptr {
					newValue = reflect.New(field.Type().Elem()).Elem()
				}

				// 判断是匿名字段
				name := fieldType.Tag.Get(string(format))

				if name == "" && fieldType.Anonymous {
					struct2MapHelper(format, newValue, confMap)
				} else {
					if name == "" {
						name = fieldType.Name
						tagSettings := ParseInitTagSettings(fieldType.Tag.Get(initTagName))
						if tagSettings.ConfigName != "" {
							name = stringsi.UpperCaseFirst(tagSettings.ConfigName)
						}
					}

					newconfMap := make(map[string]any)
					confMap[name] = newconfMap
					struct2MapHelper(format, newValue, newconfMap)
					if len(newconfMap) == 0 {
						delete(confMap, name)
					}
				}
			}
			continue
		}

		if field.CanInterface() {
			name := fieldType.Tag.Get(string(format))
			if name == "" {
				name = fieldType.Name
			}
			confMap[name] = field.Interface()
		}
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
			struct2MapHelper(format, reflect.ValueOf(field.Addr().Interface().(DaoField).Config()).Elem(), newconfMap)
		}
	}
}
