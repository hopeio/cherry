package iter

import "iter"

func Filter2[K, V any](seq iter.Seq2[K, V], test Predicate2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if test(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

func MapKeys[K, V, R any](seq iter.Seq2[K, V], f Function2[K, V, R]) iter.Seq2[R, V] {
	return func(yield func(R, V) bool) {
		for k, v := range seq {
			if !yield(f(k, v), v) {
				return
			}
		}
	}
}

func MapValues[K, V, R any](seq iter.Seq2[K, V], f Function2[K, V, R]) iter.Seq2[K, R] {
	return func(yield func(K, R) bool) {
		for k, v := range seq {
			if !yield(k, f(k, v)) {
				return
			}
		}
	}
}

func Peek2[K, V any](seq iter.Seq2[K, V], accept Consumer2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			accept(k, v)
			if !yield(k, v) {
				return
			}
		}
	}
}

func Distinct2[K, V any, Cmp comparable](seq iter.Seq2[K, V], f Function2[K, V, Cmp]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var set = make(map[Cmp]struct{})
		for k, v := range seq {
			c := f(k, v)
			_, ok := set[c]
			set[c] = struct{}{}
			if !ok && !yield(k, v) {
				return
			}
		}
	}
}
