package orderedobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	json "github.com/go-json-experiment/json"
)

func TestMarshal(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		object   *Object[any]
		expected string
	}{
		{
			name: "Empty object",
			object: func() *Object[any] {
				return NewObject[any](0)
			}(),
			expected: `{}`,
		},
		{
			name: "Single key-value pair",
			object: func() *Object[any] {
				obj := NewObject[any](1)
				obj.Set("key", "value")
				return obj
			}(),
			expected: `{"key":"value"}`,
		},
		{
			name: "Multiple key-value pairs",
			object: func() *Object[any] {
				obj := NewObject[any](3)
				obj.Set("name", "John")
				obj.Set("age", 30)
				obj.Set("city", "New York")
				return obj
			}(),
			expected: `{"name":"John","age":30,"city":"New York"}`,
		},
		{
			name: "Nested objects",
			object: func() *Object[any] {
				address := NewObject[any](2)
				address.Set("street", "123 Main St")
				address.Set("city", "London")

				person := NewObject[any](3)
				person.Set("name", "Alice")
				person.Set("age", 28)
				person.Set("address", address)

				return person
			}(),
			expected: `{"name":"Alice","age":28,"address":{"street":"123 Main St","city":"London"}}`,
		},
		{
			name: "Array of objects",
			object: func() *Object[any] {
				person1 := NewObject[any](2)
				person1.Set("name", "Bob")
				person1.Set("age", 35)

				person2 := NewObject[any](2)
				person2.Set("name", "Charlie")
				person2.Set("age", 40)

				people := NewObject[any](1)
				people.Set("people", []any{person1, person2})
				return people
			}(),
			expected: `{"people":[{"name":"Bob","age":35},{"name":"Charlie","age":40}]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := json.Marshal(tc.object) // Use v2 Marshal
			require.NoError(t, err)
			assert.Equal(t, tc.expected, string(encoded))
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](3)
	obj.Set("string", "value")
	obj.Set("int", 42)
	obj.Set("bool", true)

	tests := []struct {
		name      string
		key       string
		wantValue any
		wantFound bool
	}{
		{
			name:      "String value",
			key:       "string",
			wantValue: "value",
			wantFound: true,
		},
		{
			name:      "Int value",
			key:       "int",
			wantValue: 42,
			wantFound: true,
		},
		{
			name:      "Bool value",
			key:       "bool",
			wantValue: true,
			wantFound: true,
		},
		{
			name:      "Non-existent key",
			key:       "missing",
			wantValue: nil,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotFound := obj.Get(tt.key)
			assert.Equal(t, tt.wantFound, gotFound)
			assert.Equal(t, tt.wantValue, gotValue)
		})
	}
}

func TestHas(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](2)
	obj.Set("key", "value")

	assert.True(t, obj.Has("key"))
	assert.False(t, obj.Has("missing"))
}

func TestDelete(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](3)
	obj.Set("a", 1)
	obj.Set("b", 2)
	obj.Set("c", 3)

	// Delete middle element and verify order is maintained
	obj.Delete("b")

	expectedJSON := `{"a":1,"c":3}`
	encoded, err := json.Marshal(obj)
	require.NoError(t, err)
	assert.Equal(t, expectedJSON, string(encoded))

	// Check key no longer exists
	assert.False(t, obj.Has("b"))

	// Delete non-existent key shouldn't affect object
	obj.Delete("missing")
	assert.Equal(t, 2, obj.Length())
}

func TestForEach(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](3)
	obj.Set("a", 1)
	obj.Set("b", 2)
	obj.Set("c", 3)

	keySum := ""
	valueSum := 0

	obj.ForEach(func(key string, value any) {
		keySum += key
		valueSum += value.(int)
	})

	assert.Equal(t, "abc", keySum)
	assert.Equal(t, 6, valueSum)

	// Verify values through Get
	value, found := obj.Get("a")
	assert.True(t, found)
	assert.Equal(t, 1, value)

	value, found = obj.Get("b")
	assert.True(t, found)
	assert.Equal(t, 2, value)

	value, found = obj.Get("c")
	assert.True(t, found)
	assert.Equal(t, 3, value)
}

func TestLength(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](0)
	assert.Equal(t, 0, obj.Length())

	obj.Set("a", 1)
	obj.Set("b", 2)
	assert.Equal(t, 2, obj.Length())

	obj.Delete("a")
	assert.Equal(t, 1, obj.Length())
}

func TestChaining(t *testing.T) {
	t.Parallel()

	// Test method chaining
	obj := NewObject[any](0).
		Set("a", 1).
		Set("b", 2).
		Set("c", 3)

	assert.Equal(t, 3, obj.Length())

	sum := 0
	obj.ForEach(func(key string, value any) {
		sum += value.(int)
	})

	assert.Equal(t, 6, sum)
}

func TestFromMap(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"name": "Alice",
		"age":  28,
		"city": "London",
	}

	obj := FromMap(m)

	// Check all key-value pairs exist
	for k, v := range m {
		assert.True(t, obj.Has(k))
		got, found := obj.Get(k)
		assert.True(t, found)
		assert.Equal(t, v, got)
	}

	// Check length
	assert.Equal(t, len(m), obj.Length())
}

func TestFromJSON(t *testing.T) {
	t.Parallel()

	jsonData := []byte(`{"name":"John","age":30,"city":"New York"}`)

	obj, err := FromJSON[any](jsonData)
	require.NoError(t, err)

	// Check all values
	name, found := obj.Get("name")
	assert.True(t, found)
	assert.Equal(t, "John", name)

	age, found := obj.Get("age")
	assert.True(t, found)
	assert.Equal(t, float64(30), age)

	city, found := obj.Get("city")
	assert.True(t, found)
	assert.Equal(t, "New York", city)

	// Check length
	assert.Equal(t, 3, obj.Length())
}

func TestClone(t *testing.T) {
	t.Parallel()

	original := NewObject[any](3)
	original.Set("a", 1)
	original.Set("b", 2)
	original.Set("c", 3)

	// Clone and check equal
	clone := original.Clone()

	// Check all values
	originalJSON, err := json.Marshal(original)
	require.NoError(t, err)
	cloneJSON, err := json.Marshal(clone)
	require.NoError(t, err)
	assert.Equal(t, string(originalJSON), string(cloneJSON))

	// Modifying clone shouldn't affect original
	clone.Set("b", 99)
	clone.Delete("c")
	clone.Set("d", 4)

	// Original should be unchanged
	value, found := original.Get("b")
	assert.True(t, found)
	assert.Equal(t, 2, value)
	assert.True(t, original.Has("c"))
	assert.False(t, original.Has("d"))
}

func TestEntries(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](3)
	obj.Set("a", 1)
	obj.Set("b", 2)
	obj.Set("c", 3)

	entries := obj.Entries()

	// Check all entries
	expected := []Entry[any]{
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
		{Key: "c", Value: 3},
	}

	assert.Equal(t, expected, entries)

	// Modifying returned entries shouldn't affect the original object
	entries[0].Key = "x"
	entries[0].Value = 99

	value, found := obj.Get("a")
	assert.True(t, found)
	assert.Equal(t, 1, value)
}

func TestCapacity(t *testing.T) {
	t.Parallel()

	// Create object with capacity 3
	obj := NewObject[any](3)

	// Verify initial state
	assert.Equal(t, 3, cap(obj.entries))
	assert.Equal(t, 0, obj.Length())

	// Add entries up to capacity
	obj.Set("a", 1)
	obj.Set("b", 2)
	obj.Set("c", 3)

	// Verify state at capacity
	assert.Equal(t, 3, cap(obj.entries))
	assert.Equal(t, 3, obj.Length())

	// Add more entries beyond capacity
	obj.Set("d", 4)
	obj.Set("e", 5)
	obj.Set("f", 6)
	obj.Set("g", 7)

	// Verify state after exceeding capacity
	assert.Equal(t, 7, obj.Length())
	assert.GreaterOrEqual(t, cap(obj.entries), 7)

	// Verify all values are correct
	expected := map[string]any{
		"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7,
	}
	for k, v := range expected {
		got, found := obj.Get(k)
		assert.True(t, found)
		assert.Equal(t, v, got)
	}

	// Verify order is preserved
	expectedOrder := []string{"a", "b", "c", "d", "e", "f", "g"}
	entries := obj.Entries()
	assert.Equal(t, len(expectedOrder), len(entries))
	for i, entry := range entries {
		assert.Equal(t, expectedOrder[i], entry.Key)
	}
}

func TestJSONRoundtrip(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		jsonStr string
	}{
		{
			name:    "Simple object",
			jsonStr: `{"name":"John","age":30,"active":true,"score":95.5}`,
		},
		// Removed: Nested object test - new design doesn't preserve order in JSON-parsed nested objects
		{
			name:    "Empty object",
			jsonStr: `{}`,
		},
		{
			name:    "Object with null values",
			jsonStr: `{"a":null,"b":null}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse JSON to Object
			obj, err := FromJSON[any]([]byte(tc.jsonStr))
			require.NoError(t, err)

			// Marshal back to JSON
			jsonBytes, err := json.Marshal(obj)
			require.NoError(t, err)

			// Compare JSON strings directly
			assert.Equal(t, tc.jsonStr, string(jsonBytes))
		})
	}
}

