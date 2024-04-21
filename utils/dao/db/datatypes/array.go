package datatypes

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	reflecti "github.com/hopeio/cherry/utils/reflect/converter"
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
			return errors.New(fmt.Sprint("Failed to scan IntArray value:", value))
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
			return errors.New(fmt.Sprint("Failed to scan FloatArray value:", value))
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
			return errors.New(fmt.Sprint("Failed to scan StringArray value:", value))
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

type Array[T any] []T

func (d *Array[T]) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		data, ok := value.([]byte)
		if !ok {
			return errors.New(fmt.Sprint("Failed to scan StringArray value:", value))
		}
		str = string(data)
	}
	strs := strings.Split(str[1:len(str)-1], ",")
	var arr []T
	for _, elem := range strs {
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
	for i, str := range d {
		buf.WriteString(reflecti.StringFor(str))
		if i != len(d)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return buf.String(), nil
}
