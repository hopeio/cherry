package engine

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	id2 "github.com/hopeio/cherry/utils/datastructure/idgen/id"
	"github.com/hopeio/cherry/utils/log"
	synci "github.com/hopeio/cherry/utils/sync"
	"sync/atomic"
	"time"
)

func (e *Engine[KEY]) Run(tasks ...*Task[KEY]) {
	e.lock.Lock()
	if e.isRunning {
		if len(tasks) > 0 {
			e.AddTasks(0, tasks...)
		}
		e.lock.Unlock()
		return
	}
	if !e.isRan {
		go func() {
			for task := range e.errChan {
				e.taskErrCount++
				e.errHandler(task)
			}
		}()
		e.addWorker()
		e.isRan = true
	}

	if !e.isRunning {
		e.isRunning = true
		e.wg.Add(1)
		go func() {
			timer := time.NewTimer(5 * time.Second)
			defer timer.Stop()
			var emptyTimes uint
			var readyTaskCh chan *Task[KEY]
			var readyTask *Task[KEY]
			taskChan := e.taskChanProducer
		loop:
			for {
				if len(e.taskReadyHeap) > 0 && readyTask == nil {
					readyTask, _ = e.taskReadyHeap.Pop()
					readyTaskCh = e.taskChanConsumer
				}

				// go 无法动态的设置case,但是可以动态的把channel置为nil
				if len(e.taskReadyHeap) >= int(e.limitWaitTaskCount) {
					taskChan = nil
				} else {
					taskChan = e.taskChanProducer
				}
				select {
				case readyTaskTmp := <-taskChan:
					e.taskReadyHeap.Push(readyTaskTmp)
				case readyTaskCh <- readyTask:
					readyTaskCh = nil
					readyTask = nil
				case <-timer.C:
					//检测任务是否已空
					if uint(e.workerCount) == 0 && len(e.taskReadyHeap) == 0 {
						e.lock.Lock()
						counter, _ := synci.WaitGroupState(&e.wg)
						if counter == 1 {
							emptyTimes++
							if emptyTimes > 2 {
								log.Debug("the task is about to end.")
								e.wg.Done()
								e.isRunning = false
								e.lock.Unlock()
								break loop
							}
						}
						e.lock.Unlock()
					}
					log.Debugf("[Running] task:%d/%d/%d,worker: %d/%d\r", e.taskDoneCount, e.taskTotalCount, e.taskFailedCount, e.workerCount, e.currentWorkerCount)
					timer.Reset(e.monitorInterval)
				case <-e.ctx.Done():
					if err := e.ctx.Err(); err != nil {
						log.Error(err)
					}
					break loop
				}

			}
		}()
	}

	e.lock.Unlock()
	if len(tasks) > 0 {
		e.AddTasks(0, tasks...)
	}
	e.wg.Wait()
	e.isFinished = true
	log.Infof("[END] total:%d,done:%d,failed:%d", e.taskTotalCount, e.taskDoneCount, e.taskFailedCount)
}

func (e *Engine[KEY]) newWorker(readyTask *Task[KEY]) {
	atomic.AddUint64(&e.currentWorkerCount, 1)
	//id := c.currentWorkerCount
	// 这里考虑回复多channel,worker数量多起来的时候,channel维护的goroutine数量太多
	worker := &Worker[KEY]{Id: uint(e.currentWorkerCount)}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.ErrorStack(r)
				log.Info(spew.Sdump(readyTask))
				atomic.AddUint64(&e.taskFailedCount, 1)
				e.wg.Done()
				// 创建一个新的
				e.newWorker(nil)
			}
			atomic.AddUint64(&e.currentWorkerCount, ^uint64(0))
		}()

		if readyTask != nil {
			e.ExecTask(worker, readyTask)
		}
		for {
			select {
			case readyTask = <-e.taskChanConsumer:
				e.ExecTask(worker, readyTask)
			case <-e.ctx.Done():
				return
			}
		}
	}()
	e.workers = append(e.workers, worker)
}

