package time

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
	"strings"
	"time"
)

// 毫秒
type Timestamp int64

func NewTimeStamp(t time.Time) Timestamp {
	return Timestamp(t.UnixMilli())
}

// Scan scan time.
func (t *Timestamp) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*t = Timestamp(nullTime.Time.UnixMilli())
	return
}

// Value get time value.
func (t Timestamp) Value() (driver.Value, error) {
	return time.UnixMilli(int64(t)), nil
}

// Time get time.
func (t Timestamp) Time() time.Time {
	return time.UnixMilli(int64(t))
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return t.Time().MarshalJSON()
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(data) == 0 || str == "null" {
		return nil
	}
	// 2018-08-08 00:00:00
	if strings.Contains(str, "-") {
		var st time.Time
		if err := st.UnmarshalJSON(data); err != nil {
			return err
		}
		*t = Timestamp(st.UnixMilli())
		return nil
	}
	parseInt, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*t = Timestamp(parseInt)
	return nil
}

type SecondTimeStamp int64

func (t SecondTimeStamp) Time() time.Time {
	return time.Unix(int64(t), 0)
}

func (ts *SecondTimeStamp) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ts = SecondTimeStamp(nullTime.Time.Unix())
	return
}

func (ts SecondTimeStamp) Value() (driver.Value, error) {
	return time.Unix(int64(ts), 0), nil
}

func (ts SecondTimeStamp) Format(foramt string) string {
	return time.Unix(int64(ts), 0).Format(foramt)
}

// GormDataType gorm common data type
func (ts SecondTimeStamp) GormDataType() string {
	return "datetime"
}

func (t SecondTimeStamp) MarshalJSON() ([]byte, error) {
	return t.Time().MarshalJSON()
}

func (t *SecondTimeStamp) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(data) == 0 || str == "null" {
		return nil
	}
	// 2018-08-08 00:00:00
	if strings.Contains(str, "-") {
		var st time.Time
		if err := st.UnmarshalJSON(data); err != nil {
			return err
		}
		*t = SecondTimeStamp(st.Second())
		return nil
	}
	parseInt, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*t = SecondTimeStamp(parseInt)
	return nil
}

type NanoTimeStamp int64

func (t NanoTimeStamp) Time() time.Time {
	return time.Unix(0, int64(t))
}

func (ts *NanoTimeStamp) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ts = NanoTimeStamp(nullTime.Time.UnixNano())
	return
}

func (ts NanoTimeStamp) Value() (driver.Value, error) {
	return time.Unix(0, int64(ts)), nil
}

func (ts NanoTimeStamp) Format(foramt string) string {
	return time.Unix(0, int64(ts)).Format(foramt)
}

// GormDataType gorm common data type
func (ts NanoTimeStamp) GormDataType() string {
	return "datetime"
}

func (t NanoTimeStamp) MarshalJSON() ([]byte, error) {
	return t.Time().MarshalJSON()
}

func (t *NanoTimeStamp) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(data) == 0 || str == "null" {
		return nil
	}
	// 2018-08-08 00:00:00
	if strings.Contains(str, "-") {
		var st time.Time
		if err := st.UnmarshalJSON(data); err != nil {
			return err
		}
		*t = NanoTimeStamp(st.UnixNano())
		return nil
	}
	parseInt, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*t = NanoTimeStamp(parseInt)
	return nil
}
