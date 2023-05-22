package collections

import (
	"github.com/go-park/stream/internal/helper"
	"github.com/go-park/stream/support/function"
)

var (
	_ Set[int] = (*set[map[int]helper.Empty, int])(nil)
	_ Set[any] = (*xet[[]any, any, int])(nil)
)

type (
	Set[T any] interface {
		// Elements returns a slice containing all elements in the set.
		Elements() []T
		// Size returns the number of elements in the set.
		Size() int
		// Clear removes all elements from the set.
		Clear()
		// Add adds an element to the set.
		Add(element ...T) Set[T]
		// Remove removes an element from the set.
		Remove(element ...T) Set[T]
		// Contains checks if the set contains a specific element.
		Contains(element T) bool
		// ContainsAny checks if the set contains any of the specified values.
		ContainsAny(values ...T) bool
		// ContainsAll checks if the set contains all of the specified values.
		ContainsAll(values ...T) bool
		// Intersection returns a new set containing the elements that are common to both the current set and the input set `ts`.
		Intersection(ts Set[T]) Set[T]
		// Union returns a new set containing all the elements that are present in both the current set and the input set `ts`.
		Union(ts Set[T]) Set[T]
		// Difference  returns a new set containing the elements that are present in the current set but not in the input set `ts`.
		Difference(ts Set[T]) Set[T]
		// SymmetricDifference returns a new set containing the elements that are present in either the current set or the input set `ts`, but not in both.
		SymmetricDifference(ts Set[T]) Set[T]
	}
	set[M map[V]helper.Empty, V comparable] struct {
		elements M
	}
)

func NewSet[V comparable](values ...V) Set[V] {
	res := &set[map[V]helper.Empty, V]{}
	elements := make(map[V]helper.Empty, len(values))
	for _, v := range values {
		elements[v] = helper.Empty{}
	}
	res.elements = elements
	return res
}

func toSet[M map[K]helper.Empty, K comparable](m M) Set[K] {
	return &set[M, K]{elements: m}
}

func (s *set[L, V]) Elements() []V {
	return Keys(s.elements)
}

func (s *set[L, V]) Size() int {
	return len(s.elements)
}

func (s *set[L, V]) Clear() {
	for v := range s.elements {
		delete(s.elements, v)
	}
}

func (s *set[L, V]) Add(values ...V) Set[V] {
	for _, v := range values {
		s.elements[v] = helper.Empty{}
	}
	return s
}

func (s *set[L, V]) Remove(values ...V) Set[V] {
	for _, v := range values {
		delete(s.elements, v)
	}
	return s
}

func (s *set[L, V]) Contains(v V) bool {
	_, exists := s.elements[v]
	return exists
}

func (s *set[L, V]) ContainsAny(values ...V) bool {
	for _, v := range values {
		if _, exists := s.elements[v]; exists {
			return true
		}
	}
	return false
}

func (s *set[L, V]) ContainsAll(values ...V) bool {
	for _, v := range values {
		if _, exists := s.elements[v]; !exists {
			return false
		}
	}
	return true
}

func (s set[M, V]) Intersection(ts Set[V]) Set[V] {
	res := NewSet[V]()
	for v := range s.elements {
		if ts.Contains(v) {
			res.Add(v)
		}
	}
	return res
}

func (s *set[M, V]) Union(ts Set[V]) Set[V] {
	res := NewSet[V]()
	for v := range s.elements {
		res.Add(v)
	}
	for _, v := range ts.Elements() {
		res.Add(v)
	}
	return res
}

func (s *set[L, V]) Difference(ts Set[V]) Set[V] {
	res := NewSet[V]()
	for v := range s.elements {
		if !ts.Contains(v) {
			res.Add(v)
		}
	}
	return res
}

func (s *set[L, V]) SymmetricDifference(ts Set[V]) Set[V] {
	return s.Difference(ts).Union(ts.Difference(s))
}

type xet[L []I, I any, O comparable] struct {
	elements L
	indicats map[O]int
	mapper   function.Func[I, O]
}

func NewXet[I any, O comparable](mapper function.Func[I, O], values ...I) Set[I] {
	helper.RequireCanButNonNil(mapper)
	s := &xet[[]I, I, O]{}
	s.mapper = mapper
	s.indicats = make(map[O]int)
	for _, v := range values {
		s.Add(v)
	}
	return s
}

func (s *xet[L, I, _]) Elements() []I {
	return s.elements
}

func (s *xet[L, I, _]) Size() int {
	return len(s.elements)
}

func (s *xet[L, I, _]) Clear() {
	for v := range s.indicats {
		delete(s.indicats, v)
	}
	s.elements = s.elements[:0]
}

func (s *xet[L, I, _]) Add(values ...I) Set[I] {
	for _, v := range values {
		o := s.mapper.Apply(v)
		if _, exists := s.indicats[o]; !exists {
			index := len(s.elements)
			s.elements = append(s.elements, v)
			s.indicats[o] = index
		}
	}
	return s
}

func (s *xet[L, I, _]) Remove(values ...I) Set[I] {
	for _, v := range values {
		o := s.mapper.Apply(v)
		if index, exists := s.indicats[o]; exists {
			lastIndex := len(s.elements) - 1
			if index != lastIndex {
				tail := s.elements[lastIndex]
				key := s.mapper.Apply(tail)
				s.elements[index] = tail
				s.indicats[key] = index
			}
			s.elements = s.elements[:lastIndex]
			delete(s.indicats, o)
		}
	}
	return s
}

func (s *xet[L, I, _]) Contains(v I) bool {
	o := s.mapper.Apply(v)
	_, exists := s.indicats[o]
	return exists
}

func (s *xet[L, I, _]) ContainsAny(values ...I) bool {
	for _, v := range values {
		o := s.mapper.Apply(v)
		if _, exists := s.indicats[o]; exists {
			return true
		}
	}
	return false
}

func (s *xet[L, I, _]) ContainsAll(values ...I) bool {
	for _, v := range values {
		o := s.mapper.Apply(v)
		if _, exists := s.indicats[o]; !exists {
			return false
		}
	}
	return true
}

func (s *xet[L, I, _]) Intersection(ts Set[I]) Set[I] {
	res := NewXet(s.mapper)
	for _, v := range s.elements {
		if ts.Contains(v) {
			res.Add(v)
		}
	}
	return res
}

func (s *xet[L, I, _]) Union(ts Set[I]) Set[I] {
	res := NewXet(s.mapper)
	for _, v := range s.elements {
		res.Add(v)
	}
	for _, v := range ts.Elements() {
		res.Add(v)
	}
	return res
}

func (s *xet[L, I, _]) Difference(ts Set[I]) Set[I] {
	res := NewXet(s.mapper)
	for _, v := range s.elements {
		if !ts.Contains(v) {
			res.Add(v)
		}
	}
	return res
}

func (s *xet[L, I, _]) SymmetricDifference(ts Set[I]) Set[I] {
	return s.Difference(ts).Union(ts.Difference(s))
}
