package stamp

import (
	"database/sql/driver"
	"strconv"
	"time"
)

// 毫秒
type TimeStamp int64

// Scan scan time.
func (t *TimeStamp) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case time.Time:
		*t = TimeStamp(sc.UnixMilli())
	case string:
		var i int64
		i, err = strconv.ParseInt(sc, 10, 64)
		*t = TimeStamp(i)
	}
	return
}

// Value get time value.
func (t TimeStamp) Value() (driver.Value, error) {
	return time.UnixMilli(int64(t)), nil
}

// Time get time.
func (t TimeStamp) Time() time.Time {
	return time.UnixMilli(int64(t))
}

func Time(t time.Time) TimeStamp {
	return TimeStamp(t.UnixMilli())
}

func (t TimeStamp) MarshalJSON() ([]byte, error) {
	return t.Time().MarshalJSON()
}

func (t *TimeStamp) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(data) == 0 || str == "null" {
		return nil
	}
	parseInt, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*t = TimeStamp(parseInt)
	return nil
}
