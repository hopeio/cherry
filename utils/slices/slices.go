package slices

type Slice[T any] []T

func (slice Slice[T]) Len() int { return len(slice) }

func (slice Slice[T]) ForEach(fn func(T)) {
	for _, t := range slice {
		fn(t)
	}
}

func (slice Slice[T]) Filter(fn func(T) bool) []T {
	var newSlices []T
	for _, t := range slice {
		if fn(t) {
			newSlices = append(newSlices, t)
		}
	}
	return newSlices
}

func (slice Slice[T]) Every(fn func(T) bool) bool {
	for _, t := range slice {
		if !fn(t) {
			return false
		}
	}
	return true
}

func (slice Slice[T]) Some(fn func(T) bool) bool {
	for _, t := range slice {
		if fn(t) {
			return true
		}
	}
	return false
}

func (slice Slice[T]) Zip(s []T) [][2]T {
	var newSlices [][2]T
	for i := range slice {
		newSlices = append(newSlices, [2]T{slice[i], s[i]})
	}
	return newSlices
}

func (slice Slice[T]) Reduce(fn func(T, T) T) T {
	ret := fn(slice[0], slice[1])
	for i := 2; i < len(slice); i++ {
		ret = fn(ret, slice[i])
	}
	return ret
}

type MapSlices[T, V any] Slice[T]

func (slices MapSlices[T, V]) Map(fn func(T) V) []V {
	ret := make([]V, 0, len(slices))
	for _, s := range slices {
		ret = append(ret, fn(s))
	}
	return ret
}

// 学学kotlin的定义
type Array[S, T any] []S

//type Function[T any] func[T]()

func (a Array[S, T]) Map(fn func(S) T) []T {
	ret := make([]T, 0, len(a))
	for _, s := range a {
		ret = append(ret, fn(s))
	}
	return ret
}

type ComparableSlices[T comparable] []T

// 去重
func (slices ComparableSlices[T]) Deduplicate() ComparableSlices[T] {
	if len(slices) < SmallArrayLen {
		newslices := make(ComparableSlices[T], 0, 2)
		for i := 0; i < len(slices); i++ {
			if !In(slices[i], newslices) {
				newslices = append(newslices, slices[i])
			}
		}
		return newslices
	}
	set := make(map[T]struct{})
	for i := 0; i < len(slices); i++ {
		set[slices[i]] = struct{}{}
	}
	newslices := make(ComparableSlices[T], 0, len(slices))
	for k := range set {
		newslices = append(newslices, k)
	}
	return newslices
}

// 去重
func Deduplicate[S ~[]T, T comparable](slices S) S {
	if len(slices) < SmallArrayLen {
		newslices := make(S, 0, 2)
		for i := 0; i < len(slices); i++ {
			if !In(slices[i], newslices) {
				newslices = append(newslices, slices[i])
			}
		}
		return newslices
	}
	set := make(map[T]struct{})
	for i := 0; i < len(slices); i++ {
		set[slices[i]] = struct{}{}
	}
	newslices := make(S, 0, len(slices))
	for k := range set {
		newslices = append(newslices, k)
	}
	return newslices
}
