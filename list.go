package stream

import "sort"

type Stream[L ~[]T, T any] struct {
	items L
}

func (l Stream[L, T]) Len() int {
	return len(l.items)
}

func (l Stream[L, T]) List() []T {
	l = l.copy()
	return l.items
}

func (l Stream[L, T]) Filter(f func(t T) bool) Stream[L, T] {
	items := make([]T, 0, l.Len())
	for _, v := range l.items {
		if f(v) {
			items = append(items, v)
		}
	}
	l.items = items
	return l
}

func (l Stream[L, T]) Sort(less func(i, j T) bool) Stream[L, T] {
	l = l.copy()
	sort.Slice(l.items, func(i int, j int) bool {
		return less(l.items[i], l.items[j])
	})
	return l
}

func (l Stream[L, T]) copy() Stream[L, T] {
	items := make([]T, l.Len())
	copy(items, l.items)
	l.items = items
	return l
}

func (l Stream[L, T]) Limit(i int) Stream[L, T] {
	if i > l.Len() {
		i = l.Len()
	}
	items := make([]T, i)
	copy(items, l.items[:i])
	l.items = items
	return l
}

func (l Stream[L, T]) ForEach(f func(t T)) Stream[L, T] {
	l = l.copy()
	for _, v := range l.items {
		f(v)
	}
	return l
}
