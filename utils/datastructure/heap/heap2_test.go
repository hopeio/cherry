package heap

import (
	"fmt"
	"math/rand"
	"testing"
)

type Foo2 struct {
	P int
}

func (receiver Foo2) Less(v Foo2) bool {
	return receiver.P < v.P
}

func TestNewHeap2(t *testing.T) {
	var arr []Foo2
	for i := 0; i < 10; i++ {
		arr = append(arr, Foo2{P: rand.Intn(10000)})
	}
	heap := NewHeap2(arr)
	fmt.Println(heap)

}
