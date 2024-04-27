package types

type ICompare[T any] interface {
	Compare(T) bool
}

type ICompareKey[T comparable] interface {
	CompareKey() T
}
