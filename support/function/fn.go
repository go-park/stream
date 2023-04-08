package function

import "github.com/go-park/stream/internal/helper"

type Fn[T, R any] func(T) R

func (fn Fn[T, R]) Apply(t T) R {
	return fn(t)
}

func (Fn[T, R]) Identify() Fn[T, T] {
	return func(t T) T { return t }
}

func Compose[T, R, V any](source Fn[T, R], before Fn[V, T]) Fn[V, R] {
	helper.RequireNonNil(source)
	helper.RequireNonNil(before)
	return func(t V) R { return source.Apply(before.Apply(t)) }
}

func AndThen[T, R, V any](source Fn[T, R], after Fn[R, V]) Fn[T, V] {
	helper.RequireNonNil(source)
	helper.RequireNonNil(after)
	return func(t T) V { return after.Apply(source.Apply(t)) }
}

type BiFn[T, U, R any] func(T, U) R

func (fn BiFn[T, U, R]) Apply(t T, u U) R {
	return fn(t, u)
}

func AndBiThen[T, U, R, V any](source BiFn[T, U, R], after Fn[R, V]) BiFn[T, U, V] {
	helper.RequireNonNil(source)
	helper.RequireNonNil(after)
	return func(t T, u U) V { return after.Apply(source.Apply(t, u)) }
}
