package stream

import (
	"context"

	"github.com/go-park/stream/support/collections"
)

func Builder[T any]() builder[T] {
	return builder[T]{}
}

type builder[T any] struct {
	iter     collections.Iterator[T]
	reusable bool
}

// func (b builder[T]) Reusable() builder[T] {
// 	b.reusable = true
// 	return b
// }

func (b builder[T]) Source(t ...T) builder[T] {
	b.iter = collections.ToIterator(t...)
	return b
}

func (b builder[T]) Build() Stream[T] {
	if b.reusable {
		return reusable(b.iter)
	}
	return b.buildSimple()
}

func (b builder[T]) buildSimple() Stream[T] {
	target := make(chan T)
	ctx, cancelFn := context.WithCancel(context.Background())
	go func() {
		defer close(target)
		b.iter.ForEachRemaining(
			func(v T) {
				select {
				case <-ctx.Done():
					return
				default:
					target <- v
				}
			})
	}()
	return SimplePipline[T]{upstream: target, cancel: cancelFn}
}

func (b builder[T]) buildReusable() Stream[T] {
	return reusable(b.iter)
}

func FromMap[M ~map[K]V, K comparable, V any](m M) Stream[collections.Entry[K, V]] {
	return From(collections.GetEntrySet(m)...)
}

func From[T any](list ...T) Stream[T] {
	return Builder[T]().Source(list...).Build()
}

func reusable[T any](iter collections.Iterator[T]) Stream[T] {
	// target := make(chan T)
	// ctx, cancelFn := context.WithCancel(context.Background())
	// go func() {
	// 	defer close(target)
	// 	for _, v := range list {
	// 		select {
	// 		case <-ctx.Done():
	// 			return
	// 		default:
	// 			target <- v
	// 		}
	// 	}
	// }()
	return ReusablePipline[T]{}
}

// func Generate[T any](sup function.Supplier[T]) Stream[T] {
// 	return FromList(t)
// }
