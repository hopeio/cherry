package types

import (
	"context"
)

type GRPCServiceMethod[REQ, RES any] func(context.Context, REQ) (RES, error)

type Func func()

type FuncWithErr func() error

type TaskFunc func(context.Context)
type TaskFuncWithErr func(context.Context) error
