package constraints

import "golang.org/x/exp/constraints"

type CompareFunc[T any] func(T, T) bool

type Compare[T any] interface {
	Compare(T) bool
}

// comparable 只能比较是否相等,不能比较大小
type OrderKey[T constraints.Ordered] interface {
	OrderKey() T
}

type CmpKey[T comparable] interface {
	CmpKey() T
}
