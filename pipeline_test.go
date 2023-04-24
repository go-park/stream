package stream_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/go-park/stream"
	"github.com/go-park/stream/support/collections"
	"github.com/go-park/stream/support/function"
	"github.com/stretchr/testify/assert"
)

type TestData[T any] struct {
	name         string
	list         []T
	count        bool
	filter       function.Predicate[T]
	max          function.BiPredicate[T, T]
	min          function.BiPredicate[T, T]
	equals       function.BiPredicate[T, T]
	less         function.BiPredicate[T, T]
	reverse      bool
	mapper       function.Func[T, T]
	reducer      function.BiFunc[T, T, T]
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
	wantCount      int
	wantList       []T
	wantAnyList    []any
	wantStringList []string
	wantIntList    []int
	wantFloatList  []float64
	wantValue      T
	wantBool       bool
}

func TestDemo(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 61, 7, 8, 9, 10, 11, 19}
	val := stream.From(slice...).
		Filter(func(t int) bool { return t > 2 }).
		Skip(2).Limit(2).
		Map(func(i int) int {
			return i + 1
		}).
		Reduce(func(i1, i2 int) int { return i1 + i2 })
	val.IfNotEmptyOrElse(
		func(v int) { assert.Equal(t, v, 5+1+61+1) },
		func() { t.Error("empty") })
}

func TestPipline(t *testing.T) {
	t.Run("fast-sequential", func(t *testing.T) {
		testPipline(t, false, false)
	})
	t.Run("simple-sequential", func(t *testing.T) {
		testPipline(t, true, false)
	})
	t.Run("simple-parallel", func(t *testing.T) {
		testPipline(t, true, true)
	})
}

func testPipline(t *testing.T, simple, parallel bool) {
	intTests := []TestData[int]{
		{
			name:      "count",
			list:      []int{1, 2, 3, 4, 5},
			count:     true,
			wantCount: 5,
			wantList:  nil,
		},
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
			name: "distinct",
			list: []int{2, 2, 3, 2, 4, 2, 6, 1, 7, 0, 1, 8, 5, math.MaxInt, math.MinInt},
			equals: func(i, j int) bool {
				return i == j
			},
			wantList: []int{2, 3, 4, 6, 1, 7, 0, 8, 5, math.MaxInt, math.MinInt},
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
			name:     "reverse",
			list:     []int{2, 3, 4, 6, 1, 7, 0, 1, 8, 5, math.MaxInt, math.MinInt},
			reverse:  true,
			wantList: []int{math.MinInt, math.MaxInt, 5, 8, 1, 0, 7, 1, 6, 4, 3, 2},
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
			less: func(i, j int) bool {
				return i < j
			},
			wantList: []int{1, 2, 2, 3, 4, 5, 6, 7, 8, 9},
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
			wantAnyList: []any{"1", "10", "2", "3", "4", "5", "6", "7", "8", "9"},
		},
		{
			name: "mapperString",
			list: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			mapperString: func(t int) string {
				return fmt.Sprintf("%d", t)
			},
			wantStringList: []string{"1", "10", "2", "3", "4", "5", "6", "7", "8", "9"},
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
			builder := stream.Builder[int]().Source(v.list...)
			pStream := builder.Build()
			if simple {
				pStream = builder.Simple()
			}
			if parallel {
				pStream = pStream.Parallel()
			}
			if v.count {
				assert.Equal(t, len(v.list), pStream.Count())
			}
			if v.filter != nil {
				pStream = pStream.Filter(v.filter)
			}
			if v.limit.ok {
				pStream = pStream.Limit(v.limit.num)
			}
			if v.skip.ok {
				pStream = pStream.Skip(v.skip.num)
			}
			if v.max != nil {
				val := pStream.Max(v.max)
				assert.Equal(t, false, val.IsEmpty())
				assert.Equal(t, false, val.IsNil())
				assert.Equal(t, v.wantValue, val.Get())
			}
			if v.min != nil {
				val := pStream.Min(v.min)
				assert.Equal(t, false, val.IsEmpty())
				assert.Equal(t, false, val.IsNil())
				assert.Equal(t, v.wantValue, val.Get())
			}
			if v.equals != nil {
				pStream = pStream.Distinct(v.equals)
			}
			if v.mapper != nil {
				pStream = pStream.Map(v.mapper)
			}
			if v.reducer != nil {
				val := pStream.Reduce(v.reducer)
				assert.Equal(t, v.wantValue, val.Get())
			}
			if v.mapperAny != nil {
				anyStream := pStream.MapToAny(v.mapperAny).
					Sort(func(t, u any) bool { return t.(string) < u.(string) })
				assert.Equal(t, v.wantAnyList, anyStream.ToSlice())
			}
			if v.mapperString != nil {
				strStream := pStream.MapToString(v.mapperString).
					Sort(func(t, u string) bool { return t < u })
				assert.Equal(t, v.wantStringList, strStream.ToSlice())
			}
			if v.mapperInt != nil {
				intStream := pStream.MapToInt(v.mapperInt).
					Sort(func(t, u int) bool { return t < u })
				assert.Equal(t, v.wantIntList, intStream.ToSlice())
			}
			if v.mapperFloat != nil {
				floatStream := pStream.MapToFloat(v.mapperFloat).
					Sort(func(t, u float64) bool { return t < u })
				assert.Equal(t, v.wantFloatList, floatStream.ToSlice(), v.wantFloatList)
			}
			if v.less != nil {
				pStream = pStream.Sort(v.less)
			}
			if v.reverse {
				pStream = pStream.Reverse()
			}
			if v.anyMatcher != nil {
				b := pStream.AnyMatch(v.anyMatcher)
				assert.Equal(t, v.wantBool, b)
			}
			if v.allMatcher != nil {
				b := pStream.AllMatch(v.allMatcher)
				assert.Equal(t, v.wantBool, b)
			}
			if v.noneMatcher != nil {
				b := pStream.NoneMatch(v.noneMatcher)
				assert.Equal(t, v.wantBool, b)
			}
			assert.Equal(t, v.wantList, pStream.ToSlice())
		})
	}
}

func BenchmarkPipeline(b *testing.B) {
	var slice []int
	for i := range make([]struct{}, 1000) {
		slice = append(slice, i)
	}

	b.ResetTimer()
	b.Run("direct-sum", func(b *testing.B) {
		var r int
		for n := 0; n < b.N; n++ {
			collections.IterableSlice(slice...).ForEachRemaining(
				func(t int) {
					r += t
				},
			)
		}
	})

	b.Run("fast-serial-sum", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			stream.From(slice...).Reduce(func(i1, i2 int) int { return i1 + i2 })
		}
	})

	// b.Run("fast-parallel-sum", func(b *testing.B) {
	// 	for n := 0; n < b.N; n++ {
	// 		stream.From(slice...).Parallel().Reduce(func(i1, i2 int) int { return i1 + i2 })
	// 	}
	// })

	b.Run("simple-serial-sum", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			stream.Builder[int]().Source(slice...).Simple().Reduce(func(i1, i2 int) int { return i1 + i2 })
		}
	})

	b.Run("simple-parallel-sum", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			stream.Builder[int]().Source(slice...).Simple().Parallel().Reduce(func(i1, i2 int) int { return i1 + i2 })
		}
	})
}
