package orderedobject

import (
	"fmt"
	"strings"
	"testing"

	json "github.com/go-json-experiment/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := json.Marshal(tc.obj)
			require.NoError(t, err, "json.Marshal() returned unexpected error")
			assert.Equal(t, tc.want, string(got))
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, found := obj.Get(tc.key)
			assert.Equal(t, tc.wantFound, found, "Get(%q) found", tc.key)
			assert.Equal(t, tc.wantValue, got, "Get(%q)", tc.key)
		})
	}
}

func TestHas(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](2).Set("key", "value")

	assert.True(t, obj.Has("key"))
	assert.False(t, obj.Has("missing"))
}

func TestDelete(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](3).
		Set("a", 1).Set("b", 2).Set("c", 3)

	// Delete middle element and verify order is maintained.
	obj.Delete("b")

	got, err := json.Marshal(obj)
	require.NoError(t, err, "json.Marshal() returned unexpected error")
	assert.Equal(t, `{"a":1,"c":3}`, string(got), "json.Marshal() after Delete")
	assert.False(t, obj.Has("b"), "Has(\"b\") after Delete")

	// Delete non-existent key should be a no-op.
	obj.Delete("missing")
	assert.Equal(t, 2, obj.Len(), "Len() after no-op Delete")
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

	assert.Equal(t, "abc", keys, "ForEach keys")
	assert.Equal(t, 6, sum, "ForEach sum")
}

func TestLen(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](0)
	assert.Equal(t, 0, obj.Len(), "Len() on empty")

	obj.Set("a", 1).Set("b", 2)
	assert.Equal(t, 2, obj.Len(), "Len() after two Sets")

	obj.Delete("a")
	assert.Equal(t, 1, obj.Len(), "Len() after Delete")
}

func TestChaining(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](0).
		Set("a", 1).Set("b", 2).Set("c", 3)

	assert.Equal(t, 3, obj.Len(), "Len()")

	got, err := json.Marshal(obj)
	require.NoError(t, err, "json.Marshal() returned unexpected error")
	assert.Equal(t, `{"a":1,"b":2,"c":3}`, string(got), "json.Marshal()")
}

func TestFromMap(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"name": "Alice",
		"age":  28,
		"city": "London",
	}

	obj := FromMap(m)
	assert.Equal(t, len(m), obj.Len(), "Len()")

	for k, want := range m {
		got, found := obj.Get(k)
		require.True(t, found, "Get(%q) not found", k)
		assert.Equal(t, want, got, "Get(%q)", k)
	}
}

