package timei

import (
	"context"
	"time"
)

// Duration be used toml unmarshal string time, like 1s, 500ms.
type Duration time.Duration

// UnmarshalText unmarshal text to duration.
func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := time.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}

// Shrink will decrease the duration by comparing with context's timeout duration
// and return new timeout\context\CancelFunc.
func (d Duration) Shrink(c context.Context) (Duration, context.Context, context.CancelFunc) {
	if deadline, ok := c.Deadline(); ok {
		if ctimeout := time.Until(deadline); ctimeout < time.Duration(d) {
			// deliver small timeout
			return Duration(ctimeout), c, func() {}
		}
	}
	ctx, cancel := context.WithTimeout(c, time.Duration(d))
	return d, ctx, cancel
}

const (
	Day        = time.Hour * 24
	MonthDay30 = Day * 30
	MonthDay31 = Day * 31
	MonthDay28 = Day * 28
	MonthDay29 = Day * 29
	Month      = MonthDay30
	YearDay365 = Day * 365
	YearDay366 = Day * 366
	Year       = YearDay365
)
