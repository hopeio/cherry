package constraints

import "golang.org/x/exp/constraints"

// comparable 只能比较是否相等,不能比较大小
type OrderKey[T constraints.Ordered] interface {
	OrderKey() T
}

type CompareKey[T comparable] interface {
	CompareKey() T
}