func TestFromJSON(t *testing.T) {
	t.Parallel()

	obj, err := FromJSON[any]([]byte(`{"name":"John","age":30,"city":"New York"}`))
	require.NoError(t, err, "FromJSON() returned unexpected error")
	assert.Equal(t, 3, obj.Len(), "Len()")

	tests := []struct {
		key  string
		want any
	}{
		{key: "name", want: "John"},
		{key: "age", want: float64(30)},
		{key: "city", want: "New York"},
	}
	for _, tc := range tests {
		t.Run(tc.key, func(t *testing.T) {
			t.Parallel()
			got, found := obj.Get(tc.key)
			require.True(t, found, "Get(%q) not found", tc.key)
			assert.Equal(t, tc.want, got, "Get(%q)", tc.key)
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
	require.NoError(t, err, "json.Marshal(original) returned unexpected error")
	cloneJSON, err := json.Marshal(clone)
	require.NoError(t, err, "json.Marshal(clone) returned unexpected error")
	assert.Equal(t, string(origJSON), string(cloneJSON), "Clone JSON")

	// Modifying clone must not affect original.
	clone.Set("b", 99).Delete("c").Set("d", 4)

	val, found := original.Get("b")
	require.True(t, found, "original.Get(\"b\") not found after clone modification")
	assert.Equal(t, 2, val, "original.Get(\"b\")")
	assert.True(t, original.Has("c"), "original.Has(\"c\") after clone Delete")
	assert.False(t, original.Has("d"), "original.Has(\"d\") after clone Set")
}

func TestEntries(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](3).
		Set("a", 1).Set("b", 2).Set("c", 3)

	got := obj.Entries()
	wantKeys := []string{"a", "b", "c"}
	wantVals := []any{1, 2, 3}

	require.Equal(t, len(wantKeys), len(got), "len(Entries())")
	for i, entry := range got {
		assert.Equal(t, wantKeys[i], entry.Key, "Entries()[%d].Key", i)
		assert.Equal(t, wantVals[i], entry.Value, "Entries()[%d].Value", i)
	}

	// Modifying returned entries must not affect the object.
	got[0].Key = "x"
	got[0].Value = 99

	val, found := obj.Get("a")
	require.True(t, found, "Get(\"a\") not found after modifying returned entries")
	assert.Equal(t, 1, val, "Get(\"a\") (entries modification leaked)")
}

func TestCapacity(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](3)
	assert.Equal(t, 0, obj.Len(), "Len() on new object")

	// Add entries up to initial capacity.
	obj.Set("a", 1).Set("b", 2).Set("c", 3)
	assert.Equal(t, 3, obj.Len(), "Len() after 3 Sets")

	// Add entries beyond initial capacity.
	obj.Set("d", 4).Set("e", 5).Set("f", 6).Set("g", 7)
	assert.Equal(t, 7, obj.Len(), "Len() after 7 Sets")

	// Verify all values are correct.
	want := map[string]any{
		"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7,
	}
	for k, wantVal := range want {
		got, found := obj.Get(k)
		require.True(t, found, "Get(%q) not found", k)
		assert.Equal(t, wantVal, got, "Get(%q)", k)
	}

	// Verify insertion order is preserved.
	wantKeys := []string{"a", "b", "c", "d", "e", "f", "g"}
	entries := obj.Entries()
	require.Equal(t, len(wantKeys), len(entries), "len(Entries())")
	for i, entry := range entries {
		assert.Equal(t, wantKeys[i], entry.Key, "Entries()[%d].Key", i)
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			obj, err := FromJSON[any]([]byte(tc.json))
			require.NoError(t, err, "FromJSON() returned unexpected error")

			got, err := json.Marshal(obj)
			require.NoError(t, err, "json.Marshal() returned unexpected error")
			assert.Equal(t, tc.json, string(got), "roundtrip")
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
		require.NoError(t, err, "ToJSON() returned unexpected error")
		assert.Equal(t, `{"settings":{"theme":"dark","notifications":true},"version":"1.0"}`, string(got), "ToJSON()")
	})

	t.Run("JSON parsing uses standard behavior", func(t *testing.T) {
		t.Parallel()
		obj, err := FromJSON[any]([]byte(`{"settings":{"theme":"dark","notifications":true},"version":"1.0"}`))
		require.NoError(t, err, "FromJSON() returned unexpected error")

		// Top-level object preserves order.
		entries := obj.Entries()
		require.Equal(t, 2, len(entries), "len(Entries())")
		assert.Equal(t, "settings", entries[0].Key, "Entries()[0].Key")
		assert.Equal(t, "version", entries[1].Key, "Entries()[1].Key")

		// Nested objects are plain maps — order is not guaranteed.
		settings, ok := entries[0].Value.(map[string]any)
		require.True(t, ok, "Entries()[0].Value type should be map[string]any")
		assert.Contains(t, settings, "theme", "nested map missing key \"theme\"")
		assert.Contains(t, settings, "notifications", "nested map missing key \"notifications\"")
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			obj := NewObject[testStruct](1).Set("user", tc.input)

			got, err := json.Marshal(obj)
			require.NoError(t, err, "json.Marshal() returned unexpected error")
			want := `{"user":` + tc.want + `}`
			assert.Equal(t, want, string(got), "json.Marshal()")

			// Roundtrip.
			var decoded Object[testStruct]
			err = json.Unmarshal(got, &decoded)
			require.NoError(t, err, "json.Unmarshal() returned unexpected error")

			val, found := decoded.Get("user")
			require.True(t, found, "decoded.Get(\"user\") not found")
			// IsActive is tagged json:"-", so it should be zero after decode.
			expected := tc.input
			expected.IsActive = false
			assert.Equal(t, expected, val, "decoded.Get(\"user\")")
		})
	}
}

func TestToMap(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](3).
		Set("name", "John").Set("age", 30).Set("city", "New York")

	m := obj.ToMap()
	assert.Equal(t, 3, len(m), "len(ToMap())")

	wantMap := map[string]any{
		"name": "John", "age": 30, "city": "New York",
	}
	for k, want := range wantMap {
		got, ok := m[k]
		require.True(t, ok, "ToMap() missing key %q", k)
		assert.Equal(t, want, got, "ToMap()[%q]", k)
	}

	// Modifying map must not affect original object.
	m["age"] = 31
	got, found := obj.Get("age")
	require.True(t, found, "Get(\"age\") not found after map modification")
	assert.Equal(t, 30, got, "Get(\"age\") (map modification leaked)")
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := tc.obj.ToJSON()
			require.NoError(t, err, "ToJSON() returned unexpected error")
			assert.Equal(t, tc.want, string(got), "ToJSON()")
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
		require.NoError(t, err, "ToJSON() iteration %d returned unexpected error", i)
		if i == 0 {
			first = string(got)
		}
		assert.Equal(t, first, string(got), "ToJSON() iteration %d", i)
	}

	// Keys must appear in sorted order.
	appleIdx := strings.Index(first, `"apple"`)
	bananaIdx := strings.Index(first, `"banana"`)
	cherryIdx := strings.Index(first, `"cherry"`)

	assert.Less(t, appleIdx, bananaIdx, "apple index should be less than banana index")
	assert.Less(t, bananaIdx, cherryIdx, "banana index should be less than cherry index")
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := FromJSON[any]([]byte(tc.input))
			require.Error(t, err, "FromJSON() should return error")
			if tc.wantError != nil {
				assert.ErrorIs(t, err, tc.wantError)
			}
		})
	}
}

