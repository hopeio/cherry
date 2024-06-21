package reflect

import "reflect"

type ValueHandler func(reflect.Value)
type ValueRecursionHandler func(reflect.Value, ValueHandler)

type KindHandler [reflect.UnsafePointer + 1][]ValueRecursionHandler

func (k *KindHandler) FillUint(v ValueRecursionHandler) {
	k[reflect.Uint] = append(k[reflect.Uint], v)
	k[reflect.Uint64] = append(k[reflect.Uint64], v)
	k[reflect.Uint32] = append(k[reflect.Uint32], v)
	k[reflect.Uint16] = append(k[reflect.Uint16], v)
	k[reflect.Uint8] = append(k[reflect.Uint8], v)
}

func (k *KindHandler) FillInt(v ValueRecursionHandler) {
	k[reflect.Int] = append(k[reflect.Int], v)
	k[reflect.Int64] = append(k[reflect.Int64], v)
	k[reflect.Int32] = append(k[reflect.Int32], v)
	k[reflect.Int16] = append(k[reflect.Int16], v)
	k[reflect.Int8] = append(k[reflect.Int8], v)
}

func (k *KindHandler) FillFloat(v ValueRecursionHandler) {
	k[reflect.Float64] = append(k[reflect.Float64], v)
	k[reflect.Float32] = append(k[reflect.Float32], v)
}

func (k *KindHandler) AddHandler(kind reflect.Kind, v ValueRecursionHandler) {
	k[kind] = append(k[kind], v)
}

func (k *KindHandler) Handle(value reflect.Value) {
	for _, f := range k[value.Kind()] {
		f(value, k.Handle)
	}
	for _, f := range k[0] {
		f(value, k.Handle)
	}
}
