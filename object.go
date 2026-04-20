// Package orderedobject provides an ordered JSON object that preserves
// insertion order of keys, designed to work with go-json-experiment/json.
package orderedobject

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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
	n := 0
	if len(capacity) > 0 {
		n = capacity[0]
	}
	return &Object[V]{
		entries: make([]Entry[V], 0, n),
	}
}

// FromMap creates an ordered object from a map.
// Key order is determined by Go's map iteration order, which is randomized.
func FromMap[V any](m map[string]V) *Object[V] {
	obj := NewObject[V](len(m))
	for k, v := range m {
		obj.entries = append(obj.entries, Entry[V]{Key: k, Value: v})
	}
	return obj
}

// FromJSON parses JSON data into an ordered object.
// Key order is preserved as it appears in the JSON input.
// Returns an error if data is not valid JSON or not an object.
func FromJSON[V any](data []byte) (*Object[V], error) {
	obj := NewObject[V]()
	if err := obj.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return obj, nil
}

// findKeyIndex returns the index of the key in the entries slice, or -1 if not found.
func (o *Object[V]) findKeyIndex(key string) int {
	return slices.IndexFunc(o.entries, func(e Entry[V]) bool {
		return e.Key == key
	})
}

// Set associates value with key in the ordered object.
// If key already exists, its value is updated in place without changing position.
// If key is new, the key-value pair is appended to the end.
// Returns the object to enable method chaining.
func (o *Object[V]) Set(key string, value V) *Object[V] {
	if idx := o.findKeyIndex(key); idx >= 0 {
		o.entries[idx].Value = value
		return o
	}
	o.entries = append(o.entries, Entry[V]{Key: key, Value: value})
	return o
}

// Get returns the value associated with key and a boolean indicating whether
// the key was found. If the key does not exist, returns the zero value for V
// and false.
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

// Delete removes the key-value pair associated with key from the ordered object.
// If key does not exist, this is a no-op.
// Returns the object to enable method chaining.
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
	if err := o.UnmarshalJSONFrom(dec); err != nil {
		return err
	}

	_, err := dec.ReadToken()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return io.ErrUnexpectedEOF
}

// UnmarshalJSONFrom decodes a JSON object from a decoder into the ordered object.
func (o *Object[V]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	clear(o.entries)
	o.entries = o.entries[:0]

	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != '{' {
		return fmt.Errorf("expected object start '{', got %v: %w", tok.Kind(), ErrExpectedObjectStart)
	}

	for dec.PeekKind() != '}' {
		if kind := dec.PeekKind(); kind != '"' && kind != 0 {
			return fmt.Errorf("expected string key, got %v: %w", kind, ErrExpectedStringKey)
		}

		tok, err := dec.ReadToken()
		if err != nil {
			return err
		}
		if tok.Kind() != '"' {
			return fmt.Errorf("expected string key, got %v: %w", tok.Kind(), ErrExpectedStringKey)
		}
		key := tok.String()

		var value V
		if err := json.UnmarshalDecode(dec, &value); err != nil {
			return err
		}

		o.entries = append(o.entries, Entry[V]{Key: key, Value: value})
	}

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
