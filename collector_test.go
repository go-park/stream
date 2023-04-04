package stream_test

import (
	"testing"

	"github.com/go-park/stream"
	"gotest.tools/assert"
)

func TestMap(t *testing.T) {
	type P struct {
		Name string
		Age  int
	}
	tests := []struct {
		name     string
		kvs      map[string]P
		filter   func(s stream.Entry[string, P]) bool
		less     func(i, j stream.Entry[string, P]) bool
		wantLen  int
		wantList stream.EntrySet[string, P]
	}{
		{
			name: "filter by key",
			kvs: map[string]P{
				"bar": {
					Name: "bar",
					Age:  2,
				},
				"foo": {
					Name: "foo",
					Age:  1,
				},
			},
			filter: func(s stream.Entry[string, P]) bool {
				return s.Key() == "foo"
			},
			wantLen: 1,
			wantList: stream.GetEntrySet(map[string]P{
				"foo": {
					Name: "foo",
					Age:  1,
				},
			}),
		},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			kvStream := stream.FromMap(v.kvs)
			if v.filter != nil {
				kvStream = kvStream.Filter(v.filter)
			}
			list := kvStream.ToSlice()
			for i, vv := range v.wantList {
				assert.Equal(t, vv.Key(), list[i].Key())
				assert.DeepEqual(t, vv.Value(), list[i].Value())
			}
		})
	}
}

func TestToMap(t *testing.T) {
	type P struct {
		Name string
		Age  int
	}
	tests := []struct {
		name     string
		list     []P
		filter   func(s stream.Entry[string, P]) bool
		less     func(i, j stream.Entry[string, P]) bool
		key      func(p P) string
		value    func(p P) int
		wantLen  int
		wantList stream.EntrySet[string, P]
		wantMap  any
	}{
		{
			name: "toMap",
			list: []P{
				{
					Name: "bar",
					Age:  2,
				},
				{
					Name: "foo",
					Age:  1,
				},
			},
			key:     func(p P) string { return p.Name },
			value:   func(p P) int { return p.Age },
			wantLen: 1,
			wantMap: map[string]int{
				"foo": 1,
				"bar": 2,
			},
		},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			s := stream.From(v.list...)
			if v.key != nil && v.value != nil {
				hash := stream.ToMap(s, v.key, v.value)
				assert.DeepEqual(t, hash, v.wantMap)
			}
		})
	}
}
