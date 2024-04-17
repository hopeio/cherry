package rate

import (
	"math/rand"
	"time"
)

type SpeedLimiter interface {
	Reset(time.Duration) bool
	Stop() bool

	Wait()

	Channel() <-chan time.Time
}

type Ticker time.Ticker

func (t *Ticker) Stop() bool {
	(*time.Ticker)(t).Stop()
	return true
}

func (t *Ticker) Reset(d time.Duration) bool {
	(*time.Ticker)(t).Reset(d)
	return true
}

func (t *Ticker) Wait() {
	<-t.C
}

func (t *Ticker) Channel() <-chan time.Time {
	return t.C
}

func NewSpeedLimiter(interval time.Duration) SpeedLimiter {
	return (*Ticker)(time.NewTicker(interval))
}

var _ SpeedLimiter = &RandTimer{}

type RandTimer struct {
	timer                 *time.Timer
	limitBase, limitRange time.Duration
}

func (t *RandTimer) Stop() bool {
	return t.timer.Stop()
}

// 设置最小间隔
func (t *RandTimer) Reset(d time.Duration) bool {
	t.limitBase = d
	return t.reset()
}

func (t *RandTimer) reset() bool {
	if t.limitRange == 0 {
		return t.timer.Reset(t.limitBase)
	}
	return t.timer.Reset(t.limitBase + time.Duration(rand.Intn(int(t.limitRange))))
}

func (t *RandTimer) Wait() {
	<-t.timer.C
	t.reset()
}

func (t *RandTimer) Channel() <-chan time.Time {
	return t.timer.C
}

// minInterval:最小等待时间
// maxInterval：最大等待时间
// maxInterval-minInterval: 等待范围
func NewRandSpeedLimiter(minInterval, maxInterval time.Duration) SpeedLimiter {
	limitRange := maxInterval - minInterval
	if limitRange == 0 {
		return (*Ticker)(time.NewTicker(minInterval))
	}
	return &RandTimer{
		timer:      time.NewTimer(minInterval + time.Duration(rand.Intn(int(limitRange)))),
		limitBase:  minInterval,
		limitRange: limitRange,
	}
}
