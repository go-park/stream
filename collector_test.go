package stream_test

import (
	"testing"

	"github.com/go-park/stream"
	"github.com/go-park/stream/support/collections"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	type P struct {
		Name string
		Age  int
	}
	tests := []struct {
		name     string
		kvs      map[string]P
		filter   func(s collections.Entry[string, P]) bool
		less     func(i, j collections.Entry[string, P]) bool
		wantLen  int
		wantList collections.EntrySet[string, P]
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
			filter: func(s collections.Entry[string, P]) bool {
				return s.Key() == "foo"
			},
			wantLen: 1,
			wantList: collections.GetEntrySet(map[string]P{
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
				assert.Equal(t, list[i].Key(), vv.Key())
				assert.Equal(t, list[i].Value(), vv.Value())
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
		filter   func(s collections.Entry[string, P]) bool
		less     func(i, j collections.Entry[string, P]) bool
		key      func(p P) string
		value    func(p P) int
		wantLen  int
		wantList collections.EntrySet[string, P]
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
				assert.Equal(t, v.wantMap, hash)
			}
		})
	}
}

func TestConstraintOp(t *testing.T) {
	// distinct
	t.Run("distinct", func(t *testing.T) {
		list := []int{1, 2, 3, 1, 2, 3, 4, 4, 5, 6, 7, 8, 9, 9}
		s := stream.Distinct(stream.From(list...))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, s.ToSlice())
	})
	// sort
	t.Run("sort", func(t *testing.T) {
		list := []int{1, 4, 7, 2, 5, 8, 3, 6, 9}
		s := stream.Sort(stream.From(list...))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, s.ToSlice())
	})

	// reverse
	t.Run("reverse", func(t *testing.T) {
		list := []int{1, 4, 7, 2, 5, 8, 3, 6, 9}
		s := stream.Sort(stream.From(list...)).Reverse()
		assert.Equal(t, []int{9, 8, 7, 6, 5, 4, 3, 2, 1}, s.ToSlice())
	})

	// max
	t.Run("max", func(t *testing.T) {
		list := []int{1, 4, 7, 2, 5, 8, 3, 6, 9}
		s := stream.Max(stream.From(list...))
		assert.Equal(t, 9, s.Get())
	})

	// min
	t.Run("min", func(t *testing.T) {
		list := []int{1, 4, 7, 2, 5, 8, 3, 6, 9}
		s := stream.Min(stream.From(list...))
		assert.Equal(t, 1, s.Get())
	})
}
