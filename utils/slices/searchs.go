package slices

import (
	constraints2 "github.com/hopeio/cherry/utils/constraints"
	"golang.org/x/exp/constraints"
)

// BinarySearch 二分查找
func BinarySearch[V constraints.Ordered](arr []constraints2.OrderKey[V], x constraints2.OrderKey[V]) int {
	l, r := 0, len(arr)-1
	for l <= r {
		mid := (l + r) / 2
		if arr[mid].OrderKey() == x.OrderKey() {
			return mid
		} else if x.OrderKey() > arr[mid].OrderKey() {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return -1
}
