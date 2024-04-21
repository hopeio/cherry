package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JSONB json.RawMessage

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *JSONB) Scan(value interface{}) error {
	switch bytes := value.(type) {
	case []byte:
		*j = append((*j)[0:0], bytes...)
		return nil
	case string:
		*j = JSONB(bytes)
		return nil
	default:
		return errors.New(fmt.Sprint("Failed to scan JSONB value:", value))
	}

}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j JSONB) Value() (driver.Value, error) {
	return json.RawMessage(j).MarshalJSON()
}

func (*JSONB) GormDataType() string {
	return "jsonb"
}

/*
func (JSONB) ConvertValue(v any) (driver.Value, error) {

}
*/
type JSON map[string]any

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Json
func (j *JSON) Scan(value interface{}) error {
	switch bytes := value.(type) {
	case []byte:
		*j = map[string]any{}
		return json.Unmarshal(bytes, j)
	case string:
		*j = map[string]any{}
		return json.Unmarshal([]byte(bytes), j)
	default:
		return errors.New(fmt.Sprint("Failed to scan JSON value:", value))
	}
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return json.Marshal(j)
}

func (JSON) GormDataType() string {
	return "jsonb"
}

type JSONArray []map[string]any

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Json
func (j *JSONArray) Scan(value interface{}) error {
	switch bytes := value.(type) {
	case []byte:
		*j = make([]map[string]any, 0)
		return json.Unmarshal(bytes, j)
	case string:
		*j = make([]map[string]any, 0)
		return json.Unmarshal([]byte(bytes), j)
	default:
		return errors.New(fmt.Sprint("Failed to scan JSONArray value:", value))
	}
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return []byte("null"), nil
	}
	if len(j) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(j)
}

func (*JSONArray) GormDataType() string {
	return "jsonb"
}

type JSONStr string

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 JsonStr
func (j *JSONStr) Scan(value interface{}) error {
	switch bytes := value.(type) {
	case []byte:
		*j = JSONStr(bytes)
		return nil
	case string:
		*j = JSONStr(bytes)
		return nil
	default:
		return errors.New(fmt.Sprint("Failed to scan JSONStr value:", value))
	}
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j JSONStr) Value() (driver.Value, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return string(j), nil
}

func (*JSONStr) GormDataType() string {
	return "jsonb"
}

type JSONT[T any] struct {
	JSON *T
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Json
func (j *JSONT[T]) Scan(value interface{}) error {
	switch bytes := value.(type) {
	case []byte:
		j.JSON = new(T)
		return json.Unmarshal(bytes, j)
	case string:
		j.JSON = new(T)
		return json.Unmarshal([]byte(bytes), j)
	default:
		return errors.New(fmt.Sprint("Failed to scan JSON value:", value))
	}
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j JSONT[T]) Value() (driver.Value, error) {
	if j.JSON == nil {
		return []byte("null"), nil
	}
	return json.Marshal(j)
}

func (*JSONT[T]) GormDataType() string {
	return "jsonb"
}

type JSONArrayT[T any] []T

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Json
func (j *JSONArrayT[T]) Scan(value interface{}) error {
	switch bytes := value.(type) {
	case []byte:
		*j = make(JSONArrayT[T], 0)
		return json.Unmarshal(bytes, j)
	case string:
		*j = make(JSONArrayT[T], 0)
		return json.Unmarshal([]byte(bytes), j)
	default:
		return errors.New(fmt.Sprint("Failed to scan JSON value:", value))
	}
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j JSONArrayT[T]) Value() (driver.Value, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return json.Marshal(j)
}

func (*JSONArrayT[T]) GormDataType() string {
	return "jsonb"
}
