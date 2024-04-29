package maps

func Map[M ~map[K]V, K comparable, V, T any](m M, subValue func(K, V) T) []T {
	r := make([]T, 0, len(m))
	for k, v := range m {
		r = append(r, subValue(k, v))
	}
	return r
}

func MapKey[M ~map[K]V, K comparable, V, T any](m M, subValue func(K) T) []T {
	r := make([]T, 0, len(m))
	for k, _ := range m {
		r = append(r, subValue(k))
	}
	return r
}

func MapValue[M ~map[K]V, K comparable, V, T any](m M, subValue func(V) T) []T {
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

func ForEachValue[M ~map[K]V, K comparable, V any](m M, handle func(v V)) {
	for _, v := range m {
		handle(v)
	}
}

func ForEachKey[M ~map[K]V, K comparable, V any](m M, handle func(v K)) {
	for k, _ := range m {
		handle(k)
	}
}

func Keys[M ~map[K]V, K comparable, V any](maps ...M) []K {
	r := make([]K, 0, len(maps))
	for _, m := range maps {
		for k := range m {
			r = append(r, k)
		}
	}
	return r
}

func Values[M ~map[K]V, K comparable, V any](maps ...M) []V {
	r := make([]V, 0, len(maps))
	for _, m := range maps {
		for _, v := range m {
			r = append(r, v)
		}
	}
	return r
}

func Merge[M ~map[K]V, K comparable, V any](maps ...M) M {
	r := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			r[k] = v
		}
	}
	return r
}
