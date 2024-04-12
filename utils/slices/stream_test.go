package slices

import (
	"testing"
)

func TestStream(t *testing.T) {
	s := []int{1, 2, 3}
	stream := Stream(s)
	stream.ForEach(func(s int) {
		t.Log(s)
	})
	ret := stream.Map(func(s int) int {
		return s + 1
	})
	t.Log(ret)
	ret = stream.Filter(func(s int) bool {
		return s == 1
	})
	t.Log(ret)
}
