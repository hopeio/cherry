package engine

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgraph-io/ristretto"
	"github.com/hopeio/cherry/utils/datastructure/heap"
	"github.com/hopeio/cherry/utils/datastructure/list/list"
	"github.com/hopeio/cherry/utils/io/fs"
	"github.com/hopeio/cherry/utils/log"
	rate2 "github.com/hopeio/cherry/utils/scheduler/rate"
	"github.com/hopeio/cherry/utils/slices"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type Config[KEY comparable] struct {
	WorkerCount uint
}

func (c *Config[KEY]) NewEngine() *Engine[KEY] {
	return NewEngine[KEY](c.WorkerCount)
}

type Engine[KEY comparable] struct {
	limitWorkerCount, currentWorkerCount uint64
	limitWaitTaskCount                   uint
	workerChan                           chan *Worker[KEY]
	workers                              []*Worker[KEY]
	workerReadyList                      list.List[*Worker[KEY]]
	taskChan                             chan *Task[KEY]
	taskReadyHeap                        heap.Heap[*Task[KEY], Tasks[KEY]]
	ctx                                  context.Context
	cancel                               context.CancelFunc // 手动停止执行
	wg                                   sync.WaitGroup     // 控制确保所有任务执行完
	fixedWorkers                         []*Worker[KEY]     // 固定只执行一种任务的worker,避免并发问题
	speedLimit                           rate2.SpeedLimiter
	rateLimiter                          *rate.Limiter
	//TODO
	monitorInterval              time.Duration // 全局检测定时器间隔时间，任务的卡住检测，worker panic recover都可以用这个检测
	isRunning, isFinished, isRan bool
	lock                         sync.RWMutex
	EngineStatistics
	done         *ristretto.Cache
	kindHandlers []*KindHandler[KEY]
	errHandler   func(task *Task[KEY])
	errChan      chan *Task[KEY]
	stopCallBack []func()
}

type KindHandler[KEY comparable] struct {
	Skip        bool
	speedLimit  rate2.SpeedLimiter
	rateLimiter *rate.Limiter
	// TODO 指定Kind的Handler
	HandleFun TaskFunc[KEY]
}

func NewEngine[KEY comparable](workerCount uint) *Engine[KEY] {
	return NewEngineWithContext[KEY](workerCount, context.Background())
}

func NewEngineWithContext[KEY comparable](workerCount uint, ctx context.Context) *Engine[KEY] {
	ctx, cancel := context.WithCancel(ctx)
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters:        1e4,   // number of keys to track frequency of (10M).
		MaxCost:            1e3,   // maximum cost of cache (MaxCost * 1MB).
		BufferItems:        64,    // number of keys per Get buffer.
		Metrics:            false, // number of keys per Get buffer.
		IgnoreInternalCost: true,
	})
	return &Engine[KEY]{
		limitWorkerCount:   uint64(workerCount),
		limitWaitTaskCount: workerCount * 10,
		ctx:                ctx,
		cancel:             cancel,
		workerChan:         make(chan *Worker[KEY]),
		taskChan:           make(chan *Task[KEY]),
		workerReadyList:    list.New[*Worker[KEY]](),
		taskReadyHeap:      heap.Heap[*Task[KEY], Tasks[KEY]]{},
		monitorInterval:    time.Second,
		done:               cache,
		errHandler: func(task *Task[KEY]) {
			log.Error(task.errs)
		},
		lock:    sync.RWMutex{},
		errChan: make(chan *Task[KEY]),
	}
}

func (e *Engine[KEY]) Context() context.Context {
	return e.ctx
}

func (e *Engine[KEY]) SkipKind(kinds ...Kind) *Engine[KEY] {
	length := slices.Max(kinds) + 1
	if e.kindHandlers == nil {
		e.kindHandlers = make([]*KindHandler[KEY], length)
	}
	if int(length) > len(e.kindHandlers) {
		e.kindHandlers = append(e.kindHandlers, make([]*KindHandler[KEY], int(length)-len(e.kindHandlers))...)
	}
	for _, kind := range kinds {
		if e.kindHandlers[kind] == nil {
			e.kindHandlers[kind] = &KindHandler[KEY]{Skip: true}
		} else {
			e.kindHandlers[kind].Skip = true
		}

	}
	return e
}

func (e *Engine[KEY]) MonitorInterval(interval time.Duration) {
	e.monitorInterval = interval
}

func (e *Engine[KEY]) ErrHandler(errHandler func(task *Task[KEY])) *Engine[KEY] {
	e.errHandler = errHandler
	return e
}

func (e *Engine[KEY]) ErrHandlerUtilSuccess() *Engine[KEY] {
	return e.ErrHandler(func(task *Task[KEY]) {
		task.errs = task.errs[:0]
		e.AsyncAddTasks(task.Priority, task)
	})
}

func (e *Engine[KEY]) ErrHandlerRetryTimes(times int) *Engine[KEY] {
	return e.ErrHandler(func(task *Task[KEY]) {
		if task.errTimes < times {
			task.errs = task.errs[:0]
			e.AsyncAddTasks(task.Priority, task)
		} else {
			log.Error(task.errs)
		}

	})
}

