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

type IDelete[K constraints.Key, V any] interface {
	Delete(key K)
}

type IInit interface {
	Init()
}

type StoreWithExpire[K constraints.Key, V any] interface {
	Set(k K, v V, expire time.Duration)
	IGet[K, V]
	IDelete[K, V]
}

type Store[K constraints.Key, V any] interface {
	ISet[K, V]
	IGet[K, V]
	IDelete[K, V]
}
