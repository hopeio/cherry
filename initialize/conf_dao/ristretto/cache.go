package ristretto

import (
	"github.com/dgraph-io/ristretto"
	"github.com/hopeio/cherry/utils/log"
)

type Config ristretto.Config

func (c *Config) InitBeforeInject() {
}
func (c *Config) InitAfterInject() {
	if c.NumCounters == 0 {
		c.NumCounters = 10000000
	}
	if c.MaxCost == 0 {
		c.MaxCost = 1000000
	}
	if c.BufferItems == 0 {
		c.BufferItems = 64
	}
}

func (c *Config) Build() *ristretto.Cache {
	c.InitAfterInject()
	cache, err := ristretto.NewCache((*ristretto.Config)(c))
	if err != nil {
		log.Fatal(err)
	}
	return cache
}

// 考虑换cache，ristretto存一个值，循环取居然还会miss(没开IgnoreInternalCost的原因),某个issue提要内存占用过大，直接初始化1.5MB
// freecache不能存对象，可能要为每个对象写UnmarshalBinary 和 MarshalBinary
// go-cache

type Cache struct {
	*ristretto.Cache
	Conf Config
}

func (c *Cache) Config() any {
	return &c.Conf
}

func (c *Cache) SetEntity() {
	c.Cache = c.Conf.Build()
}

func (e *Cache) Close() error {
	if e.Cache != nil {
		e.Cache.Close()
	}
	return nil
}
