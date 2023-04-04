package stream_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/go-park/stream"
	"github.com/go-park/stream/function"
	"gotest.tools/assert"
)

type TestData[T any] struct {
	name         string
	list         []T
	filter       function.Predicate[T]
	less         function.BiPredicate[T, T]
	max          function.BiPredicate[T, T]
	min          function.BiPredicate[T, T]
	mapper       function.Fn[T, T]
	reducer      function.BiFn[T, T, T]
	mapperAny    func(t T) any
	mapperString func(t T) string
	mapperInt    func(t T) int
	mapperFloat  func(t T) float64
	anyMatcher   function.Predicate[T]
	allMatcher   function.Predicate[T]
	noneMatcher  function.Predicate[T]
	limit        struct {
		num uint
		ok  bool
	}
	skip struct {
		num uint
		ok  bool
	}
	wantList       []T
	wantAnyList    []any
	wantStringList []string
	wantIntList    []int
	wantFloatList  []float64
	wantValue      T
	wantBool       bool
}

func TestSimplePipline(t *testing.T) {
	intTests := []TestData[int]{
		{
			name: "filter",
			list: []int{1, 2, 3, 4, 5},
			filter: func(i int) bool {
				return i > 2
			},
			wantList: []int{3, 4, 5},
		},
		{
			name: "filter",
			list: []int{},
			filter: func(i int) bool {
				return i < 1
			},
			wantList: nil,
		},
		{
			name: "sort",
			list: []int{2, 3, 4, 6, 1, 7, 0, 1, 8, 5, math.MaxInt, math.MinInt},
			less: func(i, j int) bool {
				return i < j
			},
			wantList: []int{math.MinInt, 0, 1, 1, 2, 3, 4, 5, 6, 7, 8, math.MaxInt},
		},
		{
			name: "sort",
			list: []int{2, 3, 4, 6, 1, 7, 0, 1, 8, 5, math.MaxInt, math.MinInt},
			less: func(i, j int) bool {
				return i > j
			},
			wantList: []int{math.MaxInt, 8, 7, 6, 5, 4, 3, 2, 1, 1, 0, math.MinInt},
		},
		{
			name: "limit",
			list: []int{2, 3, 4, 6, 1, 7, 0, 1, 8, 5, math.MaxInt, math.MinInt},
			limit: struct {
				num uint
				ok  bool
			}{num: 0, ok: true},
			wantList: nil,
		},
		{
			name: "limit",
			list: []int{2, 3, 4, 6, 1, 7, 0, 1, 8, 5, math.MaxInt, math.MinInt},
			limit: struct {
				num uint
				ok  bool
			}{num: 5, ok: true},
			wantList: []int{2, 3, 4, 6, 1},
		},
		{
			name: "skip",
			list: []int{2, 3, 4, 6, 1, 7, 0, 1, 8, 5, math.MaxInt, math.MinInt},
			skip: struct {
				num uint
				ok  bool
			}{num: 0, ok: true},
			wantList: []int{2, 3, 4, 6, 1, 7, 0, 1, 8, 5, math.MaxInt, math.MinInt},
		},
		{
			name: "skip",
			list: []int{2, 3, 4, 6, 1, 7, 0, 1, 8, 5, math.MaxInt, math.MinInt},
			skip: struct {
				num uint
				ok  bool
			}{num: 5, ok: true},
			wantList: []int{7, 0, 1, 8, 5, math.MaxInt, math.MinInt},
		},
		{
			name: "max",
			list: []int{2, 3, 4, 6, 1, 7, 0, 1, 8, 5, math.MaxInt, math.MinInt},
			max: func(t, u int) bool {
				return t < u
			},
			wantValue: math.MaxInt,
		},
		{
			name: "min",
			list: []int{2, 3, 4, 6, 1, 7, 0, 1, 8, 5, math.MaxInt, math.MinInt},
			min: func(t, u int) bool {
				return t < u
			},
			wantValue: math.MinInt,
		},
		{
			name: "map",
			list: []int{2},
			mapper: func(t int) int {
				return t + 1
			},
			wantList: []int{3},
		},
		{
			name: "map",
			list: []int{2, 3, 4, 6, 1, 7, 0, 1, 8, 5},
			mapper: func(t int) int {
				return t + 1
			},
			wantList: []int{3, 4, 5, 7, 2, 8, 1, 2, 9, 6},
		},
		{
			name: "reduce",
			list: []int{3},
			reducer: func(t, u int) int {
				return t + u
			},
			wantValue: 3,
		},
		{
			name: "reduce",
			list: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			reducer: func(t, u int) int {
				return t + u
			},
			wantValue: (1 + 10) * 10 / 2,
		},
		{
			name: "mapperAny",
			list: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			mapperAny: func(t int) any {
				return fmt.Sprintf("%d", t)
			},
			wantAnyList: []any{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
		},
		{
			name: "mapperString",
			list: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			mapperString: func(t int) string {
				return fmt.Sprintf("%d", t)
			},
			wantStringList: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
		},
		{
			name: "mapperInt",
			list: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			mapperInt: func(t int) int {
				return t
			},
			wantIntList: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name: "mapperFloat",
			list: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			mapperFloat: func(t int) float64 {
				return float64(t) / 10
			},
			wantFloatList: []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
		},
		{
			name: "anyMatch",
			list: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			anyMatcher: func(t int) bool {
				return t == 7
			},
			wantBool: true,
		},
		{
			name: "anyMatch",
			list: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			anyMatcher: func(t int) bool {
				return t < 1
			},
			wantBool: false,
		},
		{
			name: "allMatch",
			list: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			allMatcher: func(t int) bool {
				return t > 0
			},
			wantBool: true,
		},
		{
			name: "allMatch",
			list: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			allMatcher: func(t int) bool {
				return t > 1
			},
			wantBool: false,
		},
		{
			name: "noneMatch",
			list: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			noneMatcher: func(t int) bool {
				return t < 0
			},
			wantBool: true,
		},
		{
			name: "noneMatch",
			list: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			noneMatcher: func(t int) bool {
				return t == 1
			},
			wantBool: false,
		},
	}
	for _, v := range intTests {
		t.Run(v.name, func(t *testing.T) {
			pStream := stream.From(v.list...)
			if v.filter != nil {
				pStream = pStream.Filter(v.filter)
			}
			if v.less != nil {
				pStream = pStream.SortBy(v.less)
			}
			if v.limit.ok {
				pStream = pStream.Limit(v.limit.num)
			}
			if v.skip.ok {
				pStream = pStream.Skip(v.skip.num)
			}
			if v.max != nil {
				val := pStream.Max(v.max)
				assert.Equal(t, val.IsEmpty(), false)
				assert.Equal(t, val.IsNil(), false)
				assert.DeepEqual(t, val.Get(), v.wantValue)
			}
			if v.min != nil {
				val := pStream.Min(v.min)
				assert.Equal(t, val.IsEmpty(), false)
				assert.Equal(t, val.IsNil(), false)
				assert.DeepEqual(t, val.Get(), v.wantValue)
			}
			if v.mapper != nil {
				pStream = pStream.Map(v.mapper)
			}
			if v.reducer != nil {
				val := pStream.Reduce(v.reducer)
				assert.Equal(t, val.IsEmpty(), false)
				assert.Equal(t, val.IsNil(), false)
				assert.DeepEqual(t, val.Get(), v.wantValue)
			}
			if v.mapperAny != nil {
				anyStream := pStream.MapToAny(v.mapperAny)
				assert.DeepEqual(t, anyStream.ToSlice(), v.wantAnyList)
			}
			if v.mapperString != nil {
				strStream := pStream.MapToString(v.mapperString)
				assert.DeepEqual(t, strStream.ToSlice(), v.wantStringList)
			}
			if v.mapperInt != nil {
				intStream := pStream.MapToInt(v.mapperInt)
				assert.DeepEqual(t, intStream.ToSlice(), v.wantIntList)
			}
			if v.mapperFloat != nil {
				floatStream := pStream.MapToFloat(v.mapperFloat)
				assert.DeepEqual(t, floatStream.ToSlice(), v.wantFloatList)
			}
			if v.anyMatcher != nil {
				b := pStream.AnyMatch(v.anyMatcher)
				assert.Equal(t, b, v.wantBool)
			}
			if v.allMatcher != nil {
				b := pStream.AllMatch(v.allMatcher)
				assert.Equal(t, b, v.wantBool)
			}
			if v.noneMatcher != nil {
				b := pStream.NoneMatch(v.noneMatcher)
				assert.Equal(t, b, v.wantBool)
			}
			assert.DeepEqual(t, pStream.ToSlice(), v.wantList)
		})
	}
}

func TestXxx(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	v := stream.From(slice...).
		Filter(func(t int) bool { return t > 2 }).
		Skip(2).Limit(2).
		Map(func(i int) int {
			return i + 1
		}).
		Reduce(func(i1, i2 int) int {
			return i1 + i2
		})
	v.IfNotEmpty(func(t int) { println(t) })
}
