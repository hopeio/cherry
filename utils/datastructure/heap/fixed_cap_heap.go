package heap

type FixedCapHeap[T any, V Interface[T]] Heap[T, V]

func NewFixedCapHeap[T any, V Interface[T]](arr V) FixedCapHeap[T, V] {
	heap := FixedCapHeap[T, V](NewHeap(arr))
	return heap
}

func (heap FixedCapHeap[T, V]) Put(val T) {
	if len(heap) < cap(heap) {
		heap = append(heap, val)
		for i := 1; i < len(heap); i++ {
			Heap[T, V](heap).up(i)
		}
		return
	}
	temp := V{heap[0], val}
	if temp.Less(1, 0) {
		return
	}
	heap[0] = val
	Heap[T, V](heap).down(0, len(heap))
}