func TestExplicitVsImplicitOrdering(t *testing.T) {
	t.Parallel()

	t.Run("Explicit creation preserves order", func(t *testing.T) {
		// Explicitly create nested ordered objects
		nested := NewObject[any]().
			Set("theme", "dark").
			Set("notifications", true)

		main := NewObject[any]().
			Set("settings", nested).
			Set("version", "1.0")

		data, err := main.ToJSON()
		require.NoError(t, err)

		// Order should be preserved for explicitly created objects
		expected := `{"settings":{"theme":"dark","notifications":true},"version":"1.0"}`
		assert.Equal(t, expected, string(data))
	})

	t.Run("JSON parsing uses standard behavior", func(t *testing.T) {
		// Parse JSON - nested objects may not preserve order
		jsonStr := `{"settings":{"theme":"dark","notifications":true},"version":"1.0"}`
		obj, err := FromJSON[any]([]byte(jsonStr))
		require.NoError(t, err)

		// Top-level object preserves order
		entries := obj.Entries()
		assert.Equal(t, "settings", entries[0].Key)
		assert.Equal(t, "version", entries[1].Key)

		// But nested objects (maps) may not preserve order - this is expected
		settings := entries[0].Value.(map[string]any)
		assert.Contains(t, settings, "theme")
		assert.Contains(t, settings, "notifications")
		// We don't assert order for nested map since it's not guaranteed
	})
}

