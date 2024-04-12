package slices

import "reflect"

type stream struct {
	kind    reflect.Kind
	entries reflect.Value
	expr    any
}

func Stream(array any) stream {
	value := reflect.ValueOf(array)
	kind := value.Kind()
	if kind == reflect.Pointer {
		value = value.Elem()
		kind = value.Kind()
	}
	if kind != reflect.Array && kind != reflect.Slice && kind != reflect.Map && kind != reflect.Chan && kind != reflect.String {
		panic("参数必须为数组,切片,地图,管道,字符串")
	}
	return stream{
		kind:    kind,
		entries: value,
	}
}

func (s stream) ForEach(express any) {
	call := reflect.ValueOf(express)
	kind := call.Kind()
	if kind != reflect.Func {
		panic("参数必须为函数")
	}

	switch s.kind {
	case reflect.Slice, reflect.String, reflect.Array:
		for i := 0; i < s.entries.Len(); i++ {
			call.Call([]reflect.Value{s.entries.Index(i)})
		}
	case reflect.Map:
		iter := s.entries.MapRange()
		for iter.Next() {
			call.Call([]reflect.Value{iter.Value()})
		}
	case reflect.Chan: // blocked
		item, ok := s.entries.Recv()
		if ok {
			call.Call([]reflect.Value{item})
		}
	}
}

func (s stream) Map(express any) any {
	call := reflect.ValueOf(express)
	kind := call.Kind()
	if kind != reflect.Func {
		panic("参数必须为函数")
	}
	if call.Type().NumOut() != 1 {
		panic("函数必须有1个返回值")
	}

	ret := make([]interface{}, 0, s.entries.Len())
	switch s.kind {
	case reflect.Slice, reflect.String, reflect.Array:
		for i := 0; i < s.entries.Len(); i++ {
			out := call.Call([]reflect.Value{s.entries.Index(i)})
			ret = append(ret, out[0].Interface())
		}
	case reflect.Map:
		iter := s.entries.MapRange()
		for iter.Next() {
			out := call.Call([]reflect.Value{iter.Value()})
			ret = append(ret, out[0].Interface())
		}
	case reflect.Chan: // blocked
		item, ok := s.entries.Recv()
		if ok {
			out := call.Call([]reflect.Value{item})
			ret = append(ret, out[0].Interface())
		}
	}
	return interface{}(ret)
}

func (s stream) Filter(express any) any {
	call := reflect.ValueOf(express)
	kind := call.Kind()
	if kind != reflect.Func {
		panic("参数必须为函数")
	}
	typ := call.Type()
	if typ.NumOut() != 1 || typ.Out(0).Kind() != reflect.Bool {
		panic("函数返回值必须为bool")
	}
	ret := make([]interface{}, 0, s.entries.Len())
	switch s.kind {
	case reflect.Slice, reflect.String, reflect.Array:
		for i := 0; i < s.entries.Len(); i++ {
			out := call.Call([]reflect.Value{s.entries.Index(i)})
			if !out[0].Bool() {
				ret = append(ret, s.entries.Index(i).Interface())
			}
		}
	case reflect.Map:
		iter := s.entries.MapRange()
		for iter.Next() {
			out := call.Call([]reflect.Value{iter.Value()})
			if !out[0].Bool() {
				ret = append(ret, iter.Value().Interface())
			}
		}
	case reflect.Chan: // blocked
		item, ok := s.entries.Recv()
		if ok {
			out := call.Call([]reflect.Value{item})
			if !out[0].Bool() {
				ret = append(ret, item.Interface())
			}
		}
	}
	return interface{}(ret)
}

func (s stream) Every(express any) bool {
	call := reflect.ValueOf(express)
	kind := call.Kind()
	if kind != reflect.Func {
		panic("参数必须为函数")
	}
	typ := call.Type()
	if typ.NumOut() != 1 || typ.Out(0).Kind() != reflect.Bool {
		panic("函数返回值必须为bool")
	}
	switch s.kind {
	case reflect.Slice, reflect.String, reflect.Array:
		for i := 0; i < s.entries.Len(); i++ {
			out := call.Call([]reflect.Value{s.entries.Index(i)})
			if !out[0].Bool() {
				return false
			}
		}
	case reflect.Map:
		iter := s.entries.MapRange()
		for iter.Next() {
			out := call.Call([]reflect.Value{iter.Value()})
			if !out[0].Bool() {
				return false
			}
		}
	case reflect.Chan: // blocked
		item, ok := s.entries.Recv()
		if ok {
			out := call.Call([]reflect.Value{item})
			if !out[0].Bool() {
				return false
			}
		}
	}
	return false
}

func (s stream) Some(express any) bool {
	call := reflect.ValueOf(express)
	kind := call.Kind()
	if kind != reflect.Func {
		panic("参数必须为函数")
	}
	typ := call.Type()
	if typ.NumOut() != 1 || typ.Out(0).Kind() != reflect.Bool {
		panic("函数返回值必须为bool")
	}

	switch s.kind {
	case reflect.Slice, reflect.String, reflect.Array:
		for i := 0; i < s.entries.Len(); i++ {
			out := call.Call([]reflect.Value{s.entries.Index(i)})
			if !out[0].Bool() {
				return false
			}
		}
	case reflect.Map:
		iter := s.entries.MapRange()
		for iter.Next() {
			out := call.Call([]reflect.Value{iter.Value()})
			if !out[0].Bool() {
				return false
			}
		}
	case reflect.Chan: // blocked
		item, ok := s.entries.Recv()
		if ok {
			out := call.Call([]reflect.Value{item})
			if !out[0].Bool() {
				return false
			}
		}
	}
	return false
}
