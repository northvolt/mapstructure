package mapstructure

import (
	"reflect"
	"testing"
)

type Test struct {
	Foo *int
	Bar *Test
}

// Tests that the mapstructure library handles fields that are pointers to other structures and scalars properly.
func TestPointerFieldDecoding(t *testing.T) {
	s := Test{}
	if err := Decode(map[string]interface{}{}, &s); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(s, Test{}) {
		t.Fatalf("bad: %#v", s)
	}

	if err := Decode(map[string]interface{}{"foo": 1, "bar": map[string]interface{}{}}, &s); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(s, Test{Foo: pInt(1), Bar: &Test{}}) {
		t.Fatalf("bad: %#v", s)
	}
}

func pInt(i int) *int {
	return &i
}

// Test that error in field decoding doesn't affect the structure being assigned in a field.
func TestStructFieldDecodingError(t *testing.T) {
	s := Test{}
	if err := Decode(map[string]interface{}{"bar": map[string]interface{}{"bar": "foo", "foo": 1}}, &s); err == nil {
		t.Fatal("expected one error")
	}
	if !reflect.DeepEqual(s, Test{Foo: nil, Bar: &Test{Foo: pInt(1)}}) {
		t.Fatalf("bad: %#v", s)
	}
}

func TestArrayElementDecodingError(t *testing.T) {
	type Foo struct {
		Foo int
	}
	type Array struct {
		A []*Foo
	}
	a := Array{}
	encoded := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{"foo": 1},
			"foo",
		},
	}
	if err := Decode(encoded, &a); err == nil {
		t.Fatal("expected one error")
	}
	if !reflect.DeepEqual(a, Array{A: []*Foo{{Foo: 1}, nil}}) {
		t.Fatalf("bad: %#v", a)
	}
}

func TestDecode2(t *testing.T) {
	type Foo struct {
		Foo int
		Bar int
	}

	var s Foo
	res := Decode2(map[string]interface{}{"foo": "foo", "baz": 1}, &s)
	if len(res.Errors) != 1 {
		t.Fatalf("expected one error: %v", res.Errors)
	}
	if !res.TotalFailure {
		t.Fatal("expected total failure")
	}

	res = Decode2(map[string]interface{}{"foo": "foo", "bar": 1}, &s)
	if len(res.Errors) != 1 {
		t.Fatalf("expected one error: %v", res.Errors)
	}
	if res.TotalFailure {
		t.Fatal("expected not total failure")
	}
}
