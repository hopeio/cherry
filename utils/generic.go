package utils

func ZeroValue[T any]() T {
	return *new(T)
}

// can compile,but will panic
func nilValue[T any]() T {
	return *(*T)(nil)
}

func ZeroValue2[T any]() T {
	var zero T
	return zero
}
