package types

type FCompare[T any] func(T, T) bool

type FCompareByIndex func(i, j int) bool
