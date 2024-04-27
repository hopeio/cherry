// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package converter

import (
	"encoding"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"reflect"
	"strconv"
	"strings"
)

type StringConverter func(string) any
type StringConverterE func(string) (any, error)

func (c StringConverterE) Converter() StringConverter {
	if c == nil {
		return nil
	}
	return func(value string) any {
		r, _ := c(value)
		return r
	}
}

var (
	invalidValue = reflect.Value{}
)

// Default converters for basic types.
/*var ConverterMaps = map[reflect.Kind]StringConverter{
	reflect.Bool:    stringConvertBool,
	reflect.Float32: stringConvertFloat32,
	reflect.Float64: stringConvertFloat64,
	reflect.Int:     stringConvertInt,
	reflect.Int8:    stringConvertInt8,
	reflect.Int16:   stringConvertInt16,
	reflect.Int32:   stringConvertInt32,
	reflect.Int64:   stringConvertInt64,
	reflect.String:  stringConvertString,
	reflect.Uint:    stringConvertUint,
	reflect.Uint8:   stringConvertUint8,
	reflect.Uint16:  stringConvertUint16,
	reflect.Uint32:  stringConvertUint32,
	reflect.Uint64:  stringConvertUint64,
}*/

// Deprecated: unsupported slices array map
var StringConverterArrays = []StringConverterE{
	reflect.Bool:    stringConvertBool,
	reflect.Int:     stringConvertInt,
	reflect.Int8:    stringConvertInt8,
	reflect.Int16:   stringConvertInt16,
	reflect.Int32:   stringConvertInt32,
	reflect.Int64:   stringConvertInt64,
	reflect.Uint:    stringConvertUint,
	reflect.Uint8:   stringConvertUint8,
	reflect.Uint16:  stringConvertUint16,
	reflect.Uint32:  stringConvertUint32,
	reflect.Uint64:  stringConvertUint64,
	reflect.Float32: stringConvertFloat32,
	reflect.Float64: stringConvertFloat64,
	reflect.String:  stringConvertString,
}

const (
	array        = 100
	ArrayBool    = reflect.Bool + array
	ArrayInt     = reflect.Int + array
	ArrayInt8    = reflect.Int8 + array
	ArrayInt16   = reflect.Int16 + array
	ArrayInt32   = reflect.Int32 + array
	ArrayInt64   = reflect.Int64 + array
	ArrayUint    = reflect.Uint + array
	ArrayUint8   = reflect.Uint8 + array
	ArrayUint16  = reflect.Uint16 + array
	ArrayUint32  = reflect.Uint32 + array
	ArrayUint64  = reflect.Uint64 + array
	ArrayFloat32 = reflect.Float32 + array
	ArrayFloat64 = reflect.Float64 + array
	ArrayString  = reflect.String + array
)

const (
	slice        = 200
	SliceBool    = reflect.Bool + slice
	SliceInt     = reflect.Int + slice
	SliceInt8    = reflect.Int8 + slice
	SliceInt16   = reflect.Int16 + slice
	SliceInt32   = reflect.Int32 + slice
	SliceInt64   = reflect.Int64 + slice
	SliceUint    = reflect.Uint + slice
	SliceUint8   = reflect.Uint8 + slice
	SliceUint16  = reflect.Uint16 + slice
	SliceUint32  = reflect.Uint32 + slice
	SliceUint64  = reflect.Uint64 + slice
	SliceFloat32 = reflect.Float32 + slice
	SliceFloat64 = reflect.Float64 + slice
	SliceString  = reflect.String + slice
)

func stringConvertBool(value string) (any, error) {
	if value == "on" {
		return true, nil
	}
	return strconv.ParseBool(value)
}

func stringConvertFloat32(value string) (any, error) {
	return strconv.ParseFloat(value, 32)
}

func stringConvertFloat64(value string) (any, error) {
	return strconv.ParseFloat(value, 64)
}

func stringConvertInt(value string) (any, error) {
	return strconv.ParseInt(value, 10, 0)
}

func stringConvertInt8(value string) (any, error) {
	return strconv.ParseInt(value, 10, 8)
}

func stringConvertInt16(value string) (any, error) {
	return strconv.ParseInt(value, 10, 16)
}

func stringConvertInt32(value string) (any, error) {
	return strconv.ParseInt(value, 10, 32)
}

