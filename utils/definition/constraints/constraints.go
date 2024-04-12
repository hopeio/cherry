package constraints

import (
	"golang.org/x/exp/constraints"
	"time"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Callback[T any] interface {
	~func() | ~func() error | ~func(T) | ~func(T) error
}

type ID interface {
	constraints.Integer | ~string | ~[]byte | ~[8]byte | ~[16]byte
}

type Range interface {
	constraints.Ordered | time.Time | ~*time.Time | ~string
}
