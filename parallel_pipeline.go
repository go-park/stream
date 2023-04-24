package stream

import (
	"runtime"

	"github.com/go-park/stream/support/function"
	"github.com/go-park/stream/support/optional"
)

var parallelism = runtime.NumCPU()

func setParallelism(n int) {
	parallelism = n
}

func GetParallelism() int {
	return parallelism
}

type (
	ParallelPipline[T any] struct {
		sp SimplePipline[T]
	}
)

func (p ParallelPipline[T]) Close() {
	return
}

func (p ParallelPipline[T]) Count() int {
	return p.sp.Count()
}

func (p ParallelPipline[T]) ToSlice() []T {
	return p.sp.ToSlice()
}

func (p ParallelPipline[T]) chunk() []chan T {
	list := p.ToSlice()
	parallelism := 1
	if p.sp.parallel {
		parallelism = GetParallelism()
	}
	chunkSize := len(list) / parallelism
	var chs []chan T
	for i := 0; i < len(list); i += chunkSize {
		ch := make(chan T)
		chs = append(chs, ch)
		go func(sub []T) {
			defer close(ch)
			for _, v := range sub {
				ch <- v
			}
		}(list[i : i+chunkSize])
	}
	return chs
}

func (p ParallelPipline[T]) ForEach(consumer function.Consumer[T]) {
	for _, source := range p.chunk() {
		go func(in chan T) {
			for {
				select {
				case v, ok := <-in:
					if ok {
						consumer(v)
					} else {
						return
					}
				}
			}
		}(source)
	}
	return
}

func (p ParallelPipline[T]) Parallel() Stream[T] {
	return nil
}

func (p ParallelPipline[T]) Sequential() Stream[T] {
	return p
}

func (p ParallelPipline[T]) Filter(pred function.Predicate[T]) Stream[T] {
	return p.sp.Filter(pred)
}

func (p ParallelPipline[T]) Limit(i uint) Stream[T] {
	return p.sp.Limit(i)
}

func (p ParallelPipline[T]) Skip(i uint) Stream[T] {
	return p.sp.Skip(i)
}

func (p ParallelPipline[T]) Distinct(equals function.BiPredicate[T, T]) Stream[T] {
	return p.sp.Distinct(equals)
}

func (p ParallelPipline[T]) Sort(less function.BiPredicate[T, T]) Stream[T] {
	return p.sp.Sort(less)
}

func (p ParallelPipline[T]) Reverse() Stream[T] {
	return p.sp.Reverse()
}

func (p ParallelPipline[T]) Max(less function.BiPredicate[T, T]) optional.Value[T] {
	return p.sp.Max(less)
}

func (p ParallelPipline[T]) Min(less function.BiPredicate[T, T]) optional.Value[T] {
	return p.sp.Min(less)
}

func (p ParallelPipline[T]) Map(mapper function.Func[T, T]) Stream[T] {
	return p.sp.Map(mapper)
}

func (p ParallelPipline[T]) Reduce(acc function.BiFunc[T, T, T]) optional.Value[T] {
	return optional.EmptyVal[T]()
}

func (p ParallelPipline[T]) MapToAny(mapper function.Func[T, any]) Stream[any] {
	return p.sp.MapToAny(mapper)
}

func (p ParallelPipline[T]) MapToString(mapper function.Func[T, string]) Stream[string] {
	return p.sp.MapToString(mapper)
}

func (p ParallelPipline[T]) MapToInt(mapper function.Func[T, int]) Stream[int] {
	return p.sp.MapToInt(mapper)
}

func (p ParallelPipline[T]) MapToFloat(mapper function.Func[T, float64]) Stream[float64] {
	return p.sp.MapToFloat(mapper)
}

func (p ParallelPipline[T]) AnyMatch(pred function.Predicate[T]) bool {
	return p.sp.AnyMatch(pred)
}

func (p ParallelPipline[T]) AllMatch(pred function.Predicate[T]) bool {
	return p.sp.AllMatch(pred)
}

func (p ParallelPipline[T]) NoneMatch(pred function.Predicate[T]) bool {
	return p.sp.NoneMatch(pred)
}

func (p ParallelPipline[T]) FindAny() optional.Value[T] {
	return p.sp.FindAny()
}
