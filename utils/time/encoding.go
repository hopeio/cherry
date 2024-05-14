package time

import (
	"strconv"
	"time"
)

type jsonType int8

const (
	JsonTypeLayout jsonType = iota
	JsonTypeUnixSeconds
	JsonTypeUnixMilliseconds
	JsonTypeUnixMicroseconds
	JsonTypeUnixNanoseconds
)

func SetJsonType(typ jsonType) {
	encoding.JsonType(typ)
}

func SetLayout(l string) {
	encoding.Layout(l)
}

var encoding = &Encoding{
	layout: time.RFC3339Nano,
}

func MarshalJSON(t time.Time) ([]byte, error) {
	return encoding.marshalJSON(t)
}

func UnmarshalJSON(t *time.Time, data []byte) error {
	return encoding.unmarshalJSON(t, data)
}

func MarshalText(t time.Time) ([]byte, error) {
	return encoding.marshalText(t)
}

func UnmarshalText(t *time.Time, data []byte) error {
	return encoding.unmarshalText(t, data)
}

type Encoding struct {
	jsonType
	layout string
}

func (u *Encoding) JsonType(typ jsonType) {
	u.jsonType = typ
}

func (u *Encoding) Layout(l string) {
	u.layout = l
}

func (u *Encoding) marshalText(t time.Time) ([]byte, error) {
	switch u.jsonType {
	case JsonTypeLayout:
		if u.layout == time.RFC3339Nano {
			return t.MarshalText()
		}
		return []byte(t.Format(u.layout)), nil
	case JsonTypeUnixSeconds:
		return strconv.AppendInt(nil, t.Unix(), 10), nil
	case JsonTypeUnixMilliseconds:
		return strconv.AppendInt(nil, t.UnixMilli(), 10), nil
	case JsonTypeUnixMicroseconds:
		return strconv.AppendInt(nil, t.UnixMicro(), 10), nil
	case JsonTypeUnixNanoseconds:
		return strconv.AppendInt(nil, t.UnixNano(), 10), nil
	}
	return t.MarshalText()
}

func (u *Encoding) unmarshalText(t *time.Time, data []byte) error {
	tstr := string(data)
	if tstr == "" {
		return nil
	}

	if u.jsonType == JsonTypeLayout {
		data = data[1 : len(data)-1]
		if u.layout == time.RFC3339Nano {
			return t.UnmarshalText(data)
		} else {
			var err error
			*t, err = time.ParseInLocation(u.layout, string(data), time.Local)
			return err
		}
	} else {
		parseInt, err := strconv.ParseInt(tstr, 10, 64)
		if err != nil {
			return err
		}
		switch u.jsonType {
		case JsonTypeUnixSeconds:
			*t = time.Unix(parseInt, 0)
			return nil
		case JsonTypeUnixMilliseconds:

			*t = time.UnixMilli(parseInt)
			return nil
		case JsonTypeUnixMicroseconds:
			*t = time.UnixMicro(parseInt)
			return nil
		case JsonTypeUnixNanoseconds:
			*t = time.Unix(0, parseInt)
			return nil
		}
	}

	return t.UnmarshalText(data)
}

func (u *Encoding) marshalJSON(t time.Time) ([]byte, error) {
	if u.jsonType == JsonTypeLayout {
		if u.layout == time.RFC3339Nano {
			return t.MarshalJSON()
		}
		return []byte(`"` + t.Format(u.layout) + `"`), nil
	} else {
		switch u.jsonType {
		case JsonTypeUnixSeconds:
			return strconv.AppendInt(nil, t.Unix(), 10), nil
		case JsonTypeUnixMilliseconds:
			return strconv.AppendInt(nil, t.UnixMilli(), 10), nil
		case JsonTypeUnixMicroseconds:
			return strconv.AppendInt(nil, t.UnixMicro(), 10), nil
		case JsonTypeUnixNanoseconds:
			return strconv.AppendInt(nil, t.UnixNano(), 10), nil
		}
	}

	return t.MarshalJSON()
}

func (u *Encoding) unmarshalJSON(t *time.Time, data []byte) error {
	tstr := string(data)
	if tstr == "null" {
		return nil
	}
	if u.jsonType == JsonTypeLayout {
		data = data[1 : len(data)-1]
		if u.layout == time.RFC3339Nano {
			return t.UnmarshalJSON(data)
		} else {
			var err error
			*t, err = time.ParseInLocation(`"`+u.layout+`"`, string(data), time.Local)
			return err
		}
	} else {
		parseInt, err := strconv.ParseInt(tstr, 10, 64)
		if err != nil {
			return err
		}
		switch u.jsonType {
		case JsonTypeUnixSeconds:
			*t = time.Unix(parseInt, 0)
			return nil
		case JsonTypeUnixMilliseconds:
			*t = time.UnixMilli(parseInt)
			return nil
		case JsonTypeUnixMicroseconds:
			*t = time.UnixMicro(parseInt)
			return nil
		case JsonTypeUnixNanoseconds:
			*t = time.Unix(0, parseInt)
			return nil
		}
	}

	return t.UnmarshalJSON(data)
}
