package reflect

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"reflect"
	"testing"
)

type Foo1 struct {
	A    int
	B    int
	Foo2 *Foo2
}

type Foo2 struct {
	A    int
	B    int
	Foo3 *Foo3
}

type Foo3 struct {
	A int
	B int
}

func TestStream(t *testing.T) {
	var handlers KindHandler
	handlers.FillInt(func(value reflect.Value, handler ValueHandler) {
		value.SetInt(rand.Int64N(128))
	})
	handlers.FillUint(func(value reflect.Value, handler ValueHandler) {
		value.SetUint(rand.Uint64N(256))
	})
	handlers.AddHandler(reflect.Pointer,
		func(value reflect.Value, handler ValueHandler) {
			if value.IsNil() && value.CanSet() {
				value.Set(reflect.New(value.Type().Elem()))
			}
			handler(value.Elem())
		})

	handlers.AddHandler(reflect.Struct,
		func(value reflect.Value, handler ValueHandler) {
			for i := 0; i < value.NumField(); i++ {
				field := value.Field(i)
				handler(field)
			}
		})

	var foo Foo1
	handlers.Handle(reflect.ValueOf(&foo))
	data, err := json.Marshal(&foo)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(data))
}
