// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package converter

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

type Converter func(string) any

var (
	invalidValue = reflect.Value{}
)

// Default converters for basic types.
/*var ConverterMaps = map[reflect.Kind]Converter{
	reflect.Bool:    convertBool,
	reflect.Float32: convertFloat32,
	reflect.Float64: convertFloat64,
	reflect.Int:     convertInt,
	reflect.Int8:    convertInt8,
	reflect.Int16:   convertInt16,
	reflect.Int32:   convertInt32,
	reflect.Int64:   convertInt64,
	reflect.String:  convertString,
	reflect.Uint:    convertUint,
	reflect.Uint8:   convertUint8,
	reflect.Uint16:  convertUint16,
	reflect.Uint32:  convertUint32,
	reflect.Uint64:  convertUint64,
}*/

// TODO: support slices map
var ConverterArrays = []Converter{
	reflect.Bool:    convertBool,
	reflect.Int:     convertInt,
	reflect.Int8:    convertInt8,
	reflect.Int16:   convertInt16,
	reflect.Int32:   convertInt32,
	reflect.Int64:   convertInt64,
	reflect.Uint:    convertUint,
	reflect.Uint8:   convertUint8,
	reflect.Uint16:  convertUint16,
	reflect.Uint32:  convertUint32,
	reflect.Uint64:  convertUint64,
	reflect.Float32: convertFloat32,
	reflect.Float64: convertFloat64,
	reflect.String:  convertString,
}

func convertBool(value string) any {
	if value == "on" {
		return true
	} else if v, err := strconv.ParseBool(value); err == nil {
		return v
	}
	return nil
}

func convertFloat32(value string) any {
	if v, err := strconv.ParseFloat(value, 32); err == nil {
		return v
	}
	return nil
}

func convertFloat64(value string) any {
	if v, err := strconv.ParseFloat(value, 64); err == nil {
		return v
	}
	return nil
}

func convertInt(value string) any {
	if v, err := strconv.ParseInt(value, 10, 0); err == nil {
		return v
	}
	return nil
}

func convertInt8(value string) any {
	if v, err := strconv.ParseInt(value, 10, 8); err == nil {
		return v
	}
	return nil
}

func convertInt16(value string) any {
	if v, err := strconv.ParseInt(value, 10, 16); err == nil {
		return v
	}
	return nil
}

func convertInt32(value string) any {
	if v, err := strconv.ParseInt(value, 10, 32); err == nil {
		return v
	}
	return nil
}

func convertInt64(value string) any {
	if v, err := strconv.ParseInt(value, 10, 64); err == nil {
		return v
	}
	return nil
}

func convertString(value string) any {
	return value
}

func convertUint(value string) any {
	if v, err := strconv.ParseUint(value, 10, 0); err == nil {
		return v
	}
	return nil
}

func convertUint8(value string) any {
	if v, err := strconv.ParseUint(value, 10, 8); err == nil {
		return v
	}
	return nil
}

func convertUint16(value string) any {
	if v, err := strconv.ParseUint(value, 10, 16); err == nil {
		return reflect.ValueOf(uint16(v))
	}
	return nil
}

func convertUint32(value string) any {
	if v, err := strconv.ParseUint(value, 10, 32); err == nil {
		return v
	}
	return nil
}

func convertUint64(value string) any {
	if v, err := strconv.ParseUint(value, 10, 64); err == nil {
		return v
	}
	return nil
}

func ConvertInt64(v interface{}) int64 {
	switch v := v.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case json.Number:
		f, _ := v.Int64()
		return f
	default:
		return 0
	}
}

func SetStructFieldByString(field, value string, dst any) error {
	if value == "" {
		return nil
	}
	fieldValue := reflect.ValueOf(dst).Elem().FieldByName(field)
	return SetFieldByString(value, fieldValue)
}

func SetFieldByString(value string, field reflect.Value) error {
	if value == "" {
		return nil
	}
	converter := ConverterArrays[field.Kind()]
	if converter != nil {
		if v := converter(value); v != nil {
			field.Set(reflect.ValueOf(v))
			return nil
		}
	}
	return errors.New("unsupported kind")
}

func StringConvert(value string, kind reflect.Kind) (any, error) {
	converter := ConverterArrays[kind]
	if converter != nil {
		if v := converter(value); v != nil {
			return v, nil
		}
	}
	return nil, errors.New("unsupported kind")
}

func StringConvertFor[T any](value string) (T, error) {
	kind := reflect.TypeFor[T]().Kind()
	converter := ConverterArrays[kind]
	if converter != nil {
		if v := converter(value); v != nil {
			return v.(T), nil
		}
	}
	return *new(T), errors.New("unsupported kind")
}

func AnyConvert[T any](v any) T {
	return v.(T)
}

func String(value reflect.Value) string {
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
	}
	return ""
}

func StringFor[T any](t T) string {
	v := reflect.ValueOf(t)
	return String(v)
}
