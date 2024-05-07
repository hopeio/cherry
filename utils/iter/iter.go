package iter

import (
	"golang.org/x/exp/constraints"
	"unicode/utf8"
)

type Iterator[T any] interface {
	Next() (v T, ok bool)
}

type Iterable[T any] interface {
	Iter() Iterator[T]
}

type stringIter struct {
	str string
}

func (it *stringIter) Next() (rune, bool) {
	if len(it.str) == 0 {
		return 0, false
	}
	value, width := utf8.DecodeRuneInString(it.str)
	it.str = it.str[width:]
	return value, true
}

// String returns an Iterator yielding runes from the supplied string.
func StringIterOf(input string) Iterator[rune] {
	return &stringIter{
		str: input,
	}
}

type rangeIter[T constraints.Integer] struct {
	begin, end, step, idx T
}

func (it *rangeIter[T]) Next() (T, bool) {
	v := it.begin + it.step*it.idx
	if it.step > 0 {
		if v >= it.end {
			return *new(T), false
		}
	} else {
		if v <= it.end {
			return *new(T), false
		}
	}
	it.idx++
	return v, true
}

// RangeIterOf returns an Iterator over a range of integers.
func RangeIterOf[T constraints.Integer](begin, end, step T) Iterator[T] {
	return &rangeIter[T]{
		begin: begin,
		end:   end,
		step:  step,
		idx:   0,
	}
}

type sliceIter[T any] struct {
	index  int
	source []T
}

func (a *sliceIter[T]) Next() (T, bool) {
	if a.index < len(a.source)-1 {
		a.index++
		return a.source[a.index], true
	}
	return *new(T), false
}

func SliceIterOf[S ~[]T, T any](source S) Iterator[T] {
	return &sliceIter[T]{0, source}
}

type IterStream[T any] struct {
	iter Iterator[T]
}

func (it *IterStream[T]) Map(transform func(T) T) *IterStream[T] {
	it.iter = Map(it.iter, transform)
	return it
}

func (it *IterStream[T]) Filter(f func(T) bool) *IterStream[T] {
	it.iter = Filter(it.iter, f)
	return it
}

func (it *IterStream[T]) Count() uint64 {
	return Count(it.iter)
}

func (it *IterStream[T]) ForEach(f func(T)) {
	ForEach(it.iter, f)
}

func (it *IterStream[T]) All(f func(T) bool) bool {
	return AllMatch(it.iter, f)
}

func (it *IterStream[T]) Any(f func(T) bool) bool {
	return AnyMatch(it.iter, f)
}

func (it *IterStream[T]) None(f func(T) bool) bool {
	return NoneMatch(it.iter, f)
}

func (it *IterStream[T]) Skip(n int) *IterStream[T] {
	it.iter = Skip(it.iter, n)
	return it
}

func (it *IterStream[T]) Limit(n int) *IterStream[T] {
	it.iter = Limit(it.iter, n)
	return it
}

func (it *IterStream[T]) Reduce(operation func(T, T) T) (T, bool) {
	return Reduce(it.iter, operation)
}

func (it *IterStream[T]) Fold(initVal T, operation func(T, T) T) T {
	return Fold(it.iter, initVal, operation)
}
