package engine

import (
	"context"
	"github.com/hopeio/cherry/utils/definition/types"
	"time"
)

type Kind uint8

const (
	KindNormal = iota
)

var (
	stdTimeout time.Duration = 0
)

type TaskMetaNew[T types.Key[KEY], KEY comparable] struct{}

type TaskMeta[KEY comparable] struct {
	id          uint64
	Kind        Kind
	Key         KEY
	Priority    int
	Describe    string
	createdAt   time.Time
	execBeginAt time.Time
	execEndAt   time.Time
	timeout     time.Duration
	TaskStatistics
}

func (t *TaskMeta[KEY]) OrderKey() int {
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
	reDoTimes uint
	errTimes  int
}

type Task[KEY comparable] struct {
	ctx context.Context
	TaskMeta[KEY]
	TaskFunc[KEY]
	errs []error
}

func (t *Task[KEY]) Errs() []error {
	return t.errs
}

type TaskInterface[KEY comparable] interface {
	Do(ctx context.Context) ([]*Task[KEY], error)
}

type Tasks[KEY comparable] []*Task[KEY]

func (tasks Tasks[KEY]) Less(i, j int) bool {
	return tasks[i].Priority > tasks[j].Priority
}

// ---------------

type ErrHandle func(context.Context, error)

type TaskFunc[KEY comparable] func(ctx context.Context) ([]*Task[KEY], error)

func (t TaskFunc[KEY]) Do(ctx context.Context) ([]*Task[KEY], error) {
	return t(ctx)
}

func emptyTaskFunc[KEY comparable](ctx context.Context) ([]*Task[KEY], error) {
	return nil, nil
}
