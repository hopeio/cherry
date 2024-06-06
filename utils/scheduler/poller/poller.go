package poller

import (
	"context"
	time2 "github.com/hopeio/cherry/utils/time"
	"time"
)

type TaskFunc = func(context.Context)

type Poller struct {
	times uint
}

func NewPoller() *Poller {
	return &Poller{}
}

func (task *Poller) Times() uint {
	return task.times
}

func (task *Poller) Run(ctx context.Context, interval time.Duration, do TaskFunc) {
	timer := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		default:
			<-timer.C
			task.times++
			do(ctx)
		}
	}
}

func (task *Poller) RandRun(ctx context.Context, minInterval, maxInterval time.Duration, do TaskFunc) {

	timer := time2.NewRandTicker(minInterval, maxInterval)
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		default:
			task.times++
			do(ctx)
			timer.Wait()
		}
	}
}
