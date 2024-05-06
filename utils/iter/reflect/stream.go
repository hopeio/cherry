package iter

import "reflect"

type stream struct {
	kind    reflect.Kind
	entries reflect.Value
	exprs   []expr
}

type exprKind int

type expr struct {
	kind exprKind
	expr any
}

const (
	mapExpr exprKind = iota
	filterExpr
	forEachExpr
)

func StreamOf(array any) *stream {
	value := reflect.ValueOf(array)
	kind := value.Kind()
	if kind == reflect.Pointer {
		value = value.Elem()
		kind = value.Kind()
	}
	if kind != reflect.Array && kind != reflect.Slice && kind != reflect.Map && kind != reflect.Chan && kind != reflect.String {
		panic("参数必须为数组,切片,地图,管道,字符串")
	}
	return &stream{
		kind:    kind,
		entries: value,
	}
}

func (s *stream) ForEach(express any) {
	zcall := reflect.ValueOf(express)
	kind := zcall.Kind()
	if kind != reflect.Func {
		panic("参数必须为函数")
	}

	switch s.kind {
	case reflect.Slice, reflect.String, reflect.Array:
	Loop:
		for i := 0; i < s.entries.Len(); i++ {
			v := s.entries.Index(i)
			for _, expr := range s.exprs {
				switch expr.kind {
				case mapExpr:
					call := reflect.ValueOf(expr.expr)
					v = call.Call([]reflect.Value{v})[0]
				case filterExpr:
					call := reflect.ValueOf(expr.expr)
					if !call.Call([]reflect.Value{v})[0].Bool() {
						continue Loop
					}
				}

			}
			zcall.Call([]reflect.Value{v})
		}
	case reflect.Map:
		iter := s.entries.MapRange()
	LoopM:
		for iter.Next() {
			k, v := iter.Key(), iter.Value()
			for _, expr := range s.exprs {
				switch expr.kind {
				case mapExpr:
					call := reflect.ValueOf(expr.expr)
					res := call.Call([]reflect.Value{k, v})
					k, v = res[0], res[1]
				case filterExpr:
					call := reflect.ValueOf(expr.expr)
					if !call.Call([]reflect.Value{k, v})[0].Bool() {
						continue LoopM
					}
				}

			}
			zcall.Call([]reflect.Value{k, v})
		}
	case reflect.Chan: // blocked
		item, ok := s.entries.Recv()
		if ok {
			zcall.Call([]reflect.Value{item})
		}
	}
}

func (s *stream) Map(express any) *stream {
	call := reflect.ValueOf(express)
	kind := call.Kind()
	if kind != reflect.Func {
		panic("参数必须为函数")
	}
	if call.Type().NumOut() != 1 {
		panic("函数必须有1个返回值")
	}
	s.exprs = append(s.exprs, expr{mapExpr, express})
	return s
}

func (s *stream) Filter(express any) *stream {
	call := reflect.ValueOf(express)
	kind := call.Kind()
	if kind != reflect.Func {
		panic("参数必须为函数")
	}
	typ := call.Type()
	if typ.NumOut() != 1 || typ.Out(0).Kind() != reflect.Bool {
		panic("函数返回值必须为bool")
	}
	s.exprs = append(s.exprs, expr{filterExpr, express})
	return s
}

func (s *stream) Every(express any) bool {
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

func (s *stream) Some(express any) bool {
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
			if out[0].Bool() {
				return true
			}
		}
	case reflect.Map:
		iter := s.entries.MapRange()
		for iter.Next() {
			out := call.Call([]reflect.Value{iter.Value()})
			if out[0].Bool() {
				return true
			}
		}
	case reflect.Chan: // blocked
		item, ok := s.entries.Recv()
		if ok {
			out := call.Call([]reflect.Value{item})
			if out[0].Bool() {
				return true
			}
		}
	}
	return false
}
