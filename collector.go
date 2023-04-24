package stream

import (
	"github.com/go-park/stream/internal/helper"
	"github.com/go-park/stream/support/function"
	"github.com/go-park/stream/support/optional"
	"golang.org/x/exp/constraints"
)

func ToList[T, V any](s Stream[T], conv function.Func[T, V]) []V {
	helper.RequireCanButNonNil(conv)
	list := make([]V, 0)
	s.ForEach(func(t T) {
		value := conv.Apply(t)
		list = append(list, value)
	})
	return list
}

func ToMap[T, V any, R comparable](s Stream[T], source function.Func[T, R], after function.Func[T, V]) map[R]V {
	helper.RequireCanButNonNil(source)
	helper.RequireCanButNonNil(after)
	hash := make(map[R]V)
	s.ForEach(func(t T) {
		key := source.Apply(t)
		value := after.Apply(t)
		hash[key] = value
	})
	return hash
}

func Distinct[T comparable](s Stream[T]) Stream[T] {
	helper.RequireCanButNonNil(s)
	s = s.Distinct(func(t, u T) bool {
		return t == u
	})
	return s
}

func Sort[T constraints.Ordered](s Stream[T]) Stream[T] {
	helper.RequireCanButNonNil(s)
	return s.Sort(func(t, u T) bool { return t < u })
}

func Max[T constraints.Ordered](s Stream[T]) optional.Value[T] {
	helper.RequireCanButNonNil(s)
	return s.Max(func(t, u T) bool { return t < u })
}

func Min[T constraints.Ordered](s Stream[T]) optional.Value[T] {
	helper.RequireCanButNonNil(s)
	return s.Min(func(t, u T) bool { return t < u })
}
