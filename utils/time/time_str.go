package time

import "time"

type TimeStr[T Layout] string

func (t TimeStr[T]) Time() (time.Time, error) {
	var v T
	parse, err := time.Parse(v.Layout(), string(t))
	if err != nil {
		return time.Time{}, err
	}
	return parse, nil
}

// time.DateTime
type DateTimeStr = TimeStr[dateTime]

func NewDateTimeStr(t time.Time) DateTimeStr {
	return DateTimeStr(t.Format(time.DateTime))
}

type DateStr = TimeStr[date]

func NewDateStr(t time.Time) DateStr {
	return DateStr(t.Format(time.DateOnly))
}
