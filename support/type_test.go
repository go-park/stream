package support_test

import (
	"sync"
	"testing"

	"github.com/go-park/stream/support"
)

func TestType(t *testing.T) {
	v := support.ValOf(&sync.Mutex{})
	println(v.IsEmpty())
	println(v.IsNil())
}
