package types

import "golang.org/x/exp/constraints"

type FloatGTValue float64

func (a FloatGTValue) Compare(b FloatGTValue) bool {
	return a > b
}

type FloatLTValue float64

func (a FloatLTValue) Compare(b FloatLTValue) bool {
	return a < b
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
