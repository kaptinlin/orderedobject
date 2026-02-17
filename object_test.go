package orderedobject

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	json "github.com/go-json-experiment/json"
)

func TestMarshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		obj  *Object[any]
		want string
	}{
		{
			name: "empty object",
			obj:  NewObject[any](0),
			want: `{}`,
		},
		{
			name: "single key-value pair",
			obj:  NewObject[any](1).Set("key", "value"),
			want: `{"key":"value"}`,
		},
		{
			name: "multiple key-value pairs",
			obj: NewObject[any](3).
				Set("name", "John").
				Set("age", 30).
				Set("city", "New York"),
			want: `{"name":"John","age":30,"city":"New York"}`,
		},
		{
			name: "nested objects",
			obj: func() *Object[any] {
				address := NewObject[any](2).
					Set("street", "123 Main St").
					Set("city", "London")
				return NewObject[any](3).
					Set("name", "Alice").
					Set("age", 28).
					Set("address", address)
			}(),
			want: `{"name":"Alice","age":28,"address":{"street":"123 Main St","city":"London"}}`,
		},
		{
			name: "array of objects",
			obj: func() *Object[any] {
				p1 := NewObject[any](2).Set("name", "Bob").Set("age", 35)
				p2 := NewObject[any](2).Set("name", "Charlie").Set("age", 40)
				return NewObject[any](1).Set("people", []any{p1, p2})
			}(),
			want: `{"people":[{"name":"Bob","age":35},{"name":"Charlie","age":40}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := json.Marshal(tt.obj)
			if err != nil {
				t.Fatalf("json.Marshal() returned unexpected error: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("json.Marshal() = %q, want %q", string(got), tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](3).
		Set("string", "value").
		Set("int", 42).
		Set("bool", true)

	tests := []struct {
		name      string
		key       string
		wantValue any
		wantFound bool
	}{
		{name: "string value", key: "string", wantValue: "value", wantFound: true},
		{name: "int value", key: "int", wantValue: 42, wantFound: true},
		{name: "bool value", key: "bool", wantValue: true, wantFound: true},
		{name: "non-existent key", key: "missing", wantValue: nil, wantFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, found := obj.Get(tt.key)
			if found != tt.wantFound {
				t.Errorf("Get(%q) found = %v, want %v", tt.key, found, tt.wantFound)
			}
			if got != tt.wantValue {
				t.Errorf("Get(%q) = %v, want %v", tt.key, got, tt.wantValue)
			}
		})
	}
}

func TestHas(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](2).Set("key", "value")

	if !obj.Has("key") {
		t.Error("Has(\"key\") = false, want true")
	}
	if obj.Has("missing") {
		t.Error("Has(\"missing\") = true, want false")
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](3).
		Set("a", 1).Set("b", 2).Set("c", 3)

	// Delete middle element and verify order is maintained.
	obj.Delete("b")

	got, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("json.Marshal() returned unexpected error: %v", err)
	}
	if want := `{"a":1,"c":3}`; string(got) != want {
		t.Errorf("json.Marshal() after Delete = %q, want %q", string(got), want)
	}
	if obj.Has("b") {
		t.Error("Has(\"b\") after Delete = true, want false")
	}

	// Delete non-existent key should be a no-op.
	obj.Delete("missing")
	if got := obj.Len(); got != 2 {
		t.Errorf("Len() after no-op Delete = %d, want 2", got)
	}
}

func TestForEach(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](3).
		Set("a", 1).Set("b", 2).Set("c", 3)

	var keys string
	var sum int

	obj.ForEach(func(key string, value any) {
		keys += key
		sum += value.(int)
	})

	if keys != "abc" {
		t.Errorf("ForEach keys = %q, want \"abc\"", keys)
	}
	if sum != 6 {
		t.Errorf("ForEach sum = %d, want 6", sum)
	}
}

func TestLen(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](0)
	if got := obj.Len(); got != 0 {
		t.Errorf("Len() on empty = %d, want 0", got)
	}

	obj.Set("a", 1).Set("b", 2)
	if got := obj.Len(); got != 2 {
		t.Errorf("Len() after two Sets = %d, want 2", got)
	}

	obj.Delete("a")
	if got := obj.Len(); got != 1 {
		t.Errorf("Len() after Delete = %d, want 1", got)
	}
}

func TestChaining(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](0).
		Set("a", 1).Set("b", 2).Set("c", 3)

	if got := obj.Len(); got != 3 {
		t.Errorf("Len() = %d, want 3", got)
	}

	got, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("json.Marshal() returned unexpected error: %v", err)
	}
	if want := `{"a":1,"b":2,"c":3}`; string(got) != want {
		t.Errorf("json.Marshal() = %q, want %q", string(got), want)
	}
}

func TestFromMap(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"name": "Alice",
		"age":  28,
		"city": "London",
	}

	obj := FromMap(m)
	if got := obj.Len(); got != len(m) {
		t.Errorf("Len() = %d, want %d", got, len(m))
	}

	for k, want := range m {
		got, found := obj.Get(k)
		if !found {
			t.Errorf("Get(%q) not found", k)
			continue
		}
		if got != want {
			t.Errorf("Get(%q) = %v, want %v", k, got, want)
		}
	}
}

func TestFromJSON(t *testing.T) {
	t.Parallel()

	obj, err := FromJSON[any]([]byte(`{"name":"John","age":30,"city":"New York"}`))
	if err != nil {
		t.Fatalf("FromJSON() returned unexpected error: %v", err)
	}
	if got := obj.Len(); got != 3 {
		t.Errorf("Len() = %d, want 3", got)
	}

	tests := []struct {
		key  string
		want any
	}{
		{key: "name", want: "John"},
		{key: "age", want: float64(30)},
		{key: "city", want: "New York"},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			t.Parallel()
			got, found := obj.Get(tt.key)
			if !found {
				t.Fatalf("Get(%q) not found", tt.key)
			}
			if got != tt.want {
				t.Errorf("Get(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestClone(t *testing.T) {
	t.Parallel()

	original := NewObject[any](3).
		Set("a", 1).Set("b", 2).Set("c", 3)

	clone := original.Clone()

	// Verify clone matches original.
	origJSON, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal(original) returned unexpected error: %v", err)
	}
	cloneJSON, err := json.Marshal(clone)
	if err != nil {
		t.Fatalf("json.Marshal(clone) returned unexpected error: %v", err)
	}
	if string(origJSON) != string(cloneJSON) {
		t.Errorf("Clone JSON = %q, want %q", string(cloneJSON), string(origJSON))
	}

	// Modifying clone must not affect original.
	clone.Set("b", 99).Delete("c").Set("d", 4)

	val, found := original.Get("b")
	if !found {
		t.Fatal("original.Get(\"b\") not found after clone modification")
	}
	if val != 2 {
		t.Errorf("original.Get(\"b\") = %v, want 2", val)
	}
	if !original.Has("c") {
		t.Error("original.Has(\"c\") after clone Delete = false, want true")
	}
	if original.Has("d") {
		t.Error("original.Has(\"d\") after clone Set = true, want false")
	}
}

func TestEntries(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](3).
		Set("a", 1).Set("b", 2).Set("c", 3)

	got := obj.Entries()
	wantKeys := []string{"a", "b", "c"}
	wantVals := []any{1, 2, 3}

	if len(got) != len(wantKeys) {
		t.Fatalf("len(Entries()) = %d, want %d", len(got), len(wantKeys))
	}
	for i, entry := range got {
		if entry.Key != wantKeys[i] {
			t.Errorf("Entries()[%d].Key = %q, want %q", i, entry.Key, wantKeys[i])
		}
		if entry.Value != wantVals[i] {
			t.Errorf("Entries()[%d].Value = %v, want %v", i, entry.Value, wantVals[i])
		}
	}

	// Modifying returned entries must not affect the object.
	got[0].Key = "x"
	got[0].Value = 99

	val, found := obj.Get("a")
	if !found {
		t.Fatal("Get(\"a\") not found after modifying returned entries")
	}
	if val != 1 {
		t.Errorf("Get(\"a\") = %v, want 1 (entries modification leaked)", val)
	}
}

func TestCapacity(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](3)
	if got := obj.Len(); got != 0 {
		t.Errorf("Len() on new object = %d, want 0", got)
	}

	// Add entries up to initial capacity.
	obj.Set("a", 1).Set("b", 2).Set("c", 3)
	if got := obj.Len(); got != 3 {
		t.Errorf("Len() after 3 Sets = %d, want 3", got)
	}

	// Add entries beyond initial capacity.
	obj.Set("d", 4).Set("e", 5).Set("f", 6).Set("g", 7)
	if got := obj.Len(); got != 7 {
		t.Errorf("Len() after 7 Sets = %d, want 7", got)
	}

	// Verify all values are correct.
	want := map[string]any{
		"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7,
	}
	for k, wantVal := range want {
		got, found := obj.Get(k)
		if !found {
			t.Errorf("Get(%q) not found", k)
			continue
		}
		if got != wantVal {
			t.Errorf("Get(%q) = %v, want %v", k, got, wantVal)
		}
	}

	// Verify insertion order is preserved.
	wantKeys := []string{"a", "b", "c", "d", "e", "f", "g"}
	entries := obj.Entries()
	if len(entries) != len(wantKeys) {
		t.Fatalf("len(Entries()) = %d, want %d", len(entries), len(wantKeys))
	}
	for i, entry := range entries {
		if entry.Key != wantKeys[i] {
			t.Errorf("Entries()[%d].Key = %q, want %q", i, entry.Key, wantKeys[i])
		}
	}
}

func TestJSONRoundtrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		json string
	}{
		{name: "simple object", json: `{"name":"John","age":30,"active":true,"score":95.5}`},
		{name: "empty object", json: `{}`},
		{name: "object with null values", json: `{"a":null,"b":null}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			obj, err := FromJSON[any]([]byte(tt.json))
			if err != nil {
				t.Fatalf("FromJSON() returned unexpected error: %v", err)
			}

			got, err := json.Marshal(obj)
			if err != nil {
				t.Fatalf("json.Marshal() returned unexpected error: %v", err)
			}
			if string(got) != tt.json {
				t.Errorf("roundtrip = %q, want %q", string(got), tt.json)
			}
		})
	}
}

