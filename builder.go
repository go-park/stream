package stream

import (
	"context"
)

func Builder[T any]() builder[T] {
	return builder[T]{}
}

type builder[T any] struct {
	items    []T
	reusable bool
}

// func (b builder[T]) Reusable() builder[T] {
// 	b.reusable = true
// 	return b
// }

func (b builder[T]) Append(t ...T) builder[T] {
	b.items = append(b.items, t...)
	return b
}

func (b builder[T]) Build() Stream[T] {
	if b.reusable {
		return reusable(b.items...)
	}
	return From(b.items...)
}

func FromMap[M ~map[K]V, K comparable, V any](m M) Stream[Entry[K, V]] {
	return From(GetEntrySet(m)...)
}

func From[T any](list ...T) Stream[T] {
	target := make(chan T)
	ctx, cancelFn := context.WithCancel(context.Background())
	go func() {
		defer close(target)
		for _, v := range list {
			select {
			case <-ctx.Done():
				return
			default:
				target <- v
			}
		}
	}()
	return SimplePipline[T]{upstream: target, cancel: cancelFn}
}

func reusable[T any](list ...T) Stream[T] {
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
