// Package orderedobject provides an ordered JSON object that preserves
// insertion order of keys, designed to work with go-json-experiment/json.
package orderedobject

import (
	"bytes"
	"errors"
	"fmt"
	"slices"

	json "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

var (
	// ErrExpectedObjectStart is returned when the JSON token is not '{'.
	ErrExpectedObjectStart = errors.New("expected object start")
	// ErrExpectedStringKey is returned when the JSON token is not a string key.
	ErrExpectedStringKey = errors.New("expected string key")
)

// OrderedMarshaler is an interface for types that marshal themselves to JSON
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

// FromJSON parses a JSON byte slice into an ordered object.
// The order of keys is preserved as they appear in the JSON input.
func FromJSON[V any](data []byte) (*Object[V], error) {
	obj := NewObject[V]()
	if err := obj.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return obj, nil
}

// findKeyIndex returns the index of the key in the entries slice, or -1 if not found.
func (o *Object[V]) findKeyIndex(key string) int {
	for i, entry := range o.entries {
		if entry.Key == key {
			return i
		}
	}
	return -1
}

// Set sets the value for a key in the ordered object.
// If the key already exists, its value is updated in place.
// Otherwise, the key-value pair is appended to the end.
// Returns the object for method chaining.
func (o *Object[V]) Set(key string, value V) *Object[V] {
	if idx := o.findKeyIndex(key); idx >= 0 {
		o.entries[idx].Value = value
		return o
	}
	o.entries = append(o.entries, Entry[V]{Key: key, Value: value})
	return o
}

// Get returns the value for a key and whether the key exists.
// If the key does not exist, it returns the zero value and false.
func (o *Object[V]) Get(key string) (V, bool) {
	if idx := o.findKeyIndex(key); idx >= 0 {
		return o.entries[idx].Value, true
	}
	var zero V
	return zero, false
}

// Has reports whether the key exists in the ordered object.
func (o *Object[V]) Has(key string) bool {
	return o.findKeyIndex(key) >= 0
}

// Delete removes a key-value pair from the ordered object.
// If the key does not exist, it is a no-op.
// Returns the object for method chaining.
func (o *Object[V]) Delete(key string) *Object[V] {
	if idx := o.findKeyIndex(key); idx >= 0 {
		o.entries = slices.Delete(o.entries, idx, idx+1)
	}
	return o
}

// Len returns the number of key-value pairs in the ordered object.
func (o *Object[V]) Len() int {
	return len(o.entries)
}

// Keys returns all keys in insertion order.
func (o *Object[V]) Keys() []string {
	keys := make([]string, len(o.entries))
	for i, entry := range o.entries {
		keys[i] = entry.Key
	}
	return keys
}

// Values returns all values in insertion order.
func (o *Object[V]) Values() []V {
	values := make([]V, len(o.entries))
	for i, entry := range o.entries {
		values[i] = entry.Value
	}
	return values
}

// Entries returns a copy of all key-value pairs in insertion order.
func (o *Object[V]) Entries() []Entry[V] {
	return slices.Clone(o.entries)
}

// ForEach calls fn for each key-value pair in insertion order.
func (o *Object[V]) ForEach(fn func(key string, value V)) {
	for _, entry := range o.entries {
		fn(entry.Key, entry.Value)
	}
}

// Clone returns a shallow copy of the ordered object.
func (o *Object[V]) Clone() *Object[V] {
	return &Object[V]{entries: slices.Clone(o.entries)}
}

// MarshalJSON encodes the ordered object as JSON.
func (o *Object[V]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	enc := jsontext.NewEncoder(&buf)
	if err := o.MarshalJSONTo(enc); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJSONTo encodes the ordered object to a JSON encoder.
func (o *Object[V]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	for _, entry := range o.entries {
		if err := enc.WriteToken(jsontext.String(entry.Key)); err != nil {
			return err
		}

		// Preserve key order for nested ordered objects.
		if m, ok := any(entry.Value).(OrderedMarshaler); ok {
			if err := m.MarshalJSONTo(enc); err != nil {
				return err
			}
			continue
		}
		if err := json.MarshalEncode(enc, entry.Value, json.Deterministic(true)); err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndObject)
}

// UnmarshalJSON decodes a JSON object into the ordered object.
func (o *Object[V]) UnmarshalJSON(data []byte) error {
	dec := jsontext.NewDecoder(bytes.NewReader(data))
	return o.UnmarshalJSONFrom(dec)
}

// UnmarshalJSONFrom decodes a JSON object from a decoder into the ordered object.
func (o *Object[V]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	// Reset the object and clear old references for GC.
	clear(o.entries)
	o.entries = o.entries[:0]

	// Check for object start.
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != '{' {
		return fmt.Errorf("%w, got %v", ErrExpectedObjectStart, tok.Kind())
	}

	// Parse key-value pairs.
	for dec.PeekKind() != '}' {
		tok, err := dec.ReadToken()
		if err != nil {
			return err
		}
		if tok.Kind() != '"' {
			return fmt.Errorf("%w, got %v", ErrExpectedStringKey, tok.Kind())
		}
		key := tok.String()

		var value V
		if err := json.UnmarshalDecode(dec, &value); err != nil {
			return err
		}

		o.entries = append(o.entries, Entry[V]{Key: key, Value: value})
	}

	// Read the closing '}'.
	_, err = dec.ReadToken()
	return err
}

// ToMap converts the ordered object to a standard Go map.
// The returned map does not preserve insertion order.
func (o *Object[V]) ToMap() map[string]V {
	m := make(map[string]V, len(o.entries))
	for _, entry := range o.entries {
		m[entry.Key] = entry.Value
	}
	return m
}

// ToJSON converts the ordered object to a JSON byte slice.
func (o *Object[V]) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}
