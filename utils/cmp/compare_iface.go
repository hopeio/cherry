package cmp

import "golang.org/x/exp/constraints"

type CompareKey[T comparable] interface {
	CompareKey() T
}

type Compare[T any] interface {
	Compare(T) int
}

// comparable 只能比较是否相等,不能比较大小
type SortKey[T constraints.Ordered] interface {
	SortKey() T
}

type Sort[T any] interface {
	Less(T) bool
}
