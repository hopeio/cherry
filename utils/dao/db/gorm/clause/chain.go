package clause

import (
	"context"
	dbi "github.com/hopeio/cherry/utils/dao/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type ChainDB func(db *gorm.DB) *gorm.DB

func noOp(db *gorm.DB) *gorm.DB {
	return db
}

func (c ChainDB) ById(id int) ChainDB {
	if id > 0 {
		return c.ByIdNoCheck(id)
	}
	return noOp
}

func (c ChainDB) ByIdNoCheck(id int) ChainDB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (c ChainDB) ByName(name string) ChainDB {
	if name != "" {
		return c.ByNameNoCheck(name)
	}
	return noOp
}

func (c ChainDB) ByNameNoCheck(name string) ChainDB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func (c ChainDB) List(dest any) ChainDB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Find(dest)
	}
}

func NewChainDB(ctx context.Context) ChainDB {
	return func(db *gorm.DB) *gorm.DB {
		return db.WithContext(ctx)
	}
}

type Expression dbi.FilterExpr

func (e *Expression) Clause() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(e.Field+(*dbi.FilterExpr)(e).Operation.SQL(), e.Value...)
	}
}

func NewClause(field string, op dbi.Operation, args ...interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+op.SQL(), args...)
	}
}

type ChainClause []func(db *gorm.DB) *gorm.DB

// db.Scope(ById(1),ByName("a")).First(v)
func (c ChainClause) ById(id int) ChainClause {
	if id > 0 {
		return c.ByIdNoCheck(id)
	}
	return c
}

func (c ChainClause) ByIdNoCheck(id int) ChainClause {
	return append(c, NewClause("id", dbi.Equal, id))
}

func (c ChainClause) ByName(name string) ChainClause {
	if name != "" {
		return c.ByNameNoCheck(name)
	}
	return c
}

func (c ChainClause) ByNameNoCheck(name string) ChainClause {
	return append(c, func(db *gorm.DB) *gorm.DB {
		return db.Where(`name = ?`, name)
	})
}

func (c ChainClause) Or(clauses ...ChainClause) ChainClause {
	return append(c, func(db *gorm.DB) *gorm.DB {
		db = db.Where(strings.Join(c.Exprs, " OR "), c.Vars...)
		return db
	})
}

func (c ChainClause) Exec(db *gorm.DB) *gorm.DB {
	db = db.Scopes(c...)
	return db
}

type Clause3 []clause.Interface

func (c Clause3) ByIdNoCheck(id int) Clause3 {
	return append(c, clause.Where{Exprs: []clause.Expression{clause.Eq{Column: "id", Value: id}}})
}
