package collections

import (
	"sort"

	"github.com/go-park/stream/internal/helper"
	"github.com/go-park/stream/support/function"
	"golang.org/x/exp/constraints"
)

type lessSwap[L ~[]V, V any] struct {
	l    L
	less func(i V, j V) bool
}

// Len is the number of elements in the collection.
func (ls lessSwap[_, _]) Len() int {
	return len(ls.l)
}

// Less reports whether the element with index i
// must sort before the element with index j.
func (ls lessSwap[_, _]) Less(i int, j int) bool {
	return ls.less(ls.l[i], ls.l[j])
}

// Swap swaps the elements with indexes i and j.
func (ls lessSwap[_, _]) Swap(i int, j int) {
	ls.l[i], ls.l[j] = ls.l[j], ls.l[i]
}

func orderedLessSwap[L ~[]V, V constraints.Ordered](list L) lessSwapOrdered[L, V] {
	return lessSwapOrdered[L, V]{lessSwap[L, V]{l: list}}
}

type lessSwapOrdered[L ~[]V, V constraints.Ordered] struct {
	lessSwap[L, V]
}

func (ls lessSwapOrdered[_, _]) Less(i int, j int) bool {
	return ls.l[i] < ls.l[j]
}

// ForEach performing a loop over the elements of the slice and applying a function to each element.
func ForEach[L ~[]V, V any](list L, consumer function.Consumer[V]) {
	for _, v := range list {
		consumer.Accept(v)
	}
}

// Filter returns a new slice of type `L` containing only the elements of `list` for which `pred` returns `true`.
func Filter[L ~[]V, V any](list L, pred function.Predicate[V]) L {
	res := make(L, 0, len(list))
	for _, v := range list {
		if pred.Test(v) {
			res = append(res, v)
		}
	}
	return res
}

// Contains returns a boolean value indicating whether `v` is present in the slice `list`.
func Contains[L ~[]V, V comparable](list L, v V) bool {
	for _, u := range list {
		if u == v {
			return true
		}
	}
	return false
}

// ConstainsF returns a boolean value indicating whether there is an element in `list` for which `comp` returns `true` when called with `v` and that element as arguments.
func ConstainsF[L ~[]V, V any](list L, v V, comp function.BiPredicate[V, V]) bool {
	for _, u := range list {
		if comp.Test(v, u) {
			return true
		}
	}
	return false
}

// ContainsAny returns a boolean value indicating whether any of the values passed as arguments are present in the slice `list`.
func ContainsAny[L ~[]V, V comparable](list L, values ...V) bool {
	res := make(map[V]helper.Empty, len(list))
	for _, v := range list {
		res[v] = helper.Empty{}
	}
	for _, v := range values {
		if _, ok := res[v]; ok {
			return true
		}
	}
	return false
}

// ContainsAnyF returns a boolean value indicating whether there is an element in `list` for which
// `comp` returns `true` when called with `v` and that element as arguments, where `v` is any of the
// values passed as arguments in `values`.
func ContainsAnyF[L ~[]V, V any](list L, comp function.BiPredicate[V, V], values ...V) bool {
	for _, v := range values {
		if ConstainsF(list, v, comp) {
			return true
		}
	}
	return false
}

// Sort is sorting the elements of a slice `list` of type `L` containing elements of
// type `V`, where `V` can be any type.
func Sort[L ~[]V, V any](list L, less func(i, j V) bool) {
	ls := lessSwap[L, V]{l: list, less: less}
	sort.Sort(ls)
}

// Asc is sorting the elements of a slice `list` in ascending order.
func Asc[L ~[]V, V constraints.Ordered](list L) {
	ls := orderedLessSwap(list)
	sort.Sort(ls)
}

// Desc is sorting the elements of a slice `list` in descending order.
func Desc[L ~[]V, V constraints.Ordered](list L) {
	ls := orderedLessSwap(list)
	sort.Sort(sort.Reverse(ls))
}

// Distinct returns a new slice of type `L` containing only the distinct elements of `list`.
func Distinct[L ~[]V, V comparable](list L) L {
	if len(list) <= 1 {
		return list
	}
	res := make(L, 0, len(list))
	m := make(map[V]helper.Empty, len(list))
	for _, v := range list {
		if _, ok := m[v]; !ok {
			m[v] = helper.Empty{}
			res = append(res, v)
		}
	}
	return res
}

// DistinctF returns a new slice of type `L` containing only the distinct
// elements of `list`, where distinct means that no two elements are equal according to the `comp` function.
func DistinctF[L ~[]V, V any](list L, comp function.BiPredicate[V, V]) L {
	if len(list) <= 1 {
		return list
	}
	res := make(L, 0, len(list))
	for _, v := range list {
		if !ConstainsF(res, v, comp) {
			res = append(res, v)
		}
	}
	return res
}

// Map returns a new
// slice of type `LO` containing the result of applying the `mapper` function to each element of the
// input slice `in`.
func Map[LI ~[]I, LO []O, I, O any](in LI, mapper function.Func[I, O]) LO {
	out := make(LO, 0, len(in))
	for _, v := range in {
		out = append(out, mapper.Apply(v))
	}
	return out
}
