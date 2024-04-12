//go:build !gccgo

package reflecti

import (
	"fmt"
	"github.com/hopeio/cherry/utils/definition/types/foo"
	"github.com/modern-go/reflect2"
	"testing"

	"github.com/hopeio/cherry/utils/log"
)

type Foo struct {
	A int
	B string
}
type Bar struct {
	Foo Foo
	C   string
}

func TestGetExpectTypeValue(t *testing.T) {
	a := Bar{Foo: Foo{A: 1}}
	b := Foo{}
	v := CopyFieldValueByType(&a, &b)
	if v {
		log.Info(b)
	}
}

func TestReflect2(t *testing.T) {
	a := foo.Foo{}
	typ := reflect2.TypeByName("foo.foo")
	structType := typ.(*reflect2.UnsafeStructType)
	new2 := structType.New()
	structType.Field(0).Set(new2, pInt(1))

	fmt.Println(a, new2)
}

var pInt = func(val int) *int {
	return &val
}
