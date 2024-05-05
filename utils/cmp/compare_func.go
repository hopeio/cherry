package cmp

type SortFunc[T any] func(T, T) bool

type CompareFunc[T any] func(T, T) int
