package rate

import (
	"math/rand"
	"time"
)

type SpeedLimiter interface {
	Reset()
	Stop() bool

	Wait()

	Channel() <-chan time.Time
}

type Ticker time.Ticker

func (t *Ticker) Stop() bool {
	(*time.Ticker)(t).Stop()
	return true
}

func (t *Ticker) Reset() {

}

func (t *Ticker) Wait() {
	<-t.C
}

func (t *Ticker) Channel() <-chan time.Time {
	return t.C
}

type RandTimer struct {
	*time.Timer
	randSpeedLimitBase, randSpeedLimitRange time.Duration
}

func NewSpeedLimiter(interval time.Duration) SpeedLimiter {
	return (*Ticker)(time.NewTicker(interval))
}

// minInterval:最小等待时间
// maxInterval：最大等待时间
// maxInterval-minInterval: 等待范围
func NewRandSpeedLimiter(minInterval, maxInterval time.Duration) SpeedLimiter {
	randSpeedLimitRange := maxInterval - minInterval
	if randSpeedLimitRange == 0 {
		return (*Ticker)(time.NewTicker(minInterval))
	}
	return &RandTimer{
		Timer:               time.NewTimer(minInterval + time.Duration(rand.Intn(int(randSpeedLimitRange)))),
		randSpeedLimitBase:  minInterval,
		randSpeedLimitRange: randSpeedLimitRange,
	}
}

func (t *RandTimer) Stop() bool {
	return t.Timer.Stop()
}

func (t *RandTimer) Reset() {
	if t.randSpeedLimitRange == 0 {
		t.Timer.Reset(t.randSpeedLimitBase)
	}
	t.Timer.Reset(t.randSpeedLimitBase + time.Duration(rand.Intn(int(t.randSpeedLimitRange))))
}

func (t *RandTimer) Wait() {
	<-t.Timer.C
	t.Reset()
}

func (t *RandTimer) Channel() <-chan time.Time {
	return t.C
}
