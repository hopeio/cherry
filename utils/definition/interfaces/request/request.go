package request

import "github.com/hopeio/cherry/utils/definition/types/request"

type PageSortReq interface {
	PageReq
	SortReq
}

type PageReq interface {
	PageNo() int
	PageSize() int
}

type SortReq interface {
	Column() string
	Type() request.SortType
}
