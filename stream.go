package stream

import (
	"github.com/go-park/stream/support/function"
	"github.com/go-park/stream/support/optional"
)

var (
	_ Stream[any] = SimplePipline[any]{}
	_ Stream[any] = ParallelPipline[any]{}
)

type Stream[T any] interface {
	Close()
	Count() int
	ToSlice() []T
	ForEach(consumer function.Consumer[T])
	Parallel() Stream[T]
	Sequential() Stream[T]
	Filter(pred function.Predicate[T]) Stream[T]
	Limit(i uint) Stream[T]
	Skip(i uint) Stream[T]
	Distinct(equals function.BiPredicate[T, T]) Stream[T]
	Sort(less function.BiPredicate[T, T]) Stream[T]
	Reverse() Stream[T]
	Max(less function.BiPredicate[T, T]) optional.Value[T]
	Min(less function.BiPredicate[T, T]) optional.Value[T]
	Map(mapper function.Func[T, T]) Stream[T]
	Reduce(acc function.BiFunc[T, T, T]) optional.Value[T]
	MapToAny(mapper function.Func[T, any]) Stream[any]
	MapToString(mapper function.Func[T, string]) Stream[string]
	MapToInt(mapper function.Func[T, int]) Stream[int]
	MapToFloat(mapper function.Func[T, float64]) Stream[float64]
	AnyMatch(pred function.Predicate[T]) bool
	AllMatch(pred function.Predicate[T]) bool
	NoneMatch(pred function.Predicate[T]) bool
	FindAny() optional.Value[T]
}
