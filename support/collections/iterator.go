package collections

import (
	"github.com/go-park/stream/internal/helper"
	"github.com/go-park/stream/support/function"
)

type Iterator[T any] interface {
	HasNext() bool
	Next() T
	ForEachRemaining(function.Consumer[T])
}

type iterator[T any] func() (func() bool, function.Supplier[T])

func (iter iterator[T]) HasNext() bool {
	hasNext, _ := iter()
	return hasNext()
}

func (iter iterator[T]) Next() T {
	_, next := iter()
	return next()
}

func (iter iterator[T]) ForEachRemaining(fn function.Consumer[T]) {
	helper.RequireCanButNonNil(fn)
	hasNext, next := iter()
	for hasNext() {
		fn(next())
	}
}

type iterableSlice[T any] struct {
	slice []T
}

func (s *iterableSlice[T]) HasNext() bool {
	return len(s.slice) > 0
}

func (s *iterableSlice[T]) Next() T {
	var v T
	if len(s.slice) > 1 {
		v, s.slice = s.slice[0], s.slice[1:]
		return v
	}
	if len(s.slice) == 1 {
		v, s.slice = s.slice[0], nil
		return v
	}
	return v
}

func (s *iterableSlice[T]) ForEachRemaining(fn function.Consumer[T]) {
	for s.HasNext() {
		fn(s.Next())
	}
}

func IterableSlice[T any](list ...T) Iterator[T] {
	return &iterableSlice[T]{list}
}

func IterableChan[T any](ch chan T) Iterator[T] {
	// escape to heap
	ready := false
	var value T
	return iterator[T](
		func() (func() bool, function.Supplier[T]) {
			hasNext := func() bool {
				if ready {
					return true
				}
				v, ok := <-ch
				if ok {
					ready = true
					value = v
				}
				return false
			}
			next := func() T {
				defer func() { ready = false }()
				var v T
				if hasNext() {
					v = value
				}
				return v
			}
			return hasNext, next
		},
	)
}

func Iterable[T any](iter iterator[T]) Iterator[T] {
	return iter
}