func (e *Engine[KEY]) addWorker() {
	if e.currentWorkerCount != 0 {
		return
	}
	e.newWorker(nil)
	go func() {
		for {
			select {
			case readyTask := <-e.taskChanConsumer:
				if e.currentWorkerCount < e.limitWorkerCount {
					e.newWorker(readyTask)
				} else {
					log.Info("worker count is full")
					e.taskChanProducer <- readyTask
					return
				}
			case <-e.ctx.Done():
				return
			}
		}
	}()

}

func (e *Engine[KEY]) AddNoPriorityTasks(tasks ...*Task[KEY]) {
	e.AddTasks(0, tasks...)
}

func (e *Engine[KEY]) AddTasks(generation int, tasks ...*Task[KEY]) {
	l := len(tasks)
	atomic.AddUint64(&e.taskTotalCount, uint64(l))
	e.wg.Add(l)
	for _, task := range tasks {
		if task == nil || task.TaskFunc == nil {
			continue
		}
		task.Priority += generation
		task.id = id2.GenOrderID()
		e.taskChanProducer <- task
	}
}

func (e *Engine[KEY]) AsyncAddTasks(generation int, tasks ...*Task[KEY]) {
	if len(tasks) > 0 {
		go e.AddTasks(generation, tasks...)
	}
}

func (e *Engine[KEY]) AddWorker(num int) {
	atomic.AddUint64(&e.limitWorkerCount, uint64(num))
}

func (e *Engine[KEY]) NewFixedWorker(interval time.Duration) int {
	taskChan := make(chan *Task[KEY])
	worker := &Worker[KEY]{Id: uint(e.currentWorkerCount), taskCh: taskChan}
	e.fixedWorkers = append(e.fixedWorkers, worker)
	e.newFixedWorker(worker, interval)
	return len(e.fixedWorkers) - 1
}

func (e *Engine[KEY]) newFixedWorker(worker *Worker[KEY], interval time.Duration) {
	go func() {
		var task *Task[KEY]
		defer func() {
			if r := recover(); r != nil {
				log.ErrorStack(r)
				log.Info(spew.Sdump(task))
				atomic.AddUint64(&e.taskFailedCount, 1)
				e.wg.Done()
				// 创建一个新的
				e.newFixedWorker(worker, interval)
			}
			atomic.AddUint64(&e.currentWorkerCount, ^uint64(0))
		}()
		var timer *time.Ticker
		// 如果有任务时间间隔,
		if interval > 0 {
			timer = time.NewTicker(interval)
		}
		for task = range worker.taskCh {
			if interval > 0 {
				<-timer.C
			}
			e.ExecTask(worker, task)
		}
	}()
}

func (e *Engine[KEY]) AddFixedTasks(workerId int, generation int, tasks ...*Task[KEY]) error {

	if workerId > len(e.fixedWorkers)-1 {
		return fmt.Errorf("不存在workId为%d的worker,请调用NewFixedWorker添加", workerId)
	}
	worker := e.fixedWorkers[workerId]
	l := len(tasks)
	atomic.AddUint64(&e.taskTotalCount, uint64(l))
	e.wg.Add(l)
	for _, task := range tasks {
		if task == nil || task.TaskFunc == nil {
			continue
		}
		task.Priority += generation
		task.id = id2.GenOrderID()
		worker.taskCh <- task
	}
	return nil
}

func (e *Engine[KEY]) RunSingleWorker(tasks ...*Task[KEY]) {
	e.limitWorkerCount = 1
	e.Run(tasks...)
}

