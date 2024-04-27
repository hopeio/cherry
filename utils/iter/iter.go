package iter

type Iterator[T any] interface {
	Next() (v T, ok bool)
}

type Iterable[T any] interface {
	Iter() Iterator[T]
}
