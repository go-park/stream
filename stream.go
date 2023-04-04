package stream

import (
	"github.com/go-park/stream/function"
	"github.com/go-park/stream/support"
)

var (
	_ Stream[any] = SimplePipline[any]{}
	_ Stream[any] = ReusablePipline[any]{}
)

type Stream[T any] interface {
	Close()
	Count() int
	ToSlice() []T
	ForEach(f function.Consumer[T])
	Filter(p function.Predicate[T]) Stream[T]
	SortBy(less function.BiPredicate[T, T]) Stream[T]
	Limit(i uint) Stream[T]
	Skip(i uint) Stream[T]
	Max(less function.BiPredicate[T, T]) support.Value[T]
	Min(less function.BiPredicate[T, T]) support.Value[T]
	Map(f function.Fn[T, T]) Stream[T]
	Reduce(acc function.BiFn[T, T, T]) support.Value[T]
	MapToAny(f function.Fn[T, any]) Stream[any]
	MapToString(f function.Fn[T, string]) Stream[string]
	MapToInt(fn function.Fn[T, int]) Stream[int]
	MapToFloat(fn function.Fn[T, float64]) Stream[float64]
	AnyMatch(pred function.Predicate[T]) bool
	AllMatch(pred function.Predicate[T]) bool
	NoneMatch(pred function.Predicate[T]) bool
}