func TestExplicitVsImplicitOrdering(t *testing.T) {
	t.Parallel()

	t.Run("explicit creation preserves order", func(t *testing.T) {
		t.Parallel()
		nested := NewObject[any]().
			Set("theme", "dark").
			Set("notifications", true)

		obj := NewObject[any]().
			Set("settings", nested).
			Set("version", "1.0")

		got, err := obj.ToJSON()
		if err != nil {
			t.Fatalf("ToJSON() returned unexpected error: %v", err)
		}
		if want := `{"settings":{"theme":"dark","notifications":true},"version":"1.0"}`; string(got) != want {
			t.Errorf("ToJSON() = %q, want %q", string(got), want)
		}
	})

	t.Run("JSON parsing uses standard behavior", func(t *testing.T) {
		t.Parallel()
		obj, err := FromJSON[any]([]byte(`{"settings":{"theme":"dark","notifications":true},"version":"1.0"}`))
		if err != nil {
			t.Fatalf("FromJSON() returned unexpected error: %v", err)
		}

		// Top-level object preserves order.
		entries := obj.Entries()
		if len(entries) != 2 {
			t.Fatalf("len(Entries()) = %d, want 2", len(entries))
		}
		if entries[0].Key != "settings" {
			t.Errorf("Entries()[0].Key = %q, want \"settings\"", entries[0].Key)
		}
		if entries[1].Key != "version" {
			t.Errorf("Entries()[1].Key = %q, want \"version\"", entries[1].Key)
		}

		// Nested objects are plain maps — order is not guaranteed.
		settings, ok := entries[0].Value.(map[string]any)
		if !ok {
			t.Fatalf("Entries()[0].Value type should be map[string]any")
		}
		if _, found := settings["theme"]; !found {
			t.Error("nested map missing key \"theme\"")
		}
		if _, found := settings["notifications"]; !found {
			t.Error("nested map missing key \"notifications\"")
		}
	})
}

