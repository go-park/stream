package function

import (
	"reflect"

	"github.com/go-park/stream/internal/helper"
)

type Predicate[T any] func(t T) bool

func (fn Predicate[T]) Negate() Predicate[T] {
	return func(t T) bool { return !fn.Test(t) }
}

func (fn Predicate[T]) Test(t T) bool {
	return fn(t)
}

func (fn Predicate[T]) And(other Predicate[T]) Predicate[T] {
	helper.RequireNonNil(other)
	return func(t T) bool {
		return fn.Test(t) && other.Test(t)
	}
}

func (fn Predicate[T]) Or(other Predicate[T]) Predicate[T] {
	helper.RequireNonNil(other)
	return func(t T) bool {
		return fn.Test(t) || other.Test(t)
	}
}

func (fn Predicate[T]) Not(other Predicate[T]) Predicate[T] {
	helper.RequireNonNil(other)
	return other.Negate()
}

// DeepEqual is a function that takes an object of any type and returns a Predicate[T]
// which checks if the given object is deeply equal to the argument passed in by reflect.DeepEqual.
func (fn Predicate[T]) DeepEqual(obj any) Predicate[T] {
	return func(t T) bool { return reflect.DeepEqual(t, obj) }
}

type BiPredicate[T, U any] func(t T, u U) bool

func (fn BiPredicate[T, U]) Negate() BiPredicate[T, U] {
	return func(t T, u U) bool { return !fn.Test(t, u) }
}

func (fn BiPredicate[T, U]) Test(t T, u U) bool {
	return fn(t, u)
}

func (fn BiPredicate[T, U]) And(other BiPredicate[T, U]) BiPredicate[T, U] {
	helper.RequireNonNil(other)
	return func(t T, u U) bool {
		return fn.Test(t, u) && other.Test(t, u)
	}
}

func (fn BiPredicate[T, U]) Or(other BiPredicate[T, U]) BiPredicate[T, U] {
	helper.RequireNonNil(other)
	return func(t T, u U) bool {
		return fn.Test(t, u) || other.Test(t, u)
	}
}
