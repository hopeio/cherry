package ristretto

import (
	"github.com/dgraph-io/ristretto"
	"github.com/hopeio/cherry/utils/log"
)

type Config[K ristretto.Key, V any] ristretto.Config[K, V]

func (c *Config[K, V]) InitBeforeInject() {
}
func (c *Config[K, V]) Init() {
	if c.NumCounters == 0 {
		c.NumCounters = 1e7
	}
	if c.MaxCost == 0 {
		c.MaxCost = 1e6
	}
	if c.BufferItems == 0 {
		c.BufferItems = 64
	}
}

func (c *Config[K, V]) Build() *ristretto.Cache[K, V] {
	c.Init()
	cache, err := ristretto.NewCache((*ristretto.Config[K, V])(c))
	if err != nil {
		log.Fatal(err)
	}
	return cache
}

// 考虑换cache，ristretto存一个值，循环取居然还会miss(没开IgnoreInternalCost的原因),某个issue提要内存占用过大，直接初始化1.5MB
// freecache不能存对象，可能要为每个对象写UnmarshalBinary 和 MarshalBinary
// go-cache

type Cache[K ristretto.Key, V any] struct {
	*ristretto.Cache[K, V]
	Conf Config[K, V]
}

func (c *Cache[K, V]) Config() any {
	return &c.Conf
}

func (c *Cache[K, V]) SetEntity() {
	c.Cache = c.Conf.Build()
}

func (e *Cache[K, V]) Close() error {
	if e.Cache != nil {
		e.Cache.Close()
	}
	return nil
}
