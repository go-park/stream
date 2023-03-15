package stream

func OfList[L ~[]T, T any](list L) Stream[L, T] {
	return Stream[L, T]{list}
}

func OfMap[M ~map[K]V, K comparable, V any](m M) Stream[[]Entry[K, V], Entry[K, V]] {
	return Stream[[]Entry[K, V], Entry[K, V]]{GetEntrySet(m)}
}

func Of[T any](t ...T) Stream[[]T, T] {
	return Stream[[]T, T]{t}
}

func Builder[T any]() builder[T] {
	return builder[T]{}
}

type builder[T any] struct {
	items []T
}

func (b builder[T]) Append(t ...T) builder[T] {
	b.items = append(b.items, t...)
	return b
}

func (b builder[T]) Build() Stream[[]T, T] {
	return OfList(b.items)
}
