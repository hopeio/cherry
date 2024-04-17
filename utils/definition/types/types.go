package types

import (
	constraintsi "github.com/hopeio/cherry/utils/definition/constraints"
	"golang.org/x/exp/constraints"
)

type Enum[T constraints.Integer] int

type String string

func (s String) Key() string {
	return string(s)
}

type Int int

func (s Int) Key() int {
	return int(s)
}

type Basic struct {
}

func ZeroValue[T any]() T {
	var zero T
	return zero
}

type ID[T constraintsi.ID] struct {
	Id T `json:"id"`
}
