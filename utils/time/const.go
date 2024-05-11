package time

import "time"

const (
	SecondsOfDay    = 24 * 60 * 60
	SecondsOfMinute = 60
	TimeDay         = SecondsOfDay * time.Second
)

const (
	TimeFormat            = "2006-01-02 15:04:05.999999"
	TimeFormatDisplay     = "2006-01-02 15:04:05"
	TimeFormatPostgresDB  = time.RFC3339
	TimeFormatNoDate      = "15:04:05"
	DayEndTime            = "23:59:59"
	DayEndTimeWithSpace   = " 23:59:59"
	DayBeginTime          = "00:00:00"
	DayBeginTimeWithSpace = " 00:00:00"
	DateFormat            = "2006-01-02"
	TimeFormatCompact     = "20060102150405"
	TimeFormatRFC1        = "2006/01/02 - 15:04:05"
)

const (
	Day        = time.Hour * 24
	MonthDay30 = Day * 30
	MonthDay31 = Day * 31
	MonthDay28 = Day * 28
	MonthDay29 = Day * 29
	Month      = MonthDay30
	YearDay365 = Day * 365
	YearDay366 = Day * 366
	Year       = YearDay365
)

const (
	January   = "January"
	February  = "February"
	March     = "March"
	April     = "April"
	May       = "May"
	June      = "June"
	July      = "July"
	August    = "August"
	September = "September"
	October   = "October"
	November  = "November"
	December  = "December"
)

const (
	Monday    = "Monday"
	Tuesday   = "Tuesday"
	Wednesday = "Wednesday"
	Thursday  = "Thursday"
	Friday    = "Friday"
	Saturday  = "Saturday"
	Sunday    = "Sunday"
)
