package engine

import (
	"context"
	"github.com/hopeio/cherry/utils/log"
	"strconv"
	"strings"
	"time"
)

type Kind uint32

const (
	KindNormal = iota
)

var (
	stdTimeout time.Duration = 0
)

type execLog struct {
	execBeginAt time.Time
	execEndAt   time.Time
	err         error
}

type TaskStatistics struct {
	ReExecTimes int
	ErrTimes    int
}

type Task[KEY Key] struct {
	context.Context
	Kind      Kind
	Key       KEY
	Priority  int
	Describe  string
	id        uint64
	createdAt time.Time
	execLog
	reExecLogs []*execLog // 多数任务只会执行一次
	deadline   time.Time
	timeout    time.Duration
	TaskStatistics
	TaskFunc[KEY]
}

func NewTask[KEY Key](task TaskFunc[KEY]) *Task[KEY] {
	return &Task[KEY]{
		TaskFunc: task,
	}
}

func (t *Task[KEY]) SetContext(ctx context.Context) *Task[KEY] {
	t.Context = ctx
	return t
}

func (t *Task[KEY]) SetPriority(priority int) *Task[KEY] {
	t.Priority = priority
	return t
}

func (t *Task[KEY]) SetKind(k Kind) *Task[KEY] {
	t.Kind = k
	return t
}

func (t *Task[KEY]) SetKey(key KEY) *Task[KEY] {
	t.Key = key
	return t
}

func (t *Task[KEY]) SetDescribe(describe string) *Task[KEY] {
	t.Describe = describe
	return t
}

func (t *Task[KEY]) Id() uint64 {
	return t.id
}

func (t *Task[KEY]) Compare(t2 *Task[KEY]) int {
	return t.Priority - t2.Priority
}

func (t *Task[KEY]) Errs() []error {
	var errs []error
	if t.err != nil {
		errs = append(errs, t.err)
	}
	for _, log := range t.reExecLogs {
		errs = append(errs, log.err)
	}
	return errs
}

func (t *Task[KEY]) ErrLog() {
	builder := strings.Builder{}
	if t.err != nil {
		builder.WriteString("[1]{")
		builder.WriteString(t.err.Error())
		builder.WriteString("}\n")
	}
	for i, log := range t.reExecLogs {
		if log.err != nil {
			builder.WriteString("[" + strconv.Itoa(i+2) + "]{")
			builder.WriteString(log.err.Error())
			builder.WriteString("}\n")
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
