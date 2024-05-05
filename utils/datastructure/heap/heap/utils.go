package heap

import (
	"github.com/hopeio/cherry/utils/cmp"
	"golang.org/x/exp/constraints"
)

func HeapInit[T cmp.OrderKey[V], V constraints.Ordered](heap []T, less cmp.SortFunc[V]) {
	// heapify
	n := len(heap)
	for i := n/2 - 1; i >= 0; i-- {
		Down(heap, i, n, less)
	}
}

func Down[T cmp.OrderKey[V], V constraints.Ordered](heap []T, i0, n int, less cmp.SortFunc[V]) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && less(heap[j2].OrderKey(), heap[j1].OrderKey()) {
			j = j2 // = 2*i + 2  // right child
		}
		if !less(heap[j].OrderKey(), heap[i].OrderKey()) {
			break
		}
		heap[i], heap[j] = heap[j], heap[i]
		i = j
	}
	return i > i0
}

func Up[T cmp.OrderKey[V], V constraints.Ordered](heap []T, j int, less cmp.SortFunc[V]) {

	for {
		i := (j - 1) / 2 // parent
		if i == j || !less(heap[j].OrderKey(), heap[i].OrderKey()) {
			break
		}
		heap[i], heap[j] = heap[j], heap[i]
		j = i
	}
}

func Fix[T cmp.OrderKey[V], V constraints.Ordered](heap []T, i int, less cmp.SortFunc[V]) {
	if !Down(heap, i, len(heap), less) {
		Up(heap, i, less)
	}
}

func HeapInitForBase[T any](heap []T, less cmp.SortFunc[T]) {
	// heapify
	n := len(heap)
	for i := n/2 - 1; i >= 0; i-- {
		DownForBase(heap, i, n, less)
	}
}

// 标准库写法
func DownForBase[T any](heap []T, i0, n int, less cmp.SortFunc[T]) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && less(heap[j2], heap[j1]) {
			j = j2 // = 2*i + 2  // right child
		}
		if !less(heap[j], heap[i]) {
			break
		}
		heap[i], heap[j] = heap[j], heap[i]
		i = j
	}
	return i > i0
}

// 与Down一致，不同的写法
func AdjustDownForBase[T any](heap []T, i int, less cmp.SortFunc[T]) {
	child := leftChild(i)
	for child < len(heap) {
		if child+1 < len(heap) && less(heap[child+1], heap[child]) {
			child++
		}
		if !less(heap[child], heap[i]) {
			break
		}
		heap[i], heap[child] = heap[child], heap[i]
		i = child
		child = leftChild(i)
	}
}

func UpForBase[T any](heap []T, j int, less cmp.SortFunc[T]) {

	for {
		i := (j - 1) / 2 // parent
		if i == j || !less(heap[j], heap[i]) {
			break
		}
		heap[i], heap[j] = heap[j], heap[i]
		j = i
	}
}

// 与Up一致，不同的写法
func AdjustUpForBase[T any](heap []T, i int, less cmp.SortFunc[T]) {
	p := parent(i)
	for p >= 0 && less(heap[i], heap[p]) {
		heap[i], heap[p] = heap[p], heap[i]
		i = p
		p = parent(i)
	}
}

func FixForBase[T any](heap []T, i int, less cmp.SortFunc[T]) {
	if !DownForBase(heap, i, len(heap), less) {
		UpForBase(heap, i, less)
	}
}

func parent(i int) int {
	return (i - 1) / 2
}
func leftChild(i int) int {
	return i*2 + 1
}
