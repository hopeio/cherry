package datatypes

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding"
	"errors"
	"fmt"
	reflecti "github.com/hopeio/cherry/utils/reflect/converter"
	"time"

	"golang.org/x/exp/constraints"
	"strconv"
	"strings"
)

// adpter postgres
type IntArray[T constraints.Integer] []T

func (d *IntArray[T]) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		data, ok := value.([]byte)
		if !ok {
			return errors.New(fmt.Sprint("failed to scan int array value:", value))
		}
		str = string(data)
	}
	strs := strings.Split(str[1:len(str)-1], ",")
	var arr []T
	for _, numstr := range strs {
		num, err := strconv.Atoi(numstr)
		if err != nil {
			return err
		}
		arr = append(arr, T(num))
	}
	*d = arr
	return nil
}

func (d IntArray[T]) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, num := range d {
		buf.WriteString(strconv.Itoa(int(num)))
		if i != len(d)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return buf.String(), nil
}

type FloatArray[T constraints.Float] []T

func (d *FloatArray[T]) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		data, ok := value.([]byte)
		if !ok {
			return errors.New(fmt.Sprint("failed to scan float array value:", value))
		}
		str = string(data)
	}
	strs := strings.Split(str[1:len(str)-1], ",")
	var arr []T
	for _, numstr := range strs {
		num, err := strconv.ParseFloat(numstr, 64)
		if err != nil {
			return err
		}
		arr = append(arr, T(num))
	}
	*d = arr
	return nil
}

func (d FloatArray[T]) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, num := range d {
		buf.WriteString(strconv.FormatFloat(float64(num), 'g', -1, 64))
		if i != len(d)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return buf.String(), nil
}

type StringArray []string

func (d *StringArray) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		data, ok := value.([]byte)
		if !ok {
			return errors.New(fmt.Sprint("failed to scan string array value:", value))
		}
		str = string(data)
	}
	strs := strings.Split(str[1:len(str)-1], ",")
	var arr []string
	for _, elem := range strs {
		arr = append(arr, elem)
	}
	*d = arr
	return nil
}

func (d StringArray) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, str := range d {
		buf.WriteString(str)
		if i != len(d)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return buf.String(), nil
}

// Array represents a PostgreSQL array for T. It implements the ArrayGetter and ArraySetter interfaces. It preserves
// PostgreSQL dimensions and custom lower bounds. Use FlatArray if these are not needed.
// 只支持一维数组,unsupported box
type Array[T any] []T

func (d *Array[T]) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		data, ok := value.([]byte)
		if !ok {
			return errors.New(fmt.Sprint("failed to scan array value:", value))
		}
		str = string(data)
	}
	strs := strings.Split(str[1:len(str)-1], ",")
	var arr []T
	for _, elem := range strs {
		var t T
		a := any(t)
		if v, ok := a.(sql.Scanner); ok {
			err := v.Scan(elem)
			if err != nil {
				return err
			}
			arr = append(arr, t)
			continue
		}
		if v, ok := a.(encoding.TextUnmarshaler); ok {
			err := v.UnmarshalText([]byte(elem))
			if err != nil {
				return err
			}
			arr = append(arr, t)
			continue
		}
		v, err := reflecti.StringConvertFor[T](elem)
		if err != nil {
			return err
		}
		arr = append(arr, v)
	}
	*d = arr
	return nil
}

func (d Array[T]) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, v := range d {
		a := any(v)
		if vv, ok := a.(driver.Valuer); ok {
			v, err := vv.Value()
			if err != nil {
				return nil, err
			}
			buf.WriteString(reflecti.StringFor(v))
			continue
		}
		if vv, ok := a.(encoding.TextMarshaler); ok {
			v, err := vv.MarshalText()
			if err != nil {
				return nil, err
			}
			buf.Write(v)
			continue
		}
		buf.WriteString(reflecti.StringFor(v))
		if i != len(d)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return buf.String(), nil
}

type TimeArray []time.Time

func (d *TimeArray) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		data, ok := value.([]byte)
		if !ok {
			return errors.New(fmt.Sprint("failed to scan string array value:", value))
		}
		str = string(data)
	}
	strs := strings.Split(str[1:len(str)-1], ",")
	var arr []time.Time
	for _, elem := range strs {
		t, err := time.Parse(time.RFC3339Nano, elem)
		if err != nil {
			return err
		}
		arr = append(arr, t)
	}
	*d = arr
	return nil
}

func (d TimeArray) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, t := range d {
		buf.WriteString(t.Format(time.RFC3339Nano))
		if i != len(d)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return buf.String(), nil
}