func (e *Engine[KEY]) ErrHandlerWriteToFile(path string) *Engine[KEY] {
	file, err := fs.Create(path)
	if err != nil {
		panic(err)
	}
	e.StopCallBack(func() {
		file.Close()
	})
	return e.ErrHandler(func(task *Task[KEY]) {
		spew.Fdump(file, task)
	})
}

func (e *Engine[KEY]) StopCallBack(callBack func()) *Engine[KEY] {
	e.stopCallBack = append(e.stopCallBack, callBack)
	return e
}

func (e *Engine[KEY]) SpeedLimited(interval time.Duration) *Engine[KEY] {
	e.speedLimit = rate2.NewSpeedLimiter(interval)
	return e
}

func (e *Engine[KEY]) RandSpeedLimited(minInterval, maxInterval time.Duration) *Engine[KEY] {
	e.speedLimit = rate2.NewRandSpeedLimiter(minInterval, maxInterval)
	return e
}

func (e *Engine[KEY]) KindSpeedLimit(kind Kind, interval time.Duration) *Engine[KEY] {
	limiter := rate2.NewRandSpeedLimiter(interval, interval)
	e.kindSpeedLimit(kind, limiter)
	return e
}

func (e *Engine[KEY]) KindRandSpeedLimit(kind Kind, minInterval, maxInterval time.Duration) *Engine[KEY] {
	limiter := rate2.NewRandSpeedLimiter(minInterval, maxInterval)
	e.kindSpeedLimit(kind, limiter)
	return e
}

func (e *Engine[KEY]) kindSpeedLimit(kind Kind, limiter rate2.SpeedLimiter) *Engine[KEY] {
	if e.kindHandlers == nil {
		e.kindHandlers = make([]*KindHandler[KEY], int(kind)+1)
	}
	if int(kind)+1 > len(e.kindHandlers) {
		e.kindHandlers = append(e.kindHandlers, make([]*KindHandler[KEY], int(kind)+1-len(e.kindHandlers))...)
	}
	if e.kindHandlers[kind] == nil {
		e.kindHandlers[kind] = &KindHandler[KEY]{speedLimit: limiter}
	} else {
		e.kindHandlers[kind].speedLimit = limiter
	}
	return e
}

// 多个kind共用一个timer
func (e *Engine[KEY]) KindGroupSpeedLimit(interval time.Duration, kinds ...Kind) *Engine[KEY] {
	limiter := rate2.NewRandSpeedLimiter(interval, interval)
	for _, kind := range kinds {
		e.kindSpeedLimit(kind, limiter)
	}
	return e
}

func (e *Engine[KEY]) KindGroupRandSpeedLimit(minInterval, maxInterval time.Duration, kinds ...Kind) *Engine[KEY] {
	limiter := rate2.NewRandSpeedLimiter(minInterval, maxInterval)
	for _, kind := range kinds {
		e.kindSpeedLimit(kind, limiter)
	}
	return e
}

func (e *Engine[KEY]) Limiter(r rate.Limit, b int) *Engine[KEY] {
	e.rateLimiter = rate.NewLimiter(r, b)
	return e
}

func (e *Engine[KEY]) KindLimiter(kind Kind, r rate.Limit, b int) *Engine[KEY] {
	e.kindLimiter(kind, r, b)
	return e
}

func (e *Engine[KEY]) kindLimiter(kind Kind, r rate.Limit, b int) {
	if e.kindHandlers == nil {
		e.kindHandlers = make([]*KindHandler[KEY], int(kind)+1)
	}
	if int(kind)+1 > len(e.kindHandlers) {
		e.kindHandlers = append(e.kindHandlers, make([]*KindHandler[KEY], int(kind)+1-len(e.kindHandlers))...)
	}
	if e.kindHandlers[kind] == nil {
		e.kindHandlers[kind] = &KindHandler[KEY]{rateLimiter: rate.NewLimiter(r, b)}
	} else {
		e.kindHandlers[kind].rateLimiter = rate.NewLimiter(r, b)
	}
}

// TaskSourceChannel 任务源,参数是一个channel,channel关闭时，代表任务源停止发送任务
func (e *Engine[KEY]) TaskSourceChannel(taskSourceChannel <-chan *Task[KEY]) {
	e.wg.Add(1)
	go func() {
		for task := range taskSourceChannel {
			if task == nil || task.TaskFunc == nil {
				continue
			}
			e.AddTasks(0, task)
		}
		e.wg.Done()
	}()
}

// TaskSourceFunc,参数为添加任务的函数，直到该函数运行结束，任务引擎才会检测任务是否结束
func (e *Engine[KEY]) TaskSourceFunc(taskSourceFunc func(*Engine[KEY])) {
	e.wg.Add(1)
	go func() {
		taskSourceFunc(e)
		e.wg.Done()
	}()
}
