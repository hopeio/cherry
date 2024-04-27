package iter

type mapIter[T any, R any] struct {
	transform func(T) R
	iterator  Iterator[T]
}

func (a *mapIter[T, R]) Next() (R, bool) {
	var null R
	if v, ok := a.iterator.Next(); ok {
		return a.transform(v), true
	}
	return null, false
}

func MapOf[T any, R any](transform func(T) R, iter Iterator[T]) Iterator[R] {
	return &mapIter[T, R]{transform, iter}
}
