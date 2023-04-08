package stream

import (
	"github.com/go-park/stream/support/function"
	"github.com/go-park/stream/support/optional"
)

var (
	_ Stream[any] = SimplePipline[any]{}
	_ Stream[any] = ReusablePipline[any]{}
)

type Stream[T any] interface {
	Close()
	Count() int
	ToSlice() []T
	ForEach(consumer function.Consumer[T])
	Filter(pred function.Predicate[T]) Stream[T]
	Limit(i uint) Stream[T]
	Skip(i uint) Stream[T]
	Distinct(equals function.BiPredicate[T, T]) Stream[T]
	Sort(less function.BiPredicate[T, T]) Stream[T]
	Reverse() Stream[T]
	Max(less function.BiPredicate[T, T]) optional.Value[T]
	Min(less function.BiPredicate[T, T]) optional.Value[T]
	Map(mapper function.Fn[T, T]) Stream[T]
	Reduce(acc function.BiFn[T, T, T]) optional.Value[T]
	MapToAny(mapper function.Fn[T, any]) Stream[any]
	MapToString(mapper function.Fn[T, string]) Stream[string]
	MapToInt(mapper function.Fn[T, int]) Stream[int]
	MapToFloat(mapper function.Fn[T, float64]) Stream[float64]
	AnyMatch(pred function.Predicate[T]) bool
	AllMatch(pred function.Predicate[T]) bool
	NoneMatch(pred function.Predicate[T]) bool
}
