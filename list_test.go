package stream_test

import (
	"testing"

	"github.com/go-park/stream"
	"gotest.tools/assert"
)

func TestStream(t *testing.T) {
	tests := []struct {
		name     string
		list     []P
		filter   func(p P) bool
		less     func(i, j P) bool
		foreach  func(p P)
		limit    int
		wantLen  int
		wantList []P
	}{
		{
			name: "filter/list/len",
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
			filter:  func(p P) bool { return p.Age < 2 },
			wantLen: 1,
			wantList: []P{
				{
					Name: "foo",
					Age:  1,
				},
			},
		},
		{
			name: "sort",
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
			less: func(i, j P) bool {
				return i.Age < j.Age
			},
			wantLen: 2,
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
		{
			name: "limit",
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
			limit:   1,
			wantLen: 1,
			wantList: []P{
				{
					Name: "bar",
					Age:  2,
				},
			},
		},
		{
			name: "limit",
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
			limit:   3,
			wantLen: 2,
			wantList: []P{
				{
					Name: "bar",
					Age:  2,
				},
				{
					Name: "foo",
					Age:  1,
				},
			},
		},
		{
			name: "foreach",
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
			foreach: func(p P) {
				p.Age = 100
			},
			limit:   3,
			wantLen: 2,
			wantList: []P{
				{
					Name: "bar",
					Age:  2,
				},
				{
					Name: "foo",
					Age:  1,
				},
			},
		},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			pStream := stream.OfList(v.list)
			if v.filter != nil {
				pStream = pStream.Filter(v.filter)
			}
			if v.less != nil {
				pStream = pStream.Sort(v.less)
			}
			if v.limit > 0 {
				pStream = pStream.Limit(v.limit)
			}
			if v.foreach != nil {
				_ = pStream.ForEach(v.foreach)
			}
			assert.Equal(t, pStream.Len(), len(pStream.List()))
			assert.Equal(t, pStream.Len(), v.wantLen)
			assert.DeepEqual(t, pStream.List(), v.wantList)
		})
	}
}

func TestStreamP(t *testing.T) {
	tests := []struct {
		name     string
		list     []*P
		filter   func(p *P) bool
		less     func(i, j *P) bool
		foreach  func(p *P)
		limit    int
		wantLen  int
		wantList []*P
	}{
		{
			name: "filter/list/len",
			list: []*P{
				{
					Name: "foo",
					Age:  1,
				},
				{
					Name: "bar",
					Age:  2,
				},
			},
			filter:  func(p *P) bool { return p.Age < 2 },
			wantLen: 1,
			wantList: []*P{
				{
					Name: "foo",
					Age:  1,
				},
			},
		},
		{
			name: "sort",
			list: []*P{
				{
					Name: "bar",
					Age:  2,
				},
				{
					Name: "foo",
					Age:  1,
				},
			},
			less: func(i, j *P) bool {
				return i.Age < j.Age
			},
			wantLen: 2,
			wantList: []*P{
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
		{
			name: "limit",
			list: []*P{
				{
					Name: "bar",
					Age:  2,
				},
				{
					Name: "foo",
					Age:  1,
				},
			},
			limit:   1,
			wantLen: 1,
			wantList: []*P{
				{
					Name: "bar",
					Age:  2,
				},
			},
		},
		{
			name: "limit",
			list: []*P{
				{
					Name: "bar",
					Age:  2,
				},
				{
					Name: "foo",
					Age:  1,
				},
			},
			limit:   3,
			wantLen: 2,
			wantList: []*P{
				{
					Name: "bar",
					Age:  2,
				},
				{
					Name: "foo",
					Age:  1,
				},
			},
		},
		{
			name: "foreach",
			list: []*P{
				{
					Name: "bar",
					Age:  2,
				},
				{
					Name: "foo",
					Age:  1,
				},
			},
			foreach: func(p *P) {
				p.Age = 100
			},
			limit:   3,
			wantLen: 2,
			wantList: []*P{
				{
					Name: "bar",
					Age:  100,
				},
				{
					Name: "foo",
					Age:  100,
				},
			},
		},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			pStream := stream.OfList(v.list)
			if v.filter != nil {
				pStream = pStream.Filter(v.filter)
			}
			if v.less != nil {
				pStream = pStream.Sort(v.less)
			}
			if v.limit > 0 {
				pStream = pStream.Limit(v.limit)
			}
			if v.foreach != nil {
				_ = pStream.ForEach(v.foreach)
			}
			assert.Equal(t, pStream.Len(), len(pStream.List()))
			assert.Equal(t, pStream.Len(), v.wantLen)
			assert.DeepEqual(t, pStream.List(), v.wantList)
		})
	}
}
