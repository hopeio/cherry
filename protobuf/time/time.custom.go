package time

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/hopeio/cherry/utils/encoding/binary"
	timei "github.com/hopeio/cherry/utils/time"
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
	return x != nil && timei.Check(x) == 0
}

// CheckValid returns an error if the timestamp is invalid.
// In particular, it checks whether the value represents a date that is
// in the range of 0001-01-01T00:00:00Z to 9999-12-31T23:59:59Z inclusive.
// An error is reported for a nil Timestamp.
func (x *Time) CheckValid() error {
	switch timei.Check(x) {
	case timei.InvalidNil:
		return protoimpl.X.NewError("invalid nil time")
	case timei.InvalidUnderflow:
		return protoimpl.X.NewError("time (%v) before 0001-01-01", x)
	case timei.InvalidOverflow:
		return protoimpl.X.NewError("time (%v) after 9999-12-31", x)
	case timei.InvalidNanos:
		return protoimpl.X.NewError("time (%v) has out-of-range nanos", x)
	default:
		return nil
	}
}

func (ts *Time) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ts = Time{Seconds: nullTime.Time.Unix(), Nanos: int32(nullTime.Time.Nanosecond())}
	return
}

func (ts *Time) Value() (driver.Value, error) {
	return time.Unix(ts.Seconds, 0), nil
}

func (ts *Time) GormDataType() string {
	return "time"
}

func (ts *Time) Time() time.Time {
	return time.Unix(ts.Seconds, int64(ts.Nanos))
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
	if ts == nil {
		return []byte("null"), nil
	}
	return timei.MarshalJSON(ts.Time())
}

func (ts *Time) UnmarshalJSON(data []byte) error {
	var t time.Time
	err := timei.UnmarshalJSON(&t, data)
	if err != nil {
		return err
	}
	ts.Seconds, ts.Nanos = t.Unix(), int32(t.Nanosecond())
	return nil
}

func (x *Time) MarshalGQL(w io.Writer) {
	text, _ := timei.MarshalText(x.Time())
	w.Write(text)
}

func (x *Time) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(string); ok {
		var t time.Time
		err := timei.UnmarshalText(&t, []byte(i))
		if err != nil {
			return err
		}
		x.Seconds, x.Nanos = t.Unix(), int32(t.Nanosecond())
	}
	return errors.New("enum need integer type")
}

type TimeInput = Time
type DateInput = Date
