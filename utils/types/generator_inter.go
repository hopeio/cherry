package types

import (
	"github.com/hopeio/cherry/utils/constraints"
	"time"
)

type IIdGenerator[T constraints.ID] interface {
	Id() T
}

type IDurationGenerator interface {
	Duration() time.Duration
}
