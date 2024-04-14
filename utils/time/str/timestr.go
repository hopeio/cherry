package str

import "time"

// time.DateTime
type TimeStr string

func (t TimeStr) Time() (time.Time, error) {
	parse, err := time.Parse(time.DateTime, string(t))
	if err != nil {
		return time.Time{}, err
	}
	return parse, nil
}

func Time(t time.Time) TimeStr {
	return TimeStr(t.Format(time.DateTime))
}
