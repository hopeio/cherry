package initialize

import (
	"bytes"
	"github.com/hopeio/cherry/utils/encoding"
	"github.com/hopeio/cherry/utils/log"
	"github.com/hopeio/cherry/utils/slices"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"reflect"
	"strings"
)

func (gc *globalConfig) UnmarshalAndSet(data []byte) {
	gc.lock.Lock()
	err := gc.Viper.MergeConfig(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	/*	tmpConfig := gc.cacheConf
		if tmpConfig == nil {
			gc.cacheConf = tmpConfig
		}*/
	tmpConfig := gc.newStruct()
	/*	format, err := common.Unmarshal(gc.ConfigCenter.Format, data, tmpConfig)
		if err != nil {
			log.Fatal(err)
		}*/
	//gc.ConfigCenter.Format = format
	commandLine := newCommandLine()
	injectFlagConfig(commandLine, reflect.ValueOf(tmpConfig).Elem())

	err = gc.Viper.BindPFlags(commandLine)
	if err != nil {
		log.Fatal(err)
	}
	parseFlag(commandLine)

	err = gc.Viper.Unmarshal(tmpConfig, decoderConfigOptions...)
	if err != nil {
		log.Fatal(err)
	}
	gc.inject(tmpConfig)
	gc.lock.Unlock()
	log.Debugf("Configuration:  %+v", tmpConfig)
}

func (gc *globalConfig) newStruct() any {
	nameValueMap := make(map[string]reflect.Value)
	var structFields []reflect.StructField
	confValue := reflect.ValueOf(gc.conf).Elem()
	confType := confValue.Type()
	for i := 0; i < confValue.NumField(); i++ {
		field := confValue.Field(i).Addr()
		if field.CanInterface() {
			inter := field.Interface()
			if c, ok := inter.(InitBeforeInject); ok {
				c.InitBeforeInject()
			}
		}
		structField := confType.Field(i)
		name := structField.Name
		tagSettings := ParseInitTagSettings(structField.Tag.Get(initTagName))
		if tagSettings.ConfigName != "" {
			name = stringsi.UpperCaseFirst(tagSettings.ConfigName)
		}

		structFields = append(structFields, reflect.StructField{
			Name:      name,
			Type:      field.Type(),
			Tag:       structField.Tag,
			Anonymous: structField.Anonymous,
		})

		nameValueMap[name] = field
	}
	// 不进行二次注入,无法确定业务中是否仍然使用,除非每次加锁,或者说每次业务中都交给一个零时变量?需要规范去控制
	if !gc.initialized && gc.dao != nil {

		daoValue := reflect.ValueOf(gc.dao).Elem()
		daoType := daoValue.Type()
		for i := 0; i < daoValue.NumField(); i++ {
			field := daoValue.Field(i)
			if field.Addr().CanInterface() {
				inter := field.Addr().Interface()
				if daofield, ok := inter.(DaoField); ok {
					structField := daoType.Field(i)
					// TODO: 加强校验,必须不为nil
					daoConfig := daofield.Config()
					if daoConfig == nil {
						log.Fatalf("dao %s Config() return nil", structField.Name)
					}

					if c, ok := daoConfig.(InitBeforeInject); ok {
						c.InitBeforeInject()
					}

					name := structField.Name
					daoConfigValue := reflect.ValueOf(daoConfig)
					daoConfigType := reflect.TypeOf(daoConfig)
					tagSettings := ParseInitTagSettings(structField.Tag.Get(initTagName))
					if tagSettings.ConfigName != "" {
						name = stringsi.UpperCaseFirst(tagSettings.ConfigName)
					}
					structFields = append(structFields, reflect.StructField{
						Name: name,
						Type: daoConfigType,
						Tag:  structField.Tag,
					})
					nameValueMap[name] = daoConfigValue
				}
			}
		}
	}
	typ := reflect.StructOf(structFields)
	newStruct := reflect.New(typ)
	gc.setNewStruct(newStruct.Elem(), nameValueMap)
	return newStruct.Interface()
}

func (gc *globalConfig) setNewStruct(value reflect.Value, typValueMap map[string]reflect.Value) {
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		structField := typ.Field(i)
		name := structField.Name
		tagSettings := ParseInitTagSettings(structField.Tag.Get(initTagName))
		if tagSettings.ConfigName != "" {
			name = stringsi.UpperCaseFirst(tagSettings.ConfigName)
		}

		field := value.Field(i)
		field.Set(typValueMap[name])
	}
}

// 注入配置及生成DAO
func (gc *globalConfig) inject(tmpConfig any) {
	confAfterInjectCall(tmpConfig)
	gc.conf.InitAfterInject()
	if !gc.initialized && gc.dao != nil {
		gc.dao.InitAfterInjectConfig()
		injectDao(gc.dao)
	}
}

func confAfterInjectCall(tmpConfig any) {
	v := reflect.ValueOf(tmpConfig).Elem()
	if !v.IsValid() {
		return
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.CanInterface() {
			inter := field.Interface()
			if subconf, ok := inter.(InitAfterInject); ok {
				subconf.InitAfterInject()
			}
		}
	}
}

func injectDao(dao Dao) {
	v := reflect.ValueOf(dao).Elem()
	if !v.IsValid() {
		return
	}
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structFiled := typ.Field(i)
		if field.Addr().CanInterface() {
			inter := field.Addr().Interface()

			if field.Kind() != reflect.Struct {
				log.Debug("ignore inject pointer type: ", field.Type().String())
				continue
			}
			confName := strings.ToUpper(structFiled.Name)
			if slices.Contains(globalConfig1.NoInject, confName) {
				continue
			}

			// 根据DaoField接口实现获取配置和要注入的类型
			if daofield, ok := inter.(DaoField); ok {
				daofield.SetEntity()
			}
		}
	}
	dao.InitAfterInject()
}

// get field name, return filed config name and skip flag
func getFieldConfigName(v reflect.StructField, format encoding.Format) (string, bool) {
	tag := v.Tag.Get(string(format))
	if tag == "" {
		return v.Name, false
	}
	if tag == "-" {
		return "", true
	}
	name, _ := parseTag(tag)
	return name, false
}

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	tag, opt, _ := strings.Cut(tag, ",")
	return tag, tagOptions(opt)
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		if name == optionName {
			return true
		}
	}
	return false
}

// 整合viper,提供单个注入,viper实现也挺简单的,原来的方案就能实现啊
/*func Inject(v Config, path string) {

}
*/
