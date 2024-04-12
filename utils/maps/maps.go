package maps

func SubValues[M ~map[K]V, K comparable, V, T any](m M, subValue func(V) T) []T {
	r := make([]T, 0, len(m))
	for _, v := range m {
		r = append(r, subValue(v))
	}
	return r
}

func ForEach[M ~map[K]V, K comparable, V any](m M, handle func(K, V)) {
	for k, v := range m {
		handle(k, v)
	}
}

func ForEachValue[S ~[]V, V any](slices S, handle func(v V)) {
	for _, v := range slices {
		handle(v)
	}
}
