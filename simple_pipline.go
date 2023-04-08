package stream

import (
	"sort"

	"github.com/go-park/stream/internal/helper"
	"github.com/go-park/stream/support/function"
	"github.com/go-park/stream/support/optional"
)

type SimplePipline[T any] struct {
	upstream chan T
	cancel   func()
}

func (p SimplePipline[T]) Close() {
	p.cancel()
	p.ForEach(func(_ T) {})
}

func (p SimplePipline[T]) Count() int {
	var count int
	p.ForEach(
		func(_ T) {
			count++
		})
	return count
}

func (p SimplePipline[T]) ToSlice() []T {
	var slice []T
	p.ForEach(
		func(t T) {
			slice = append(slice, t)
		})
	return slice
}

func (p SimplePipline[T]) ForEach(fn function.Consumer[T]) {
	helper.RequireNonNil(fn)
	source := p.upstream
	for {
		select {
		case v, ok := <-source:
			if ok {
				fn(v)
			} else {
				return
			}
		}
	}
}

func (p SimplePipline[T]) Filter(pred function.Predicate[T]) Stream[T] {
	helper.RequireNonNil(pred)
	source := p.upstream
	target := make(chan T)
	p.upstream = target
	go func() {
		defer close(target)
		for {
			select {
			case v, ok := <-source:
				if ok {
					if pred.Test(v) {
						target <- v
					}
				} else {
					return
				}
			}
		}
	}()
	return p
}

func (p SimplePipline[T]) Limit(i uint) Stream[T] {
	source := p.upstream
	target := make(chan T)
	p.upstream = target
	go func() {
		defer close(target)
		var num uint = 0
		for {
			select {
			case v, ok := <-source:
				if ok {
					if num < i {
						target <- v
						num++
					}
				} else {
					return
				}
			}
		}
	}()
	return p
}

func (p SimplePipline[T]) Skip(i uint) Stream[T] {
	source := p.upstream
	target := make(chan T)
	p.upstream = target
	go func() {
		defer close(target)
		var num uint = 0
		for {
			select {
			case v, ok := <-source:
				if ok {
					if num < i {
						num++
					} else {
						target <- v
					}
				} else {
					return
				}
			}
		}
	}()

	return p
}

func (p SimplePipline[T]) Distinct(equals function.BiPredicate[T, T]) Stream[T] {
	helper.RequireNonNil(equals)
	source := p.upstream
	target := make(chan T)
	p.upstream = target
	go func() {
		defer close(target)
		var list []T
		for {
			select {
			case v, ok := <-source:
				if ok {
					noneMatch := From(list...).NoneMatch(
						func(t T) bool {
							return equals(v, t)
						})
					if noneMatch {
						list = append(list, v)
					}
				} else {
					for _, v := range list {
						target <- v
					}
					return
				}
			}
		}
	}()
	return p
}

func (p SimplePipline[T]) Sort(less function.BiPredicate[T, T]) Stream[T] {
	helper.RequireNonNil(less)
	source := p.upstream
	target := make(chan T)
	p.upstream = target
	go func() {
		defer close(target)
		var list []T
		for {
			select {
			case v, ok := <-source:
				if ok {
					list = append(list, v)
				} else {
					sort.Slice(list, func(i, j int) bool {
						return less(list[i], list[j])
					})
					for _, v := range list {
						target <- v
					}
					return
				}
			}
		}
	}()
	return p
}

func (p SimplePipline[T]) Reverse() Stream[T] {
	source := p.upstream
	target := make(chan T)
	p.upstream = target
	go func() {
		defer close(target)
		var list []T
		for {
			select {
			case v, ok := <-source:
				if ok {
					list = append(list, v)
				} else {
					for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
						list[i], list[j] = list[j], list[i]
					}
					for _, v := range list {
						target <- v
					}
					return
				}
			}
		}
	}()
	return p
}

func (p SimplePipline[T]) Max(less function.BiPredicate[T, T]) optional.Value[T] {
	helper.RequireNonNil(less)
	var max T
	val := optional.EmptyVal[T]()
	hasPre := false
	p.ForEach(func(t T) {
		if !hasPre {
			hasPre = true
			max = t
		} else {
			if less.Test(max, t) {
				max = t
			}
		}
		val = optional.ValOf(max)
	})
	return val
}

