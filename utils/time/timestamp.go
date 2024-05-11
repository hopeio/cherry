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
	return t.Time(), nil
}

// Time get time.
func (t Timestamp) Time() time.Time {
	return time.UnixMilli(int64(t))
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return strconv.AppendInt(nil, int64(t), 10), nil
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
