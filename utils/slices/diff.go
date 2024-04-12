package slices

import (
	_interface "github.com/hopeio/cherry/utils/definition/constraints"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

// 没有泛型，范例，实际需根据不同类型各写一遍,用CmpKey，基本类型又用不了，go需要能给基本类型实现方法不能给外部类型实现方法
// 1.20以后字段均是comparable的结构体也是comparable的
// 判断是否有重合元素
func HasCoincide[S ~[]T, T comparable](s1, s2 S) bool {
	for i := range s1 {
		for j := range s2 {
			if s1[i] == s2[j] {
				return true
			}
		}
	}
	return false
}

func HasCoincideByKey[S ~[]_interface.CmpKey[T], T comparable](s1, s2 S) bool {
	for i := range s1 {
		for j := range s2 {
			if s1[i].CmpKey() == s2[j].CmpKey() {
				return true
			}
		}
	}
	return false
}

func RemoveDuplicates[S ~[]T, T comparable](s S) S {
	var m = make(map[T]struct{})
	for _, i := range s {
		m[i] = struct{}{}
	}
	s = s[:0]
	for k, _ := range m {
		s = append(s, k)
	}
	return s
}

func RemoveDuplicatesByKey[S ~[]_interface.CmpKey[T], T comparable](s S) S {
	var m = make(map[T]_interface.CmpKey[T])
	for _, i := range s {
		m[i.CmpKey()] = i
	}
	s = s[:0]
	for _, i := range m {
		s = append(s, i)
	}
	return s
}

// 取交集
func Intersection[S ~[]T, T comparable](a S, b S) S {
	if len(a) < SmallArrayLen && len(b) < SmallArrayLen {
		if len(a) > len(b) {
			return smallArrayIntersection(a, b)
		}
		return smallArrayIntersection(b, a)
	}
	return intersection(a, b)
}

func smallArrayIntersection[S ~[]T, T comparable](a S, b S) S {
	var ret S
	for _, x := range a {
		if In(x, b) {
			ret = append(ret, x)
		}
	}
	return ret
}

func intersection[S ~[]T, T comparable](a S, b S) S {
	return maps.Keys(IntersectionMap(a, b))
}

func IntersectionMap[S ~[]T, T comparable](a S, b S) map[T]struct{} {
	intersectionMap := make(map[T]struct{})
Loop:
	for _, i := range a {
		for _, j := range b {
			if i == j {
				intersectionMap[i] = struct{}{}
				continue Loop
			}
		}
	}
	return intersectionMap
}

func IntersectionByKey[S ~[]_interface.CmpKey[T], T comparable](a S, b S) S {
	if len(a) < SmallArrayLen && len(b) < SmallArrayLen {
		if len(a) > len(b) {
			return intersectionByKey(a, b)
		}
		return intersectionByKey(b, a)
	}
	panic("TODO:大数组利用map取并集")
}

func intersectionByKey[S ~[]_interface.CmpKey[T], T comparable](a S, b S) S {
	var ret S
	for _, x := range a {
		if InByKey(x.CmpKey(), b) {
			ret = append(ret, x)
		}
	}
	return ret
}

// 有序数组取交集
func OrderedArrayIntersection[S ~[]T, T constraints.Ordered](a S, b S) S {
	var ret S
	if len(a) == 0 || len(b) == 0 {
		return nil
	}
	var idx int
	for _, x := range a {
		if x > b[len(b)-1] {
			return ret
		}
		for j := idx; idx < len(b)-1; j++ {
			if a[len(a)-1] < b[idx] {
				return ret
			}
			if x == b[idx] {
				ret = append(ret, x)
				idx = j
			}
		}
	}
	return ret
}

// 取并集
func Union[S ~[]T, T comparable](a S, b S) S {
	var set = make(map[T]struct{}, len(a)+len(b))
	for _, x := range a {
		set[x] = struct{}{}
	}
	for _, x := range b {
		set[x] = struct{}{}
	}
	return maps.Keys(set)
}

func UnionByKey[S ~[]_interface.CmpKey[T], T comparable](a S, b S) S {
	var m = make(map[T]_interface.CmpKey[T], len(a)+len(b))
	for _, x := range a {
		m[x.CmpKey()] = x
	}
	for _, x := range b {
		m[x.CmpKey()] = x
	}
	return maps.Values(m)
}

// 取差集,返回为A-B
func DifferenceSet[S ~[]T, T comparable](a S, b S) S {
	if len(a) == 0 {
		return S{}
	}
	if len(b) == 0 {
		return a
	}
	if len(a) < SmallArrayLen && len(b) < SmallArrayLen {
		return smallArrayDifferenceSet(a, b)
	}
	return differenceSet(a, b)
}

func smallArrayDifferenceSet[S ~[]T, T comparable](a S, b S) S {
	var diff S
	for _, x := range a {
		if !In(x, b) {
			diff = append(diff, x)
		}
	}
	return diff
}

func differenceSet[S ~[]T, T comparable](a S, b S) S {
	var diff S
	if len(b)/len(a) >= 2 {
		aMap := make(map[T]bool)
		for _, x := range a {
			aMap[x] = false
		}
		for _, x := range b {
			if _, ok := aMap[x]; ok {
				aMap[x] = true
			}
		}
		for v, exits := range aMap {
			if !exits {
				diff = append(diff, v)
			}
		}
	} else {
		bMap := make(map[T]struct{})
		for _, x := range b {
			bMap[x] = struct{}{}
		}
		for _, x := range a {
			if _, ok := bMap[x]; !ok {
				diff = append(diff, x)
			}
		}
	}
	return diff
}

// 指定key取差集,返回为A-B
func DifferenceSetByKey[S ~[]_interface.CmpKey[T], T comparable](a S, b S) S {
	if len(a) == 0 {
		return S{}
	}
	if len(b) == 0 {
		return a
	}
	if len(a) < SmallArrayLen && len(b) < SmallArrayLen {
		return smallArrayDifferenceSetByKey(a, b)
	}
	aMap := make(map[T]struct{})
	for _, x := range a {
		aMap[x.CmpKey()] = struct{}{}
	}
	var diff S
	for _, x := range b {
		if _, ok := aMap[x.CmpKey()]; !ok {
			diff = append(diff, x)
		}
	}
	return diff
}

func smallArrayDifferenceSetByKey[S ~[]_interface.CmpKey[T], T comparable](a S, b S) S {
	var diff S
	for _, x := range a {
		if !InByKey(x.CmpKey(), b) {
			diff = append(diff, x)
		}
	}
	return diff
}

func differenceSetByKey[S ~[]_interface.CmpKey[T], T comparable](a S, b S) S {
	var diff S
	if len(b)/len(a) >= 2 {
		aMap := make(map[T]bool)
		for _, x := range a {
			aMap[x.CmpKey()] = false
		}
		for _, x := range b {
			if _, ok := aMap[x.CmpKey()]; ok {
				aMap[x.CmpKey()] = true
			}
		}
		for _, x := range a {
			if _, exits := aMap[x.CmpKey()]; !exits {
				diff = append(diff, x)
			}
		}

	} else {
		bMap := make(map[T]struct{})
		for _, x := range b {
			bMap[x.CmpKey()] = struct{}{}
		}
		for _, x := range a {
			if _, ok := bMap[x.CmpKey()]; !ok {
				diff = append(diff, x)
			}
		}
	}
	return diff
}

// 两个数组各自相对的差集,返回为A-B,B-A
func Difference[S ~[]T, T comparable](a, b S) (S, S) {
	var diff1, diff2 S
	intersectionMap := make(map[T]struct{})
Loop:
	for _, i := range a {
		for _, j := range b {
			if i == j {
				intersectionMap[i] = struct{}{}
				continue Loop
			}
		}
		diff1 = append(diff1, i)
	}
	for _, i := range b {
		if _, ok := intersectionMap[i]; !ok {
			diff2 = append(diff2, i)
		}
	}
	return diff1, diff2
}

// 取差集，通过循环比较key
func DifferenceByKey[S ~[]_interface.CmpKey[T], T comparable](a, b S) (S, S) {
	var diff1, diff2 S
	intersectionMap := make(map[T]struct{})
Loop:
	for _, i := range a {
		for _, j := range b {
			if i.CmpKey() == j.CmpKey() {
				intersectionMap[i.CmpKey()] = struct{}{}
				continue Loop
			}
		}
		diff1 = append(diff1, i)
	}
	for _, i := range b {
		if _, ok := intersectionMap[i.CmpKey()]; !ok {
			diff2 = append(diff2, i)
		}
	}
	return diff1, diff2
}

// 并集和差集，返回AUB,A-B,B-A
func IntersectionAndDifference[S ~[]T, T comparable](a, b S) (S, S, S) {
	var diff1, diff2 S
	intersectionMap := make(map[T]struct{})
Loop:
	for _, i := range a {
		for _, j := range b {
			if i == j {
				intersectionMap[i] = struct{}{}
				continue Loop
			}
		}
		diff1 = append(diff1, i)
	}
	for _, i := range b {
		if _, ok := intersectionMap[i]; !ok {
			diff2 = append(diff2, i)
		}
	}

	return maps.Keys(intersectionMap), diff1, diff2
}