func stringConvertInt64(value string) (any, error) {
	return strconv.ParseInt(value, 10, 64)
}

func stringConvertString(value string) (any, error) {
	return value, nil
}

func stringConvertArray(value string) (any, error) {
	return value, nil
}

func stringConvertUint(value string) (any, error) {
	return strconv.ParseUint(value, 10, 0)
}

func stringConvertUint8(value string) (any, error) {
	return strconv.ParseUint(value, 10, 8)
}

func stringConvertUint16(value string) (any, error) {
	return strconv.ParseUint(value, 10, 16)
}

func stringConvertUint32(value string) (any, error) {
	return strconv.ParseUint(value, 10, 32)
}

func stringConvertUint64(value string) (any, error) {
	return strconv.ParseUint(value, 10, 64)
}

func ConvertInt64(v interface{}) int64 {
	return cast.ToInt64(v)
}

func SetFieldByString(dst any, field, value string) error {
	if value == "" {
		return nil
	}

	fieldValue := reflect.ValueOf(dst).Elem().FieldByName(field)
	return SetValueByString(fieldValue, value)
}

func SetValueByString(field reflect.Value, value string) error {
	if value == "" {
		return nil
	}

	v := field.Interface()
	if t, ok := v.(encoding.TextUnmarshaler); ok {
		return t.UnmarshalText([]byte(value))
	}
	kind := field.Kind()
	if kind == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
		v = field.Interface()
		if t, ok := v.(encoding.TextUnmarshaler); ok {
			return t.UnmarshalText([]byte(value))
		}
		kind = field.Kind()
	}
	switch kind {
	case reflect.Array, reflect.Slice:
		subType := field.Type().Elem()
		eKind := subType.Kind()
		if eKind == reflect.Array || eKind == reflect.Slice || eKind == reflect.Map {
			return fmt.Errorf("unsupported sub type %v", subType)
		}
		strs := strings.Split(value, ",")
		if kind == reflect.Slice {
			field.Set(reflect.MakeSlice(field.Type(), len(strs), len(strs)))
		}
		for i := 0; i < field.Len(); i++ {
			if err := SetValueByString(field.Index(i), strs[i]); err != nil {
				return err
			}
		}
	case reflect.Map:
		subType := field.Type().Elem()
		eKind := subType.Kind()
		if eKind == reflect.Array || eKind == reflect.Slice || eKind == reflect.Map {
			return fmt.Errorf("unsupported sub type %v", subType)
		}
		strs := strings.Split(value, ",")
		field.Set(reflect.MakeMapWithSize(field.Type(), len(strs)/2))
		for i := 0; i < len(strs)/2; i += 2 {
			key := reflect.New(field.Type().Key())
			err := SetValueByString(key, strs[i])
			if err != nil {
				return err
			}
			v := reflect.New(field.Type().Elem())
			err = SetValueByString(v, strs[i+1])
			if err != nil {
				return err
			}
			field.SetMapIndex(key, v)
		}
	}

	v, err := StringConvert(kind, value)
	if err == nil {
		field.Set(reflect.ValueOf(v))
		return nil
	}
	return err
}

func StringConvert(kind reflect.Kind, value string) (any, error) {
	converter := StringConverterArrays[kind]
	if converter != nil {
		return converter(value)
	}
	return nil, errors.New("unsupported kind")
}

func StringConvertFor[T any](value string) (T, error) {
	kind := reflect.TypeFor[T]().Kind()
	converter := StringConverterArrays[kind]
	if converter != nil {
		if v, err := converter(value); err != nil {
			return *new(T), err
		} else {
			return v.(T), nil
		}
	}
	return *new(T), errors.New("unsupported kind")
}

func String(value reflect.Value) string {
	v := value.Interface()
	if t, ok := v.(encoding.TextMarshaler); ok {
		s, _ := t.MarshalText()
		return string(s)
	}

	kind := value.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Pointer, reflect.UnsafePointer:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.String:
		return value.String()
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float64, reflect.Float32:
		return strconv.FormatFloat(value.Float(), 'g', -1, 64)
	case reflect.Array, reflect.Slice:
		var strs []string
		for i := 0; i < value.Len(); i++ {
			strs = append(strs, String(value.Index(i)))
		}
		return strings.Join(strs, ",")
	}
	return ""
}

func StringFor[T any](t T) string {
	v := reflect.ValueOf(t)
	return String(v)
}