func (e *Engine[KEY]) Stop() {
	e.cancel()
	close(e.taskChanConsumer)
	close(e.taskChanProducer)
	close(e.errChan)
	for _, worker := range e.fixedWorkers {
		close(worker.taskCh)
	}
	if e.speedLimit != nil {
		e.speedLimit.Stop()
	}
	e.done.Close()
	for _, kindHandler := range e.kindHandlers {
		if kindHandler != nil {
			if kindHandler.speedLimit != nil {
				kindHandler.speedLimit.Stop()
			}
			if kindHandler.rateLimiter != nil {
				kindHandler.rateLimiter = nil
			}
		}
	}

	for _, callback := range e.stopCallBack {
		callback()
	}
}

func (e *Engine[KEY]) ExecTask(worker *Worker[KEY], task *Task[KEY]) {
	atomic.AddUint64(&e.workerCount, 1)
	worker.isExecuting = true
	worker.currentTask = task
	if task != nil {
		if task.TaskFunc != nil {
			if task.ctx == nil {
				task.ctx = e.ctx
			}
			if !e.execTask(task) {
				atomic.AddUint64(&e.taskDoneCount, ^uint64(0))
			}
		}
	}
	atomic.AddUint64(&e.taskDoneCount, 1)
	atomic.AddUint64(&e.workerCount, ^uint64(0))
	e.wg.Done()
	worker.isExecuting = false
}

func (e *Engine[KEY]) execTask(task *Task[KEY]) bool {

	if task.Key != e.zeroKey {
		if _, ok := e.done.Get(task.Key); ok {
			return false
		}
	}

	if e.speedLimit != nil {
		e.speedLimit.Wait()
	}

	if e.rateLimiter != nil {
		e.rateLimiter.Wait(task.ctx)
	}

	var kindHandler *KindHandler[KEY]
	if e.kindHandlers != nil && int(task.Kind) < len(e.kindHandlers) {
		kindHandler = e.kindHandlers[task.Kind]
	}

	if kindHandler != nil {
		if kindHandler.Skip {
			return false
		}

		if kindHandler.speedLimit != nil {
			kindHandler.speedLimit.Wait()
		}
		if kindHandler.rateLimiter != nil {
			_ = kindHandler.rateLimiter.Wait(task.ctx)
		}
	}

	if task.reExecTimes > 0 {
		task.reExecLogs = append(task.reExecLogs, &ExecLog{
			execBeginAt: time.Now(),
		})
	} else {
		task.execBeginAt = time.Now()
	}
	tasks, err := task.TaskFunc.Do(task.ctx)
	if task.reExecTimes > 0 {
		task.reExecLogs[len(task.reExecLogs)-1].execEndAt = time.Now()
	} else {
		task.execEndAt = time.Now()
	}

	if err != nil {
		task.errTimes++
		if task.reExecTimes > 0 {
			task.reExecLogs[len(task.reExecLogs)-1].err = err
		} else {
			task.err = err
		}

		if task.errTimes < 5 {
			task.reExecTimes++
			log.Warnf("%v执行失败:%v,将第%d次执行", task.Key, err, task.reExecTimes)
			e.AsyncAddTasks(task.Priority+1, task)
		} else {
			log.Warn(task.Key, "多次执行失败:", err, ",将执行错误处理")
			e.errChan <- task
		}
		return false
	}
	if task.Key != e.zeroKey {
		e.done.SetWithTTL(task.Key, struct{}{}, 1, time.Hour)
	}
	if len(tasks) > 0 {
		e.AsyncAddTasks(task.Priority+1, tasks...)
	}
	return true
}

func (e *Engine[KEY]) Cancel() {
	log.Info("任务取消")
	e.cancel()
	synci.WaitGroupStopWait(&e.wg)

}

func (e *Engine[KEY]) CancelAfter(interval time.Duration) *Engine[KEY] {
	time.AfterFunc(interval, e.Cancel)
	return e
}

func (e *Engine[KEY]) StopAfter(interval time.Duration) *Engine[KEY] {
	time.AfterFunc(interval, e.Stop)
	return e
}
