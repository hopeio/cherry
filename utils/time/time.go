package time

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"strconv"
	"time"
)

type Constraints interface {
	DefaultTime | DateTime | UnixTime | UnixMilliTime | UnixNanoTime | DisplayTime
}

type DefaultTime struct{}
type DateTime struct{}
type UnixTime struct{}
type UnixMilliTime struct{}
type UnixMicroTime struct{}
type UnixNanoTime struct{}
type DisplayTime struct{}

// TODO: 这是不优雅的实现
type Time[T any] time.Time

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
	t := time.Time(dt)
	switch any(v).(type) {
	case DisplayTime:
		b := make([]byte, 0, len(TimeFormatDisplay)+2)
		b = append(b, '"')
		b = t.AppendFormat(b, TimeFormatDisplay)
		b = append(b, '"')
		return b, nil
	case DateTime:
		b := make([]byte, 0, len(DateFormat)+2)
		b = append(b, '"')
		b = t.AppendFormat(b, DateFormat)
		b = append(b, '"')
		return b, nil
	case UnixTime:
		return strconv.AppendInt(nil, time.Time(dt).Unix(), 10), nil
	case UnixMilliTime:
		return strconv.AppendInt(nil, time.Time(dt).UnixMilli(), 10), nil
	case UnixMicroTime:
		return strconv.AppendInt(nil, time.Time(dt).UnixMicro(), 10), nil
	case UnixNanoTime:
		return strconv.AppendInt(nil, time.Time(dt).UnixNano(), 10), nil
	default:
		return t.MarshalJSON()
	}
}

func (dt *Time[T]) UnmarshalJSON(data []byte) error {
	str := string(data)
	// Ignore null, like in the main JSON package.
	if str == "null" {
		return nil
	}
	var v T
	switch any(v).(type) {
	case DisplayTime:
		t, err := time.ParseInLocation(`"`+TimeFormatDisplay+`"`, string(data), time.Local)
		*dt = (Time[T])(t)
		return err
	case DateTime:
		t, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)
		*dt = (Time[T])(t)
		return err
	case UnixTime:
		ts, err := strconv.ParseInt(str, 10, 64)
		*dt = (Time[T])(time.Unix(ts, 0))
		return err
	case UnixMilliTime:
		ts, err := strconv.ParseInt(str, 10, 64)
		*dt = (Time[T])(time.UnixMilli(ts))
		return err
	case UnixMicroTime:
		ts, err := strconv.ParseInt(str, 10, 64)
		*dt = (Time[T])(time.UnixMicro(ts))
		return err
	case UnixNanoTime:
		ts, err := strconv.ParseInt(str, 10, 64)
		*dt = (Time[T])(time.Unix(0, ts))
		return err
	default:
		return (*time.Time)(dt).UnmarshalJSON(data)
	}
}

// 对应数据库datetime或timestamp,或date
// typ 0 序列化为 "2006-01-02 15:04:05",typ 1序列化为"2006-01-02",typ 2 序列化为秒时间戳, typ 3序列化为毫秒时间戳, typ 4 序列化为纳秒时间戳
// 序列化,反序列化前需设置typ
type UnionTime struct {
	time.Time
	typ uint8
}

func NewUnionTime(t time.Time, typ uint8) UnionTime {
	return UnionTime{Time: t, typ: typ}
}

func ZeroUnionTime(typ uint8) UnionTime {
	return UnionTime{typ: typ}
}
func (ut *UnionTime) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ut = UnionTime{Time: nullTime.Time}
	return
}

func (ut UnionTime) Value() (driver.Value, error) {
	if ut.typ == 1 {
		return ut.Format(DateFormat), nil
	}
	return ut.Time, nil
}

func (ut UnionTime) MarshalJSON() ([]byte, error) {
	t := ut.Time
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	switch ut.typ {
	case 0:
		b := make([]byte, 0, len(TimeFormatDisplay)+2)
		b = append(b, '"')
		b = t.AppendFormat(b, TimeFormatDisplay)
		b = append(b, '"')
		return b, nil
	case 1:
		b := make([]byte, 0, len(DateFormat)+2)
		b = append(b, '"')
		b = t.AppendFormat(b, DateFormat)
		b = append(b, '"')
		return b, nil
	case 2:
		return stringsi.ToBytes(strconv.Itoa(int(t.Unix()))), nil
	case 3:
		return stringsi.ToBytes(strconv.Itoa(int(t.UnixNano()))), nil
	}

	return nil, nil
}

func (ut *UnionTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	var err error
	var t time.Time
	switch ut.typ {
	case 0:
		t, err = time.ParseInLocation(`"`+TimeFormatDisplay+`"`, string(data), time.Local)
	case 1:
		t, err = time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)
	case 2:
		str, err := strconv.Atoi(stringsi.BytesToString(data))
		if err != nil {
			return err
		}
		t = time.Unix(int64(str), 0)
	case 3:
		str, err := strconv.Atoi(stringsi.BytesToString(data))
		if err != nil {
			return err
		}
		t = time.Unix(0, int64(str))
	}
	*ut = UnionTime{Time: t, typ: ut.typ}
	return err
}

func (ut *UnionTime) Type(typ uint8) UnionTime {
	ut.typ = typ
	return *ut
}
