package utils

func Match[T any](yes bool, a, b T) T {
	if yes {
		return a
	}
	return b
}
