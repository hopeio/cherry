package clause

import (
	dbi "github.com/hopeio/cherry/utils/dao/db"
	"github.com/hopeio/cherry/utils/definition/types/request"
	"gorm.io/gorm/clause"
)

type RangeReq[T request.Ordered] request.RangeReq[T]

func (req *RangeReq[T]) Clause() clause.Expression {
	if req == nil || req.RangeField == "" {
		return new(EmptyClause)
	}
	// 泛型还很不好用，这种方式代替原来的interface{}
	zero := *new(T)
	operation := dbi.Between
	if req.RangeEnd == zero && req.RangeStart != zero {
		operation = dbi.Greater
		if req.Include {
			operation = dbi.GreaterOrEqual
		}
		return NewWhereClause(req.RangeField, operation, req.RangeStart)
	}
	if req.RangeStart == zero && req.RangeEnd != zero {
		operation = dbi.Less
		if req.Include {
			operation = dbi.LessOrEqual
		}
		return NewWhereClause(req.RangeField, operation, req.RangeStart)
	}
	if req.RangeStart != zero && req.RangeEnd != zero {
		if req.Include {
			return NewWhereClause(req.RangeField, operation, req.RangeStart, req.RangeEnd)
		} else {
			return clause.Where{Exprs: []clause.Expression{NewWhereClause(req.RangeField, dbi.Greater, req.RangeStart), NewWhereClause(req.RangeField, dbi.Less, req.RangeStart)}}
		}
	}
	return new(EmptyClause)
}
