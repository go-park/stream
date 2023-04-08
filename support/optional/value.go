package optional

import (
	"github.com/go-park/stream/internal/helper"
	"github.com/go-park/stream/support/function"
)

type Value[T any] struct {
	v        T
	nonEmpty bool
}

func ValOf[T any](v T) Value[T] {
	helper.RequireNonNil(v)
	return Value[T]{v: v, nonEmpty: true}
}

func EmptyVal[T any]() Value[T] {
	return Value[T]{}
}

func (v Value[T]) IsEmpty() bool {
	return !v.nonEmpty
}

func (v Value[T]) IfNotEmpty(fn function.Consumer[T]) {
	if !v.IsEmpty() {
		fn(v.v)
	}
	return
}

func (v Value[T]) IsNil() bool {
	canNil, isNil := helper.IsNil(v.v)
	return canNil && isNil
}

func (v Value[T]) Get() T {
	return v.v
}
