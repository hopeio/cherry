package _interface

import "time"

type IdGenerator interface {
	Id() int64
}

type DurationGenerator interface {
	Duration() time.Duration
}