func TestDuplicateKeys_Rejection(t *testing.T) {
	t.Parallel()

	_, err := FromJSON[any]([]byte(`{"key": "value1", "key": "value2"}`))
	require.Error(t, err, "FromJSON() with duplicate keys should return error")
	assert.Contains(t, err.Error(), "duplicate", "error should contain 'duplicate'")
}

func TestKeys(t *testing.T) {
	t.Parallel()

	t.Run("populated object", func(t *testing.T) {
		t.Parallel()
		obj := NewObject[any]().Set("a", 1).Set("b", 2).Set("c", 3)
		got := obj.Keys()
		want := []string{"a", "b", "c"}
		require.Equal(t, len(want), len(got), "len(Keys())")
		for i := range want {
			assert.Equal(t, want[i], got[i], "Keys()[%d]", i)
		}
	})

	t.Run("empty object", func(t *testing.T) {
		t.Parallel()
		obj := NewObject[any]()
		assert.Empty(t, obj.Keys(), "Keys() on empty")
	})
}

func TestValues(t *testing.T) {
	t.Parallel()

	t.Run("populated object", func(t *testing.T) {
		t.Parallel()
		obj := NewObject[any]().Set("a", 1).Set("b", 2).Set("c", 3)
		got := obj.Values()
		want := []any{1, 2, 3}
		require.Equal(t, len(want), len(got), "len(Values())")
		for i := range want {
			assert.Equal(t, want[i], got[i], "Values()[%d]", i)
		}
	})

	t.Run("empty object", func(t *testing.T) {
		t.Parallel()
		obj := NewObject[any]()
		assert.Empty(t, obj.Values(), "Values() on empty")
	})
}

func TestMarshalJSON_Direct(t *testing.T) {
	t.Parallel()

	obj := NewObject[any]().Set("x", 1).Set("y", 2)
	got, err := obj.MarshalJSON()
	require.NoError(t, err, "MarshalJSON() returned unexpected error")
	assert.Equal(t, `{"x":1,"y":2}`+"\n", string(got), "MarshalJSON()")
}

func TestMarshalJSONTo_NestedOrderedMarshaler(t *testing.T) {
	t.Parallel()

	inner := NewObject[any]().Set("z", 3).Set("w", 4)
	outer := NewObject[any]().Set("nested", inner)

	got, err := outer.MarshalJSON()
	require.NoError(t, err, "MarshalJSON() returned unexpected error")
	assert.Equal(t, `{"nested":{"z":3,"w":4}}`+"\n", string(got), "MarshalJSON()")
}

func TestMarshalJSONTo_UnmarshalableValue(t *testing.T) {
	t.Parallel()

	// Functions cannot be marshaled to JSON.
	obj := NewObject[any]().Set("fn", func() {})
	_, err := obj.MarshalJSON()
	require.Error(t, err, "MarshalJSON() with unmarshalable value should return error")
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var obj Object[any]
			err := obj.UnmarshalJSON([]byte(tc.input))
			require.Error(t, err, "UnmarshalJSON() should return error")
			assert.ErrorIs(t, err, tc.wantError)
		})
	}
}

func TestUnmarshalJSON_ClearsExistingEntries(t *testing.T) {
	t.Parallel()

	obj := NewObject[any]().Set("old", "data")
	err := obj.UnmarshalJSON([]byte(`{"new":"value"}`))
	require.NoError(t, err, "UnmarshalJSON() returned unexpected error")
	assert.False(t, obj.Has("old"), "Has(\"old\") after UnmarshalJSON")
	got, found := obj.Get("new")
	require.True(t, found, "Get(\"new\") not found")
	assert.Equal(t, "value", got, "Get(\"new\")")
}

func TestNewObject_NoCapacity(t *testing.T) {
	t.Parallel()

	obj := NewObject[any]()
	obj.Set("a", 1)
	assert.Equal(t, 1, obj.Len(), "Len()")
}

func TestSetUpdateExistingKey(t *testing.T) {
	t.Parallel()

	obj := NewObject[any]().Set("a", 1).Set("b", 2).Set("a", 99)

	got, found := obj.Get("a")
	require.True(t, found, "Get(\"a\") not found")
	assert.Equal(t, 99, got, "Get(\"a\")")
	// Order must be preserved: "a" stays at index 0.
	keys := obj.Keys()
	assert.Equal(t, "a", keys[0], "Keys()[0]")
	assert.Equal(t, 2, obj.Len(), "Len()")
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
