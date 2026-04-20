package orderedobject

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	json "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/google/go-cmp/cmp"
)

func TestNewObject(t *testing.T) {
	t.Parallel()

	t.Run("without capacity", func(t *testing.T) {
		t.Parallel()

		obj := NewObject[int]()
		if got := obj.Len(); got != 0 {
			t.Fatalf("Len() = %d, want 0", got)
		}

		obj.Set("a", 1)
		if got := obj.Len(); got != 1 {
			t.Fatalf("Len() after Set = %d, want 1", got)
		}
	})

	t.Run("with capacity", func(t *testing.T) {
		t.Parallel()

		obj := NewObject[int](3)
		obj.Set("a", 1).Set("b", 2).Set("c", 3).Set("d", 4)
		if got := obj.Len(); got != 4 {
			t.Fatalf("Len() = %d, want 4", got)
		}

		wantKeys := []string{"a", "b", "c", "d"}
		if diff := cmp.Diff(wantKeys, obj.Keys()); diff != "" {
			t.Errorf("Keys() mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestFromMap(t *testing.T) {
	t.Parallel()

	t.Run("copies all entries", func(t *testing.T) {
		t.Parallel()

		input := map[string]int{"a": 1, "b": 2, "c": 3}
		obj := FromMap(input)
		if got := obj.Len(); got != len(input) {
			t.Fatalf("Len() = %d, want %d", got, len(input))
		}

		for key, want := range input {
			got, found := obj.Get(key)
			if !found {
				t.Fatalf("Get(%q) not found", key)
			}
			if got != want {
				t.Fatalf("Get(%q) = %d, want %d", key, got, want)
			}
		}
	})

	t.Run("does not alias input map", func(t *testing.T) {
		t.Parallel()

		input := map[string]int{"count": 1}
		obj := FromMap(input)
		input["count"] = 2

		got, found := obj.Get("count")
		if !found {
			t.Fatal("Get(\"count\") not found")
		}
		if got != 1 {
			t.Fatalf("Get(\"count\") = %d, want 1", got)
		}
	})
}

func TestSetGetHasDeleteLen(t *testing.T) {
	t.Parallel()

	obj := NewObject[any]()
	if got := obj.Len(); got != 0 {
		t.Fatalf("Len() = %d, want 0", got)
	}
	if obj.Has("missing") {
		t.Fatal("Has(\"missing\") = true, want false")
	}
	if got, found := obj.Get("missing"); found || got != nil {
		t.Fatalf("Get(\"missing\") = (%v, %v), want (nil, false)", got, found)
	}

	returned := obj.Set("a", 1)
	if returned != obj {
		t.Fatal("Set() did not return receiver")
	}
	obj.Set("b", 2).Set("c", 3)
	if got := obj.Len(); got != 3 {
		t.Fatalf("Len() = %d, want 3", got)
	}
	if !obj.Has("b") {
		t.Fatal("Has(\"b\") = false, want true")
	}

	obj.Set("a", 99)
	if got, found := obj.Get("a"); !found || got != 99 {
		t.Fatalf("Get(\"a\") = (%v, %v), want (99, true)", got, found)
	}
	wantKeys := []string{"a", "b", "c"}
	if diff := cmp.Diff(wantKeys, obj.Keys()); diff != "" {
		t.Errorf("Keys() after update mismatch (-want +got):\n%s", diff)
	}

	returned = obj.Delete("b")
	if returned != obj {
		t.Fatal("Delete() did not return receiver")
	}
	if obj.Has("b") {
		t.Fatal("Has(\"b\") = true after Delete")
	}
	if got := obj.Len(); got != 2 {
		t.Fatalf("Len() after Delete = %d, want 2", got)
	}

	obj.Delete("missing")
	if got := obj.Len(); got != 2 {
		t.Fatalf("Len() after deleting missing key = %d, want 2", got)
	}

	wantKeys = []string{"a", "c"}
	if diff := cmp.Diff(wantKeys, obj.Keys()); diff != "" {
		t.Errorf("Keys() after Delete mismatch (-want +got):\n%s", diff)
	}
}

func TestEmptyStringKey(t *testing.T) {
	t.Parallel()

	obj := NewObject[int]().Set("", 1).Set("named", 2)

	if got, found := obj.Get(""); !found || got != 1 {
		t.Fatalf("Get(\"\") = (%v, %v), want (1, true)", got, found)
	}
	if !obj.Has("") {
		t.Fatal("Has(\"\") = false, want true")
	}

	obj.Set("", 3)
	if diff := cmp.Diff([]string{"", "named"}, obj.Keys()); diff != "" {
		t.Errorf("Keys() after empty-key update mismatch (-want +got):\n%s", diff)
	}
	if got, found := obj.Get(""); !found || got != 3 {
		t.Fatalf("Get(\"\") after update = (%v, %v), want (3, true)", got, found)
	}

	obj.Delete("")
	if obj.Has("") {
		t.Fatal("Has(\"\") = true after Delete")
	}
	if diff := cmp.Diff([]string{"named"}, obj.Keys()); diff != "" {
		t.Errorf("Keys() after empty-key Delete mismatch (-want +got):\n%s", diff)
	}
}

func TestKeysValuesEntries(t *testing.T) {
	t.Parallel()

	t.Run("empty object", func(t *testing.T) {
		t.Parallel()

		obj := NewObject[int]()
		if diff := cmp.Diff([]string{}, obj.Keys()); diff != "" {
			t.Errorf("Keys() mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff([]int{}, obj.Values()); diff != "" {
			t.Errorf("Values() mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff([]Entry[int]{}, obj.Entries()); diff != "" {
			t.Errorf("Entries() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("preserve insertion order and return copies", func(t *testing.T) {
		t.Parallel()

		obj := NewObject[int]().Set("a", 1).Set("b", 2).Set("c", 3)
		wantKeys := []string{"a", "b", "c"}
		wantValues := []int{1, 2, 3}
		wantEntries := []Entry[int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "c", Value: 3}}

		if diff := cmp.Diff(wantKeys, obj.Keys()); diff != "" {
			t.Errorf("Keys() mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff(wantValues, obj.Values()); diff != "" {
			t.Errorf("Values() mismatch (-want +got):\n%s", diff)
		}
		gotEntries := obj.Entries()
		if diff := cmp.Diff(wantEntries, gotEntries); diff != "" {
			t.Errorf("Entries() mismatch (-want +got):\n%s", diff)
		}

		gotEntries[0] = Entry[int]{Key: "x", Value: 99}
		if diff := cmp.Diff(wantEntries, obj.Entries()); diff != "" {
			t.Errorf("Entries() leaked mutation (-want +got):\n%s", diff)
		}
	})
}

func TestForEach(t *testing.T) {
	t.Parallel()

	obj := NewObject[int]().Set("a", 1).Set("b", 2).Set("c", 3)
	var gotKeys []string
	var sum int

	obj.ForEach(func(key string, value int) {
		gotKeys = append(gotKeys, key)
		sum += value
	})

	wantKeys := []string{"a", "b", "c"}
	if diff := cmp.Diff(wantKeys, gotKeys); diff != "" {
		t.Errorf("ForEach keys mismatch (-want +got):\n%s", diff)
	}
	if sum != 6 {
		t.Fatalf("ForEach sum = %d, want 6", sum)
	}
}

func TestLargeValueOperations(t *testing.T) {
	t.Parallel()

	type largeValue struct {
		Name    string
		Payload [64]int
	}

	want := largeValue{Name: "alpha"}
	for i := range want.Payload {
		want.Payload[i] = i + 1
	}

	obj := NewObject[largeValue]().Set("first", want)

	got, found := obj.Get("first")
	if !found {
		t.Fatal("Get(\"first\") not found")
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("Get(\"first\") mismatch (-want +got):\n%s", diff)
	}
	if !obj.Has("first") {
		t.Fatal("Has(\"first\") = false, want true")
	}

	if diff := cmp.Diff([]largeValue{want}, obj.Values()); diff != "" {
		t.Errorf("Values() mismatch (-want +got):\n%s", diff)
	}

	var gotEntries []Entry[largeValue]
	obj.ForEach(func(key string, value largeValue) {
		gotEntries = append(gotEntries, Entry[largeValue]{Key: key, Value: value})
	})
	wantEntries := []Entry[largeValue]{{Key: "first", Value: want}}
	if diff := cmp.Diff(wantEntries, gotEntries); diff != "" {
		t.Errorf("ForEach() mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(map[string]largeValue{"first": want}, obj.ToMap()); diff != "" {
		t.Errorf("ToMap() mismatch (-want +got):\n%s", diff)
	}
}

func TestClone(t *testing.T) {
	t.Parallel()

	original := NewObject[any]().
		Set("a", 1).
		Set("nested", []int{1, 2}).
		Set("c", 3)

	clone := original.Clone()
	if diff := cmp.Diff(original.Entries(), clone.Entries()); diff != "" {
		t.Errorf("Clone() mismatch (-want +got):\n%s", diff)
	}

	clone.Set("a", 99).Delete("c").Set("d", 4)

	if got, found := original.Get("a"); !found || got != 1 {
		t.Fatalf("original.Get(\"a\") = (%v, %v), want (1, true)", got, found)
	}
	if !original.Has("c") {
		t.Fatal("original lost key \"c\" after clone mutation")
	}
	if original.Has("d") {
		t.Fatal("original unexpectedly gained key \"d\"")
	}

	sliceValue, found := original.Get("nested")
	if !found {
		t.Fatal("original.Get(\"nested\") not found")
	}
	cloneSlice, found := clone.Get("nested")
	if !found {
		t.Fatal("clone.Get(\"nested\") not found")
	}
	if diff := cmp.Diff(sliceValue, cloneSlice); diff != "" {
		t.Errorf("Clone() shallow copy mismatch (-want +got):\n%s", diff)
	}
}

func TestToMap(t *testing.T) {
	t.Parallel()

	obj := NewObject[int]().Set("a", 1).Set("b", 2)
	got := obj.ToMap()
	want := map[string]int{"a": 1, "b": 2}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ToMap() mismatch (-want +got):\n%s", diff)
	}

	got["a"] = 99
	value, found := obj.Get("a")
	if !found {
		t.Fatal("Get(\"a\") not found after ToMap mutation")
	}
	if value != 1 {
		t.Fatalf("Get(\"a\") = %d, want 1", value)
	}
}

func TestMarshalJSONAndToJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		obj  *Object[any]
		want string
	}{
		{
			name: "empty object",
			obj:  NewObject[any](),
			want: `{}` + "\n",
		},
		{
			name: "ordered values",
			obj: NewObject[any]().
				Set("name", "John").
				Set("age", 30).
				Set("city", "New York"),
			want: `{"name":"John","age":30,"city":"New York"}` + "\n",
		},
		{
			name: "nested ordered object",
			obj: NewObject[any]().
				Set("name", "Alice").
				Set("address", NewObject[any]().Set("street", "123 Main St").Set("zip", "10001")),
			want: `{"name":"Alice","address":{"street":"123 Main St","zip":"10001"}}` + "\n",
		},
		{
			name: "array of ordered objects",
			obj: NewObject[any]().Set("people", []any{
				NewObject[any]().Set("name", "Bob").Set("age", 35),
				NewObject[any]().Set("name", "Charlie").Set("age", 40),
			}),
			want: `{"people":[{"name":"Bob","age":35},{"name":"Charlie","age":40}]}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.obj.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("MarshalJSON() = %q, want %q", string(got), tt.want)
			}

			toJSON, err := tt.obj.ToJSON()
			if err != nil {
				t.Fatalf("ToJSON() error = %v", err)
			}
			if string(toJSON)+"\n" != tt.want {
				t.Fatalf("ToJSON() = %q, want %q", string(toJSON), strings.TrimSuffix(tt.want, "\n"))
			}
		})
	}
}

func TestMarshalJSONTo(t *testing.T) {
	t.Parallel()

	t.Run("writes deterministic map ordering", func(t *testing.T) {
		t.Parallel()

		obj := NewObject[any]().Set("fruits", map[string]any{
			"zebra":  1,
			"apple":  2,
			"mango":  3,
			"banana": 4,
		})

		var first string
		for i := range 5 {
			var buf bytes.Buffer
			enc := jsontext.NewEncoder(&buf)
			if err := obj.MarshalJSONTo(enc); err != nil {
				t.Fatalf("MarshalJSONTo() iteration %d error = %v", i, err)
			}
			got := buf.String()
			if i == 0 {
				first = got
				continue
			}
			if got != first {
				t.Fatalf("MarshalJSONTo() iteration %d = %q, want %q", i, got, first)
			}
		}

		if strings.Index(first, `"apple"`) >= strings.Index(first, `"banana"`) ||
			strings.Index(first, `"banana"`) >= strings.Index(first, `"mango"`) ||
			strings.Index(first, `"mango"`) >= strings.Index(first, `"zebra"`) {
			t.Fatalf("MarshalJSONTo() did not sort nested map keys: %q", first)
		}
	})

	t.Run("returns error for unmarshalable value", func(t *testing.T) {
		t.Parallel()

		obj := NewObject[any]().Set("fn", func() {})
		var buf bytes.Buffer
		enc := jsontext.NewEncoder(&buf)
		if err := obj.MarshalJSONTo(enc); err == nil {
			t.Fatal("MarshalJSONTo() error = nil, want non-nil")
		}
	})
}

func TestOrderedMarshalerInterface(t *testing.T) {
	t.Parallel()

	inner := NewObject[any]().Set("z", 3).Set("w", 4)
	var marshaler OrderedMarshaler = inner

	var buf bytes.Buffer
	enc := jsontext.NewEncoder(&buf)
	if err := marshaler.MarshalJSONTo(enc); err != nil {
		t.Fatalf("MarshalJSONTo() error = %v", err)
	}
	if got := buf.String(); got != `{"z":3,"w":4}`+"\n" {
		t.Fatalf("MarshalJSONTo() = %q, want %q", got, `{"z":3,"w":4}`+"\n")
	}
}

func TestFromJSON(t *testing.T) {
	t.Parallel()

	t.Run("preserves top level order", func(t *testing.T) {
		t.Parallel()

		obj, err := FromJSON[any]([]byte(`{"name":"John","age":30,"city":"New York"}`))
		if err != nil {
			t.Fatalf("FromJSON() error = %v", err)
		}

		want := []Entry[any]{
			{Key: "name", Value: "John"},
			{Key: "age", Value: float64(30)},
			{Key: "city", Value: "New York"},
		}
		if diff := cmp.Diff(want, obj.Entries()); diff != "" {
			t.Errorf("Entries() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("nested values decode as standard maps for any", func(t *testing.T) {
		t.Parallel()

		obj, err := FromJSON[any]([]byte(`{"settings":{"theme":"dark","notifications":true},"version":"1.0"}`))
		if err != nil {
			t.Fatalf("FromJSON() error = %v", err)
		}

		entries := obj.Entries()
		wantTopLevelKeys := []string{"settings", "version"}
		gotTopLevelKeys := make([]string, 0, len(entries))
		for _, entry := range entries {
			gotTopLevelKeys = append(gotTopLevelKeys, entry.Key)
		}
		if diff := cmp.Diff(wantTopLevelKeys, gotTopLevelKeys); diff != "" {
			t.Errorf("top-level keys mismatch (-want +got):\n%s", diff)
		}

		settings, ok := entries[0].Value.(map[string]any)
		if !ok {
			t.Fatalf("entries[0].Value type = %T, want map[string]any", entries[0].Value)
		}
		if settings["theme"] != "dark" {
			t.Fatalf("settings[\"theme\"] = %v, want dark", settings["theme"])
		}
		if settings["notifications"] != true {
			t.Fatalf("settings[\"notifications\"] = %v, want true", settings["notifications"])
		}
	})
}

func TestFromJSONErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantError error
	}{
		{name: "malformed JSON", input: `{"key":"value"`, wantError: nil},
		{name: "array instead of object", input: `["a","b"]`, wantError: ErrExpectedObjectStart},
		{name: "non-string key", input: `{1:"value"}`, wantError: ErrExpectedStringKey},
		{name: "primitive instead of object", input: `"hello"`, wantError: ErrExpectedObjectStart},
		{name: "empty input", input: ``, wantError: nil},
		{name: "duplicate keys", input: `{"key":"value1","key":"value2"}`, wantError: nil},
		{name: "trailing garbage", input: `{"ok":1} trailing`, wantError: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := FromJSON[any]([]byte(tt.input))
			if err == nil {
				t.Fatal("FromJSON() error = nil, want non-nil")
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("errors.Is(%v, %v) = false", err, tt.wantError)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("clears existing entries", func(t *testing.T) {
		t.Parallel()

		obj := NewObject[any]().Set("old", "data")
		if err := obj.UnmarshalJSON([]byte(`{"new":"value"}`)); err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		if obj.Has("old") {
			t.Fatal("old key still present after UnmarshalJSON()")
		}
		if got, found := obj.Get("new"); !found || got != "value" {
			t.Fatalf("Get(\"new\") = (%v, %v), want (value, true)", got, found)
		}
	})

	t.Run("returns ErrExpectedObjectStart for non objects", func(t *testing.T) {
		t.Parallel()

		for _, input := range []string{`42`, `true`, `null`} {
			var obj Object[any]
			err := obj.UnmarshalJSON([]byte(input))
			if err == nil {
				t.Fatalf("UnmarshalJSON(%q) error = nil, want non-nil", input)
			}
			if !errors.Is(err, ErrExpectedObjectStart) {
				t.Fatalf("errors.Is(%v, ErrExpectedObjectStart) = false", err)
			}
		}
	})
}

func TestUnmarshalJSONFrom(t *testing.T) {
	t.Parallel()

	t.Run("decodes into typed object", func(t *testing.T) {
		t.Parallel()

		type user struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		dec := jsontext.NewDecoder(strings.NewReader(`{"first":{"name":"Alice","age":30},"second":{"name":"Bob","age":35}}`))
		var obj Object[user]
		if err := obj.UnmarshalJSONFrom(dec); err != nil {
			t.Fatalf("UnmarshalJSONFrom() error = %v", err)
		}

		want := []Entry[user]{
			{Key: "first", Value: user{Name: "Alice", Age: 30}},
			{Key: "second", Value: user{Name: "Bob", Age: 35}},
		}
		if diff := cmp.Diff(want, obj.Entries()); diff != "" {
			t.Errorf("Entries() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("leaves following top level value for caller", func(t *testing.T) {
		t.Parallel()

		dec := jsontext.NewDecoder(strings.NewReader(`{"ok":1} true`))
		var obj Object[any]
		if err := obj.UnmarshalJSONFrom(dec); err != nil {
			t.Fatalf("UnmarshalJSONFrom() error = %v", err)
		}

		tok, err := dec.ReadToken()
		if err != nil {
			t.Fatalf("ReadToken() error = %v", err)
		}
		if tok.Kind() != 't' {
			t.Fatalf("next token kind = %q, want %q", tok.Kind(), 't')
		}
	})
}

func TestJSONRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{name: "simple object", input: `{"name":"John","age":30,"active":true,"score":95.5}`},
		{name: "empty object", input: `{}`},
		{name: "null values", input: `{"a":null,"b":null}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			obj, err := FromJSON[any]([]byte(tt.input))
			if err != nil {
				t.Fatalf("FromJSON() error = %v", err)
			}

			got, err := obj.ToJSON()
			if err != nil {
				t.Fatalf("ToJSON() error = %v", err)
			}
			if string(got) != tt.input {
				t.Fatalf("ToJSON() = %q, want %q", string(got), tt.input)
			}
		})
	}
}

func TestJSONTags(t *testing.T) {
	t.Parallel()

	type user struct {
		Name      string `json:"full_name"`
		Age       int    `json:"age"`
		Email     string `json:"email,omitempty"`
		IsActive  bool   `json:"-"`
		SecretKey string `json:"secret_key,omitempty"`
	}

	tests := []struct {
		name  string
		input user
		want  string
	}{
		{
			name: "all fields populated",
			input: user{
				Name:      "John Doe",
				Age:       30,
				Email:     "john@example.com",
				IsActive:  true,
				SecretKey: "abc123",
			},
			want: `{"user":{"full_name":"John Doe","age":30,"email":"john@example.com","secret_key":"abc123"}}`,
		},
		{
			name: "omitempty fields omitted",
			input: user{
				Name:     "Jane Smith",
				Age:      25,
				IsActive: true,
			},
			want: `{"user":{"full_name":"Jane Smith","age":25}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			obj := NewObject[user]().Set("user", tt.input)
			got, err := obj.ToJSON()
			if err != nil {
				t.Fatalf("ToJSON() error = %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("ToJSON() = %q, want %q", string(got), tt.want)
			}

			var decoded Object[user]
			if err := json.Unmarshal(got, &decoded); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			decodedUser, found := decoded.Get("user")
			if !found {
				t.Fatal("decoded.Get(\"user\") not found")
			}
			wantUser := tt.input
			wantUser.IsActive = false
			if diff := cmp.Diff(wantUser, decodedUser); diff != "" {
				t.Errorf("decoded user mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func BenchmarkObjectSet(b *testing.B) {
	obj := NewObject[any](100)
	b.ResetTimer()
	for b.Loop() {
		obj.Set("key", 42)
	}
}

func BenchmarkObjectGet(b *testing.B) {
	obj := NewObject[any](100)
	for i := range 100 {
		obj.Set(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for b.Loop() {
		obj.Get("key50")
	}
}

func BenchmarkObjectHas(b *testing.B) {
	obj := NewObject[any](100)
	for i := range 100 {
		obj.Set(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for b.Loop() {
		obj.Has("key50")
	}
}

func BenchmarkObjectDelete(b *testing.B) {
	obj := NewObject[any](100)
	for i := range 100 {
		obj.Set(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for b.Loop() {
		obj.Delete("key50")
		obj.Set("key50", 50)
	}
}

func BenchmarkObjectMarshalJSON(b *testing.B) {
	obj := NewObject[any](10)
	obj.Set("name", "John").
		Set("age", 30).
		Set("city", "New York").
		Set("active", true)
	b.ResetTimer()
	for b.Loop() {
		_, _ = obj.MarshalJSON()
	}
}
