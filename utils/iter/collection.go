package iter

import (
	"github.com/hopeio/cherry/utils/types"
	"golang.org/x/exp/constraints"
)

type Collection[T any] interface {
	Iterable[T]

	Count() int
}

// Ruturns true if the count of collection is 0.
func IsEmpty[T any](c Collection[T]) bool {
	return c.Count() == 0
}

// Ruturns true if the count of collection is 0.
func IsNotEmpty[T any](c Collection[T]) bool {
	return c.Count() != 0
}

// Converts a collection to a Slice.
func ToSlice[T any](c Collection[T]) Slice[T] {
	var arr = make([]T, 0, c.Count())
	ForEach(func(t T) {
		arr = append(arr, t)
	}, c.Iter())
	return arr
}

// Returns true if the target is included in the iterator.
func Contains[T comparable](target T, it Iterator[T]) bool {
	for {
		if v, ok := it.Next(); ok {
			if v == target {
				return true
			}
		} else {
			break
		}
	}
	return false
}

// Returns the sum of all the elements in the iterator.
func Sum[T constraints.Integer | constraints.Float](it Iterator[T]) T {
	return Fold(0, func(a, b T) T {
		return a + b
	}, it)
}

// Returns the product of all the elements in the iterator.
func Product[T constraints.Integer | constraints.Float](it Iterator[T]) T {
	return Fold(1, func(a, b T) T {
		return a * b
	}, it)
}

// Returns the average of all the elements in the iterator.
func Average[T constraints.Integer | constraints.Float](it Iterator[T]) float64 {
	return Fold(0.0, func(result float64, item *types.Pair[int, T]) float64 {
		return result + (float64(item.Second)-result)/float64(item.First+1)
	}, Enumerate(it))
}

// Return the total number of iterators.
func Count[T any](it Iterator[T]) int {
	return Fold(0, func(v int, _ T) int {
		return v + 1
	}, it)
}

// Return the maximum value of all elements of the iterator.
func Max[T constraints.Ordered](it Iterator[T]) (T, bool) {
	return Reduce(func(a T, b T) T {
		if a > b {
			return a
		} else {
			return b
		}
	}, it)
}

// Return the maximum value of all elements of the iterator.
func MaxBy[T any](greater func(T, T) bool, it Iterator[T]) (T, bool) {
	return Reduce(func(a T, b T) T {
		if greater(a, b) {
			return a
		} else {
			return b
		}
	}, it)
}

// Return the minimum value of all elements of the iterator.
func Min[T constraints.Ordered](it Iterator[T]) (T, bool) {
	return Reduce(func(a T, b T) T {
		if a < b {
			return a
		} else {
			return b
		}
	}, it)
}

// Return the minimum value of all elements of the iterator.
func MinBy[T any](less func(T, T) bool, it Iterator[T]) (T, bool) {
	return Reduce(func(a T, b T) T {
		if less(a, b) {
			return a
		} else {
			return b
		}
	}, it)
}

// The action is executed for each element of the iterator, and the argument to the action is the element.
func ForEach[T any](action func(T), it Iterator[T]) {
	for {
		if v, ok := it.Next(); ok {
			action(v)
		} else {
			break
		}
	}
}

// Returns true if all elements in the iterator match the condition.
func AllMatch[T any](predicate func(T) bool, it Iterator[T]) bool {
	for {
		if v, ok := it.Next(); ok {
			if !predicate(v) {
				return false
			}
		} else {
			break
		}
	}
	return true
}

// Returns true if none elements in the iterator match the condition.
func NoneMatch[T any](predicate func(T) bool, it Iterator[T]) bool {
	for {
		if v, ok := it.Next(); ok {
			if predicate(v) {
				return false
			}
		} else {
			break
		}
	}
	return true
}

// Returns true if any elements in the iterator match the condition.
func AnyMatch[T any](predicate func(T) bool, it Iterator[T]) bool {
	for {
		if v, ok := it.Next(); ok {
			if predicate(v) {
				return true
			}
		} else {
			break
		}
	}
	return false
}

// Return the first element.
func First[T any](it Iterator[T]) (T, bool) {
	return it.Next()
}

// Return the last element.
func Last[T any](it Iterator[T]) (T, bool) {
	return Fold(types.None[T](), func(_ *types.Option[T], next T) *types.Option[T] {
		return types.Some(next)
	}, it).Val()
}

// Return the element at index.
func At[T any](index int, it Iterator[T]) (T, bool) {
	var result, ok = it.Next()
	var i = 0
	for i < index && ok {
		result, ok = it.Next()
		i++
	}
	return result, ok
}

// Return the value of the final composite, operates on the iterator from front to back.
func Reduce[T any](operation func(T, T) T, it Iterator[T]) (T, bool) {
	if v, ok := it.Next(); ok {
		return Fold(v, operation, it), true
	}
	return *new(T), false
}

// Return the value of the final composite, operates on the iterator from back to front.
func Fold[T any, R any](initial R, operation func(R, T) R, it Iterator[T]) R {
	var result = initial
	for {
		if v, ok := it.Next(); ok {
			result = operation(result, v)
		} else {
			break
		}
	}
	return result
}

// Splitting an iterator whose elements are pair into two lists.
func Unzip[A any, B any](it Iterator[types.Pair[A, B]]) ([]A, []B) {
	var arrA = make([]A, 0)
	var arrB = make([]B, 0)
	for {
		if v, ok := it.Next(); ok {
			arrA = append(arrA, v.First)
			arrB = append(arrB, v.Second)
		} else {
			break
		}
	}
	return arrA, arrB
}

type Collector[S any, T any, R any] interface {
	Builder() S
	Append(builder S, element T)
	Finish(builder S) R
}

// Collecting via Collector.
func Collect[T any, S any, R any](collector Collector[S, T, R], it Iterator[T]) R {
	var s = collector.Builder()
	for {
		if v, ok := it.Next(); ok {
			collector.Append(s, v)
		} else {
			break
		}
	}
	return collector.Finish(s)
}

// Collect to built-in map.
func CollectToMap[K comparable, V any](it Iterator[types.Pair[K, V]]) map[K]V {
	var r = make(map[K]V, 0)
	for {
		if v, ok := it.Next(); ok {
			r[v.First] = v.Second
		} else {
			break
		}
	}
	return r
}
