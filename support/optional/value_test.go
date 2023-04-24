package optional_test

import (
	"sync"
	"testing"

	"github.com/go-park/stream/support/optional"
	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	v := optional.EmptyVal[*sync.Mutex]()
	assert.Equal(t, true, v.IsEmpty())
	assert.Equal(t, true, v.IsNil())

	v = optional.ValOf(&sync.Mutex{})
	assert.Equal(t, false, v.IsEmpty())
	assert.Equal(t, false, v.IsNil())

	w := optional.ValOf(sync.Mutex{})
	assert.Equal(t, false, w.IsEmpty())
	assert.Equal(t, false, w.IsNil())

	w = optional.EmptyVal[sync.Mutex]()
	assert.Equal(t, true, w.IsEmpty())
	assert.Equal(t, false, w.IsNil())
}
