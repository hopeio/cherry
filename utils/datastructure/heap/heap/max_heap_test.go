package heap

import "testing"

type Foo struct {
	A int
	B string
}

func (f *Foo) OrderKey() int {
	return f.A
}

func TestHeap(t *testing.T) {
	heap := MaxHeap[*Foo, int]{}
	heap.Init()
	heap.Push(&Foo{10, "10"})
	heap.Push(&Foo{5, "5"})
	heap.Push(&Foo{8, "8"})
	heap.Push(&Foo{2, "2"})
	heap.Push(&Foo{5, "51"})
	heap.Push(&Foo{26, "26"})
	heap.Push(&Foo{6, "6"})
	heap.Push(&Foo{9, "9"})
	heap.Push(&Foo{5, "52"})
	heap.Push(&Foo{1, "1"})
	for _, foo := range heap {
		t.Log(foo)
	}
	t.Log("------------------")
	n := len(heap)
	for i := 0; i < n; i++ {
		t.Log(heap.Pop())
	}
}
