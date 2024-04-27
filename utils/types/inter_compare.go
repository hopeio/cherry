package types

import constraintsi "github.com/hopeio/cherry/utils/constraints"

type ICompare[T any] interface {
	Compare(T) bool
}

type IKey[T constraintsi.Key] interface {
	Key() T
}