func TestJSONTags(t *testing.T) {
	t.Parallel()

	// Define a struct with various JSON tags
	type TestStruct struct {
		Name      string `json:"full_name"`            // Custom field name
		Age       int    `json:"age"`                  // Standard field name
		Email     string `json:"email,omitempty"`      // Omit if empty
		IsActive  bool   `json:"-"`                    // Skip this field
		SecretKey string `json:"secret_key,omitempty"` // Omit if empty
	}

	testCases := []struct {
		name     string
		input    TestStruct
		expected string
	}{
		{
			name: "All fields populated",
			input: TestStruct{
				Name:      "John Doe",
				Age:       30,
				Email:     "john@example.com",
				IsActive:  true,
				SecretKey: "abc123",
			},
			expected: `{"full_name":"John Doe","age":30,"email":"john@example.com","secret_key":"abc123"}`,
		},
		{
			name: "Empty fields should be omitted",
			input: TestStruct{
				Name:      "Jane Smith",
				Age:       25,
				Email:     "",
				IsActive:  true,
				SecretKey: "",
			},
			expected: `{"full_name":"Jane Smith","age":25}`,
		},
		{
			name: "IsActive should always be skipped",
			input: TestStruct{
				Name:      "Bob Wilson",
				Age:       35,
				Email:     "bob@example.com",
				IsActive:  false,
				SecretKey: "xyz789",
			},
			expected: `{"full_name":"Bob Wilson","age":35,"email":"bob@example.com","secret_key":"xyz789"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create an ordered object with the test struct
			obj := NewObject[TestStruct](1)
			obj.Set("user", tc.input)

			// Marshal to JSON
			data, err := json.Marshal(obj)
			require.NoError(t, err)

			// Compare with expected JSON
			expected := `{"user":` + tc.expected + `}`
			assert.Equal(t, expected, string(data))

			// Test roundtrip
			var newObj Object[TestStruct]
			err = json.Unmarshal(data, &newObj)
			require.NoError(t, err)

			// Verify the unmarshaled object
			value, found := newObj.Get("user")
			require.True(t, found, "Expected to find 'user' key in unmarshaled object")

			// Verify fields that should be preserved
			assert.Equal(t, tc.input.Name, value.Name)
			assert.Equal(t, tc.input.Age, value.Age)
			assert.Equal(t, tc.input.Email, value.Email)
			assert.Equal(t, tc.input.SecretKey, value.SecretKey)

			// Verify that IsActive is zero value (false) after unmarshal
			// This is expected behavior for fields with json:"-" tag
			assert.Equal(t, false, value.IsActive, "IsActive should be false (zero value)")
		})
	}
}

func TestToMap(t *testing.T) {
	t.Parallel()

	// Create test object
	obj := NewObject[any](3)
	obj.Set("name", "John").
		Set("age", 30).
		Set("city", "New York")

	// Convert to map
	m := obj.ToMap()

	// Verify map contents
	assert.Equal(t, 3, len(m))

	// Check all values exist
	expected := map[string]any{
		"name": "John",
		"age":  30,
		"city": "New York",
	}

	for k, v := range expected {
		got, ok := m[k]
		assert.True(t, ok)
		assert.Equal(t, v, got)
	}

	// Verify modifying map doesn't affect original object
	m["age"] = 31
	got, found := obj.Get("age")
	assert.True(t, found)
	assert.Equal(t, 30, got)
}

func TestToJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		object   *Object[any]
		expected string
	}{
		{
			name: "Simple object",
			object: func() *Object[any] {
				obj := NewObject[any](3)
				obj.Set("name", "John").
					Set("age", 30).
					Set("city", "New York")
				return obj
			}(),
			expected: `{"name":"John","age":30,"city":"New York"}`,
		},
		{
			name: "Empty object",
			object: func() *Object[any] {
				return NewObject[any](0)
			}(),
			expected: `{}`,
		},
		{
			name: "Nested object",
			object: func() *Object[any] {
				address := NewObject[any](2)
				address.Set("street", "123 Main St").
					Set("zipcode", "10001")

				person := NewObject[any](2)
				person.Set("name", "Alice").
					Set("address", address)
				return person
			}(),
			expected: `{"name":"Alice","address":{"street":"123 Main St","zipcode":"10001"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert to JSON - order should be preserved for explicitly created objects
			data, err := tc.object.ToJSON()
			require.NoError(t, err)

			// Compare with expected JSON
			assert.Equal(t, tc.expected, string(data))

			// Note: We don't test roundtrip for nested objects because
			// the new design only preserves order for explicitly created objects,
			// not for objects parsed from JSON
		})
	}
}
