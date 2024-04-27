package request

type IPageSortReq interface {
	IPageReq
	ISortReq
}

type IPageReq interface {
	PageNo() int
	PageSize() int
}

type ISortReq interface {
	Column() string
	Type() SortType
}
