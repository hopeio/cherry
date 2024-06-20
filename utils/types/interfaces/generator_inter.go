package interfaces

import (
	"github.com/hopeio/cherry/utils/types/constraints"
	"time"
)

type IdGenerator[T constraints.ID] interface {
	Id() T
}

type DurationGenerator interface {
	Duration() time.Duration
}
