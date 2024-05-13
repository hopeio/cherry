package slices

import (
	constraints2 "github.com/hopeio/cherry/utils/cmp"
	"golang.org/x/exp/constraints"
)

// BinarySearch 二分查找
func BinarySearch[V constraints.Ordered](arr []constraints2.SortKey[V], x constraints2.SortKey[V]) int {
	l, r := 0, len(arr)-1
	for l <= r {
		mid := (l + r) / 2
		if arr[mid].SortKey() == x.SortKey() {
			return mid
		} else if x.SortKey() > arr[mid].SortKey() {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return -1
}