func TestJSONTags(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Name      string `json:"full_name"`
		Age       int    `json:"age"`
		Email     string `json:"email,omitempty"`
		IsActive  bool   `json:"-"`
		SecretKey string `json:"secret_key,omitempty"`
	}

	tests := []struct {
		name  string
		input testStruct
		want  string
	}{
		{
			name: "all fields populated",
			input: testStruct{
				Name: "John Doe", Age: 30,
				Email: "john@example.com", IsActive: true,
				SecretKey: "abc123",
			},
			want: `{"full_name":"John Doe","age":30,"email":"john@example.com","secret_key":"abc123"}`,
		},
		{
			name: "empty fields omitted",
			input: testStruct{
				Name: "Jane Smith", Age: 25,
				IsActive: true,
			},
			want: `{"full_name":"Jane Smith","age":25}`,
		},
		{
			name: "skipped field excluded",
			input: testStruct{
				Name: "Bob Wilson", Age: 35,
				Email:     "bob@example.com",
				SecretKey: "xyz789",
			},
			want: `{"full_name":"Bob Wilson","age":35,"email":"bob@example.com","secret_key":"xyz789"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			obj := NewObject[testStruct](1).Set("user", tt.input)

			got, err := json.Marshal(obj)
			if err != nil {
				t.Fatalf("json.Marshal() returned unexpected error: %v", err)
			}
			want := `{"user":` + tt.want + `}`
			if string(got) != want {
				t.Errorf("json.Marshal() = %q, want %q", string(got), want)
			}

			// Roundtrip.
			var decoded Object[testStruct]
			err = json.Unmarshal(got, &decoded)
			if err != nil {
				t.Fatalf("json.Unmarshal() returned unexpected error: %v", err)
			}

			val, found := decoded.Get("user")
			if !found {
				t.Fatalf("decoded.Get(\"user\") not found")
			}
			// IsActive is tagged json:"-", so it should be zero after decode.
			expected := tt.input
			expected.IsActive = false
			if val != expected {
				t.Errorf("decoded.Get(\"user\") = %v, want %v", val, expected)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](3).
		Set("name", "John").Set("age", 30).Set("city", "New York")

	m := obj.ToMap()
	if len(m) != 3 {
		t.Fatalf("len(ToMap()) = %d, want 3", len(m))
	}

	wantMap := map[string]any{
		"name": "John", "age": 30, "city": "New York",
	}
	for k, want := range wantMap {
		got, ok := m[k]
		if !ok {
			t.Fatalf("ToMap() missing key %q", k)
		}
		if got != want {
			t.Errorf("ToMap()[%q] = %v, want %v", k, got, want)
		}
	}

	// Modifying map must not affect original object.
	m["age"] = 31
	got, found := obj.Get("age")
	if !found {
		t.Fatalf("Get(\"age\") not found after map modification")
	}
	if got != 30 {
		t.Errorf("Get(\"age\") = %v, want 30 (map modification leaked)", got)
	}
}

func TestToJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		obj  *Object[any]
		want string
	}{
		{
			name: "simple object",
			obj: NewObject[any](3).
				Set("name", "John").Set("age", 30).Set("city", "New York"),
			want: `{"name":"John","age":30,"city":"New York"}`,
		},
		{
			name: "empty object",
			obj:  NewObject[any](0),
			want: `{}`,
		},
		{
			name: "nested object",
			obj: func() *Object[any] {
				addr := NewObject[any](2).
					Set("street", "123 Main St").Set("zipcode", "10001")
				return NewObject[any](2).
					Set("name", "Alice").Set("address", addr)
			}(),
			want: `{"name":"Alice","address":{"street":"123 Main St","zipcode":"10001"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.obj.ToJSON()
			if err != nil {
				t.Fatalf("ToJSON() returned unexpected error: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("ToJSON() = %q, want %q", string(got), tt.want)
			}
		})
	}
}

func TestDeterministicMapOrdering(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](1).Set("fruits", map[string]any{
		"zebra": 1, "apple": 2, "mango": 3, "banana": 4,
		"cherry": 5, "dragon": 6, "elderberry": 7, "fig": 8,
	})

	// Marshal multiple times — output must be identical.
	var first string
	for i := range 10 {
		got, err := obj.ToJSON()
		if err != nil {
			t.Fatalf("ToJSON() iteration %d returned unexpected error: %v", i, err)
		}
		if i == 0 {
			first = string(got)
		}
		if string(got) != first {
			t.Errorf("ToJSON() iteration %d = %q, want %q", i, string(got), first)
		}
	}

	// Keys must appear in sorted order.
	appleIdx := strings.Index(first, `"apple"`)
	bananaIdx := strings.Index(first, `"banana"`)
	cherryIdx := strings.Index(first, `"cherry"`)

	if appleIdx >= bananaIdx {
		t.Errorf("apple index %d should be less than banana index %d", appleIdx, bananaIdx)
	}
	if bananaIdx >= cherryIdx {
		t.Errorf("banana index %d should be less than cherry index %d", bananaIdx, cherryIdx)
	}
}

