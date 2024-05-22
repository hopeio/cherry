package clause

import (
	dbi "github.com/hopeio/cherry/utils/dao/db"
	"github.com/hopeio/cherry/utils/types/request"
	"gorm.io/gorm/clause"
)

type RangeReq[T request.Ordered] request.RangeReq[T]

func (req *RangeReq[T]) Clause() clause.Expression {
	if req == nil || req.RangeField == "" {
		return nil
	}

	var zero T
	operation := dbi.Between
	if req.RangeEnd == zero && req.RangeBegin != zero {
		operation = dbi.Greater
		if req.Include {
			operation = dbi.GreaterOrEqual
		}
		return NewWhereClause(req.RangeField, operation, req.RangeBegin)
	}
	if req.RangeBegin == zero && req.RangeEnd != zero {
		operation = dbi.Less
		if req.Include {
			operation = dbi.LessOrEqual
		}
		return NewWhereClause(req.RangeField, operation, req.RangeBegin)
	}
	if req.RangeBegin != zero && req.RangeEnd != zero {
		if req.Include {
			return NewWhereClause(req.RangeField, operation, req.RangeBegin, req.RangeEnd)
		} else {
			return clause.Where{Exprs: []clause.Expression{NewWhereClause(req.RangeField, dbi.Greater, req.RangeBegin), NewWhereClause(req.RangeField, dbi.Less, req.RangeBegin)}}
		}
	}
	return nil
}
