package stream

import (
	"context"
	"sort"
	"sync"

	"github.com/go-park/stream/internal/helper"
	"github.com/go-park/stream/support/collections"
	"github.com/go-park/stream/support/function"
	"github.com/go-park/stream/support/optional"
	"github.com/go-park/stream/support/routine"
)

type FastPipline[T any] struct {
	source    collections.Iterator[T]
	opWrapper func(down function.Consumer[T]) function.Consumer[T]
	cancel    context.CancelFunc
	parallel  bool
}

func (p *FastPipline[T]) Close() {
	p.cancel()
}

func (p *FastPipline[T]) Parallel() Stream[T] {
	// p.parallel = true
	return p
}

func (p *FastPipline[T]) Sequential() Stream[T] {
	p.parallel = false
	return p
}

func (p *FastPipline[T]) Count() int {
	acc := func(_ T, i int) int {
		return i + 1
	}
	var list []*optional.Value[int]
	fac := func() function.Consumer[T] {
		val, op := reduceOp(acc, 0)
		list = append(list, val)
		return p.opWrapper(op)
	}
	p.exec(fac)
	val, sum := reduceOp(func(i1, i2 int) int {
		return i1 + i2
	}, 0)
	for _, item := range list {
		item.IfNotEmpty(sum)
	}
	return val.Get()
}

func (p *FastPipline[T]) ToSlice() []T {
	slices := make([][]T, 0)

	fac := func() function.Consumer[T] {
		slices = append(slices, make([]T, 0))
		index := len(slices) - 1
		return p.opWrapper(func(t T) { slices[index] = append(slices[index], t) })
	}
	p.exec(fac)
	var r []T
	for _, v := range slices {
		r = append(r, v...)
	}
	return r
}

func (p *FastPipline[T]) ForEach(fn function.Consumer[T]) {
	helper.RequireCanButNonNil(fn)
	fac := func() function.Consumer[T] { return p.opWrapper(fn) }
	p.exec(fac)
}

func (p *FastPipline[T]) exec(fac func() function.Consumer[T]) {
	if !p.parallel {
		p.source.ForEachRemaining(fac())
		return
	}
	var chunks [][]T
	var wg sync.WaitGroup
	chunkNum := GetParallelism()
	count := 0
	p.source.ForEachRemaining(
		func(v T) {
			chunk := count % chunkNum
			if len(chunks) == chunk {
				chunks = append(chunks, make([]T, 0))
			}
			chunks[chunk] = append(chunks[chunk], v)
			count++
		})
	for _, list := range chunks {
		wg.Add(1)
		op := fac()
		routine.RunArg(list, func(list []T) {
			defer wg.Done()
			for _, v := range list {
				op(v)
			}
		})
	}
	wg.Wait()
}

func (p *FastPipline[T]) Filter(pred function.Predicate[T]) Stream[T] {
	helper.RequireCanButNonNil(pred)
	wrap := p.opWrapper
	p.opWrapper = func(down function.Consumer[T]) function.Consumer[T] {
		op := func(t T) {
			if pred.Test(t) {
				down.Accept(t)
			}
		}
		return wrap(op)
	}
	return p
}

func (p *FastPipline[T]) Limit(i uint) Stream[T] {
	wrap := p.opWrapper
	p.opWrapper = func(down function.Consumer[T]) function.Consumer[T] {
		var num uint = 0
		op := func(t T) {
			if num < i {
				down.Accept(t)
				num++
			}
		}
		return wrap(op)
	}
	return p
}

func (p *FastPipline[T]) Skip(i uint) Stream[T] {
	wrap := p.opWrapper
	p.opWrapper = func(down function.Consumer[T]) function.Consumer[T] {
		var num uint = 0
		op := func(t T) {
			if num < i {
				num++
			} else {
				down.Accept(t)
			}
		}
		return wrap(op)
	}
	return p
}

func (p *FastPipline[T]) Distinct(equals function.BiPredicate[T, T]) Stream[T] {
	helper.RequireCanButNonNil(equals)
	wrap := p.opWrapper
	p.opWrapper = func(down function.Consumer[T]) function.Consumer[T] {
		var list []T
		op := func(v T) {
			exists := false
			for _, item := range list {
				if equals.Test(v, item) {
					exists = true
					break
				}
			}
			if !exists {
				list = append(list, v)
				down.Accept(v)
			}
		}
		return wrap(op)
	}
	return p
}

func (p *FastPipline[T]) Sort(less function.BiPredicate[T, T]) Stream[T] {
	helper.RequireCanButNonNil(less)
	wrap := p.opWrapper
	p.opWrapper = func(down function.Consumer[T]) function.Consumer[T] {
		var list []T
		op := func(t T) {
			list = append(list, t)
		}
		p.exec(func() function.Consumer[T] { return wrap(op) })
		sort.Slice(list, func(i, j int) bool {
			return less(list[i], list[j])
		})
		p.source = collections.IterableSlice(list...)
		return defultOpWrapper(down)
	}
	return p
}

func (p *FastPipline[T]) Reverse() Stream[T] {
	wrap := p.opWrapper
	p.opWrapper = func(down function.Consumer[T]) function.Consumer[T] {
		var list []T
		op := func(t T) {
			list = append(list, t)
		}
		p.exec(func() function.Consumer[T] { return wrap(op) })
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
		p.source = collections.IterableSlice(list...)
		return defultOpWrapper(down)
	}
	return p
}

