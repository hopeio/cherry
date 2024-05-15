package gorm

import (
	"gorm.io/gorm"
)

type ChainDao struct {
	DB, OriginDB *gorm.DB
}

func NewChainDao(db *gorm.DB) *ChainDao {
	return &ChainDao{
		db, db,
	}
}

func (c *ChainDao) ResetDB() {
	c.DB = c.OriginDB
}

func (c *ChainDao) By(field string, v any) *ChainDao {
	c.DB = c.DB.Where(field+` = ?`, v)
	return c
}

func (c *ChainDao) ByIdCheck(id int) *ChainDao {
	if id != 0 {
		return c.ById(id)
	}
	return c
}

func (c *ChainDao) ByNameCheck(name string) *ChainDao {
	if name != "" {
		return c.ByName(name)
	}
	return c
}

func (c *ChainDao) ById(id int) *ChainDao {
	c.DB = c.DB.Where(`id = ?`, id)
	return c
}

func (c *ChainDao) ByName(name string) *ChainDao {
	c.DB = c.DB.Where(`name = ?`, name)
	return c
}
