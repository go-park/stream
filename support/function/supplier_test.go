package function_test

import (
	"testing"

	"github.com/go-park/stream/support/function"
	"gotest.tools/assert"
)

func TestSupplier(t *testing.T) {
	fn := func() int {
		return 5
	}
	supplier := function.Supplier[int](fn)
	result := supplier.Get()
	assert.Equal(t, result, 5)
}
