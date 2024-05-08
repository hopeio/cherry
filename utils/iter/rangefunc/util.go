package iter

import (
	"github.com/hopeio/cherry/utils/constraints"
	"github.com/hopeio/cherry/utils/types"
)

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

// TODO: IEDA错误: 这里Seq中类型是types.Pair[K, V]会报错,*types.Pair[K, V]则不会, 2024-05-09 ,IDEA2024.1 Build #IU-241.15989.150, built on April 29, 2024,go插件241.15989.21
func HashMapSeqOf[K comparable, V any](m map[K]V) Seq[*types.Pair[K, V]] {
	return func(yield func(*types.Pair[K, V]) bool) {
		for k, v := range m {
			if !yield(types.PairOf(k, v)) {
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

func ChannelSeqOf[T any](c chan T) Seq[T] {
	return func(yield func(T) bool) {
		for v := range c {
			if !yield(v) {
				return
			}
		}
	}
}

func ChannelSeq2Of[T any](c chan T) Seq2[int, T] {
	return func(yield func(int, T) bool) {
		var count int
		for v := range c {
			if !yield(count, v) {
				return
			}
			count++
		}
	}
}

func StringSeqOf(input string) Seq[rune] {
	return func(yield func(rune) bool) {
		for _, v := range input {
			if !yield(v) {
				return
			}
		}
	}
}

func StringSeq2Of(input string) Seq2[int, rune] {
	return func(yield func(int, rune) bool) {
		for i, v := range input {
			if !yield(i, v) {
				return
			}
		}
	}
}

func RangeSeqOf[T constraints.Number](begin, end, step T) Seq[T] {
	return func(yield func(T) bool) {
		for v := begin; v <= end; v += step {
			if !yield(v) {
				return
			}
		}
	}
}

func RangeSeq2Of[T constraints.Number](begin, end, step T) Seq2[int, T] {
	return func(yield func(int, T) bool) {
		var count int
		for v := begin; v <= end; v += step {
			if !yield(count, v) {
				return
			}
			count++
		}
	}
}
