package function_test

import (
	"testing"

	"github.com/go-park/stream/function"
)

func TestSupplier(t *testing.T) {
	fn := func() int {
		return 5
	}
	supplier := function.Supplier[int](fn)
	result := supplier.Get()
	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}
}
