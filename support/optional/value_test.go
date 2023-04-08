package optional_test

import (
	"sync"
	"testing"

	"github.com/go-park/stream/support/optional"
	"gotest.tools/assert"
)

func TestValue(t *testing.T) {
	v := optional.EmptyVal[*sync.Mutex]()
	assert.Equal(t, v.IsEmpty(), true)
	assert.Equal(t, v.IsNil(), true)

	v = optional.ValOf(&sync.Mutex{})
	assert.Equal(t, v.IsEmpty(), false)
	assert.Equal(t, v.IsNil(), false)

	w := optional.ValOf(sync.Mutex{})
	assert.Equal(t, w.IsEmpty(), false)
	assert.Equal(t, w.IsNil(), false)

	w = optional.EmptyVal[sync.Mutex]()
	assert.Equal(t, w.IsEmpty(), true)
	assert.Equal(t, w.IsNil(), false)
}
