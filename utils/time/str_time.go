package time

import "time"

type Layout interface {
	Layout() string
}

type StrTime[T Layout] string

func (t StrTime[T]) Time() (time.Time, error) {
	var v T
	parse, err := time.Parse(v.Layout(), string(t))
	if err != nil {
		return time.Time{}, err
	}
	return parse, nil
}

func (dt StrTime[T]) MarshalJSON() ([]byte, error) {
	return []byte(`"` + dt + `"`), nil
}

func (dt *StrTime[T]) UnmarshalJSON(data []byte) error {
	str := string(data)
	// Ignore null, like in the main JSON package.
	if str == "null" {
		return nil
	}
	*dt = StrTime[T](str[1 : len(str)-1])
	return nil
}

// time.DateTime
type DateTimeStr = StrTime[EdateTime]

func NewDateTimeStr(t time.Time) DateTimeStr {
	return DateTimeStr(t.Format(time.DateTime))
}

type DateStr = StrTime[Edate]

func NewDateStr(t time.Time) DateStr {
	return DateStr(t.Format(time.DateOnly))
}
