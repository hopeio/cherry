package iter

import (
	"iter"
)

// Supplier 产生一个元素
type Supplier2[K, V any] func() (K, V)

// Function 将一个类型转为另一个类型
type Function2[K, V, R any] func(K, V) R

// Predicate 断言是否满足指定条件
type Predicate2[K, V any] func(K, V) bool

// Comparator 比较两个元素.
// 第一个元素大于第二个元素时，返回正数;
// 第一个元素小于第二个元素时，返回负数;
// 否则返回 0.
type Comparator2[K, V any] func(K, V, K, V) bool

// Consumer 消费一个元素
type Consumer2[K, V any] func(K, V)

type Stream2[K, V any] interface {
	Seq2() iter.Seq2[K, V]

	Filter(Predicate2[K, V]) Stream2[K, V]
	Peek(Consumer2[K, V]) Stream2[K, V]

	Distinct(Function2[K, V, int]) Stream2[K, V]
	SortedByKeys(Comparator[K]) Stream2[K, V]
	SortedByValues(Comparator[V]) Stream2[K, V]
	Limit(int64) Stream2[K, V]
	Skip(int64) Stream2[K, V]

	ForEach(Consumer2[K, V])
	CollectKeys() []K
	CollectValues() []V
	AllMatch(Predicate2[K, V]) bool
	NoneMatch(Predicate2[K, V]) bool
	AnyMatch(Predicate2[K, V]) bool

	First() (K, V)
	Count() int64
}

type Seq2[K, V any] iter.Seq2[K, V]
