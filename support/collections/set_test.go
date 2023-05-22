package collections_test

import (
	"testing"

	"github.com/go-park/stream/support/collections"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	t.Run("clear", func(t *testing.T) {
		set := collections.NewSet(list...)
		set.Clear()
		actual := set.Elements()
		expected := []int{}
		assert.Equal(t, expected, actual)
		assert.Equal(t, len(expected), set.Size())
	})
	t.Run("add", func(t *testing.T) {
		set := collections.NewSet(list...)
		set.Add(1)
		set.Add(0)
		assert.Equal(t, 10, set.Size())
		expected := append([]int{0}, list...)
		actual := set.Elements()
		collections.Asc(actual)
		assert.Equal(t, expected, actual)
	})
	t.Run("remove", func(t *testing.T) {
		set := collections.NewSet(list...)
		set.Remove(8)
		assert.Equal(t, 8, set.Size())
		actual := set.Elements()
		collections.Asc(actual)
		expected := append(list[:7:7], list[8:]...)
		assert.Equal(t, expected, actual)
	})
	t.Run("contains", func(t *testing.T) {
		set := collections.NewSet(list...)
		assert.True(t, set.Contains(8))
		set.Remove(8)
		assert.False(t, set.Contains(8))
	})
	t.Run("containsAny", func(t *testing.T) {
		set := collections.NewSet(list...)
		assert.True(t, set.ContainsAny(8, 10))
		set.Remove(8)
		assert.False(t, set.ContainsAny(8, 10))
	})
	t.Run("containsAll", func(t *testing.T) {
		set := collections.NewSet(list...)
		assert.True(t, set.ContainsAll(8, 7))
		set.Remove(8)
		assert.False(t, set.ContainsAll(8, 7))
	})
	t.Run("intersection", func(t *testing.T) {
		set1 := collections.NewSet(list...)
		set2 := collections.NewSet(0, 1, 2)
		set := set1.Intersection(set2)
		expected := []int{1, 2}
		actual := set.Elements()
		collections.Asc(actual)
		assert.Equal(t, len(expected), set.Size())
		assert.Equal(t, expected, actual)
	})
	t.Run("union", func(t *testing.T) {
		set1 := collections.NewSet(list...)
		set2 := collections.NewSet(0, 1, 2)
		set := set1.Union(set2)
		expected := append([]int{0}, list...)
		actual := set.Elements()
		collections.Asc(actual)
		assert.Equal(t, len(expected), set.Size())
		assert.Equal(t, expected, actual)
	})
	t.Run("difference", func(t *testing.T) {
		set1 := collections.NewSet(list...)
		set2 := collections.NewSet(0, 1, 2)
		set := set1.Difference(set2)
		expected := list[2:]
		actual := set.Elements()
		collections.Asc(actual)
		assert.Equal(t, len(expected), set.Size())
		assert.Equal(t, expected, actual)
	})
	t.Run("symmetricDifference", func(t *testing.T) {
		set1 := collections.NewSet(list...)
		set2 := collections.NewSet(0, 1, 2)
		set := set1.SymmetricDifference(set2)
		expected := append([]int{0}, list[2:]...)
		actual := set.Elements()
		collections.Asc(actual)
		assert.Equal(t, len(expected), set.Size())
		assert.Equal(t, expected, actual)
	})
}

func TestXet(t *testing.T) {
	type sint struct{ i int }
	list := []sint{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}}
	mapper := func(s sint) int { return s.i }
	less := func(u, v sint) bool { return u.i < v.i }
	t.Run("clear", func(t *testing.T) {
		set := collections.NewXet(mapper, list...)
		set.Clear()
		actual := set.Elements()
		expected := []sint{}
		assert.Equal(t, expected, actual)
		assert.Equal(t, len(expected), set.Size())
	})
	t.Run("add", func(t *testing.T) {
		set := collections.NewXet(mapper, list...)
		set.Add(sint{1})
		set.Add(sint{0})
		assert.Equal(t, 10, set.Size())
		expected := append([]sint{{0}}, list...)
		actual := set.Elements()
		collections.Sort(actual, less)
		assert.Equal(t, expected, actual)
	})
	t.Run("remove", func(t *testing.T) {
		set := collections.NewXet(mapper, list...)
		set.Remove(sint{8})
		assert.Equal(t, 8, set.Size())
		actual := set.Elements()
		collections.Sort(actual, less)
		expected := append(list[:7:7], list[8:]...)
		assert.Equal(t, expected, actual)
	})
	t.Run("contains", func(t *testing.T) {
		set := collections.NewXet(mapper, list...)
		assert.True(t, set.Contains(sint{8}))
		set.Remove(sint{8})
		assert.False(t, set.Contains(sint{8}))
	})
	t.Run("containsAny", func(t *testing.T) {
		set := collections.NewXet(mapper, list...)
		assert.True(t, set.ContainsAny(sint{8}, sint{10}))
		set.Remove(sint{8})
		assert.False(t, set.ContainsAny(sint{8}, sint{10}))
	})
	t.Run("containsAll", func(t *testing.T) {
		set := collections.NewXet(mapper, list...)
		assert.True(t, set.ContainsAll(sint{8}, sint{7}))
		set.Remove(sint{8})
		assert.False(t, set.ContainsAll(sint{8}, sint{7}))
	})
	t.Run("intersection", func(t *testing.T) {
		set1 := collections.NewXet(mapper, list...)
		set2 := collections.NewXet(mapper, sint{0}, sint{1}, sint{2})
		set := set1.Intersection(set2)
		expected := []sint{{1}, {2}}
		assert.Equal(t, len(expected), set.Size())
		assert.Equal(t, expected, set.Elements())
	})
	t.Run("union", func(t *testing.T) {
		set1 := collections.NewSet(list...)
		set2 := collections.NewSet(sint{0}, sint{1}, sint{2})
		set := set1.Union(set2)
		expected := append([]sint{{0}}, list...)
		actual := set.Elements()
		collections.Sort(actual, less)
		assert.Equal(t, len(expected), set.Size())
		assert.Equal(t, expected, actual)
	})
	t.Run("difference", func(t *testing.T) {
		set1 := collections.NewSet(list...)
		set2 := collections.NewSet(sint{0}, sint{1}, sint{2})
		set := set1.Difference(set2)
		expected := list[2:]
		actual := set.Elements()
		collections.Sort(actual, less)
		assert.Equal(t, len(expected), set.Size())
		assert.Equal(t, expected, actual)
	})
	t.Run("symmetricDifference", func(t *testing.T) {
		set1 := collections.NewSet(list...)
		set2 := collections.NewSet(sint{0}, sint{1}, sint{2})
		set := set1.SymmetricDifference(set2)
		expected := append([]sint{{0}}, list[2:]...)
		actual := set.Elements()
		collections.Sort(actual, less)
		assert.Equal(t, len(expected), set.Size())
		assert.Equal(t, expected, actual)
	})
}
