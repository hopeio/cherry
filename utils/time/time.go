package time

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
	"time"
)

type Date = Time[date]
type date struct{}

func (date) Layout() string {
	return LayoutDate
}

type UnixSecondTime = UnixTime[SecondTime]
type UnixMilliTime = UnixTime[MilliTime]
type UnixMicroTime = UnixTime[MicroTime]
type UnixNanoTime = UnixTime[NanoTime]

type DateTime = Time[dateTime]

type dateTime struct{}

func (dateTime) Layout() string {
	return time.DateTime
}

type DefaultTime struct{}

func (DefaultTime) Layout() string {
	return time.RFC3339
}

type Layout interface {
	Layout() string
}

type Time[T Layout] time.Time

func NewTime[T Layout](t time.Time) Time[T] {
	return Time[T](t)
}

func (dt *Time[T]) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*dt = Time[T](nullTime.Time)
	return
}

func (dt Time[T]) Value() (driver.Value, error) {
	return time.Time(dt), nil
}

func (dt Time[T]) Format(format string) string {
	return time.Time(dt).Format(format)
}

func (dt Time[T]) GormDataType() string {
	var v T
	switch any(v).(type) {
	case DateTime:
		return "date"
	default:
		return "datetime"
	}
}

func (dt Time[T]) MarshalBinary() ([]byte, error) {
	return time.Time(dt).MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (dt *Time[T]) UnmarshalBinary(data []byte) error {
	return (*time.Time)(dt).UnmarshalBinary(data)
}

func (dt Time[T]) GobEncode() ([]byte, error) {
	return dt.MarshalBinary()
}

func (dt *Time[T]) GobDecode(data []byte) error {
	return dt.UnmarshalBinary(data)
}

func (dt Time[T]) MarshalJSON() ([]byte, error) {
	var v T
	layout := v.Layout()
	if layout != "" {
		b := make([]byte, 0, len(layout)+2)
		b = append(b, '"')
		b = time.Time(dt).AppendFormat(b, layout)
		b = append(b, '"')
		return b, nil
	}
	return time.Time(dt).MarshalJSON()
}

func (dt *Time[T]) UnmarshalJSON(data []byte) error {
	str := string(data)
	// Ignore null, like in the main JSON package.
	if str == "null" {
		return nil
	}
	var v T
	layout := v.Layout()
	if layout != "" {
		t, err := time.ParseInLocation(`"`+layout+`"`, string(data), time.Local)
		*dt = (Time[T])(t)
		return err

	}
	return (*time.Time)(dt).UnmarshalJSON(data)
}

type UnixTime[T TimestampConstraints] time.Time

func NewUnixTime[T TimestampConstraints](t time.Time) UnixTime[T] {
	return UnixTime[T](t)
}

func (dt *UnixTime[T]) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*dt = UnixTime[T](nullTime.Time)
	return
}

func (dt UnixTime[T]) Value() (driver.Value, error) {
	return time.Time(dt), nil
}

func (dt UnixTime[T]) Format(format string) string {
	return time.Time(dt).Format(format)
}

func (dt UnixTime[T]) GormDataType() string {
	return "datetime"
}

func (dt UnixTime[T]) MarshalBinary() ([]byte, error) {
	return time.Time(dt).MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (dt *UnixTime[T]) UnmarshalBinary(data []byte) error {
	return (*time.Time)(dt).UnmarshalBinary(data)
}

func (dt UnixTime[T]) GobEncode() ([]byte, error) {
	return dt.MarshalBinary()
}

func (dt *UnixTime[T]) GobDecode(data []byte) error {
	return dt.UnmarshalBinary(data)
}

func (dt UnixTime[T]) MarshalJSON() ([]byte, error) {
	var v T
	switch any(v).(type) {
	case SecondTime:
		return strconv.AppendInt(nil, time.Time(dt).Unix(), 10), nil
	case MilliTime:
		return strconv.AppendInt(nil, time.Time(dt).UnixMilli(), 10), nil
	case MicroTime:
		return strconv.AppendInt(nil, time.Time(dt).UnixMicro(), 10), nil
	case NanoTime:
		return strconv.AppendInt(nil, time.Time(dt).UnixNano(), 10), nil
	default:
		return time.Time(dt).MarshalJSON()
	}
}

func (dt *UnixTime[T]) UnmarshalJSON(data []byte) error {
	str := string(data)
	// Ignore null, like in the main JSON package.
	if str == "null" {
		return nil
	}
	var v T
	ts, err := strconv.ParseInt(str, 10, 64)
	switch any(v).(type) {
	case SecondTime:
		*dt = (UnixTime[T])(time.Unix(ts, 0))
		return err
	case MilliTime:
		*dt = (UnixTime[T])(time.UnixMilli(ts))
		return err
	case MicroTime:
		*dt = (UnixTime[T])(time.UnixMicro(ts))
		return err
	case NanoTime:
		*dt = (UnixTime[T])(time.Unix(0, ts))
		return err
	default:
		return (*time.Time)(dt).UnmarshalJSON(data)
	}
}