func TestFromJSON_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantError error
	}{
		{name: "malformed JSON", input: `{"key": "value"`},
		{name: "array instead of object", input: `["a", "b"]`, wantError: ErrExpectedObjectStart},
		{name: "primitive instead of object", input: `"hello"`, wantError: ErrExpectedObjectStart},
		{name: "empty input", input: ``},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := FromJSON[any]([]byte(tt.input))
			if err == nil {
				t.Fatal("FromJSON() should return error")
			}
			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Errorf("expected error %v, got %v", tt.wantError, err)
				}
			}
		})
	}
}

func TestDuplicateKeys_Rejection(t *testing.T) {
	t.Parallel()

	_, err := FromJSON[any]([]byte(`{"key": "value1", "key": "value2"}`))
	if err == nil {
		t.Fatal("FromJSON() with duplicate keys should return error")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("error should contain 'duplicate', got: %v", err)
	}
}

func TestKeys(t *testing.T) {
	t.Parallel()

	t.Run("populated object", func(t *testing.T) {
		t.Parallel()
		obj := NewObject[any]().Set("a", 1).Set("b", 2).Set("c", 3)
		got := obj.Keys()
		want := []string{"a", "b", "c"}
		if len(got) != len(want) {
			t.Fatalf("len(Keys()) = %d, want %d", len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("Keys()[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})

	t.Run("empty object", func(t *testing.T) {
		t.Parallel()
		obj := NewObject[any]()
		if keys := obj.Keys(); len(keys) != 0 {
			t.Errorf("Keys() on empty = %v, want empty slice", keys)
		}
	})
}

func TestValues(t *testing.T) {
	t.Parallel()

	t.Run("populated object", func(t *testing.T) {
		t.Parallel()
		obj := NewObject[any]().Set("a", 1).Set("b", 2).Set("c", 3)
		got := obj.Values()
		want := []any{1, 2, 3}
		if len(got) != len(want) {
			t.Fatalf("len(Values()) = %d, want %d", len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("Values()[%d] = %v, want %v", i, got[i], want[i])
			}
		}
	})

	t.Run("empty object", func(t *testing.T) {
		t.Parallel()
		obj := NewObject[any]()
		if vals := obj.Values(); len(vals) != 0 {
			t.Errorf("Values() on empty = %v, want empty slice", vals)
		}
	})
}

func TestMarshalJSON_Direct(t *testing.T) {
	t.Parallel()

	obj := NewObject[any]().Set("x", 1).Set("y", 2)
	got, err := obj.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() returned unexpected error: %v", err)
	}
	if want := `{"x":1,"y":2}` + "\n"; string(got) != want {
		t.Errorf("MarshalJSON() = %q, want %q", string(got), want)
	}
}

func TestMarshalJSONTo_NestedOrderedMarshaler(t *testing.T) {
	t.Parallel()

	inner := NewObject[any]().Set("z", 3).Set("w", 4)
	outer := NewObject[any]().Set("nested", inner)

	got, err := outer.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() returned unexpected error: %v", err)
	}
	if want := `{"nested":{"z":3,"w":4}}` + "\n"; string(got) != want {
		t.Errorf("MarshalJSON() = %q, want %q", string(got), want)
	}
}

func TestMarshalJSONTo_UnmarshalableValue(t *testing.T) {
	t.Parallel()

	// Functions cannot be marshaled to JSON.
	obj := NewObject[any]().Set("fn", func() {})
	_, err := obj.MarshalJSON()
	if err == nil {
		t.Fatal("MarshalJSON() with unmarshalable value should return error")
	}
}

func TestUnmarshalJSONFrom_NotAnObject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantError error
	}{
		{
			name:      "number instead of object",
			input:     `42`,
			wantError: ErrExpectedObjectStart,
		},
		{
			name:      "boolean instead of object",
			input:     `true`,
			wantError: ErrExpectedObjectStart,
		},
		{
			name:      "null instead of object",
			input:     `null`,
			wantError: ErrExpectedObjectStart,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var obj Object[any]
			err := obj.UnmarshalJSON([]byte(tt.input))
			if err == nil {
				t.Fatal("UnmarshalJSON() should return error")
			}
			if !errors.Is(err, tt.wantError) {
				t.Errorf("expected error %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestUnmarshalJSON_ClearsExistingEntries(t *testing.T) {
	t.Parallel()

	obj := NewObject[any]().Set("old", "data")
	err := obj.UnmarshalJSON([]byte(`{"new":"value"}`))
	if err != nil {
		t.Fatalf("UnmarshalJSON() returned unexpected error: %v", err)
	}
	if obj.Has("old") {
		t.Error("Has(\"old\") after UnmarshalJSON = true, want false")
	}
	got, found := obj.Get("new")
	if !found {
		t.Fatal("Get(\"new\") not found")
	}
	if got != "value" {
		t.Errorf("Get(\"new\") = %v, want \"value\"", got)
	}
}

func TestNewObject_NoCapacity(t *testing.T) {
	t.Parallel()

	obj := NewObject[any]()
	obj.Set("a", 1)
	if obj.Len() != 1 {
		t.Errorf("Len() = %d, want 1", obj.Len())
	}
}

func TestSetUpdateExistingKey(t *testing.T) {
	t.Parallel()

	obj := NewObject[any]().Set("a", 1).Set("b", 2).Set("a", 99)

	got, found := obj.Get("a")
	if !found {
		t.Fatal("Get(\"a\") not found")
	}
	if got != 99 {
		t.Errorf("Get(\"a\") = %v, want 99", got)
	}
	// Order must be preserved: "a" stays at index 0.
	keys := obj.Keys()
	if keys[0] != "a" {
		t.Errorf("Keys()[0] = %q, want \"a\"", keys[0])
	}
	if obj.Len() != 2 {
		t.Errorf("Len() = %d, want 2", obj.Len())
	}
}

// Benchmark tests

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
		obj.Set("key50", 50) // Re-insert so next iteration deletes a real key.
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
