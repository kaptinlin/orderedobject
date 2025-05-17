package orderedobject

import (
	"reflect"
	"testing"

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
		tc := tc // capture
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := json.Marshal(tc.object) // Use v2 Marshal
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := string(encoded); got != tc.expected {
				t.Errorf("want %s, got %s", tc.expected, got)
			}
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
			if gotFound != tt.wantFound {
				t.Errorf("Get() found = %v, want %v", gotFound, tt.wantFound)
			}
			if gotValue != tt.wantValue {
				t.Errorf("Get() value = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestHas(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](2)
	obj.Set("key", "value")

	if !obj.Has("key") {
		t.Error("Has() = false, want true")
	}

	if obj.Has("missing") {
		t.Error("Has() = true, want false")
	}
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
	encoded, _ := json.Marshal(obj)
	if string(encoded) != expectedJSON {
		t.Errorf("After Delete() JSON = %s, want %s", encoded, expectedJSON)
	}

	// Check key no longer exists
	if obj.Has("b") {
		t.Error("Key 'b' should have been deleted")
	}

	// Delete non-existent key shouldn't affect object
	obj.Delete("missing")
	if obj.Length() != 2 {
		t.Errorf("Length after deleting non-existent key = %d, want %d", obj.Length(), 2)
	}
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

	if keySum != "abc" {
		t.Errorf("ForEach() key sum = %s, want %s", keySum, "abc")
	}

	if valueSum != 6 {
		t.Errorf("ForEach() value sum = %d, want %d", valueSum, 6)
	}

	// Verify values through Get
	if got, found := obj.Get("a"); !found || got != 1 {
		t.Errorf("Value for 'a' = %v (found: %v), want %v", got, found, 1)
	}
	if got, found := obj.Get("b"); !found || got != 2 {
		t.Errorf("Value for 'b' = %v (found: %v), want %v", got, found, 2)
	}
	if got, found := obj.Get("c"); !found || got != 3 {
		t.Errorf("Value for 'c' = %v (found: %v), want %v", got, found, 3)
	}
}

func TestLength(t *testing.T) {
	t.Parallel()

	obj := NewObject[any](0)
	if obj.Length() != 0 {
		t.Errorf("Length() = %d, want 0", obj.Length())
	}

	obj.Set("a", 1)
	obj.Set("b", 2)
	if obj.Length() != 2 {
		t.Errorf("Length() = %d, want 2", obj.Length())
	}

	obj.Delete("a")
	if obj.Length() != 1 {
		t.Errorf("Length() = %d, want 1", obj.Length())
	}
}

func TestChaining(t *testing.T) {
	t.Parallel()

	// Test method chaining
	obj := NewObject[any](0).
		Set("a", 1).
		Set("b", 2).
		Set("c", 3)

	if obj.Length() != 3 {
		t.Errorf("Chaining length = %d, want %d", obj.Length(), 3)
	}

	sum := 0
	obj.ForEach(func(key string, value any) {
		sum += value.(int)
	})

	if sum != 6 {
		t.Errorf("Chaining sum = %d, want %d", sum, 6)
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

	// Check all key-value pairs exist
	for k, v := range m {
		if !obj.Has(k) {
			t.Errorf("Key %s not found in object", k)
		}
		if got, found := obj.Get(k); !found || got != v {
			t.Errorf("Value for key %s = %v (found: %v), want %v", k, got, found, v)
		}
	}

	// Check length
	if obj.Length() != len(m) {
		t.Errorf("Object length = %d, want %d", obj.Length(), len(m))
	}
}

func TestFromJSON(t *testing.T) {
	t.Parallel()

	jsonData := []byte(`{"name":"John","age":30,"city":"New York"}`)

	obj, err := FromJSON[any](jsonData)
	if err != nil {
		t.Fatalf("FromJSON() error = %v", err)
	}

	// Check all values
	if got, found := obj.Get("name"); !found || got != "John" {
		t.Errorf("Value for 'name' = %v (found: %v), want %v", got, found, "John")
	}
	if got, found := obj.Get("age"); !found || got != float64(30) {
		t.Errorf("Value for 'age' = %v (found: %v), want %v", got, found, float64(30))
	}
	if got, found := obj.Get("city"); !found || got != "New York" {
		t.Errorf("Value for 'city' = %v (found: %v), want %v", got, found, "New York")
	}

	// Check length
	if obj.Length() != 3 {
		t.Errorf("Object length = %d, want %d", obj.Length(), 3)
	}
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
	originalJSON, _ := json.Marshal(original)
	cloneJSON, _ := json.Marshal(clone)

	if string(originalJSON) != string(cloneJSON) {
		t.Errorf("Clone JSON = %s, want %s", cloneJSON, originalJSON)
	}

	// Modifying clone shouldn't affect original
	clone.Set("b", 99)
	clone.Delete("c")
	clone.Set("d", 4)

	// Original should be unchanged
	if got, found := original.Get("b"); !found || got != 2 {
		t.Errorf("Original value for 'b' = %v (found: %v), want %v", got, found, 2)
	}
	if !original.Has("c") {
		t.Errorf("Original should still have key 'c'")
	}
	if original.Has("d") {
		t.Errorf("Original should not have key 'd'")
	}
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

	if !reflect.DeepEqual(entries, expected) {
		t.Errorf("Entries() = %v, want %v", entries, expected)
	}

	// Modifying returned entries shouldn't affect the original object
	entries[0].Key = "x"
	entries[0].Value = 99

	if got, found := obj.Get("a"); !found || got != 1 {
		t.Errorf("Original value for 'a' = %v (found: %v), want %v", got, found, 1)
	}
}

func TestCapacity(t *testing.T) {
	t.Parallel()

	// Create object with capacity 3
	obj := NewObject[any](3)

	// Verify initial state
	if cap(obj.entries) != 3 {
		t.Errorf("Initial capacity = %d, want %d", cap(obj.entries), 3)
	}
	if obj.Length() != 0 {
		t.Errorf("Initial length = %d, want %d", obj.Length(), 0)
	}

	// Add entries up to capacity
	obj.Set("a", 1)
	obj.Set("b", 2)
	obj.Set("c", 3)

	// Verify state at capacity
	if cap(obj.entries) != 3 {
		t.Errorf("Capacity at limit = %d, want %d", cap(obj.entries), 3)
	}
	if obj.Length() != 3 {
		t.Errorf("Length at limit = %d, want %d", obj.Length(), 3)
	}

	// Add more entries beyond capacity
	obj.Set("d", 4)
	obj.Set("e", 5)
	obj.Set("f", 6)
	obj.Set("g", 7)

	// Verify state after exceeding capacity
	if obj.Length() != 7 {
		t.Errorf("Length after exceeding capacity = %d, want %d", obj.Length(), 7)
	}
	if cap(obj.entries) < 7 {
		t.Errorf("Capacity after exceeding limit = %d, want >= %d", cap(obj.entries), 7)
	}

	// Verify all values are correct
	expected := map[string]any{
		"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7,
	}
	for k, v := range expected {
		if got, found := obj.Get(k); !found || got != v {
			t.Errorf("Value for key %s = %v (found: %v), want %v", k, got, found, v)
		}
	}

	// Verify order is preserved
	expectedOrder := []string{"a", "b", "c", "d", "e", "f", "g"}
	entries := obj.Entries()
	if len(entries) != len(expectedOrder) {
		t.Errorf("Number of entries = %d, want %d", len(entries), len(expectedOrder))
	}
	for i, entry := range entries {
		if entry.Key != expectedOrder[i] {
			t.Errorf("Entry[%d].Key = %s, want %s", i, entry.Key, expectedOrder[i])
		}
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
		{
			name:    "Nested object",
			jsonStr: `{"user":{"name":"Alice","age":28},"settings":{"theme":"dark","notifications":true},"tags":["go","json"]}`,
		},
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
		tc := tc // capture
		t.Run(tc.name, func(t *testing.T) {
			// Parse JSON to Object
			obj, err := FromJSON[any]([]byte(tc.jsonStr))
			if err != nil {
				t.Fatalf("FromJSON() error = %v", err)
			}

			// Marshal back to JSON
			jsonBytes, err := json.Marshal(obj)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			// Compare JSON strings directly
			if got := string(jsonBytes); got != tc.jsonStr {
				t.Errorf("JSON mismatch:\nOriginal: %s\nGot:      %s", tc.jsonStr, got)
			}
		})
	}
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
		tc := tc // capture
		t.Run(tc.name, func(t *testing.T) {
			// Create an ordered object with the test struct
			obj := NewObject[TestStruct](1)
			obj.Set("user", tc.input)

			// Marshal to JSON
			data, err := json.Marshal(obj)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			// Compare with expected JSON
			expected := `{"user":` + tc.expected + `}`
			if got := string(data); got != expected {
				t.Errorf("JSON mismatch:\nExpected: %s\nGot:      %s", expected, got)
			}

			// Test roundtrip
			var newObj Object[TestStruct]
			if err := json.Unmarshal(data, &newObj); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			// Verify the unmarshaled object
			if value, found := newObj.Get("user"); found {
				// Verify fields that should be preserved
				if value.Name != tc.input.Name {
					t.Errorf("Name = %v, want %v", value.Name, tc.input.Name)
				}
				if value.Age != tc.input.Age {
					t.Errorf("Age = %v, want %v", value.Age, tc.input.Age)
				}
				if value.Email != tc.input.Email {
					t.Errorf("Email = %v, want %v", value.Email, tc.input.Email)
				}
				if value.SecretKey != tc.input.SecretKey {
					t.Errorf("SecretKey = %v, want %v", value.SecretKey, tc.input.SecretKey)
				}

				// Verify that IsActive is zero value (false) after unmarshal
				// This is expected behavior for fields with json:"-" tag
				if value.IsActive != false {
					t.Errorf("IsActive = %v, want false (zero value)", value.IsActive)
				}
			} else {
				t.Error("Expected to find 'user' key in unmarshaled object")
			}
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
	if len(m) != 3 {
		t.Errorf("Map length = %d, want %d", len(m), 3)
	}

	// Check all values exist
	expected := map[string]any{
		"name": "John",
		"age":  30,
		"city": "New York",
	}

	for k, v := range expected {
		if got, ok := m[k]; !ok || got != v {
			t.Errorf("Map[%s] = %v (found: %v), want %v", k, got, ok, v)
		}
	}

	// Verify modifying map doesn't affect original object
	m["age"] = 31
	if got, found := obj.Get("age"); !found || got != 30 {
		t.Errorf("Original object age = %v (found: %v), want %v", got, found, 30)
	}
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
		tc := tc // capture
		t.Run(tc.name, func(t *testing.T) {
			// Convert to JSON
			data, err := tc.object.ToJSON()
			if err != nil {
				t.Fatalf("ToJSON() error = %v", err)
			}

			// Compare with expected JSON
			if got := string(data); got != tc.expected {
				t.Errorf("JSON mismatch:\nExpected: %s\nGot:      %s", tc.expected, got)
			}

			// Verify roundtrip
			newObj, err := FromJSON[any](data)
			if err != nil {
				t.Fatalf("FromJSON() error = %v", err)
			}

			// Compare original and new object
			originalJSON, _ := tc.object.ToJSON()
			newJSON, _ := newObj.ToJSON()
			if string(originalJSON) != string(newJSON) {
				t.Errorf("Roundtrip JSON mismatch:\nOriginal: %s\nNew:      %s",
					string(originalJSON), string(newJSON))
			}
		})
	}
}
