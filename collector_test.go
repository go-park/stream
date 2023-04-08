package stream_test

import (
	"testing"

	"github.com/go-park/stream"
	"github.com/go-park/stream/support/collections"
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
				assert.DeepEqual(t, hash, v.wantMap)
			}
		})
	}
}

func TestConstraintOp(t *testing.T) {
	// distinct
	t.Run("distinct", func(t *testing.T) {
		list := []int{1, 2, 3, 1, 2, 3, 4, 4, 5, 6, 7, 8, 9, 9}
		s := stream.Distinct(stream.From(list...))
		assert.DeepEqual(t, s.ToSlice(), []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	})
	// sort
	t.Run("sort", func(t *testing.T) {
		list := []int{1, 4, 7, 2, 5, 8, 3, 6, 9}
		s := stream.Sort(stream.From(list...))
		assert.DeepEqual(t, s.ToSlice(), []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	})

	// reverse
	t.Run("reverse", func(t *testing.T) {
		list := []int{1, 4, 7, 2, 5, 8, 3, 6, 9}
		s := stream.Sort(stream.From(list...)).Reverse()
		assert.DeepEqual(t, s.ToSlice(), []int{9, 8, 7, 6, 5, 4, 3, 2, 1})
	})

	// max
	t.Run("max", func(t *testing.T) {
		list := []int{1, 4, 7, 2, 5, 8, 3, 6, 9}
		s := stream.Max(stream.From(list...))
		assert.DeepEqual(t, s.Get(), 9)
	})

	// min
	t.Run("min", func(t *testing.T) {
		list := []int{1, 4, 7, 2, 5, 8, 3, 6, 9}
		s := stream.Min(stream.From(list...))
		assert.DeepEqual(t, s.Get(), 1)
	})
}
