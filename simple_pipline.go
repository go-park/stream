package stream

import (
	"context"
	"sort"
	"sync"

	"github.com/go-park/stream/internal/helper"
	"github.com/go-park/stream/support/function"
	"github.com/go-park/stream/support/optional"
	"github.com/go-park/stream/support/routine"
)

type SimplePipline[T any] struct {
	upstream chan T
	cancel   context.CancelFunc
	parallel bool
}

func (p SimplePipline[T]) Close() {
	p.cancel()
}

func (p SimplePipline[T]) Parallel() Stream[T] {
	p.parallel = true
	return p
}

func (p SimplePipline[T]) Sequential() Stream[T] {
	p.parallel = false
	return p
}

func (p SimplePipline[T]) Count() int {
	acc := func(_ T, i int) int {
		return i + 1
	}
	ch := aggregator(reduce[T, int], acc, 0, p.chunk())
	return reduce(ch, function.Sum[int], 0).Get()
}

func (p SimplePipline[T]) ToSlice() []T {
	var slice []T
	for v := range p.upstream {
		slice = append(slice, v)
	}
	return slice
}

func (p SimplePipline[T]) chunk() chan chan T {
	// serial default
	chunkNum := 1
	if p.parallel {
		chunkNum = GetParallelism()
	}
	count := 0
	chs := make(chan chan T)
	routine.Run(func() {
		var slice []chan T
		defer func() { close(chs) }()
		for v := range p.upstream {
			chunk := count % chunkNum
			if len(slice) == chunk {
				slice = append(slice, make(chan T))
				chs <- slice[chunk]
				defer close(slice[chunk])
			}
			slice[chunk] <- v
			count++
		}
	})

	return chs
}

type reducer[T, R any] func(in chan T, acc function.BiFunc[T, R, R], identify R) optional.Value[R]

func reduce[T, R any](in chan T, acc function.BiFunc[T, R, R], identify R) optional.Value[R] {
	val := optional.ValOf(identify)
	for v := range in {
		val = optional.ValOf(acc.Apply(v, val.Get()))
	}
	return val
}

func aggregator[T, R any](reduce reducer[T, R], acc function.BiFunc[T, R, R], identify R, chs chan chan T) chan R {
	helper.RequireCanButNonNil(reduce)
	ch := make(chan R)
	var wg sync.WaitGroup
	routine.Run(
		func() {
			for source := range chs {
				wg.Add(1)
				routine.RunArg(source, func(in chan T) {
					defer wg.Done()
					reduce(in, acc, identify).IfNotEmpty(
						func(t R) {
							ch <- t
						})
				})
			}
			wg.Wait()
			close(ch)
		})
	return ch
}

func (p SimplePipline[T]) ForEach(fn function.Consumer[T]) {
	helper.RequireCanButNonNil(fn)
	acc := func(t T, i struct{}) struct{} {
		fn(t)
		return struct{}{}
	}
	for range aggregator(reduce[T, struct{}], acc, struct{}{}, p.chunk()) {
	}
}

func (p SimplePipline[T]) Filter(pred function.Predicate[T]) Stream[T] {
	helper.RequireCanButNonNil(pred)
	target := make(chan T)
	acc := func(t T, _ struct{}) struct{} {
		if pred.Test(t) {
			target <- t
		}
		return struct{}{}
	}
	ch := aggregator(reduce[T, struct{}], acc, struct{}{}, p.chunk())
	routine.Run(
		func() {
			defer close(target)
			for range ch {
			}
		})
	p.upstream = target
	return p
}

func (p SimplePipline[T]) Limit(i uint) Stream[T] {
	source := p.upstream
	target := make(chan T)
	p.upstream = target
	routine.Run(func() {
		defer close(target)
		var num uint = 0
		for v := range source {
			if num < i {
				target <- v
				num++
			}
		}
	})
	return p
}

func (p SimplePipline[T]) Skip(i uint) Stream[T] {
	source := p.upstream
	target := make(chan T)
	p.upstream = target
	routine.Run(func() {
		defer close(target)
		var num uint = 0
		for v := range source {
			if num < i {
				num++
			} else {
				target <- v
			}
		}
	})

	return p
}

func (p SimplePipline[T]) Distinct(equals function.BiPredicate[T, T]) Stream[T] {
	helper.RequireCanButNonNil(equals)
	source := p.upstream
	target := make(chan T)
	p.upstream = target
	routine.Run(func() {
		defer close(target)
		var list []T
		for v := range source {
			exists := false
			for _, item := range list {
				if equals.Test(v, item) {
					exists = true
					break
				}
			}
			if !exists {
				list = append(list, v)
				target <- v
			}
		}
		return
	})
	return p
}

func (p SimplePipline[T]) Sort(less function.BiPredicate[T, T]) Stream[T] {
	helper.RequireCanButNonNil(less)
	source := p.upstream
	target := make(chan T)
	p.upstream = target
	routine.Run(func() {
		defer close(target)
		var list []T
		for v := range source {
			list = append(list, v)
		}
		sort.Slice(list, func(i, j int) bool {
			return less(list[i], list[j])
		})
		for _, v := range list {
			target <- v
		}
	})
	return p
}

