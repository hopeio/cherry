package iter

import (
	"github.com/stretchr/testify/assert"
	"iter"
	"testing"
)

func TestIter(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	seq := SliceSeqOf(s)
	t.Log(First[int](iter.Seq[int](seq)))
	for v := range seq {
		t.Log(v)
	}
	assert.Equal(t, true, SliceSeqOf(s).IsSorted(func(i int, i2 int) bool {
		return i < i2
	}))
	SliceSeqOf(s).Filter(func(i int) bool {
		return i%2 == 0
	}).Map(func(i int) int {
		return i + 10
	}).ForEach(func(i int) {
		t.Log(i)
	})
}
