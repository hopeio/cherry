package configor

import (
	"github.com/hopeio/cherry/utils/log"
	"time"
)

func DurationNotify(td time.Duration, stdTd time.Duration) {
	if td > 0 && td < stdTd {
		log.GetSkipLogger(1).Warnf("except: %s level,but got %s", stdTd, td)
	}
}

// 标准化TimeDuration
func StdDuration(td time.Duration, stdTd time.Duration) time.Duration {
	if td == 0 {
		return td
	}
	// 1/10为可容忍数量级?
	if td < stdTd/10 {
		return td * stdTd
	}
	return td
}
