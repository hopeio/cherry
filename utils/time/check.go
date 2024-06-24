package time

import (
	"errors"
	"fmt"
)

type TimeCheck interface {
	GetSeconds() int64
	GetNanos() int32
}

func IsValid(x TimeCheck) bool {
	return x != nil && Check(x) == 0
}

func CheckValid(x TimeCheck) error {
	switch Check(x) {
	case InvalidNil:
		return errors.New("invalid nil Timestamp")
	case InvalidUnderflow:
		return fmt.Errorf("time (%v) before 0001-01-01", x)
	case InvalidOverflow:
		return fmt.Errorf("time (%v) after 9999-12-31", x)
	case InvalidNanos:
		return fmt.Errorf("time (%v) has out-of-range nanos", x)
	default:
		return nil
	}
}

const (
	_ = iota
	InvalidNil
	InvalidUnderflow
	InvalidOverflow
	InvalidNanos
)

func Check(x TimeCheck) uint {
	const minTimestamp = -62135596800  // Seconds between 1970-01-01T00:00:00Z and 0001-01-01T00:00:00Z, inclusive
	const maxTimestamp = +253402300799 // Seconds between 1970-01-01T00:00:00Z and 9999-12-31T23:59:59Z, inclusive

	secs := x.GetSeconds()
	nanos := x.GetNanos()
	switch {
	case x == nil:
		return InvalidNil
	case secs < minTimestamp:
		return InvalidUnderflow
	case secs > maxTimestamp:
		return InvalidOverflow
	case nanos < 0 || nanos >= 1e9:
		return InvalidNanos
	default:
		return 0
	}
}
