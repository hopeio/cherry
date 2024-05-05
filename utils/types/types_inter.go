package types

import "github.com/hopeio/cherry/utils/constraints"

type Key[T constraints.Key] interface {
	Key() T
}