func (p SimplePipline[T]) Reverse() Stream[T] {
	source := p.upstream
	target := make(chan T)
	p.upstream = target
	routine.Run(func() {
		defer close(target)
		var list []T
		for v := range source {
			list = append(list, v)
		}
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
		for _, v := range list {
			target <- v
		}
	})
	return p
}

func (p SimplePipline[T]) Max(less function.BiPredicate[T, T]) optional.Value[T] {
	helper.RequireCanButNonNil(less)
	return p.Reduce(func(t1, t2 T) T {
		if !less.Test(t1, t2) {
			return t1
		}
		return t2
	})
}

func (p SimplePipline[T]) Min(less function.BiPredicate[T, T]) optional.Value[T] {
	helper.RequireCanButNonNil(less)
	return p.Reduce(func(t1, t2 T) T {
		if less.Test(t1, t2) {
			return t1
		}
		return t2
	})
}

func (p SimplePipline[T]) Map(mapper function.Func[T, T]) Stream[T] {
	helper.RequireCanButNonNil(mapper)
	target := make(chan T)
	acc := func(t T, _ struct{}) struct{} {
		target <- mapper(t)
		return struct{}{}
	}
	ch := aggregator(reduce[T, struct{}], acc, struct{}{}, p.chunk())
	routine.Run(
		func() {
			defer close(target)
			for range ch {
			}
		})
	p.upstream = target
	return p
}

func (p SimplePipline[T]) Reduce(acc function.BiFunc[T, T, T]) optional.Value[T] {
	helper.RequireCanButNonNil(acc)
	reduce := func(in chan T, acc function.BiFunc[T, T, T], _ T) optional.Value[T] {
		val := optional.EmptyVal[T]()
		for v := range in {
			val.IfNotEmptyOrElse(
				func(t T) { val = optional.ValOf(acc.Apply(v, t)) },
				func() { val = optional.ValOf(v) })
		}
		return val
	}
	var identify T
	ch := aggregator(reduce, acc, identify, p.chunk())
	return reduce(ch, acc, identify)
}

func (p SimplePipline[T]) MapToAny(mapper function.Func[T, any]) Stream[any] {
	helper.RequireCanButNonNil(mapper)
	target := make(chan any)
	acc := func(t T, _ struct{}) struct{} {
		target <- mapper(t)
		return struct{}{}
	}
	ch := aggregator(reduce[T, struct{}], acc, struct{}{}, p.chunk())
	routine.Run(
		func() {
			defer close(target)
			for range ch {
			}
		})
	return SimplePipline[any]{
		upstream: target,
		cancel:   p.cancel,
		parallel: p.parallel,
	}
}

func (p SimplePipline[T]) MapToString(mapper function.Func[T, string]) Stream[string] {
	helper.RequireCanButNonNil(mapper)
	target := make(chan string)
	acc := func(t T, _ struct{}) struct{} {
		target <- mapper(t)
		return struct{}{}
	}
	ch := aggregator(reduce[T, struct{}], acc, struct{}{}, p.chunk())
	routine.Run(
		func() {
			defer close(target)
			for range ch {
			}
		})
	return SimplePipline[string]{
		upstream: target,
		cancel:   p.cancel,
		parallel: p.parallel,
	}
}

func (p SimplePipline[T]) MapToInt(mapper function.Func[T, int]) Stream[int] {
	helper.RequireCanButNonNil(mapper)
	target := make(chan int)
	acc := func(t T, _ struct{}) struct{} {
		target <- mapper(t)
		return struct{}{}
	}
	ch := aggregator(reduce[T, struct{}], acc, struct{}{}, p.chunk())
	routine.Run(
		func() {
			defer close(target)
			for range ch {
			}
		})
	return SimplePipline[int]{
		upstream: target,
		cancel:   p.cancel,
		parallel: p.parallel,
	}
}

func (p SimplePipline[T]) MapToFloat(mapper function.Func[T, float64]) Stream[float64] {
	helper.RequireCanButNonNil(mapper)
	target := make(chan float64)
	acc := func(t T, _ struct{}) struct{} {
		target <- mapper(t)
		return struct{}{}
	}
	ch := aggregator(reduce[T, struct{}], acc, struct{}{}, p.chunk())
	routine.Run(
		func() {
			defer close(target)
			for range ch {
			}
		})
	return SimplePipline[float64]{
		upstream: target,
		cancel:   p.cancel,
		parallel: p.parallel,
	}
}

func (p SimplePipline[T]) AnyMatch(pred function.Predicate[T]) bool {
	helper.RequireCanButNonNil(pred)
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
	helper.RequireCanButNonNil(pred)
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
	helper.RequireCanButNonNil(pred)
	match := false
	p.ForEach(func(t T) {
		if pred.Test(t) {
			match = true
			p.cancel()
		}
	})
	return !match
}

func (p SimplePipline[T]) FindAny() optional.Value[T] {
	source := p.upstream
	r := optional.EmptyVal[T]()
	if v, ok := <-source; ok {
		r = optional.ValOf(v)
	}
	return r
}
