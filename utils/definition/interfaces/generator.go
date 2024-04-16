package interfaces

import (
	"github.com/hopeio/cherry/utils/definition/constraints"
	"time"
)

type IdGenerator[T constraints.ID] interface {
	Id() T
}

type DurationGenerator interface {
	Duration() time.Duration
}