func (p SimplePipline[T]) Min(less function.BiPredicate[T, T]) optional.Value[T] {
	helper.RequireNonNil(less)
	var min T
	val := optional.EmptyVal[T]()
	hasPre := false
	p.ForEach(func(t T) {
		if !hasPre {
			hasPre = true
			min = t
		} else {
			if less.Test(t, min) {
				min = t
			}
		}
		val = optional.ValOf(min)
	})
	return val
}

func (p SimplePipline[T]) Map(fn function.Fn[T, T]) Stream[T] {
	helper.RequireNonNil(fn)
	target := make(chan T)
	source := p.upstream
	p.upstream = target
	go func() {
		defer close(target)
		for {
			select {
			case v, ok := <-source:
				if ok {
					target <- fn(v)
				} else {
					return
				}
			}
		}
	}()
	return p
}

func (p SimplePipline[T]) Reduce(acc function.BiFn[T, T, T]) optional.Value[T] {
	helper.RequireNonNil(acc)
	var res T
	val := optional.EmptyVal[T]()
	hasPre := false
	p.ForEach(func(t T) {
		if !hasPre {
			hasPre = true
			res = t
		} else {
			res = acc.Apply(res, t)
		}
		val = optional.ValOf(res)
	})
	return val
}

func (p SimplePipline[T]) MapToAny(fn function.Fn[T, any]) Stream[any] {
	helper.RequireNonNil(fn)
	source := p.upstream
	target := make(chan any)
	anyPipe := SimplePipline[any]{
		upstream: target,
	}
	go func() {
		defer close(target)
		for {
			select {
			case v, ok := <-source:
				if ok {
					target <- fn(v)
				} else {
					return
				}
			}
		}
	}()
	return anyPipe
}

func (p SimplePipline[T]) MapToString(fn function.Fn[T, string]) Stream[string] {
	helper.RequireNonNil(fn)
	source := p.upstream
	target := make(chan string)
	strPipe := SimplePipline[string]{
		upstream: target,
	}
	go func() {
		defer close(target)
		for {
			select {
			case v, ok := <-source:
				if ok {
					target <- fn(v)
				} else {
					return
				}
			}
		}
	}()
	return strPipe
}

func (p SimplePipline[T]) MapToInt(fn function.Fn[T, int]) Stream[int] {
	helper.RequireNonNil(fn)
	source := p.upstream
	target := make(chan int)
	intPipe := SimplePipline[int]{
		upstream: target,
	}
	go func() {
		defer close(target)
		for {
			select {
			case v, ok := <-source:
				if ok {
					target <- fn(v)
				} else {
					return
				}
			}
		}
	}()
	return intPipe
}

func (p SimplePipline[T]) MapToFloat(fn function.Fn[T, float64]) Stream[float64] {
	helper.RequireNonNil(fn)
	source := p.upstream
	target := make(chan float64)
	floatPipe := SimplePipline[float64]{
		upstream: target,
	}
	go func() {
		defer close(target)
		for {
			select {
			case v, ok := <-source:
				if ok {
					target <- fn(v)
				} else {
					return
				}
			}
		}
	}()
	return floatPipe
}

func (p SimplePipline[T]) AnyMatch(pred function.Predicate[T]) bool {
	helper.RequireNonNil(pred)
	match := false
	p.ForEach(func(t T) {
		if pred.Test(t) {
			match = true
			p.cancel()
		}
	})
	return match
}

func (p SimplePipline[T]) AllMatch(pred function.Predicate[T]) bool {
	helper.RequireNonNil(pred)
	match := true
	p.ForEach(func(t T) {
		if !pred.Test(t) {
			match = false
			p.cancel()
		}
	})
	return match
}

func (p SimplePipline[T]) NoneMatch(pred function.Predicate[T]) bool {
	helper.RequireNonNil(pred)
	match := false
	p.ForEach(func(t T) {
		if pred.Test(t) {
			match = true
			p.cancel()
		}
	})
	return !match
}
