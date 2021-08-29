package lru_go

import (
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func AssertArrayEqual(t *testing.T, a, b []interface{}) {
	if len(a) != len(b) {
		t.Errorf("length is not equal, %v != %v", len(a), len(b))
	}
	for i := range a {
		if a[i] != b[i] {
			t.Errorf("In index %v, Received %v (type %v), expected %v (type %v)", i, a[i], reflect.TypeOf(a[i]), b[i], reflect.TypeOf(b[i]))
		}
	}
}

func TestLen(t *testing.T) {
	lru := New(4)
	AssertEqual(t, 0, lru.Len())
	lru.Insert(1)
	AssertEqual(t, 1, lru.Len())
	lru.Insert(2)
	AssertEqual(t, 2, lru.Len())
	lru.Insert(3)
	AssertEqual(t, 3, lru.Len())
	lru.Insert(4)
	AssertEqual(t, 4, lru.Len())
	lru.Insert(5)
	AssertEqual(t, 4, lru.Len())
}

func TestInsert(t *testing.T) {
	lru := New(4)
	lru.Insert(4)
	AssertArrayEqual(t, lru.ToArray(), []interface{}{4})
	lru.Insert(8)
	AssertArrayEqual(t, lru.ToArray(), []interface{}{8, 4})
	lru.Insert(11)
	AssertArrayEqual(t, lru.ToArray(), []interface{}{11, 8, 4})
	lru.Insert(1)
	AssertArrayEqual(t, lru.ToArray(), []interface{}{1, 11, 8, 4})
	lru.Insert(5)
	AssertArrayEqual(t, lru.ToArray(), []interface{}{5, 1, 11, 8})
	lru.Insert(132)
	AssertArrayEqual(t, lru.ToArray(), []interface{}{132, 5, 1, 11})
}

func TestFind(t *testing.T) {
	lru := New(4)
	lru.Insert(1)
	AssertArrayEqual(t, lru.ToArray(), []interface{}{1})
	lru.Insert(2)
	AssertArrayEqual(t, lru.ToArray(), []interface{}{2, 1})
	lru.Insert(3)
	AssertArrayEqual(t, lru.ToArray(), []interface{}{3, 2, 1})
	lru.Insert(4)
	AssertArrayEqual(t, lru.ToArray(), []interface{}{4, 3, 2, 1})
	AssertEqual(t, lru.Find(func(i interface{}) bool {
		return i.(int) == 3
	}), 3)
	AssertArrayEqual(t, lru.ToArray(), []interface{}{3, 4, 2, 1})
	AssertEqual(t, lru.Find(func(i interface{}) bool {
		return i.(int) == 1
	}), 1)
	AssertArrayEqual(t, lru.ToArray(), []interface{}{1, 3, 4, 2})
}
