package time

import (
	"database/sql"
	"database/sql/driver"
	"github.com/hopeio/cherry/utils/encoding/binary"
	timei "github.com/hopeio/cherry/utils/time"
	"time"
)

func (ts *Timestamp) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ts = Timestamp{Millis: nullTime.Time.UnixMilli()}
	return
}

func (ts *Timestamp) Value() (driver.Value, error) {
	return time.UnixMilli(ts.Millis), nil
}

func (ts *Timestamp) Time() time.Time {
	return time.UnixMilli(ts.Millis)
}

func (ts *Timestamp) GormDataType() string {
	return "time"
}

func (ts *Timestamp) MarshalBinary() ([]byte, error) {
	return binary.ToBinary(ts.Millis), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (ts *Timestamp) UnmarshalBinary(data []byte) error {
	ts.Millis = binary.BinaryTo[int64](data)
	return nil
}

func (ts *Timestamp) GobEncode() ([]byte, error) {
	return ts.MarshalBinary()
}

func (ts *Timestamp) GobDecode(data []byte) error {
	return ts.UnmarshalBinary(data)
}

func (ts *Timestamp) MarshalJSON() ([]byte, error) {
	if ts == nil {
		return []byte("null"), nil
	}
	t := time.UnixMilli(ts.Millis)
	return timei.MarshalJSON(t)
}

func (ts *Timestamp) UnmarshalJSON(data []byte) error {
	var t time.Time
	err := timei.UnmarshalJSON(&t, data)
	if err != nil {
		return err
	}
	ts.Millis = t.UnixMilli()
	return err
}
