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
	helper.RequireNonNil(hasNext)

	return hasNext()
}

func (iter iterator[T]) Next() T {
	_, next := iter()
	helper.RequireNonNil(next)

	return next()
}

func (iter iterator[T]) ForEachRemaining(fn function.Consumer[T]) {
	helper.RequireNonNil(fn)
	hasNext, next := iter()
	helper.RequireNonNil(hasNext)
	helper.RequireNonNil(next)

	for hasNext() {
		fn(next())
	}
}

func ToIterator[T any](list ...T) Iterator[T] {
	// escape to heap
	index, remainNum := 0, len(list)
	return iterator[T](
		func() (func() bool, function.Supplier[T]) {
			hasNext := func() bool {
				return remainNum > 0
			}
			next := func() T {
				var v T
				if hasNext() {
					remainNum--
					v = list[index]
					index++
				}
				return v
			}
			return hasNext, next
		},
	)
}

// func
// type iterator[T any] struct{

// }
