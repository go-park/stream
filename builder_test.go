package stream_test

import (
	"testing"

	"github.com/go-park/stream"
	"gotest.tools/assert"
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
			bs := stream.Builder[P]().Append(v.list...).Build()
			ls := stream.OfList(v.list)
			os := stream.Of(v.list...)
			assert.DeepEqual(t, bs.List(), ls.List())
			assert.DeepEqual(t, bs.List(), os.List())
		})
	}
}
