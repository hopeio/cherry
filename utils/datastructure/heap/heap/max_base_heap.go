package heap

import (
	constraintsi "github.com/hopeio/cherry/utils/types"
	"golang.org/x/exp/constraints"
)

type MaxBaseHeap[T constraints.Ordered] []T

func NewMaxBaseHeap[T constraints.Ordered](l int) MaxBaseHeap[T] {
	maxHeap := make(MaxBaseHeap[T], 0, l)
	return maxHeap
}

func NewMaxBaseHeapFromArray[T constraints.Ordered](arr []T) MaxBaseHeap[T] {
	heap := NewBaseHeapFromArray[T](arr, constraintsi.GreaterFunc[T])
	return MaxBaseHeap[T](heap)
}

func (h MaxBaseHeap[T]) Init() {
	BaseHeap[T](h).init(constraintsi.GreaterFunc[T])
}

func (h *MaxBaseHeap[T]) Push(x T) {
	(*BaseHeap[T])(h).push(x, constraintsi.GreaterFunc[T])
}

func (h *MaxBaseHeap[T]) Pop() T {
	return (*BaseHeap[T])(h).pop(constraintsi.GreaterFunc[T])
}

func (h *MaxBaseHeap[T]) Remove(i int) T {
	return (*BaseHeap[T])(h).remove(i, constraintsi.GreaterFunc[T])
}
