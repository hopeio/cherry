package iter

import "github.com/hopeio/cherry/utils/types"

func SliceSeqOf[T any](input []T) Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range input {
			if !yield(v) {
				return
			}
		}
	}
}

func SliceSeq2Of[T any](input []T) Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range input {
			if !yield(i, v) {
				return
			}
		}
	}
}

func HashMapSeqOf[K comparable, V any](m map[K]V) Seq[types.Pair[K, V]] {
	return func(yield func(types.Pair[K, V]) bool) {
		for k, v := range m {
			if !yield(types.Pair[K, V]{First: k, Second: v}) {
				return
			}
		}
	}
}

func HashMapSeq2Of[K comparable, V any](m map[K]V) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}
