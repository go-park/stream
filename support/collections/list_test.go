package collections_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/go-park/stream/support/collections"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	equal := func(u, v int) bool {
		return u == v
	}
	t.Run("foreach", func(t *testing.T) {
		index := 0
		collections.ForEach(list, func(v int) {
			assert.Equal(t, list[index], v)
			index++
		})
	})
	t.Run("filter", func(t *testing.T) {
		actual := collections.Filter(list, func(v int) bool {
			return v > 6
		})
		expected := []int{7, 8, 9}
		assert.Equal(t, expected, actual)
	})
	t.Run("has-one", func(t *testing.T) {
		hasTrue := collections.Contains(list, 7)
		hasFalse := collections.Contains(list, 10)
		assert.True(t, hasTrue)
		assert.False(t, hasFalse)

		hasTrue = collections.ConstainsF(list, 7, equal)
		hasFalse = collections.ConstainsF(list, 0, equal)
		assert.True(t, hasTrue)
		assert.False(t, hasFalse)
	})
	t.Run("has-any", func(t *testing.T) {
		hasAnyTrue := collections.ContainsAny(list, 1, 10, 11, 12)
		hasAnyFalse := collections.ContainsAny(list, 10, 11, 12)
		assert.True(t, hasAnyTrue)
		assert.False(t, hasAnyFalse)

		hasAnyTrue = collections.ContainsAnyF(list, equal, 1, 10, 11, 12)
		hasAnyFalse = collections.ContainsAnyF(list, equal, 10, 11, 12)
		assert.True(t, hasAnyTrue)
		assert.False(t, hasAnyFalse)
	})
	t.Run("sort", func(t *testing.T) {
		shuffle := []int{6, 7, 8, 9, 1, 5, 2, 3, 4}
		collections.Sort(shuffle, func(i, j int) bool {
			return i < j
		})
		assert.Equal(t, list, shuffle)
	})
	t.Run("ascending-order", func(t *testing.T) {
		cl := make([]int, len(list))
		copy(cl, list)
		rand.Shuffle(9, func(i, j int) {
			cl[i], cl[j] = cl[j], cl[i]
		})
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		assert.NotEqual(t, expected, cl)
		collections.Asc(cl)
		assert.Equal(t, expected, cl)
	})

	t.Run("descending-order", func(t *testing.T) {
		cl := make([]int, len(list))
		copy(cl, list)
		collections.Desc(cl)
		expected := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
		assert.Equal(t, expected, cl)
	})
	t.Run("distinct", func(t *testing.T) {
		repeated := []int{1, 1, 2, 3, 4, 5, 6, 7, 8, 8, 9}
		actual := collections.Distinct(repeated)
		assert.Equal(t, list, actual)

		repeated = []int{1, 1, 2, 3, 4, 5, 6, 7, 8, 8, 9}
		actual = collections.DistinctF(repeated, equal)
		assert.Equal(t, list, actual)
	})
	t.Run("map", func(t *testing.T) {
		expected := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
		out := collections.Map(list, func(i int) string {
			return strconv.Itoa(i)
		})
		assert.Equal(t, expected, out)
	})
}
