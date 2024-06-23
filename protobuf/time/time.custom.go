package time

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/hopeio/cherry/utils/encoding/binary"
	"google.golang.org/protobuf/runtime/protoimpl"
	"io"
	"time"
)

func Now() *Time {
	return NewTime(time.Now())
}

// New constructs a new Timestamp from the provided time.Time.
func NewTime(t time.Time) *Time {
	return &Time{Seconds: t.Unix(), Nanos: int32(t.Nanosecond())}
}

// AsTime converts x to a time.Time.
func (x *Time) AsTime() time.Time {
	return time.Unix(x.GetSeconds(), int64(x.GetNanos()))
}

// IsValid reports whether the timestamp is valid.
// It is equivalent to CheckValid == nil.
func (x *Time) IsValid() bool {
	return x != nil && x.check() == 0
}

// CheckValid returns an error if the timestamp is invalid.
// In particular, it checks whether the value represents a date that is
// in the range of 0001-01-01T00:00:00Z to 9999-12-31T23:59:59Z inclusive.
// An error is reported for a nil Timestamp.
func (x *Time) CheckValid() error {
	switch x.check() {
	case invalidNil:
		return protoimpl.X.NewError("invalid nil Timestamp")
	case invalidUnderflow:
		return protoimpl.X.NewError("timestamp (%v) before 0001-01-01", x)
	case invalidOverflow:
		return protoimpl.X.NewError("timestamp (%v) after 9999-12-31", x)
	case invalidNanos:
		return protoimpl.X.NewError("timestamp (%v) has out-of-range nanos", x)
	default:
		return nil
	}
}

const (
	_ = iota
	invalidNil
	invalidUnderflow
	invalidOverflow
	invalidNanos
)

func (x *Time) check() uint {
	const minTimestamp = -62135596800  // Seconds between 1970-01-01T00:00:00Z and 0001-01-01T00:00:00Z, inclusive
	const maxTimestamp = +253402300799 // Seconds between 1970-01-01T00:00:00Z and 9999-12-31T23:59:59Z, inclusive
	secs := x.GetSeconds()
	nanos := x.GetNanos()
	switch {
	case x == nil:
		return invalidNil
	case secs < minTimestamp:
		return invalidUnderflow
	case secs > maxTimestamp:
		return invalidOverflow
	case nanos < 0 || nanos >= 1e9:
		return invalidNanos
	default:
		return 0
	}
}

func (ts *Time) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ts = Time{Seconds: nullTime.Time.UnixMilli(), Nanos: int32(nullTime.Time.Nanosecond())}
	return
}

func (ts *Time) Value() (driver.Value, error) {
	return time.Unix(ts.Seconds, 0), nil
}

func (ts *Time) GormDataType() string {
	return "time"
}

func (ts *Time) Time() time.Time {
	return time.Unix(ts.Seconds, 0)
}

func (ts *Time) MarshalBinary() ([]byte, error) {
	return binary.ToBinary(ts.Seconds), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (ts *Time) UnmarshalBinary(data []byte) error {
	ts.Seconds = binary.BinaryTo[int64](data)
	return nil
}

func (ts *Time) GobEncode() ([]byte, error) {
	return ts.MarshalBinary()
}

func (ts *Time) GobDecode(data []byte) error {
	return ts.UnmarshalBinary(data)
}

func (ts *Time) MarshalJSON() ([]byte, error) {
	t := time.Unix(ts.Seconds, 0)
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(time.DateOnly)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, time.DateOnly)
	b = append(b, '"')
	return b, nil
}

func (ts *Time) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(str) == 0 || str == "null" {
		return nil
	}
	t, err := time.ParseInLocation(time.DateOnly, str[1:len(str)-1], time.Local)
	if err != nil {
		return err
	}
	ts.Seconds = t.Unix()
	return nil
}
func (x *Time) MarshalGQL(w io.Writer) {
	w.Write([]byte(time.Unix(x.Seconds, 0).Format(time.DateOnly)))
}

func (x *Time) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(string); ok {
		t, err := time.ParseInLocation(time.DateOnly, i, time.Local)
		if err != nil {
			return err
		}
		x.Seconds = t.Unix()
		return nil
	}
	return errors.New("enum need integer type")
}

type TimeInput = Time
type DateInput = Date
