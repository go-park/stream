package stream

import (
	"github.com/go-park/stream/function"
	"github.com/go-park/stream/internal/helper"
)

func ToList[T, V any](s Stream[T], conv function.Fn[T, V]) []V {
	helper.RequireNonNil(conv)
	list := make([]V, 0)
	s.ForEach(func(t T) {
		value := conv.Apply(t)
		list = append(list, value)
	})
	return list
}

func ToMap[T, V any, R comparable](s Stream[T], source function.Fn[T, R], after function.Fn[T, V]) map[R]V {
	helper.RequireNonNil(source)
	helper.RequireNonNil(after)
	hash := make(map[R]V)
	s.ForEach(func(t T) {
		key := source.Apply(t)
		value := after.Apply(t)
		hash[key] = value
	})
	return hash
}
