package types

import (
	"context"
)

type FGrpcServiceMethod[REQ, RES any] func(context.Context, REQ) (RES, error)

type Func func()

type FuncWithErr func() error

type FTask func(context.Context)
type FTaskWithErr func(context.Context) error
