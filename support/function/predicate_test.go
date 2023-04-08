package function_test

import (
	"testing"

	"github.com/go-park/stream/support/function"
)

// Test for Test()
func TestTest(t *testing.T) {
	fn := func(t int) bool { return t > 0 }
	predicate := function.Predicate[int](fn)
	if !predicate.Test(1) {
		t.Errorf("Expected negatePredicate to be false")
	}
	if predicate.Test(-1) {
		t.Errorf("Expected negatePredicate to be true")
	}
}

// Test for Negate()
func TestNegate(t *testing.T) {
	fn := func(t int) bool { return t > 0 }
	predicate := function.Predicate[int](fn)
	negatePredicate := predicate.Negate()
	if negatePredicate.Test(1) {
		t.Errorf("Expected negatePredicate to be false")
	}
	if !negatePredicate.Test(-1) {
		t.Errorf("Expected negatePredicate to be true")
	}
}

// Test for And()
func TestAnd(t *testing.T) {
	fn1 := func(t int) bool { return t > 0 }
	fn2 := func(t int) bool { return t < 10 }
	predicate1 := function.Predicate[int](fn1)
	predicate2 := function.Predicate[int](fn2)
	andPredicate := predicate1.And(predicate2)
	if !andPredicate.Test(5) {
		t.Errorf("Expected andPredicate to be true")
	}
	if andPredicate.Test(-1) {
		t.Errorf("Expected andPredicate to be false")
	}
}

// Test for Or()
func TestOr(t *testing.T) {
	fn1 := func(t int) bool { return t > 0 }
	fn2 := func(t int) bool { return t > 10 }
	predicate1 := function.Predicate[int](fn1)
	predicate2 := function.Predicate[int](fn2)
	orPredicate := predicate1.Or(predicate2)
	if !orPredicate.Test(15) {
		t.Errorf("Expected orPredicate to be true")
	}
	if orPredicate.Test(-1) {
		t.Errorf("Expected orPredicate to be false")
	}
}

// Test for Not()
func TestNot(t *testing.T) {
	fn := func(t int) bool { return t > 0 }
	predicate := function.Predicate[int](fn)
	notPredicate := predicate.Not(predicate)
	if notPredicate.Test(1) {
		t.Errorf("Expected notPredicate to be false")
	}
	if !notPredicate.Test(-1) {
		t.Errorf("Expected notPredicate to be true")
	}
}

// Test for DeepEqual()
func TestDeepEqual(t *testing.T) {
	obj := "test"
	fn := func(t string) bool { return t == obj }
	predicate := function.Predicate[string](fn)
	deepEqualPredicate := predicate.DeepEqual(obj)
	if !deepEqualPredicate.Test("test") {
		t.Errorf("Expected deepEqualPredicate to be true")
	}
	if deepEqualPredicate("not test") {
		t.Errorf("Expected deepEqualPredicate to be false")
	}
}

// Test for Test()
func TestBiTest(t *testing.T) {
	fn := func(t, u int) bool { return t > 0 && u > 0 }
	predicate := function.BiPredicate[int, int](fn)
	if !predicate.Test(1, 1) {
		t.Errorf("Expected negatePredicate to be false")
	}
	if predicate.Test(-1, 1) {
		t.Errorf("Expected negatePredicate to be true")
	}
}

// Test for Negate()
func TestBiNegate(t *testing.T) {
	fn := func(t, u int) bool { return t > 0 && u > 0 }
	predicate := function.BiPredicate[int, int](fn)
	negatePredicate := predicate.Negate()
	if negatePredicate.Test(1, 1) {
		t.Errorf("Expected negatePredicate to be false")
	}
	if !negatePredicate.Test(1, -1) {
		t.Errorf("Expected negatePredicate to be true")
	}
}

// Test for And()
func TestBiAnd(t *testing.T) {
	fn1 := func(t, u int) bool { return t > 0 && u > 0 }
	fn2 := func(t, u int) bool { return t < 10 && u < 10 }
	predicate1 := function.BiPredicate[int, int](fn1)
	predicate2 := function.BiPredicate[int, int](fn2)
	andPredicate := predicate1.And(predicate2)
	if !andPredicate.Test(5, 5) {
		t.Errorf("Expected andPredicate to be true")
	}
	if andPredicate.Test(-1, 5) {
		t.Errorf("Expected andPredicate to be false")
	}
	if andPredicate.Test(11, 5) {
		t.Errorf("Expected andPredicate to be false")
	}
}

// Test for Or()
func TestBiOr(t *testing.T) {
	fn1 := func(t, u int) bool { return t > 0 && u > 0 }
	fn2 := func(t, u int) bool { return t > 10 && u > 10 }
	predicate1 := function.BiPredicate[int, int](fn1)
	predicate2 := function.BiPredicate[int, int](fn2)
	orPredicate := predicate1.Or(predicate2)
	if !orPredicate.Test(1, 1) {
		t.Errorf("Expected orPredicate to be true")
	}
	if orPredicate.Test(-1, -1) {
		t.Errorf("Expected orPredicate to be false")
	}
}
