package reflecti

import (
	"errors"
	"fmt"
	"reflect"
)

// 获取引用类型的原始类型
func OriginalType(typ reflect.Type) reflect.Type {
	switch typ.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		return OriginalType(typ.Elem())
	}
	return typ
}

// TODO:
func ContainType() {

}

// SetSubField 设置字段值
// 参数均为指针,dst类型为src field的类型
func SetSubField(src any, sub any) bool {
	srcValue := reflect.ValueOf(src).Elem()
	subValue := reflect.ValueOf(sub).Elem()
	return SetSubFieldValue(srcValue, subValue)
}

// SetSubFieldValue 设置字段值
// subValue为srcValue field的类型
func SetSubFieldValue(srcValue reflect.Value, subValue reflect.Value) bool {
	for i := 0; i < srcValue.NumField(); i++ {
		if srcValue.Field(i).Type() == subValue.Type() {
			srcValue.Field(i).Set(subValue)
			return true
		}
	}
	return false
}

// CopyFieldValueByType 根据类型复制字段
// 参数均为指针,sub类型为src field的类型
func CopyFieldValueByType(src any, sub any) bool {
	srcValue := reflect.ValueOf(src).Elem()
	dstValue := reflect.ValueOf(sub).Elem()
	for i := 0; i < srcValue.NumField(); i++ {
		if srcValue.Field(i).Type() == dstValue.Type() {
			dstValue.Set(srcValue.Field(i))
			return true
		}
	}
	return false
}

func SetField(structValue reflect.Value, name string, value any) error {

	fieldValue := structValue.FieldByName(name)
	if !fieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in obj ", name)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value ", name)
	}

	fieldType := fieldValue.Type()
	val := reflect.ValueOf(value)

	valTypeKind := val.Type().Kind()
	fieldTypeKind := fieldType.Kind()
	if fieldType != val.Type() && val.CanConvert(fieldType) {
		val = val.Convert(fieldType)
	} else {
		return fmt.Errorf("provided value type %s didn't match obj field type %s", valTypeKind, fieldTypeKind)
	}
	fieldValue.Set(val)
	return nil
}

// Copy 复制结构体,浅表复制
func CopyStruct(src any, dest any) error {
	valueOfS := reflect.ValueOf(src)
	typeOfT := reflect.TypeOf(dest)

	if valueOfS.Kind() == reflect.Ptr {
		valueOfS = valueOfS.Elem()
	}

	if valueOfS.Kind() != reflect.Struct {
		return errors.New("src is not a ptr or struct")
	}

	if typeOfT.Kind() != reflect.Ptr {
		return errors.New("dest is not a ptr")
	}

	typeOfT = typeOfT.Elem()
	valueOfT := reflect.ValueOf(dest).Elem()

	for i := 0; i < typeOfT.NumField(); i++ {
		// 获取每个成员的结构体字段值
		fieldType := typeOfT.Field(i)
		// 赋值
		valueOfT.Field(i).Set(valueOfS.FieldByName(fieldType.Name))
	}

	return nil
}
