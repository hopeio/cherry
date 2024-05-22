package queue

import "github.com/hopeio/cherry/utils/cmp"

type TinyQueue[T cmp.CompareLess[T]] struct {
	length int
	data   []T
	zero   T
}

func New[T cmp.CompareLess[T]](data []T) *TinyQueue[T] {
	q := &TinyQueue[T]{}
	q.data = data
	q.length = len(data)
	if q.length > 0 {
		i := q.length >> 1
		for ; i >= 0; i-- {
			q.down(i)
		}
	}
	return q
}

func (q *TinyQueue[T]) Push(item T) {
	q.data = append(q.data, item)
	q.length++
	q.up(q.length - 1)
}
func (q *TinyQueue[T]) Pop() (T, bool) {
	if q.length == 0 {
		return q.zero, false
	}
	top := q.data[0]
	q.length--
	if q.length > 0 {
		q.data[0] = q.data[q.length]
		q.down(0)
	}
	q.data = q.data[:len(q.data)-1]
	return top, true
}
func (q *TinyQueue[T]) Peek() (T, bool) {
	if q.length == 0 {
		return q.zero, false
	}
	return q.data[0], true
}
func (q *TinyQueue[T]) Len() int {
	return q.length
}
func (q *TinyQueue[T]) down(pos int) {
	data := q.data
	halfLength := q.length >> 1
	item := data[pos]
	for pos < halfLength {
		left := (pos << 1) + 1
		right := left + 1
		best := data[left]
		if right < q.length && data[right].Less(best) {
			left = right
			best = data[right]
		}
		if !best.Less(item) {
			break
		}
		data[pos] = best
		pos = left
	}
	data[pos] = item
}

func (q *TinyQueue[T]) up(pos int) {
	data := q.data
	item := data[pos]
	for pos > 0 {
		parent := (pos - 1) >> 1
		current := data[parent]
		if !item.Less(current) {
			break
		}
		data[pos] = current
		pos = parent
	}
	data[pos] = item
}
