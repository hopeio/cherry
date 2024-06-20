package interfaces

import "github.com/hopeio/cherry/utils/types/constraints"

type Key[T constraints.Key] interface {
	Key() T
}
