package gorm

import (
	"context"
	"gorm.io/gorm"
)

type ChainDao struct {
	DB, OriginDB *gorm.DB
}

func NewChainDao(ctx context.Context, db *gorm.DB) *ChainDao {
	return &ChainDao{
		db.WithContext(ctx), db.WithContext(ctx),
	}
}

func (c *ChainDao) ResetDB() {
	c.DB = c.OriginDB
}

func (c *ChainDao) ById(id int) *ChainDao {
	if id != 0 {
		return c.ByIdNoCheck(id)
	}
	return c
}

func (c *ChainDao) ByName(name string) *ChainDao {
	if name != "" {
		return c.ByNameNoCheck(name)
	}
	return c
}

func (c *ChainDao) ByIdNoCheck(id int) *ChainDao {
	c.DB = c.DB.Where(`id = ?`, id)
	return c
}

func (c *ChainDao) ByNameNoCheck(name string) *ChainDao {
	c.DB = c.DB.Where(`name = ?`, name)
	return c
}
