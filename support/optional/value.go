package optional

import (
	"github.com/go-park/stream/internal/helper"
	"github.com/go-park/stream/support/function"
)

type Value[T any] struct {
	v             T
	canNil, isNil bool
	nonEmpty      bool
}

func ValOf[T any](v T) Value[T] {
	// canNil, isNil := helper.IsNil(v)
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
		helper.RequireCanButNonNil(fn)
		fn.Accept(v.v)
	}
	return
}

func (v Value[T]) IfNotEmptyOrElse(fn function.Consumer[T], runner function.Runner) {
	if !v.IsEmpty() {
		helper.RequireCanButNonNil(fn)
		fn.Accept(v.v)
	} else {
		helper.RequireCanButNonNil(runner)
		runner.Run()
	}
	return
}

func (v Value[T]) IsNil() bool {
	return v.canNil && v.isNil
}

func (v Value[T]) Get() T {
	return v.v
}
