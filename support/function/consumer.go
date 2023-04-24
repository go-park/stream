package function

import "github.com/go-park/stream/internal/helper"

type Consumer[T any] func(t T)

func (fn Consumer[T]) Accept(t T) {
	fn(t)
}

func (fn Consumer[T]) AndThen(after Consumer[T]) Consumer[T] {
	helper.RequireCanButNonNil(after)
	return func(t T) { fn.Accept(t); after.Accept(t) }
}
