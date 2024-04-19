package poller

import (
	"context"
	"github.com/hopeio/cherry/utils/scheduler/rate"
	"time"
)

type TaskFunc = func(context.Context)

type Poller struct {
	times         uint
	firstExecuted bool
	do            TaskFunc
}

func NewPoller() *Poller {
	return &Poller{}
}

func (task *Poller) Times() uint {
	return task.times
}

func (task *Poller) Run(ctx context.Context, interval time.Duration, do TaskFunc) {
	task.do = do
	timer := time.NewTicker(interval)
	if !task.firstExecuted {
		task.times = 1
		task.do(ctx)
		task.firstExecuted = true
	}
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			task.times++
			task.do(ctx)
		}
	}
}

func (task *Poller) RandRun(ctx context.Context, minInterval, maxInterval time.Duration, do TaskFunc) {
	task.do = do
	timer := rate.NewRandSpeedLimiter(minInterval, maxInterval)
	ch := timer.Channel()
	task.times = 1
	task.do(ctx)
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-ch:
			task.times++
			task.do(ctx)
			timer.Reset(minInterval)
		}
	}
}
