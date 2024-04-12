package heap

type Interface2[T any] interface {
	Less(v T) bool
}

type Heap2[T Interface2[T]] []T

func NewHeap2[T Interface2[T]](arr []T) Heap2[T] {
	heap := Heap2[T](arr)
	for i := 1; i < len(arr); i++ {
		heap.up(i)
	}
	return heap
}

func (heap Heap2[T]) up(j int) {

	for {
		i := (j - 1) / 2 // parent
		if i == j || !heap[j].Less(heap[i]) {
			break
		}
		heap[i], heap[j] = heap[j], heap[i]
		j = i
	}
}
