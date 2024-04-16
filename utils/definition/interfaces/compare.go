package interfaces

import constraintsi "github.com/hopeio/cherry/utils/definition/constraints"

type Compare[T any] interface {
	Compare(T) bool
}

type Key[T constraintsi.Key] interface {
	Key() T
}
