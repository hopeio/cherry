package initialize

import (
	"bytes"
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
		if gc.editTimes == 0 {
			log.Fatal(err)
		} else {
			log.Error(err)
			return
		}
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

	err = gc.Viper.Unmarshal(tmpConfig, decoderConfigOptions...)
	if err != nil {
		if gc.editTimes == 0 {
			log.Fatal(err)
		} else {
			log.Error(err)
			return
		}
	}
	applyFlagConfig(gc.Viper, tmpConfig)

	gc.inject(tmpConfig)
	gc.editTimes++
	gc.lock.Unlock()
	log.Debugf("Configuration:  %+v", tmpConfig)
}

func (gc *globalConfig) newStruct() any {
	nameValueMap := make(map[string]reflect.Value)
	var structFields []reflect.StructField

	confValue := reflect.ValueOf(&gc.BuiltinConfig).Elem()
	confType := confValue.Type()
	for i := 0; i < confValue.NumField(); i++ {
		field := confValue.Field(i).Addr()

		structField := confType.Field(i)
		name := structField.Name
		tagSettings := ParseInitTagSettings(structField.Tag.Get(initTagName))
		if tagSettings.ConfigName != "" {
			name = stringsi.UpperCaseFirst(tagSettings.ConfigName)
		}

		if field.CanInterface() {
			inter := field.Interface()
			if c, ok := inter.(InitBeforeInject); ok {
				c.InitBeforeInject()
			}
			if c, ok := inter.(InitBeforeInjectWithInitConfig); ok {
				c.InitBeforeInjectWithInitConfig(&gc.InitConfig)
			}
		}
		structFields = append(structFields, reflect.StructField{
			Name:      name,
			Type:      field.Type(),
			Tag:       structField.Tag,
			Anonymous: structField.Anonymous,
		})

		nameValueMap[name] = field
	}
	confValue = reflect.ValueOf(gc.conf).Elem()
	confType = confValue.Type()
	for i := 0; i < confValue.NumField(); i++ {
		field := confValue.Field(i)
		fieldType := field.Type()
		// panic: reflect: embedded type with methods not implemented if type is not first field // Issue 15924.
		if confValue.Field(i).Type() == EmbeddedPresetsType {
			continue
		}
		if fieldType.Kind() != reflect.Ptr && fieldType.Kind() != reflect.Map {
			field = field.Addr()
		}

		structField := confType.Field(i)
		name := structField.Name
		tagSettings := ParseInitTagSettings(structField.Tag.Get(initTagName))
		if tagSettings.ConfigName != "" {
			name = stringsi.UpperCaseFirst(tagSettings.ConfigName)
		}

		if v, ok := nameValueMap[name]; ok {
			if v.Type() == field.Type() {
				log.Fatalf(`exists builtin config field: %s, please delete the field`, name)
			} else {
				log.Fatalf(`exists builtin config field: %s, please rename or use init tag [init:"config:{{other config name}}"]`, name)
			}
		}

		if field.CanInterface() {
			inter := field.Interface()
			if c, ok := inter.(InitBeforeInject); ok {
				c.InitBeforeInject()
			}
			if c, ok := inter.(InitBeforeInjectWithInitConfig); ok {
				c.InitBeforeInjectWithInitConfig(&gc.InitConfig)
			}
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
			if field.Type().Kind() == reflect.Struct {
				field = field.Addr()
			}
			if field.CanInterface() {
				inter := field.Interface()
				if daoField, ok := inter.(DaoField); ok {

					structField := daoType.Field(i)

					// TODO: 加强校验,必须不为nil
					daoConfig := daoField.Config()
					if daoConfig == nil {
						log.Fatalf("dao %s Config() return nil", structField.Name)
					}

					name := structField.Name
					daoConfigValue := reflect.ValueOf(daoConfig)
					daoConfigType := reflect.TypeOf(daoConfig)
					tagSettings := ParseInitTagSettings(structField.Tag.Get(initTagName))
					if tagSettings.ConfigName != "" {
						name = stringsi.UpperCaseFirst(tagSettings.ConfigName)
					}

					if c, ok := daoConfig.(InitBeforeInject); ok {
						c.InitBeforeInject()
					}
					if c, ok := inter.(InitBeforeInjectWithInitConfig); ok {
						c.InitBeforeInjectWithInitConfig(&gc.InitConfig)
					}

					if _, ok := nameValueMap[name]; ok {
						log.Fatalf(`exists field: %s, please rename or use init tag [init:"{{otherConfigName}}"]`, name)
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
	gc.confAfterInjectCall(tmpConfig)
	gc.conf.InitAfterInject()
	if c, ok := gc.conf.(InitAfterInjectWithInitConfig); ok {
		c.InitAfterInjectWithInitConfig(&gc.InitConfig)
	}
	if !gc.initialized && gc.dao != nil {
		gc.dao.InitAfterInjectConfig()
		if c, ok := gc.conf.(InitAfterInjectConfigWithInitConfig); ok {
			c.InitAfterInjectConfigWithInitConfig(&gc.InitConfig)
		}
		gc.injectDao()
	}
}

func (gc *globalConfig) confAfterInjectCall(tmpConfig any) {
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
			if subconf, ok := inter.(InitAfterInjectWithInitConfig); ok {
				subconf.InitAfterInjectWithInitConfig(&gc.InitConfig)
			}
		}
	}
}

func (gc *globalConfig) injectDao() {
	v := reflect.ValueOf(gc.dao).Elem()
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
			if slices.Contains(gConfig.InitConfig.NoInject, confName) {
				continue
			}

			// 根据DaoField接口实现获取配置和要注入的类型
			if daofield, ok := inter.(DaoField); ok {
				daofield.SetEntity()
			}
		}
	}
	gc.dao.InitAfterInject()
	if c, ok := gc.dao.(InitAfterInjectWithInitConfig); ok {
		c.InitAfterInjectWithInitConfig(&gc.InitConfig)
	}
}
