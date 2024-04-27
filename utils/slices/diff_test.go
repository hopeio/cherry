package slices

import (
	"github.com/hopeio/cherry/utils/constraints"
	"log"
	"testing"
)

type Foo struct {
	ID  uint64
	Str string
}

func (f *Foo) IsEqual(v interface{}) bool {
	if f1, ok := v.(*Foo); ok {
		if f1.ID == f.ID {
			return true
		}
	}
	return false
}
func (f *Foo) CompareKey() uint64 {
	return f.ID
}

var _ constraints.CompareKey[uint64] = &Foo{}

func TestIsEqu(t *testing.T) {
	s1 := []constraints.CompareKey[uint64]{
		&Foo{1, "1"},
		&Foo{2, "2"},
		&Foo{5, "3"},
	}
	s2 := []constraints.CompareKey[uint64]{
		&Foo{4, "1"},
		&Foo{5, "1"},
		&Foo{6, "1"},
	}
	log.Println(HasCoincide(s1, s2))
	log.Println(HasCoincideByKey(s1, s2))
}

func TestDiff(t *testing.T) {
	a := []uint64{1, 2, 3, 4}
	b := []uint64{2, 3, 4, 5}
	t.Log(Difference(a, b))
	t.Log(intersection(a, b))
	t.Log(IntersectionAndDifference(a, b))
}
