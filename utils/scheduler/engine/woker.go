package engine

import (
	"context"
	"time"
)

type Type uint8

const (
	normalType Type = iota
	fixedType
)

type Worker[KEY comparable] struct {
	Id          uint
	Type        Type
	Kind        Kind
	taskCh      chan *Task[KEY]
	currentTask *Task[KEY]
	isExecuting bool
	ctx         context.Context
}

// WorkStatistics worker统计数据
type WorkStatistics struct {
	averageTimeCost                                                                  time.Duration
	taskDoneCount, taskTotalCount, taskErrorCount, taskTimeoutCount, taskFailedCount uint64
}

// EngineStatistics 基本引擎统计数据
type EngineStatistics struct {
	WorkStatistics
	taskErrCount uint64
}
