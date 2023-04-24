package stream_test

import (
	"testing"

	"github.com/go-park/stream"
	"github.com/stretchr/testify/assert"
)

type P struct {
	Name string
	Age  int
}

func TestBuilder(t *testing.T) {
	tests := []struct {
		name     string
		list     []P
		kvs      map[string]P
		wantList []P
	}{
		{
			name: "builder",
			list: []P{
				{
					Name: "foo",
					Age:  1,
				},
				{
					Name: "bar",
					Age:  2,
				},
			},
			wantList: []P{
				{
					Name: "foo",
					Age:  1,
				},
				{
					Name: "bar",
					Age:  2,
				},
			},
		},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
		})
	}
}

func TestRange(t *testing.T) {
	t.Run("0-99", func(t *testing.T) {
		var slice []int
		for i := 0; i < 100; i++ {
			slice = append(slice, i)
		}
		assert.Equal(t, stream.Range(0, 99).ToSlice(), slice)
	})
}