func (p *FastPipline[T]) Max(less function.BiPredicate[T, T]) optional.Value[T] {
	helper.RequireCanButNonNil(less)
	return p.Reduce(func(t1, t2 T) T {
		if !less.Test(t1, t2) {
			return t1
		}
		return t2
	})
}

func (p *FastPipline[T]) Min(less function.BiPredicate[T, T]) optional.Value[T] {
	helper.RequireCanButNonNil(less)
	return p.Reduce(func(t1, t2 T) T {
		if less.Test(t1, t2) {
			return t1
		}
		return t2
	})
}

func (p *FastPipline[T]) Map(mapper function.Func[T, T]) Stream[T] {
	helper.RequireCanButNonNil(mapper)
	wrap := p.opWrapper
	p.opWrapper = func(down function.Consumer[T]) function.Consumer[T] {
		op := func(t T) {
			down.Accept(mapper(t))
		}
		return wrap(op)
	}
	return p
}

func reduceOp[T, R any](acc function.BiFunc[T, R, R], identify R) (*optional.Value[R], function.Consumer[T]) {
	val := optional.EmptyVal[R]()
	return &val, func(v T) {
		val.IfNotEmptyOrElse(
			func(t R) { val = optional.ValOf(acc.Apply(v, t)) },
			func() { val = optional.ValOf(acc.Apply(v, identify)) })
	}
}

func (p *FastPipline[T]) reduceOp(acc function.BiFunc[T, T, T], identify T) (*optional.Value[T], function.Consumer[T]) {
	val := optional.EmptyVal[T]()
	return &val, func(v T) {
		val.IfNotEmptyOrElse(
			func(t T) {
				val = optional.ValOf(acc.Apply(v, t))
			},
			func() {
				val = optional.ValOf(v)
			})
	}
}

func (p *FastPipline[T]) Reduce(acc function.BiFunc[T, T, T]) optional.Value[T] {
	var list []*optional.Value[T]
	var identify T

	fac := func() function.Consumer[T] {
		val, op := p.reduceOp(acc, identify)
		list = append(list, val)
		return p.opWrapper(op)
	}
	p.exec(fac)
	v, op := p.reduceOp(acc, identify)
	for _, item := range list {
		item.IfNotEmpty(op)
	}
	return *v
}

func (p *FastPipline[T]) MapToAny(mapper function.Func[T, any]) Stream[any] {
	helper.RequireCanButNonNil(mapper)
	anyPipe := &FastPipline[any]{
		cancel:   p.cancel,
		parallel: p.parallel,
	}
	anyPipe.opWrapper = func(down function.Consumer[any]) function.Consumer[any] {
		var list []any
		op := func(t T) {
			list = append(list, mapper(t))
		}
		p.exec(func() function.Consumer[T] { return p.opWrapper(op) })
		anyPipe.source = collections.IterableSlice(list...)
		return defultOpWrapper(down)
	}
	return anyPipe
}

func (p *FastPipline[T]) MapToString(mapper function.Func[T, string]) Stream[string] {
	helper.RequireCanButNonNil(mapper)
	strPipe := &FastPipline[string]{
		cancel:   p.cancel,
		parallel: p.parallel,
	}
	var list []string
	op := func(t T) {
		list = append(list, mapper(t))
	}
	strPipe.opWrapper = func(down function.Consumer[string]) function.Consumer[string] {
		p.exec(func() function.Consumer[T] { return p.opWrapper(op) })
		strPipe.source = collections.IterableSlice(list...)
		return defultOpWrapper(down)
	}
	return strPipe
}

func (p *FastPipline[T]) MapToInt(mapper function.Func[T, int]) Stream[int] {
	intPipe := &FastPipline[int]{
		cancel:   p.cancel,
		parallel: p.parallel,
	}
	intPipe.opWrapper = func(down function.Consumer[int]) function.Consumer[int] {
		var list []int
		op := func(t T) {
			list = append(list, mapper(t))
		}
		p.exec(func() function.Consumer[T] { return p.opWrapper(op) })
		intPipe.source = collections.IterableSlice(list...)
		return defultOpWrapper(down)
	}
	return intPipe
}

func (p *FastPipline[T]) MapToFloat(mapper function.Func[T, float64]) Stream[float64] {
	floatPipe := &FastPipline[float64]{
		cancel:   p.cancel,
		parallel: p.parallel,
	}
	floatPipe.opWrapper = func(down function.Consumer[float64]) function.Consumer[float64] {
		var list []float64
		op := func(t T) {
			list = append(list, mapper(t))
		}
		p.exec(func() function.Consumer[T] { return p.opWrapper(op) })
		floatPipe.source = collections.IterableSlice(list...)
		return defultOpWrapper(down)
	}
	return floatPipe
}

func (p *FastPipline[T]) AnyMatch(pred function.Predicate[T]) bool {
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

func (p *FastPipline[T]) AllMatch(pred function.Predicate[T]) bool {
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

func (p *FastPipline[T]) NoneMatch(pred function.Predicate[T]) bool {
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

func (p *FastPipline[T]) FindAny() optional.Value[T] {
	val := optional.EmptyVal[T]()
	p.ForEach(func(t T) {
		val = optional.ValOf(t)
		p.cancel()
	})
	return val
}
