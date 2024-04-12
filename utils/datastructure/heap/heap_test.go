package heap

import (
	"fmt"
	"math/rand"
	"testing"
)

type heap []int

func (h heap) Less(i, j int) bool {
	return h[i] < h[j]
}

func TestNewHeap(t *testing.T) {
	var arr heap
	for i := 0; i < 10; i++ {
		arr = append(arr, rand.Intn(10000))
	}
	heap := NewHeap(arr)
	fmt.Println(heap)
	for i := 0; i < 10; i++ {
		heap.Push(rand.Intn(10000))
	}
	fmt.Println(heap)
	n := len(heap)
	for i := 0; i < n; i++ {
		t.Log(heap.Pop())
	}
}

type Foo struct {
	A int
	B int
}

type Foos []Foo

func (f Foos) Less(i, j int) bool {
	if f[i].A == f[j].A {
		return f[i].B < f[j].B
	}
	return f[i].A < f[j].A
}

func TestPushHeap(t *testing.T) {
	var arr Foos
	heap := NewHeap(arr)
	heap.Push(Foo{10, 10})
	heap.Push(Foo{5, 5})
	heap.Push(Foo{8, 8})
	heap.Push(Foo{2, 2})
	heap.Push(Foo{5, 51})
	heap.Push(Foo{26, 26})
	heap.Push(Foo{6, 6})
	heap.Push(Foo{9, 9})
	heap.Push(Foo{5, 52})
	heap.Push(Foo{1, 1})
	for _, foo := range heap {
		t.Log(foo)
	}
	t.Log("------------------")
	n := len(heap)
	for i := 0; i < n; i++ {
		t.Log(heap.Pop())
	}
}
