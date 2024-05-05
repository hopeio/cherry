//go:build goexperiment.rangefunc

package iter

import (
	"github.com/hopeio/cherry/utils/types"
	"iter"
)

// Supplier 产生一个元素
type Supplier[T any] func() T

// Function 将一个类型转为另一个类型
type Function[T, R any] func(T) R

// Predicate 断言是否满足指定条件
type Predicate[T any] Function[T, bool]

// UnaryOperator 对输入进行一元运算返回相同类型的结果
type UnaryOperator[T any] Function[T, T]

// BiFunction 将两个类型转为第三个类型
type BiFunction[T, R, U any] func(T, R) U

// BinaryOperator 输入两个相同类型的参数，对其做二元运算，返回相同类型的结果
type BinaryOperator[T any] BiFunction[T, T, T]

// Comparator 比较两个元素.
// 第一个元素大于第二个元素时，返回正数;
// 第一个元素小于第二个元素时，返回负数;
// 否则返回 0.
type Comparator[T any] BiFunction[T, T, int]

// Consumer 消费一个元素
type Consumer[T any] func(T)

type Stream[T any] interface {
	Seq() iter.Seq[T]

	Filter(Predicate[T]) Stream[T]
	Map(Function[T, T]) Stream[T]
	FlatMap(Function[T, iter.Seq[T]]) Stream[T]
	Peek(Consumer[T]) Stream[T]

	Distinct(Function[T, int]) Stream[T]
	Sorted(Comparator[T]) Stream[T]
	Limit(int64) Stream[T]
	Skip(int64) Stream[T]

	ForEach(Consumer[T])
	Collect() []T
	AllMatch(Predicate[T]) bool
	NoneMatch(Predicate[T]) bool
	AnyMatch(Predicate[T]) bool
	Reduce(acc BinaryOperator[T]) *types.Option[T]
	ReduceFrom(initVal T, acc BinaryOperator[T]) T
	FindFirst() *types.Option[T]
	Count() int64
}

type Seq[T any] iter.Seq[T]

func (it Seq[T]) Seq() iter.Seq[T] {
	return iter.Seq[T](it)
}

func (it Seq[T]) Filter(test Predicate[T]) Stream[T] {
	return Seq[T](Filter(iter.Seq[T](it), test))
}

func (it Seq[T]) Map(f Function[T, T]) Stream[T] {
	return Maps(it, f)
}

func (it Seq[T]) FlatMap(f Function[T, iter.Seq[T]]) Stream[T] {
	return FlatMaps(it, f)
}

func (it Seq[T]) Peek(accept Consumer[T]) Stream[T] {
	return Seq[T](Peek(iter.Seq[T](it), accept))
}

func (it Seq[T]) Distinct(f Function[T, int]) Stream[T] {
	return Seq[T](Distinct(iter.Seq[T](it), f))
}

func (it Seq[T]) Sorted(cmp Comparator[T]) Stream[T] {
	return Seq[T](Sorted(iter.Seq[T](it), cmp))
}

func (it Seq[T]) Limit(limit int64) Stream[T] {
	return Seq[T](Limit(iter.Seq[T](it), limit))
}

func (it Seq[T]) Skip(skip int64) Stream[T] {
	return Seq[T](Skip(iter.Seq[T](it), skip))
}

func (it Seq[T]) ForEach(accept Consumer[T]) {
	ForEach(iter.Seq[T](it), accept)
}

func (it Seq[T]) Collect() []T {
	return Collect(iter.Seq[T](it))
}

func (it Seq[T]) AllMatch(test Predicate[T]) bool {
	return AllMatch(iter.Seq[T](it), test)
}

func (it Seq[T]) NoneMatch(test Predicate[T]) bool {
	return NoneMatch(iter.Seq[T](it), test)
}

func (it Seq[T]) AnyMatch(test Predicate[T]) bool {
	return AnyMatch(iter.Seq[T](it), test)
}

func (it Seq[T]) Reduce(acc BinaryOperator[T]) *types.Option[T] {
	return Reduce(iter.Seq[T](it), acc)
}

func (it Seq[T]) ReduceFrom(initVal T, acc BinaryOperator[T]) T {
	return ReduceFrom(iter.Seq[T](it), initVal, acc)
}

func (it Seq[T]) FindFirst() *types.Option[T] {
	return FindFirst(iter.Seq[T](it))
}

func (it Seq[T]) Count() int64 {
	return Count(iter.Seq[T](it))
}
