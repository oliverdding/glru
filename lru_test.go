package lru

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
	AssertEqual(t, 4, lru.entries[lru.head].val)
	lru.Insert(8)
	AssertEqual(t, 8, lru.entries[lru.head].val)
	lru.Insert(11)
	AssertEqual(t, 11, lru.entries[lru.head].val)
	lru.Insert(1)
	AssertEqual(t, 1, lru.entries[lru.head].val)
	lru.Insert(5)
	AssertEqual(t, 5, lru.entries[lru.head].val)
	result := make([]interface{}, 0, 4)
	itor := lru.Itor()
	for itor.HasNext() {
		result = append(result, lru.entries[itor.GetNext()].val)
	}
	AssertArrayEqual(t, result, []interface{}{5, 1, 11, 8})
	lru.Insert(132)
	result = make([]interface{}, 0, 4)
	itor = lru.Itor()
	for itor.HasNext() {
		result = append(result, lru.entries[itor.GetNext()].val)
	}
	AssertArrayEqual(t, result, []interface{}{132, 5, 1, 11})
}

func TestFind(t *testing.T) {
	lru := New(4)
	lru.Insert(1)
	AssertEqual(t, 1, lru.Len())
	lru.Insert(2)
	AssertEqual(t, 2, lru.Len())
	lru.Insert(3)
	AssertEqual(t, 3, lru.Len())
	lru.Insert(4)
	AssertEqual(t, 4, lru.Len())
	result := make([]interface{}, 0, 4)
	itor := lru.Itor()
	for itor.HasNext() {
		result = append(result, lru.entries[itor.GetNext()].val)
	}
	AssertArrayEqual(t, result, []interface{}{4, 3, 2, 1})
	lru.Find(func(i interface{}) bool {
		return i.(int) == 3
	})
	result = make([]interface{}, 0, 4)
	itor = lru.Itor()
	for itor.HasNext() {
		result = append(result, lru.entries[itor.GetNext()].val)
	}
	AssertArrayEqual(t, result, []interface{}{3, 4, 2, 1})

	lru.Find(func(i interface{}) bool {
		return i.(int) == 1
	})
	result = make([]interface{}, 0, 4)
	itor = lru.Itor()
	for itor.HasNext() {
		result = append(result, lru.entries[itor.GetNext()].val)
	}
	AssertArrayEqual(t, result, []interface{}{1, 3, 4, 2})
}
