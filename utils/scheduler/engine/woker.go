package engine

import (
	"time"
)

type Type uint32

const (
	normalType Type = iota
	fixedType
)

type Worker[KEY Key] struct {
	id                      uint
	typ                     Type
	kind                    Kind
	taskCh                  chan *Task[KEY]
	createdAt               time.Time
	currentTask             *Task[KEY]
	isExecuting, canExecute bool
}

// workStatistics worker统计数据
type workStatistics struct {
	timeCost                                                                          time.Duration
	taskTotalCount, taskDoneCount, taskSkipCount, taskErrHandleCount, taskFailedCount uint64
	taskRepeatTimes, taskErrorTimes, taskTimeoutTimes                                 uint64
}

// EngineStatistics 基本引擎统计数据
type EngineStatistics struct {
	workStatistics
}
