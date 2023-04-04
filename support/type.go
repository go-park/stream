package support

import (
	"github.com/go-park/stream/function"
	"github.com/go-park/stream/internal/helper"
)

type Value[V any] struct {
	v        V
	nonEmpty bool
}

func ValOf[V any](v V) Value[V] {
	helper.RequireNonNil(v)
	return Value[V]{v: v, nonEmpty: true}
}

func EmptyVal[V any]() Value[V] {
	return Value[V]{}
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
	canNil, isNil := helper.IsNil(v)
	return canNil && isNil
}

func (v Value[T]) Get() T {
	return v.v
}
