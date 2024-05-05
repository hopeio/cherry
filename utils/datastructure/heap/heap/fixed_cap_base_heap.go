package heap

import (
	"github.com/hopeio/cherry/utils/cmp"
	"golang.org/x/exp/constraints"
)

type FixedCapBaseHeap[T constraints.Ordered] BaseHeap[T]

func NewFixedCapBaseHeap[T constraints.Ordered](l int) FixedCapBaseHeap[T] {
	heap := make([]T, 0, l)
	return heap
}

func NewFixedCapBaseHeapFromArray[T constraints.Ordered](arr []T, less cmp.SortFunc[T]) FixedCapBaseHeap[T] {
	return FixedCapBaseHeap[T](NewBaseHeapFromArray(arr, less))
}

func (h FixedCapBaseHeap[T]) init(less cmp.SortFunc[T]) {
	// heapify
	n := len(h)
	for i := n/2 - 1; i >= 0; i-- {
		BaseHeap[T](h).down(i, n, less)
	}
}

func (h FixedCapBaseHeap[T]) put(val T, less cmp.SortFunc[T]) {
	if len(h) < cap(h) {
		h = append(h, val)
		for i := 1; i < len(h); i++ {
			BaseHeap[T](h).up(i, less)
		}
		return
	}
	if less(val, h[0]) {
		return
	}
	h[0] = val
	BaseHeap[T](h).down(0, len(h), less)
}
