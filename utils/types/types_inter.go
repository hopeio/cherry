package types

import "github.com/hopeio/cherry/utils/constraints"

type IKey[T constraints.Key] interface {
	Key() T
}
