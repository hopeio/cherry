package cmp

import (
	constraintsi "github.com/hopeio/cherry/utils/constraints"
	"golang.org/x/exp/constraints"
)

func Less[T constraints.Ordered](a T, b T) bool {
	return a < b
}

func LessByKey[K constraints.Ordered, T SortKey[K]](a T, b T) bool {
	return a.SortKey() < b.SortKey()
}

func Greater[T constraints.Ordered](a T, b T) bool {
	return a > b
}

func GreaterByKey[K constraints.Ordered, T SortKey[K]](a T, b T) bool {
	return a.SortKey() > b.SortKey()
}

func Equal[T comparable](a T, b T) bool {
	return a == b
}

func EqualByKey[K constraints.Ordered, T SortKey[K]](a T, b T) bool {
	return a.SortKey() == b.SortKey()
}

func CompareNumber[T constraintsi.Number](a T, b T) int {
	return int(a - b)
}

func CompareByKey[K constraints.Ordered, T SortKey[K]](a T, b T) int {
	return Compare(a.SortKey(), b.SortKey())
}

func Compare[T constraints.Ordered](x, y T) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}

type GTValue[T constraints.Ordered] struct {
	Value T
}

func (a GTValue[T]) Compare(b GTValue[T]) bool {
	return a.Value > b.Value
}

type LTValue[T constraints.Ordered] struct {
	Value T
}

func (a LTValue[T]) Compare(b GTValue[T]) bool {
	return a.Value < b.Value
}
