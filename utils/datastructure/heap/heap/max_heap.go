package heap

import (
	constraints2 "github.com/hopeio/cherry/utils/constraints"
	"github.com/hopeio/cherry/utils/types"
	"golang.org/x/exp/constraints"
)

type MaxHeap[T constraints2.OrderKey[V], V constraints.Ordered] Heap[T, V]

func NewMaxHeap[T constraints2.OrderKey[V], V constraints.Ordered](l int) MaxHeap[T, V] {
	maxHeap := make([]T, 0, l)
	return MaxHeap[T, V]{
		entry: maxHeap,
		less:  types.Greater[V],
	}
}

func NewMaxHeapFromArray[T constraints2.OrderKey[V], V constraints.Ordered](arr []T) MaxHeap[T, V] {
	heap := NewHeapFromArray[T, V](arr, types.Greater[V])
	return MaxHeap[T, V](heap)
}

func (h *MaxHeap[T, V]) Init() {
	(*Heap[T, V])(h).Init()
}

func (h *MaxHeap[T, V]) Push(x T) {
	(*Heap[T, V])(h).Push(x)
}

func (h *MaxHeap[T, V]) Pop() T {
	return (*Heap[T, V])(h).Pop()
}

func (h *MaxHeap[T, V]) Remove(i int) T {

	return (*Heap[T, V])(h).Remove(i)
}
