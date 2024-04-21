package utils

func ZeroValue[T any]() T {
	return *new(T)
}

// will panic
func nilValue[T any]() T {
	return *(*T)(nil)
}
