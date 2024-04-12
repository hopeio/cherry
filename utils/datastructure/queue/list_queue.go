package queue

type Node[T any] struct {
	data     interface{}
	previous *Node[T]
	next     *Node[T]
}

func (n *Node[T]) Value() interface{} {
	return n.data
}

func (n *Node[T]) Set(value interface{}) {
	n.data = value
}

func (n *Node[T]) Previous() *Node[T] {
	return n.previous
}

func (n *Node[T]) Next() *Node[T] {
	return n.next
}

type ListQueue[T any] struct {
	head *Node[T]
	end  *Node[T]
	size int
}

func NewListQueue[T any](size int) *ListQueue[T] {
	q := &ListQueue[T]{nil, nil, size}
	return q
}

func (q *ListQueue[T]) push(data T) {
	n := &Node[T]{data: data, next: nil}

	if q.end == nil {
		q.head = n
		q.end = n
	} else {
		q.end.next = n
		q.end = n
	}

	return
}

func (q *ListQueue[T]) pop() (T, bool) {
	if q.head == nil {
		return *new(T), false
	}

	data := q.head.data
	q.head = q.head.next
	if q.head == nil {
		q.end = nil
	}

	return data, true
}
