//go:build go1.18

package clause

import (
	"github.com/hopeio/cherry/utils/types/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PageSortReq request.PageSortReq

func (req *PageSortReq) Clause() []clause.Expression {
	if req.PageNo == 0 && req.PageSize == 0 {
		return nil
	}
	if req.SortReq == nil || req.SortReq.SortField == "" {
		return []clause.Expression{Page(req.PageNo, req.PageSize)}
	}

	return []clause.Expression{Sort(req.SortField, req.SortType), Page(req.PageNo, req.PageSize)}
}

func List[T any, O request.Ordered](db *gorm.DB, req *request.ListReq[O]) ([]T, error) {
	var models []T
	clauses := ListClause(req)
	if len(clauses) > 0 {
		db = db.Clauses(clauses...)
	}
	err := db.Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func ListClause[O request.Ordered](req *request.ListReq[O]) []clause.Expression {
	return (*ListReq[O])(req).Clause()
}

type ListReq[T request.Ordered] request.ListReq[T]

func (req *ListReq[O]) Clause() []clause.Expression {
	pqc := (*PageSortReq)(&req.PageSortReq).Clause()
	rqc := (*RangeReq[O])(req.RangeReq).Clause()
	if pqc != nil && rqc != nil {
		return append((*PageSortReq)(&req.PageSortReq).Clause(), (*RangeReq[O])(req.RangeReq).Clause())
	}
	if rqc == nil {
		return pqc
	}
	if rqc != nil {
		return []clause.Expression{rqc}
	}
	return nil
}
