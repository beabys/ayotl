package config

import (
	"reflect"
	"testing"
)

func TestFlatten_Basic(t *testing.T) {
	in := map[string]interface{}{
		"a": "b",
		"c": map[string]interface{}{
			"d": "e",
		},
	}
	want := map[string]interface{}{
		"a":   "b",
		"c.d": "e",
	}

	got := Flatten(in)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Flatten() = %#v, want %#v", got, want)
	}
}

func TestFlatten_EmptyNestedMap(t *testing.T) {
	in := map[string]interface{}{
		"a": map[string]interface{}{},
	}
	want := map[string]interface{}{
		"a": map[string]interface{}{},
	}

	got := Flatten(in)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Flatten(empty nested) = %#v, want %#v", got, want)
	}
}

func TestFlatten_MultiLevel(t *testing.T) {
	in := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": 1,
			},
		},
		"x": map[string]interface{}{
			"y": 2,
		},
	}
	want := map[string]interface{}{
		"a.b.c": 1,
		"x.y":   2,
	}

	got := Flatten(in)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Flatten(multi-level) = %#v, want %#v", got, want)
	}
}

func TestFlatten_NilAndEmpty(t *testing.T) {
	var nilMap map[string]interface{}
	gotNil := Flatten(nilMap)
	if len(gotNil) != 0 {
		t.Fatalf("Flatten(nil) returned non-empty map: %#v", gotNil)
	}

	empty := map[string]interface{}{}
	gotEmpty := Flatten(empty)
	if len(gotEmpty) != 0 {
		t.Fatalf("Flatten(empty) returned non-empty map: %#v", gotEmpty)
	}
}
func TestMergeKeys_AddNonExisting(t *testing.T) {
	m1 := ConfigMap{
		"a": "1",
	}
	m2 := ConfigMap{
		"b": "2",
	}

	got := MergeKeys(m1, m2)
	want := map[string]interface{}{
		"a": "1",
		"b": "2",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MergeKeys(add non-existing) = %#v, want %#v", got, want)
	}
}

func TestMergeKeys_ReplaceScalar(t *testing.T) {
	m1 := ConfigMap{
		"x": "old",
	}
	m2 := ConfigMap{
		"x": "new",
	}

	got := MergeKeys(m1, m2)
	want := map[string]interface{}{
		"x": "new",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MergeKeys(replace scalar) = %#v, want %#v", got, want)
	}
}

func TestMergeKeys_MergeNested(t *testing.T) {
	m1 := ConfigMap{
		"s": map[string]interface{}{
			"k1": "v1",
		},
	}
	m2 := ConfigMap{
		"s": ConfigMap{
			"k2": "v2",
		},
	}

	got := MergeKeys(m1, m2)
	want := map[string]interface{}{
		"s": map[string]interface{}{
			"k1": "v1",
			"k2": "v2",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MergeKeys(merge nested) = %#v, want %#v", got, want)
	}
}

func TestMergeKeys_ReplaceWithMap(t *testing.T) {
	m1 := ConfigMap{
		"s": "primitive",
	}
	m2 := ConfigMap{
		"s": ConfigMap{
			"k": "v",
		},
	}

	got := MergeKeys(m1, m2)
	want := map[string]interface{}{
		"s": ConfigMap{
			"k": "v",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MergeKeys(replace with map) = %#v, want %#v", got, want)
	}
}
