package engine

import (
	"context"
	"github.com/hopeio/cherry/utils/log"
	"strings"
	"time"
)

type Kind uint8

const (
	KindNormal = iota
)

var (
	stdTimeout time.Duration = 0
)

type TaskMeta[KEY Key] struct {
	id        uint64
	Kind      Kind
	Key       KEY
	Priority  int
	Describe  string
	createdAt time.Time
	ExecLog
	reExecLogs []*ExecLog // 多数任务只会执行一次
	timeout    time.Duration
	TaskStatistics
}

type ExecLog struct {
	execBeginAt time.Time
	execEndAt   time.Time
	err         error
}

func (t *TaskMeta[KEY]) SortKey() int {
	return t.Priority
}

func (t *TaskMeta[KEY]) SetPriority(priority int) {
	t.Priority = priority
}

func (t *TaskMeta[KEY]) SetKind(k Kind) {
	t.Kind = k
}

func (t *TaskMeta[KEY]) SetKey(key KEY) {
	t.Key = key
}

func (t *TaskMeta[KEY]) Id() uint64 {
	return t.id
}

type TaskStatistics struct {
	reExecTimes int
	errTimes    int
}

type Task[KEY Key] struct {
	ctx context.Context
	TaskMeta[KEY]
	TaskFunc[KEY]
}

func (t *Task[KEY]) SetContext(ctx context.Context) {
	t.ctx = ctx
}

func (t *Task[KEY]) ErrLog() {
	builder := strings.Builder{}
	if t.err != nil {
		builder.WriteString(t.err.Error())
	}
	for _, log := range t.reExecLogs {
		if log.err != nil {
			builder.WriteByte(' ')
			builder.WriteString(log.err.Error())
		}
	}
	log.Error(builder.String())
}

type TaskInterface[KEY Key] interface {
	Do(ctx context.Context) ([]*Task[KEY], error)
}

type Tasks[KEY Key] []*Task[KEY]

func (tasks Tasks[KEY]) Less(i, j int) bool {
	return tasks[i].Priority > tasks[j].Priority
}

// ---------------

type ErrHandle func(context.Context, error)

type TaskFunc[KEY Key] func(ctx context.Context) ([]*Task[KEY], error)

func (t TaskFunc[KEY]) Do(ctx context.Context) ([]*Task[KEY], error) {
	return t(ctx)
}

func emptyTaskFunc[KEY Key](ctx context.Context) ([]*Task[KEY], error) {
	return nil, nil
}
