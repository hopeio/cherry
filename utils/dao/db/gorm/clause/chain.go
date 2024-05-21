package clause

import (
	dbi "github.com/hopeio/cherry/utils/dao/db"
	"gorm.io/gorm/clause"
)

type ChainClause []clause.Interface

func (c ChainClause) ById(id int) ChainClause {
	if id != 0 {
		return c.ByIdNoCheck(id)
	}
	return c
}

func (c ChainClause) ByIdNoCheck(id any) ChainClause {
	return append(c, clause.Where{Exprs: []clause.Expression{clause.Eq{Column: dbi.ColumnId, Value: id}}})
}

func (c ChainClause) ByName(name string) ChainClause {
	if name != "" {
		return c.ByNameNoCheck(name)
	}
	return c
}

func (c ChainClause) ByNameNoCheck(name string) ChainClause {
	return append(c, clause.Where{Exprs: []clause.Expression{clause.Eq{Column: dbi.ColumnName, Value: name}}})
}
