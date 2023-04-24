package function

import "github.com/go-park/stream/internal/helper"

type (
	Runner         func()
	Func[T, R any] func(T) R
)

func (fn Runner) Run() {
	helper.RequireCanButNonNil(fn)
	fn()
}

func (fn Runner) AndThen(after Runner) Runner {
	helper.RequireCanButNonNil(after)
	return func() { fn.Run(); after.Run() }
}

func (fn Func[T, R]) Apply(t T) R {
	return fn(t)
}

func (Func[T, R]) Identify() Func[T, T] {
	return func(t T) T { return t }
}

func Compose[T, R, V any](source Func[T, R], before Func[V, T]) Func[V, R] {
	helper.RequireCanButNonNil(source)
	helper.RequireCanButNonNil(before)
	return func(t V) R { return source.Apply(before.Apply(t)) }
}

func AndThen[T, R, V any](source Func[T, R], after Func[R, V]) Func[T, V] {
	helper.RequireCanButNonNil(source)
	helper.RequireCanButNonNil(after)
	return func(t T) V { return after.Apply(source.Apply(t)) }
}

type BiFunc[T, U, R any] func(T, U) R

func (fn BiFunc[T, U, R]) Apply(t T, u U) R {
	helper.RequireCanButNonNil(fn)
	return fn(t, u)
}

func AndBiThen[T, U, R, V any](source BiFunc[T, U, R], after Func[R, V]) BiFunc[T, U, V] {
	helper.RequireCanButNonNil(source)
	helper.RequireCanButNonNil(after)
	return func(t T, u U) V { return after.Apply(source.Apply(t, u)) }
}
