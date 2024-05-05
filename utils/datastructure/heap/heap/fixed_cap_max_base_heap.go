package heap

import (
	"github.com/hopeio/cherry/utils/cmp"
	"golang.org/x/exp/constraints"
)

// 最大（小）堆是指在树中，存在一个结点而且该结点有儿子结点，该结点的data域值都不小于（大于）其儿子结点的data域值

// 定容最大堆 可用于保留前n个最小元素
type FixedCapMaxBaseHeap[T constraints.Ordered] []T

func NewFixedCapMaxBaseHeap[T constraints.Ordered](cap int) FixedCapMaxBaseHeap[T] {
	maxHeap := make(FixedCapMaxBaseHeap[T], 0, cap)
	return maxHeap
}

func NewFixedCapMaxBaseHeapFromArray[T constraints.Ordered](arr []T) FixedCapMaxBaseHeap[T] {
	return FixedCapMaxBaseHeap[T](NewBaseHeapFromArray[T](arr, cmp.Greater[T]))
}

func (heap FixedCapMaxBaseHeap[T]) Put(val T) {
	FixedCapBaseHeap[T](heap).put(val, cmp.Greater[T])
}
