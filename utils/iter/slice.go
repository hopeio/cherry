package iter

type Slice[T any] []T

func (a Slice[T]) Iterator() Iterator[T] {
	return &sliceIterator[T]{-1, a}
}

func (a Slice[T]) Count() int {
	return len(a)
}

type sliceIterator[T any] struct {
	index  int
	source []T
}

func (a *sliceIterator[T]) Next() (T, bool) {
	if a.index < len(a.source)-1 {
		a.index++
		return a.source[a.index], true
	}
	return *new(T), false
}
