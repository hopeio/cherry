package time

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/hopeio/cherry/utils/encoding/binary"
	"io"
	"time"
)

func (ts *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ts = Date{Seconds: nullTime.Time.UnixMilli()}
	return
}

func (ts *Date) Value() (driver.Value, error) {
	return time.Unix(ts.Seconds, 0), nil
}

func (ts *Date) GormDataType() string {
	return "time"
}

func (ts *Date) Time() time.Time {
	return time.Unix(ts.Seconds, 0)
}

func (ts *Date) MarshalBinary() ([]byte, error) {
	return binary.ToBinary(ts.Seconds), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (ts *Date) UnmarshalBinary(data []byte) error {
	ts.Seconds = binary.BinaryTo[int64](data)
	return nil
}

func (ts *Date) GobEncode() ([]byte, error) {
	return ts.MarshalBinary()
}

func (ts *Date) GobDecode(data []byte) error {
	return ts.UnmarshalBinary(data)
}

func (ts *Date) MarshalJSON() ([]byte, error) {
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

func (ts *Date) UnmarshalJSON(data []byte) error {
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
func (x *Date) MarshalGQL(w io.Writer) {
	w.Write([]byte(time.Unix(x.Seconds, 0).Format(time.DateOnly)))
}

func (x *Date) UnmarshalGQL(v interface{}) error {
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
