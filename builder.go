package stream

import (
	"context"

	"github.com/go-park/stream/support/collections"
	"github.com/go-park/stream/support/function"
	"github.com/go-park/stream/support/routine"
	"golang.org/x/exp/constraints"
)

func Builder[T any]() builder[T] {
	return builder[T]{}
}

type builder[T any] struct {
	iter     collections.Iterator[T]
	reusable bool
	parallel bool
}

func (b builder[T]) Source(t ...T) builder[T] {
	b.iter = collections.IterableSlice(t...)
	return b
}

func (b builder[T]) iterator(iter collections.Iterator[T]) builder[T] {
	b.iter = iter
	return b
}

func (b builder[T]) Parallel() builder[T] {
	b.parallel = true
	return b
}

func (b builder[T]) Simple() Stream[T] {
	return b.buildSimple()
}

func (b builder[T]) Build() Stream[T] {
	return b.buildFast()
}

func (b builder[T]) buildSimple() SimplePipline[T] {
	target := make(chan T)
	ctx, cancelFn := context.WithCancel(context.Background())
	routine.Run(func() {
		defer cancelFn()
		defer close(target)
		for b.iter.HasNext() {
			select {
			case <-ctx.Done():
				return
			default:
				target <- b.iter.Next()
			}
		}
	})
	return SimplePipline[T]{upstream: target, cancel: cancelFn, parallel: b.parallel}
}

func (b builder[T]) buildFast() Stream[T] {
	_, cancelFn := context.WithCancel(context.Background())
	return &FastPipline[T]{source: b.iter, opWrapper: defultOpWrapper[T], cancel: cancelFn}
}

func defultOpWrapper[T any](down function.Consumer[T]) function.Consumer[T] {
	return func(t T) {
		down.Accept(t)
	}
}

func FromMap[M ~map[K]V, K comparable, V any](m M) Stream[collections.Entry[K, V]] {
	return From(collections.GetEntrySet(m)...)
}

func From[T any](list ...T) Stream[T] {
	return Builder[T]().Source(list...).Build()
}

func Range[T constraints.Integer](start, end T) Stream[T] {
	var iter collections.Iterator[T] = collections.Iterable(
		func() (func() bool, function.Supplier[T]) {
			hasNext := func() bool {
				return start <= end
			}
			next := func() T {
				cur := start
				start++
				return cur
			}
			return hasNext, next
		})
	return Builder[T]().iterator(iter).Build()
}
