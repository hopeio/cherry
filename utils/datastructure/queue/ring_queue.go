package queue

type Queue[T any] interface {
	// 获取当前链表长度。
	Len() int
	// 获取当前链表容量。
	Capacity() int
	// 获取当前链表头结点。
	Front() *Node[T]
	// 获取当前链表尾结点。
	Rear() *Node[T]
	// 入列。
	Enqueue(value T) bool
	// 出列。
	Dequeue() T
}

type RingQueue[T any] struct {
	head, tail int
	len        int
	buf        []T
}

func NewRingQueue[T any](capacity int) *RingQueue[T] {
	nodes := make([]T, capacity)
	return &RingQueue[T]{
		head: -1,
		tail: -1,
		buf:  nodes,
	}
}

func (q *RingQueue[T]) Length() int {
	return q.len
}

func (q *RingQueue[T]) Capacity() int {
	return len(q.buf)
}

func (q *RingQueue[T]) Front() T {
	if q.len == 0 {
		return *new(T)
	}

	return q.buf[q.head]
}

func (q *RingQueue[T]) Tail() T {
	if q.len == 0 {
		return *new(T)
	}

	return q.buf[q.tail]
}

func (q *RingQueue[T]) Enqueue(value T) bool {
	if q.IsFull() || value == nil {
		return false
	}

	q.tail++
	if q.tail == len(q.buf) {
		q.tail = 0
	}
	q.buf[q.tail] = value
	q.len++

	if q.len == 1 {
		q.head = q.tail
	}

	return true
}

func (q *RingQueue[T]) Dequeue() T {
	if q.len == 0 {
		return *new(T)
	}

	result := q.buf[q.head]
	q.buf[q.head] = *new(T)
	q.head++
	q.len--
	if q.head == len(q.buf) {
		q.head = 0
	}

	return result
}

// IsFull checks if the ring buffer is full
func (q *RingQueue[T]) IsFull() bool {
	return q.len == len(q.buf)
}

// LookAll reads all elements from ring buffer
// this method doesn't consume all elements
func (q *RingQueue[T]) LookAll() []interface{} {
	all := make([]interface{}, q.len)
	if q.len == 0 {
		return all
	}
	j := 0
	for i := q.head; ; i++ {
		if i == len(q.buf) {
			i = 0
		}
		all[j] = q.buf[i]
		if i == q.tail {
			break
		}
		j++
	}
	return all
}
