package reflecti

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

type Request struct {
	FuncName string        `json:"funcName"`
	Params   []interface{} `json:"params"`
}

type Response struct {
	FuncName string        `json:"funcName"`
	Result   []interface{} `json:"result"`
	ErrorMsg string        `json:"errorMsg"`
}

type ReflectInvoker struct {
	Methods map[string]reflect.Value
}

const (
	NoError           = "NoError"
	JsonError         = "JsonError"
	MethodNotFound    = "MethodNotFound"
	ParameterNotMatch = "ParameterNotMatch"
)

func NewReflectInvoker() *ReflectInvoker {
	return &ReflectInvoker{
		Methods: make(map[string]reflect.Value),
	}
}

func (r *ReflectInvoker) RegisterMethod(v interface{}) {
	reflectType := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	for i := 0; i < value.NumMethod(); i++ {
		m := reflectType.Method(i)
		f := value.Method(i)
		r.Methods[reflectType.Elem().Name()+"."+m.Name] = f
	}

}

func (r *ReflectInvoker) InvokeByReflectArgs(funcName string, par []reflect.Value) []reflect.Value {

	return r.Methods[funcName].Call(par)
}

func (r *ReflectInvoker) InvokeByInterfaceArgs(funcName string, Params []interface{}) []reflect.Value {

	paramsValue, err := convertParam(r.Methods[funcName], Params)

	if err != nil {
		return nil
	}

	return r.Methods[funcName].Call(paramsValue)
}

func (r *ReflectInvoker) InvokeByJson(byteData []byte) []byte {

	req := &Request{}
	err := json.Unmarshal(byteData, req)

	resultData := &Response{}

	if err != nil {
		resultData.ErrorMsg = JsonError
	} else {
		resultData.FuncName = req.FuncName

		methodInfo, found := r.Methods[req.FuncName]

		if found {

			paramsValue, err := convertParam(methodInfo, req.Params)

			if err != nil {
				resultData.ErrorMsg = ParameterNotMatch
			} else {
				resultData.Result = InvokeByValues(methodInfo, paramsValue)
			}

		} else {
			resultData.ErrorMsg = MethodNotFound
		}

	}
	resultData.ErrorMsg = NoError
	data, _ := json.Marshal(resultData)

	return data
}

func convertParam(method reflect.Value, params []interface{}) ([]reflect.Value, error) {

	if len(params) != method.Type().NumIn() {
		return nil, errors.New("convertParam number error" + strconv.Itoa(method.Type().NumIn()))
	}

	paramsValue := make([]reflect.Value, 0, len(params))
	//跳过 receiver
	for i := 0; i < len(params); i++ {
		paramsValue = append(paramsValue, reflect.ValueOf(params[i]))
	}

	return paramsValue, nil
}

func InvokeByValues(method reflect.Value, params []reflect.Value) []any {

	var ret []any
	result := method.Call(params)
	for _, x := range result {
		ret = append(ret, x.Interface())
	}

	return ret
}
