package cmp

import (
	constraintsi "github.com/hopeio/cherry/utils/constraints"
	"golang.org/x/exp/constraints"
)

type EqualKey[T comparable] interface {
	EqualKey() T
}

type IsEqual[T any] interface {
	Equal(T) bool
}

type CompareKey[T constraintsi.Number] interface {
	CompareKey() T
}

type compareKey[K any, T Comparable[K]] interface {
	CompareKey() T
}

type Comparable[T any] interface {
	Compare(T) int
}

// comparable 只能比较是否相等,不能比较大小
type SortKey[T constraints.Ordered] interface {
	SortKey() T
}

type Sort[T any] interface {
	Less(T) bool
}
