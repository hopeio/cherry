package time

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
	"time"
)

type TimestampConstraints interface {
	Timestamp(time.Time) int64
	Time(int64) time.Time
}

type timestamp[T TimestampConstraints] int64

func (t timestamp[T]) Time() time.Time {
	var v T
	return v.Time(int64(t))
}

// Scan scan time.
func (t *timestamp[T]) Scan(value interface{}) (err error) {
	var v T
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*t = timestamp[T](v.Timestamp(nullTime.Time))
	return
}

// Value get time value.
func (t timestamp[T]) Value() (driver.Value, error) {
	var v T
	return v.Time(int64(t)), nil
}

func (t timestamp[T]) MarshalJSON() ([]byte, error) {
	return strconv.AppendInt(nil, int64(t), 10), nil
}

func (t *timestamp[T]) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(data) == 0 || str == "null" {
		return nil
	}
	parseInt, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*t = timestamp[T](parseInt)
	return nil
}

// go生成的纳米级时间戳最后两位恒为0
type NanoTime struct{}

func (NanoTime) Timestamp(t time.Time) int64 {
	return t.UnixNano()
}
func (NanoTime) Time(t int64) time.Time {
	return time.Unix(0, t)
}

type MicroTime struct{}

func (MicroTime) Timestamp(t time.Time) int64 {
	return t.UnixMicro()
}
func (MicroTime) Time(t int64) time.Time {
	return time.UnixMicro(t)
}

type MilliTime struct{}

func (MilliTime) Timestamp(t time.Time) int64 {
	return t.UnixMilli()
}
func (MilliTime) Time(t int64) time.Time {
	return time.UnixMilli(t)
}

type SecondTime struct{}

func (SecondTime) Timestamp(t time.Time) int64 {
	return t.Unix()
}
func (SecondTime) Time(t int64) time.Time {
	return time.Unix(t, 0)
}

// 毫秒
type Timestamp = timestamp[MilliTime]

func NewTimeStamp(t time.Time) Timestamp {
	return Timestamp(t.UnixMilli())
}

type SecondTimestamp = timestamp[SecondTime]
type MilliTimestamp = timestamp[MilliTime]
type MicroTimestamp = timestamp[MicroTime]
type NanoTimestamp = timestamp[NanoTime]
