package slices

import (
	"github.com/hopeio/cherry/utils/cmp"
	"github.com/stretchr/testify/assert"
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

var _ cmp.CompareKey[uint64] = &Foo{}

func TestHasCoincide(t *testing.T) {
	s1 := []cmp.CompareKey[uint64]{
		&Foo{1, "1"},
		&Foo{2, "2"},
		&Foo{5, "3"},
	}
	s2 := []cmp.CompareKey[uint64]{
		&Foo{4, "1"},
		&Foo{5, "1"},
		&Foo{6, "1"},
	}
	assert.Equal(t, false, HasCoincide(s1, s2))
	assert.Equal(t, true, HasCoincideByKey(s1, s2))
}

func TestDifference(t *testing.T) {
	a := []uint64{1, 2, 3, 4}
	b := []uint64{2, 3, 4, 5}
	diff1, diff2 := Difference(a, b)
	assert.ElementsMatch(t, []uint64{1}, diff1)
	assert.ElementsMatch(t, []uint64{5}, diff2)
	assert.ElementsMatch(t, []uint64{2, 3, 4}, intersection(a, b))
	u, i, d1, d2 := UnionAndIntersectionAndDifference(a, b)
	assert.ElementsMatch(t, []uint64{1, 2, 3, 4, 5}, u)
	assert.ElementsMatch(t, []uint64{2, 3, 4}, i)
	assert.ElementsMatch(t, []uint64{1}, d1)
	assert.ElementsMatch(t, []uint64{5}, d2)
}
