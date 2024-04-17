package interfaces

import (
	"github.com/hopeio/cherry/utils/definition/constraints"
	"time"
)

type Get[K constraints.Key, V any] interface {
	Get(key K) V
}

type Set[K constraints.Key, V any] interface {
	Set(key K, v V)
}

type Delete[K constraints.Key, V any] interface {
	Delete(key K)
}

type Init interface {
	Init()
}

type StoreWithExpire[K constraints.Key, V any] interface {
	Set(k K, v V, expire time.Duration)
	Get[K, V]
	Delete[K, V]
}

type Store[K constraints.Key, V any] interface {
	Set[K, V]
	Get[K, V]
	Delete[K, V]
}
