package heap

type Interface[T any] interface {
	~[]T
	Less(i, j int) bool
}

type Heap[T any, V Interface[T]] []T

func NewHeap[T any, V Interface[T]](arr V) Heap[T, V] {
	heap := Heap[T, V](arr)
	for i := 1; i < len(arr); i++ {
		heap.up(i)
	}
	return heap
}

func (heap Heap[T, V]) Init() {
	// heapify
	n := len(heap)
	for i := n/2 - 1; i >= 0; i-- {
		heap.down(i, n)
	}
}

func (heap *Heap[T, V]) Push(x T) {
	h := *heap
	*heap = append(h, x)
	heap.up(len(h))
}

func (heap *Heap[T, V]) Pop() T {
	h := *heap
	n := len(h) - 1
	item := h[0]
	h[0], h[n] = h[n], h[0]
	heap.down(0, n)
	*heap = h[:n]
	return item
}

func (heap Heap[T, V]) First() T {
	return heap[0]
}

func (heap Heap[T, V]) Last() T {
	return heap[len(heap)-1]
}

func (heap *Heap[T, V]) Remove(i int) T {
	h := *heap
	n := len(h) - 1
	item := h[i]
	if n != i {
		h[i], h[n] = h[n], h[i]
		if !heap.down(i, n) {
			heap.up(i)
		}
	}
	*heap = h[:n]
	return item
}

func (heap Heap[T, V]) up(j int) {

	for {
		i := (j - 1) / 2 // parent
		if i == j || !V(heap).Less(j, i) {
			break
		}
		heap[i], heap[j] = heap[j], heap[i]
		j = i
	}
}

func (heap Heap[T, V]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && V(heap).Less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !V(heap).Less(j, i) {
			break
		}
		heap[i], heap[j] = heap[j], heap[i]
		i = j
	}
	return i > i0
}
