package types

import (
	"encoding/json"
)

type Option[T any] struct {
	value T
	ok    bool
}

func Some[T any](v T) *Option[T] {
	return &Option[T]{value: v, ok: true}
}

func None[T any]() *Option[T] {
	return &Option[T]{ok: false}
}
func Nil[T any]() *Option[T] {
	return &Option[T]{ok: false}
}

func (opt *Option[T]) Val() (T, bool) {
	return opt.value, opt.ok
}

func (opt *Option[T]) Get() (T, bool) {
	return opt.value, opt.ok
}

func (opt *Option[T]) IsNone() bool {
	return !opt.ok
}

func (opt *Option[T]) IsSome() bool {
	return opt.ok
}

func (opt *Option[T]) Unwrap() T {
	if opt.IsNone() {
		panic("Attempted to unwrap an empty Option.")
	}
	return opt.value
}

func (opt *Option[T]) UnwrapOr(def T) T {
	if opt.IsSome() {
		return opt.Unwrap()
	}
	return def
}

func (opt *Option[T]) UnwrapOrElse(fn func() T) T {
	if opt.IsSome() {
		return opt.Unwrap()
	}
	return fn()
}

func MapOption[T any, R any](opt *Option[T], fn func(T) R) *Option[R] {
	if !opt.IsSome() {
		return None[R]()
	}
	return Some(fn(opt.Unwrap()))
}

func (a *Option[T]) IfSome(action func(value T)) {
	if a.ok {
		action(a.value)
	}
}

func (a *Option[T]) IfNone(action func()) {
	if !a.ok {
		action()
	}
}

func (a *Option[T]) Next() *Option[T] {
	return a
}

func (a *Option[T]) MarshalJSON() ([]byte, error) {
	if a.ok {
		return json.Marshal(a.value)
	}
	return []byte("null"), nil
}

func (a *Option[T]) UnmarshalJSON(data []byte) error {
	if len(data) < 5 && string(data) == "null" {
		a.ok = false
		return nil
	}
	a.ok = true
	return json.Unmarshal(data, &a.value)
}
