package stream_test

import (
	"testing"
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
