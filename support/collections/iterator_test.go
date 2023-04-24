package collections_test

import (
	"testing"

	"github.com/go-park/stream/support/collections"
	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	iter := collections.IterableSlice(list...)
	for _, v := range list {
		if iter.HasNext() {
			assert.Equal(t, v, iter.Next())
		}
	}
	assert.Equal(t, 0, iter.Next())
	assert.Equal(t, false, iter.HasNext())

	index := 0
	iter2 := collections.IterableSlice(list...)
	iter2.ForEachRemaining(func(v int) {
		assert.Equal(t, v, list[index])
		index++
	})
	assert.Equal(t, 0, iter2.Next())
	assert.Equal(t, false, iter2.HasNext())
}
