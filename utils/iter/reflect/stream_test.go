package iter

import (
	"testing"
)

func TestStream(t *testing.T) {
	s := []int{1, 2, 3}
	StreamOf(s).Map(func(s int) int {
		return s + 1
	}).Filter(func(s int) bool {
		return s == 2
	}).ForEach(func(s int) {
		t.Log(s)
	})
}
