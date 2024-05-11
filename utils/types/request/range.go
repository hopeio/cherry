package request

import (
	"github.com/hopeio/cherry/utils/constraints"
	"time"
)

type SortType int

const (
	_ SortType = iota
	SortTypeASC
	SortTypeDESC
)

type PageSortReq struct {
	PageReq
	*SortReq
}

type PageReq struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}

type SortReq struct {
	SortField string   `json:"sortField,omitempty"`
	SortType  SortType `json:"sortType,omitempty"`
}

func (receiver *SortReq) Column() string {
	return receiver.SortField
}

func (receiver *SortReq) Type() SortType {
	return receiver.SortType
}

type DateRangeReq[T ~string | time.Time] RangeReq[T]

type RangeReq[T constraints.Range] struct {
	RangeField string `json:"rangeField,omitempty"`
	RangeBegin T      `json:"rangeBegin,omitempty"`
	RangeEnd   T      `json:"rangeEnd,omitempty"`
	Include    bool   `json:"include,omitempty"`
}
