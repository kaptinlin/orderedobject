// Package orderedobject provides an ordered JSON object implementation
// designed to work with github.com/go-json-experiment/json.
package orderedobject

import (
	"bytes"
	"errors"
	"fmt"

	json "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

var (
	// ErrExpectedObjectStart is returned when the JSON token is not an object start
	ErrExpectedObjectStart = errors.New("expected object start")
	// ErrExpectedStringKey is returned when the JSON token is not a string key
	ErrExpectedStringKey = errors.New("expected string key")
)

// OrderedMarshaler is an interface for objects that can marshal themselves to JSON
// while preserving key order.
type OrderedMarshaler interface {
	MarshalJSONTo(enc *jsontext.Encoder) error
}

// Entry represents a key-value pair.
type Entry[V any] struct {
	Key   string
	Value V
}

// Object is an ordered JSON object that preserves insertion order.
type Object[V any] struct {
	entries []Entry[V]
}

// NewObject returns an ordered object with optional pre-allocated capacity.
func NewObject[V any](capacity ...int) *Object[V] {
	cap := 0
	if len(capacity) > 0 {
		cap = capacity[0]
	}
	return &Object[V]{
		entries: make([]Entry[V], 0, cap),
	}
}

// FromMap creates an ordered object from a map.
// The order of the keys will be determined by the map iteration order.
func FromMap[V any](m map[string]V) *Object[V] {
	obj := NewObject[V](len(m))
	for k, v := range m {
		obj.Set(k, v)
	}
	return obj
}

// FromJSON parses a JSON string into an ordered object.
// The order of keys will be preserved as in the JSON string.
func FromJSON[V any](data []byte) (*Object[V], error) {
	obj := NewObject[V]()
	err := obj.UnmarshalJSON(data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return obj, nil
}

// findKeyIndex returns the index of the key in the entries slice, or -1 if not found.
func (object *Object[V]) findKeyIndex(key string) int {
	for i, entry := range object.entries {
		if entry.Key == key {
			return i
		}
	}
	return -1
}

// Set sets the value for a key in the ordered object.
// If the key already exists, its value is updated.
// Otherwise, the key-value pair is appended to the end.
// Returns the object for chaining.
func (object *Object[V]) Set(key string, value V) *Object[V] {
	if idx := object.findKeyIndex(key); idx >= 0 {
		object.entries[idx].Value = value
	} else {
		object.entries = append(object.entries, Entry[V]{Key: key, Value: value})
	}
	return object
}

// Get returns the value for a key and whether the key exists.
// If the key does not exist, it returns the zero value and false.
func (object *Object[V]) Get(key string) (V, bool) {
	if idx := object.findKeyIndex(key); idx >= 0 {
		return object.entries[idx].Value, true
	}
	var zero V
	return zero, false
}

// Has returns whether the key exists in the ordered object.
func (object *Object[V]) Has(key string) bool {
	return object.findKeyIndex(key) >= 0
}

// Delete removes a key-value pair from the ordered object.
// If the key does not exist, it does nothing.
// Returns the object for chaining.
func (object *Object[V]) Delete(key string) *Object[V] {
	if idx := object.findKeyIndex(key); idx >= 0 {
		object.entries = append(object.entries[:idx], object.entries[idx+1:]...)
	}
	return object
}

// Length returns the number of key-value pairs in the ordered object.
func (object *Object[V]) Length() int {
	return len(object.entries)
}

// Keys returns all keys in the ordered object.
func (object *Object[V]) Keys() []string {
	keys := make([]string, len(object.entries))
	for i, entry := range object.entries {
		keys[i] = entry.Key
	}
	return keys
}

// Values returns all values in the ordered object.
func (object *Object[V]) Values() []V {
	values := make([]V, len(object.entries))
	for i, entry := range object.entries {
		values[i] = entry.Value
	}
	return values
}

// Entries returns all key-value pairs in the ordered object.
func (object *Object[V]) Entries() []Entry[V] {
	entries := make([]Entry[V], len(object.entries))
	copy(entries, object.entries)
	return entries
}

// ForEach executes a function for each key-value pair in the ordered object.
func (object *Object[V]) ForEach(fn func(key string, value V)) {
	for _, entry := range object.entries {
		fn(entry.Key, entry.Value)
	}
}

// Clone returns a deep copy of the ordered object.
func (object *Object[V]) Clone() *Object[V] {
	entries := make([]Entry[V], len(object.entries))
	copy(entries, object.entries)
	return &Object[V]{entries: entries}
}

// MarshalJSON encodes the ordered object as JSON.
func (object *Object[V]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	enc := jsontext.NewEncoder(&buf)
	if err := object.MarshalJSONTo(enc); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJSONTo encodes the ordered object to a JSON encoder.
func (object *Object[V]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	for _, entry := range object.entries {
		if err := enc.WriteToken(jsontext.String(entry.Key)); err != nil {
			return err
		}

		// Check if value implements OrderedMarshaler and handle it specially
		if orderedMarshaler, ok := any(entry.Value).(OrderedMarshaler); ok {
			if err := orderedMarshaler.MarshalJSONTo(enc); err != nil {
				return err
			}
		} else {
			// Use Deterministic option to ensure nested maps have consistent ordering
			if err := json.MarshalEncode(enc, entry.Value, json.Deterministic(true)); err != nil {
				return err
			}
		}
	}
	return enc.WriteToken(jsontext.EndObject)
}

// UnmarshalJSON decodes a JSON object into the ordered object.
func (object *Object[V]) UnmarshalJSON(data []byte) error {
	dec := jsontext.NewDecoder(bytes.NewReader(data))
	return object.UnmarshalJSONFrom(dec)
}

// UnmarshalJSONFrom decodes a JSON object from a decoder into the ordered object.
func (object *Object[V]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	// Reset the object
	object.entries = object.entries[:0]

	// Check for object start
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != '{' {
		return fmt.Errorf("%w, got %v", ErrExpectedObjectStart, tok.Kind())
	}

	// Parse key-value pairs
	for dec.PeekKind() != '}' {
		// Read key
		tok, err := dec.ReadToken()
		if err != nil {
			return err
		}
		if tok.Kind() != '"' {
			return fmt.Errorf("%w, got %v", ErrExpectedStringKey, tok.Kind())
		}
		key := tok.String()

		// Read value
		var value V
		if err := json.UnmarshalDecode(dec, &value); err != nil {
			return err
		}

		// Add to entries
		object.entries = append(object.entries, Entry[V]{Key: key, Value: value})
	}

	// Read the closing '}'
	if _, err := dec.ReadToken(); err != nil {
		return err
	}

	return nil
}

// ToMap converts the ordered object to a standard Go map.
// The returned map will not preserve the insertion order.
func (object *Object[V]) ToMap() map[string]V {
	m := make(map[string]V, len(object.entries))
	for _, entry := range object.entries {
		m[entry.Key] = entry.Value
	}
	return m
}

// ToJSON converts the ordered object to a JSON byte slice.
// This is a convenience method that internally uses json.Marshal.
func (object *Object[V]) ToJSON() ([]byte, error) {
	return json.Marshal(object)
}
