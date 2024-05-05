package cmp

import (
	"golang.org/x/exp/constraints"
)

func Less[T constraints.Ordered](a T, b T) bool {
	return a < b
}

func Greater[T constraints.Ordered](a T, b T) bool {
	return a > b
}

func Equal[T comparable](a T, b T) bool {
	return a == b
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
