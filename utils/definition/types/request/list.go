package request

import (
	"golang.org/x/exp/constraints"
	"time"
)

type ListReq[T Ordered] struct {
	PageSortReq
	*RangeReq[T]
}

func NewListReq[T Ordered](pageNo, pageSize int) *ListReq[T] {
	return &ListReq[T]{
		PageSortReq: PageSortReq{
			PageReq: PageReq{
				PageNo:   pageNo,
				PageSize: pageSize,
			},
		},
	}
}

func (req *ListReq[T]) WithSort(field string, typ SortType) *ListReq[T] {
	req.SortReq = &SortReq{
		SortField: field,
		SortType:  typ,
	}
	return req
}

func (req *ListReq[T]) WithRange(field string, start, end T, include bool) *ListReq[T] {
	req.RangeReq = &RangeReq[T]{
		RangeField: field,
		RangeStart: start,
		RangeEnd:   end,
		Include:    include,
	}
	return req
}

type Ordered interface {
	constraints.Ordered | time.Time
}
