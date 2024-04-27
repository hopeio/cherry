package types

import (
	"github.com/hopeio/cherry/utils/constraints"
	"time"
)

type IGet[K constraints.Key, V any] interface {
	Get(key K) V
}

type ISet[K constraints.Key, V any] interface {
	Set(key K, v V)
}

type ISetWithExpire[K constraints.Key, V any] interface {
	SetWithExpire(key K, v V, expire time.Duration)
}

type IDelete[K constraints.Key, V any] interface {
	Delete(key K)
}

type IInit interface {
	Init()
}

type IStoreWithExpire[K constraints.Key, V any] interface {
	ISetWithExpire[K, V]
	IGet[K, V]
	IDelete[K, V]
}

type IStore[K constraints.Key, V any] interface {
	ISet[K, V]
	IGet[K, V]
	IDelete[K, V]
}
