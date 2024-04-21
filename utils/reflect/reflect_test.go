package reflect

import (
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
