package time

import (
	"math/rand"
	"time"
)

type Ticker interface {
	Reset(time.Duration) bool
	Stop() bool
	Wait()
}

type FixTicker time.Ticker

func (t *FixTicker) Stop() bool {
	(*time.Ticker)(t).Stop()
	return true
}

func (t *FixTicker) Reset(d time.Duration) bool {
	(*time.Ticker)(t).Reset(d)
	return true
}

func (t *FixTicker) Wait() {
	<-t.C
}

func NewTicker(interval time.Duration) Ticker {
	return (*FixTicker)(time.NewTicker(interval))
}

var _ Ticker = &RandTicker{}

type RandTicker struct {
	timer                 *time.Timer
	limitBase, limitRange time.Duration
}

// 设置最小间隔
func (t *RandTicker) Reset(d time.Duration) bool {
	t.limitBase = d
	return t.reset()
}

func (t *RandTicker) reset() bool {
	if t.limitRange == 0 {
		return t.timer.Reset(t.limitBase)
	}
	return t.timer.Reset(t.limitBase + time.Duration(rand.Intn(int(t.limitRange))))
}

func (t *RandTicker) Wait() {
	<-t.timer.C
	t.reset()
}

func (t *RandTicker) Stop() bool {
	return t.timer.Stop()
}

// minInterval:最小等待时间
// maxInterval：最大等待时间
// maxInterval-minInterval: 等待范围
func NewRandTicker(minInterval, maxInterval time.Duration) *RandTicker {
	limitRange := maxInterval - minInterval

	return &RandTicker{
		timer:      time.NewTimer(minInterval + time.Duration(rand.Intn(int(limitRange)))),
		limitBase:  minInterval,
		limitRange: limitRange,
	}
}
