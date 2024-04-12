package types

import (
	constraintsi "github.com/hopeio/cherry/utils/definition/constraints"
	"golang.org/x/exp/constraints"
)

type CompareFunc func(a, b any) bool

func SignedConvert[T, V constraints.Signed](v V) T {
	return T(v)
}

func FloatConvert[T, V constraints.Float](v V) T {
	return T(v)
}

func UnsignedConvert[T, V constraints.Unsigned](v V) T {
	return T(v)
}

func IntegerConvert[T, V constraints.Integer](v V) T {
	return T(v)
}

func NumberConvert[T, V constraintsi.Number](v V) T {
	return T(v)
}

func LessFunc[T constraints.Ordered](a T, b T) bool {
	return a < b
}

func GreaterFunc[T constraints.Ordered](a T, b T) bool {
	return a > b
}

func EqualFunc[T constraints.Ordered](a T, b T) bool {
	return a > b
}
