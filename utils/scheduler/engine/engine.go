package engine

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgraph-io/ristretto"
	"github.com/hopeio/cherry/utils/datastructure/heap/idxless"
	"github.com/hopeio/cherry/utils/io/fs"
	"github.com/hopeio/cherry/utils/log"
	"github.com/hopeio/cherry/utils/slices"
	time2 "github.com/hopeio/cherry/utils/time"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

// 目前受限于ristretto.Cache的泛型限制,考虑移除并引入lru或boolong filter
type Key interface {
	uint64 | string | byte | int | int32 | uint32 | int64
}

type Config[KEY Key] struct {
	WorkerCount uint
}

func (c *Config[KEY]) NewEngine() *Engine[KEY] {
	return NewEngine[KEY](c.WorkerCount)
}

type Engine[KEY Key] struct {
	limitWorkerCount, currentWorkerCount, workerCount uint64
	limitWaitTaskCount                                uint
	workers                                           []*Worker[KEY]
	fixedWorkers                                      []*Worker[KEY] // 固定只执行一种任务的worker,避免并发问题
	taskChanProducer                                  chan *Task[KEY]
	taskChanConsumer                                  chan *Task[KEY]
	taskReadyHeap                                     heap.Heap[*Task[KEY], Tasks[KEY]]
	ctx                                               context.Context
	cancel                                            context.CancelFunc // 手动停止执行
	wg                                                sync.WaitGroup     // 控制确保所有任务执行完
	speedLimit                                        time2.Ticker
	rateLimiter                                       *rate.Limiter
	//TODO
	monitorInterval              time.Duration // 全局检测定时器间隔时间，任务的卡住检测，worker panic recover都可以用这个检测
	enableTracing, enableMetrics bool
	isRunning, isFinished, isRan bool
	lock                         sync.RWMutex
	EngineStatistics
	done         *ristretto.Cache[KEY, struct{}]
	kindHandlers []*KindHandler[KEY]
	errHandler   func(task *Task[KEY])
	errChan      chan *Task[KEY]
	stopCallBack []func()
	zeroKey      KEY // 泛型不够强大,又为了性能妥协的字段
}

type KindHandler[KEY Key] struct {
	Skip        bool
	speedLimit  time2.Ticker
	rateLimiter *rate.Limiter
	// TODO 指定Kind的Handler
	HandleFun TaskFunc[KEY]
}

func NewEngine[KEY Key](workerCount uint) *Engine[KEY] {
	return NewEngineWithContext[KEY](workerCount, context.Background())
}

func NewEngineWithContext[KEY Key](workerCount uint, ctx context.Context) *Engine[KEY] {
	ctx, cancel := context.WithCancel(ctx)
	cache, _ := ristretto.NewCache(&ristretto.Config[KEY, struct{}]{
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
		taskChanProducer:   make(chan *Task[KEY]),
		taskChanConsumer:   make(chan *Task[KEY]),
		taskReadyHeap:      heap.Heap[*Task[KEY], Tasks[KEY]]{},
		monitorInterval:    5 * time.Second,
		done:               cache,
		errHandler:         func(task *Task[KEY]) { task.ErrLog() },
		lock:               sync.RWMutex{},
		errChan:            make(chan *Task[KEY]),
		zeroKey:            *new(KEY),
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
	if interval < time.Second {
		log.Warn("monitor interval min one second")
		interval = time.Second
	}

	e.monitorInterval = interval
}

func (e *Engine[KEY]) ErrHandler(errHandler func(task *Task[KEY])) *Engine[KEY] {
	e.errHandler = errHandler
	return e
}

func (e *Engine[KEY]) ErrHandlerUtilSuccess() *Engine[KEY] {
	log.Warn("it will clear history exec log contains err")
	return e.ErrHandler(func(task *Task[KEY]) {
		task.errTimes = 0
		task.reExecLogs = task.reExecLogs[:0]
		e.AsyncAddTasks(task.Priority, task)
	})
}

func (e *Engine[KEY]) ErrHandlerRetryTimes(times int) *Engine[KEY] {
	return e.ErrHandler(func(task *Task[KEY]) {
		if task.reExecTimes < times {
			task.errTimes = 0
			task.reExecLogs = task.reExecLogs[:0]
			e.AsyncAddTasks(task.Priority, task)
		} else {
			task.ErrLog()
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
	e.speedLimit = time2.NewTicker(interval)
	return e
}

func (e *Engine[KEY]) RandSpeedLimited(minInterval, maxInterval time.Duration) *Engine[KEY] {
	e.speedLimit = time2.NewRandTicker(minInterval, maxInterval)
	return e
}

func (e *Engine[KEY]) KindSpeedLimit(kind Kind, interval time.Duration) *Engine[KEY] {
	limiter := time2.NewRandTicker(interval, interval)
	e.kindSpeedLimit(kind, limiter)
	return e
}

func (e *Engine[KEY]) KindRandSpeedLimit(kind Kind, minInterval, maxInterval time.Duration) *Engine[KEY] {
	limiter := time2.NewRandTicker(minInterval, maxInterval)
	e.kindSpeedLimit(kind, limiter)
	return e
}

func (e *Engine[KEY]) kindSpeedLimit(kind Kind, limiter time2.Ticker) *Engine[KEY] {
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
	limiter := time2.NewRandTicker(interval, interval)
	for _, kind := range kinds {
		e.kindSpeedLimit(kind, limiter)
	}
	return e
}

func (e *Engine[KEY]) KindGroupRandSpeedLimit(minInterval, maxInterval time.Duration, kinds ...Kind) *Engine[KEY] {
	limiter := time2.NewRandTicker(minInterval, maxInterval)
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
