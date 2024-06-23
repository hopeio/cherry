package timestamp

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	timei "github.com/hopeio/cherry/utils/time"
	"google.golang.org/protobuf/runtime/protoimpl"
	"io"
	"time"
)

// Now constructs a new Timestamp from the current time.
func Now() *Timestamp {
	return New(time.Now())
}

// New constructs a new Timestamp from the provided time.Time.
func New(t time.Time) *Timestamp {
	return &Timestamp{Seconds: t.Unix(), Nanos: int32(t.Nanosecond())}
}

// AsTime converts x to a time.Time.
func (x *Timestamp) AsTime() time.Time {
	return time.Unix(x.GetSeconds(), int64(x.GetNanos()))
}

// IsValid reports whether the timestamp is valid.
// It is equivalent to CheckValid == nil.
func (x *Timestamp) IsValid() bool {
	return x != nil && timei.Check(x) == 0
}

// CheckValid returns an error if the timestamp is invalid.
// In particular, it checks whether the value represents a date that is
// in the range of 0001-01-01T00:00:00Z to 9999-12-31T23:59:59Z inclusive.
// An error is reported for a nil Timestamp.
func (x *Timestamp) CheckValid() error {
	switch timei.Check(x) {
	case timei.InvalidNil:
		return protoimpl.X.NewError("invalid nil Timestamp")
	case timei.InvalidUnderflow:
		return protoimpl.X.NewError("timestamp (%v) before 0001-01-01", x)
	case timei.InvalidOverflow:
		return protoimpl.X.NewError("timestamp (%v) after 9999-12-31", x)
	case timei.InvalidNanos:
		return protoimpl.X.NewError("timestamp (%v) has out-of-range nanos", x)
	default:
		return nil
	}
}

// Scan scan time.
func (t *Timestamp) Scan(value interface{}) error {
	nullTime := &sql.NullTime{}
	err := nullTime.Scan(value)
	if err != nil {
		return err
	}
	if nullTime.Valid {
		*t = Timestamp{Seconds: nullTime.Time.Unix(), Nanos: int32(nullTime.Time.Nanosecond())}
	}
	return nil
}

// Value get time value.
func (t *Timestamp) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return t.Time(), nil
}

func (ts *Timestamp) GormDataType() string {
	return "time"
}

// Time get time.
func (t *Timestamp) Time() time.Time {
	return time.Unix(t.Seconds, int64(t.Nanos))
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return timei.MarshalJSON(t.Time())
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var st time.Time
	if err := timei.UnmarshalJSON(&st, data); err != nil {
		return err
	}
	t.Seconds, t.Nanos = st.Unix(), int32(st.Nanosecond())
	return nil
}

func (x *Timestamp) MarshalGQL(w io.Writer) {
	data, _ := timei.MarshalText(x.Time())
	w.Write(data)
}

func (x *Timestamp) UnmarshalGQL(v interface{}) error {
	var t time.Time
	if i, ok := v.(string); ok {
		err := timei.UnmarshalText(&t, []byte(i))
		if err != nil {
			return err
		}
		*x = Timestamp{Seconds: t.Unix(), Nanos: int32(t.Nanosecond())}
		return nil
	}
	return errors.New("enum need integer type")
}

type TimestampInput = Timestamp
