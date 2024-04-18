package time

import (
	"os/exec"
	"time"

	"github.com/hopeio/cherry/utils/log"
)

func Format(t time.Time) string {
	return t.Format(TimeFormat)
}

func TimeCost(start time.Time) {
	log.Info(time.Since(start))
}

// 设置系统时间
func SetUnixSysTime(t time.Time) {
	cmd := exec.Command("date", "-s", t.Format("01/02/2006 15:04:05.999999999"))
	cmd.Run()
}

func SyncHwTime() {
	cmd := exec.Command("clock --systohc")
	cmd.Run()
}

func TodayZeroTime() time.Time {
	todayZeroTime, _ := time.ParseInLocation(DateFormat, time.Now().Format(DateFormat), time.Local)
	return todayZeroTime
}

var ZeroTime = time.Time{}

// 标准化TimeDuration
func StdDuration(td time.Duration, stdTd time.Duration) time.Duration {
	if td == 0 {
		return td
	}
	if td < stdTd {
		return td * stdTd
	}
	return td
}
