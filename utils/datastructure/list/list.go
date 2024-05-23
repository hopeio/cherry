package list

type Node[T any] struct {
	Value T
	Next  *Node[T]
}

type List[T any] struct {
	head, tail *Node[T]
	size       uint
	zero       T
}

func New[T any]() List[T] {
	l := List[T]{}
	l.head = nil //head指向头部结点
	l.tail = nil //tail指向尾部结点
	l.size = 0
	return l
}

func (l *List[T]) Len() uint {
	return l.size
}

func (l *List[T]) Head() *Node[T] {
	if l.size == 0 {
		return nil
	}
	return l.head
}

func (l *List[T]) Tail() *Node[T] {
	if l.size == 0 {
		return nil
	}
	return l.tail
}

func (l *List[T]) First() (T, bool) {
	if l.size == 0 {
		return l.zero, false
	}
	return l.head.Value, true
}

func (l *List[T]) Last() (T, bool) {
	if l.size == 0 {
		return l.zero, false
	}
	return l.tail.Value, true
}

func (l *List[T]) Pop() (T, bool) {
	if l.size == 0 {
		return l.zero, false
	}

	p := l.head
	l.head = p.Next
	if l.size == 1 {
		l.tail = nil
	}
	l.size--
	return p.Value, true
}

func (l *List[T]) Push(v T) {
	node := &Node[T]{v, nil}
	if l.size == 0 {
		l.head = node
		l.tail = node
		l.size++
		return
	}
	l.tail.Next = node
	l.tail = node
	l.size++
}
