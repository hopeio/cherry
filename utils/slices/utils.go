package slices

import (
	"fmt"
	"github.com/hopeio/cherry/utils/constraints"
	"strings"
)

func Contains[S ~[]V, V comparable](slices S, v V) bool {
	for i := 0; i < len(slices); i++ {
		if slices[i] == v {
			return true
		}
	}
	return false
}

func In[S ~[]V, V comparable](v V, slices S) bool {
	for _, x := range slices {
		if x == v {
			return true
		}
	}
	return false
}

func InByKey[S ~[]constraints.CompareKey[K], K comparable](key K, slices S) bool {
	for _, x := range slices {
		if x.CompareKey() == key {
			return true
		}
	}
	return false
}

func Reverse[S ~[]T, T any](slices S) S {
	for i, j := 0, len(slices)-1; i < j; i, j = i+1, j-1 {
		slices[i], slices[j] = slices[j], slices[i]
	}

	return slices
}

func Max[S ~[]T, T constraints.Number](slices S) T {
	if len(slices) == 0 {
		return *new(T)
	}
	max := slices[0]
	if len(slices) == 1 {
		return max
	}
	for i := 1; i < len(slices); i++ {
		if slices[i] > max {
			max = slices[i]
		}
	}

	return max
}

func Min[S ~[]T, T constraints.Number](slices S) T {
	if len(slices) == 0 {
		return *new(T)
	}
	min := slices[0]
	if len(slices) == 1 {
		return min
	}
	for i := 1; i < len(slices); i++ {
		if slices[i] < min {
			min = slices[i]
		}
	}

	return min
}

// 将切片转换为map
func SliceToMap[S ~[]T, T any, K comparable, V any](slices S, getKV func(T) (K, V)) map[K]V {
	m := make(map[K]V)
	for _, s := range slices {
		k, v := getKV(s)
		m[k] = v
	}
	return m
}

// 将切片按照某个key分类
func SliceClassify[S ~[]T, T any, K comparable, V any](slices S, getKV func(T) (K, V)) map[K][]V {
	m := make(map[K][]V)
	for _, s := range slices {
		k, v := getKV(s)
		if ms, ok := m[k]; ok {
			m[k] = append(ms, v)
		} else {
			m[k] = []V{v}
		}

	}
	return m
}

func Swap[S ~[]V, V any](slices S, i, j int) {
	slices[i], slices[j] = slices[j], slices[i]
}

func ForEach[S ~[]V, V any](slices S, handle func(idx int, v V)) {
	for i, t := range slices {
		handle(i, t)
	}
}

func ForEachValue[S ~[]V, V any](slices S, handle func(v V)) {
	for _, v := range slices {
		handle(v)
	}
}

// 遍历切片,参数为下标，利用闭包实现遍历
func ForEachIndex[S ~[]V, V any](slices S, handle func(i int)) {
	for i := range slices {
		handle(i)
	}
}

func JoinByIndex[S ~[]V, V any](slices S, toString func(i int) string, sep string) string {
	switch len(slices) {
	case 0:
		return ""
	case 1:
		return toString(0)
	}
	n := len(sep) * (len(slices) - 1)
	for i := 0; i < len(slices); i++ {
		n += len(toString(i))
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(toString(0))
	for i := range slices[1:] {
		b.WriteString(sep)
		b.WriteString(toString(i))
	}
	return b.String()
}

func JoinByValue[S ~[]V, V any](slices S, toString func(v V) string, sep string) string {
	switch len(slices) {
	case 0:
		return ""
	case 1:
		return toString(slices[0])
	}
	n := len(sep) * (len(slices) - 1)
	for i := 0; i < len(slices); i++ {
		n += len(toString(slices[i]))
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(toString(slices[0]))
	for _, s := range slices[1:] {
		b.WriteString(sep)
		b.WriteString(toString(s))
	}
	return b.String()
}

func Join[S ~[]V, V fmt.Stringer](slices S, sep string) string {
	switch len(slices) {
	case 0:
		return ""
	case 1:
		return slices[0].String()
	}
	n := len(sep) * (len(slices) - 1)
	for i := 0; i < len(slices); i++ {
		n += len(slices[i].String())
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(slices[0].String())
	for _, s := range slices[1:] {
		b.WriteString(sep)
		b.WriteString(s.String())
	}
	return b.String()
}

func ReverseForEach[S ~[]V, V any](slices S, handle func(idx int, v V)) {
	l := len(slices)
	for i := l - 1; i > 0; i-- {
		handle(i, slices[i])
	}
}

func Map[S ~[]T, T any, V any](slices S, fn func(T) V) []V {
	ret := make([]V, 0, len(slices))
	for _, s := range slices {
		ret = append(ret, fn(s))
	}
	return ret
}

func Filter[S ~[]T, T any](fn func(T) bool, src S) S {
	var dst []T
	for _, v := range src {
		if fn(v) {
			dst = append(dst, v)
		}
	}
	return dst
}

func Reduce[S ~[]T, T any](slices S, fn func(T, T) T) T {
	ret := fn(slices[0], slices[1])
	for i := 2; i < len(slices); i++ {
		ret = fn(ret, slices[i])
	}
	return ret
}
