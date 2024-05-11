package time

import "time"

// time.DateTime
type Timestr string

func (t Timestr) Time() (time.Time, error) {
	parse, err := time.Parse(time.DateTime, string(t))
	if err != nil {
		return time.Time{}, err
	}
	return parse, nil
}

func NewTimeStr(t time.Time) Timestr {
	return Timestr(t.Format(time.DateTime))
}
