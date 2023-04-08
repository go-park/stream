package stream

import (
	"github.com/go-park/stream/support/function"
	"github.com/go-park/stream/support/optional"
)

type ReusablePipline[T any] struct {
	source []T
	cancel func()
}

func (p ReusablePipline[T]) Close() {
}

func (p ReusablePipline[T]) Count() int {
	var count int
	p.ForEach(
		func(_ T) {
			count++
		})
	return count
}

func (p ReusablePipline[T]) ToSlice() []T {
	var slice []T
	p.ForEach(
		func(t T) {
			slice = append(slice, t)
		})
	return slice
}

func (p ReusablePipline[T]) initSource() (r <-chan T, w chan T) {
	// w = make(chan T)
	// if !p.notHead {
	// 	p.notHead = true
	// 	p.upstream = make(chan T)
	// 	p.listener <- p.upstream
	// }
	// return p.upstream, w
	return nil, nil
}

func (p ReusablePipline[T]) ForEach(f function.Consumer[T]) {
	source, _ := p.initSource()
	for {
		select {
		case v, ok := <-source:
			if ok {
				f(v)
			} else {
				return
			}
		}
	}
}

func (p ReusablePipline[T]) Filter(pred function.Predicate[T]) Stream[T] {
	source, target := p.initSource()
	go func() {
		for {
			select {
			case v, ok := <-source:
				if ok {
					if pred.Test(v) {
						target <- v
					}
				} else {
					close(target)
					return
				}
			}
		}
	}()
	return p
}

func (p ReusablePipline[T]) Limit(i uint) Stream[T] {
	return nil
}

func (p ReusablePipline[T]) Skip(i uint) Stream[T] {
	return nil
}

func (p ReusablePipline[T]) Distinct(less function.BiPredicate[T, T]) Stream[T] {
	return nil
}

func (p ReusablePipline[T]) Sort(less function.BiPredicate[T, T]) Stream[T] {
	return nil
}

func (p ReusablePipline[T]) Reverse() Stream[T] {
	return nil
}

func (p ReusablePipline[T]) Max(less function.BiPredicate[T, T]) optional.Value[T] {
	return optional.EmptyVal[T]()
}

func (p ReusablePipline[T]) Min(less function.BiPredicate[T, T]) optional.Value[T] {
	return optional.EmptyVal[T]()
}

func (p ReusablePipline[T]) Map(f function.Fn[T, T]) Stream[T] {
	// target := make(chan T)
	// source := p.upstream
	// go func() {
	// 	for v := range source {
	// 		target <- f(v)
	// 	}
	// }()
	// p.upstream = target
	return p
}

func (p ReusablePipline[T]) Reduce(acc function.BiFn[T, T, T]) optional.Value[T] {
	return optional.EmptyVal[T]()
}

func (p ReusablePipline[T]) MapToAny(f function.Fn[T, any]) Stream[any] {
	return nil
}

func (p ReusablePipline[T]) MapToString(f function.Fn[T, string]) Stream[string] {
	return nil
}

func (p ReusablePipline[T]) MapToInt(fn function.Fn[T, int]) Stream[int] {
	return nil
}

func (p ReusablePipline[T]) MapToFloat(fn function.Fn[T, float64]) Stream[float64] {
	return nil
}

func (p ReusablePipline[T]) AnyMatch(pred function.Predicate[T]) bool {
	return false
}

func (p ReusablePipline[T]) AllMatch(pred function.Predicate[T]) bool {
	return false
}

func (p ReusablePipline[T]) NoneMatch(pred function.Predicate[T]) bool {
	return false
}
