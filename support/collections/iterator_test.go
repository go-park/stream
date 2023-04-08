package collections_test

import (
	"testing"

	"github.com/go-park/stream/support/collections"
	"gotest.tools/assert"
)

func TestIterator(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	iter := collections.ToIterator(list...)
	for _, v := range list {
		if iter.HasNext() {
			assert.Equal(t, v, iter.Next())
		}
	}
	assert.Equal(t, iter.Next(), 0)
	assert.Equal(t, iter.HasNext(), false)

	index := 0
	iter2 := collections.ToIterator(list...)
	iter2.ForEachRemaining(func(v int) {
		assert.Equal(t, v, list[index])
		index++
	})
	assert.Equal(t, iter2.Next(), 0)
	assert.Equal(t, iter2.HasNext(), false)
}
