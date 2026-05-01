// Package orderedobject provides an ordered JSON object that preserves key
// insertion order during JSON marshaling and unmarshaling.
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
	// ErrExpectedObjectStart is returned when the next JSON token is not '{'.
	ErrExpectedObjectStart = errors.New("expected object start")
	// ErrExpectedStringKey is returned when the next JSON token is not a string key.
	ErrExpectedStringKey = errors.New("expected string key")
)

// OrderedMarshaler marshals a value to JSON while preserving key order.
type OrderedMarshaler interface {
	// MarshalJSONTo writes the JSON encoding of the value to enc.
	MarshalJSONTo(enc *jsontext.Encoder) error
}

// Entry holds one key-value pair in an Object.
type Entry[V any] struct {
	// Key is the object member name.
	Key string
	// Value is the value associated with Key.
	Value V
}

// Object stores key-value pairs in insertion order.
type Object[V any] struct {
	entries []Entry[V]
}

// NewObject returns an Object with optional initial capacity.
func NewObject[V any](capacity ...int) *Object[V] {
	n := 0
	if len(capacity) > 0 {
		n = capacity[0]
	}
	return &Object[V]{
		entries: make([]Entry[V], 0, n),
	}
}

// FromMap returns an Object containing the entries in m.
// The resulting entry order matches Go's randomized map iteration order.
func FromMap[V any](m map[string]V) *Object[V] {
	obj := NewObject[V](len(m))
	for k, v := range m {
		obj.entries = append(obj.entries, Entry[V]{Key: k, Value: v})
	}
	return obj
}

// FromJSON decodes data into an Object, preserving key order from the input.
// FromJSON returns an error if data is not valid JSON or does not encode an object.
func FromJSON[V any](data []byte) (*Object[V], error) {
	obj := NewObject[V]()
	if err := obj.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return obj, nil
}

func (o *Object[V]) findKeyIndex(key string) int {
	for i := range o.entries {
		if o.entries[i].Key == key {
			return i
		}
	}
	return -1
}

// Set stores value under key and returns o.
// If key already exists, Set updates its value without changing its position.
func (o *Object[V]) Set(key string, value V) *Object[V] {
	if idx := o.findKeyIndex(key); idx >= 0 {
		o.entries[idx].Value = value
		return o
	}
	o.entries = append(o.entries, Entry[V]{Key: key, Value: value})
	return o
}

// Get returns the value for key and whether key is present.
// If key is not present, Get returns the zero value of V and false.
func (o *Object[V]) Get(key string) (V, bool) {
	if idx := o.findKeyIndex(key); idx >= 0 {
		return o.entries[idx].Value, true
	}
	var zero V
	return zero, false
}

// Has reports whether key is present.
func (o *Object[V]) Has(key string) bool {
	return o.findKeyIndex(key) >= 0
}

// Delete removes key and returns o.
// Delete is a no-op if key is not present.
func (o *Object[V]) Delete(key string) *Object[V] {
	if idx := o.findKeyIndex(key); idx >= 0 {
		o.entries = slices.Delete(o.entries, idx, idx+1)
	}
	return o
}

// Len returns the number of entries in o.
func (o *Object[V]) Len() int {
	return len(o.entries)
}

// Keys returns a new slice of keys in insertion order.
func (o *Object[V]) Keys() []string {
	keys := make([]string, len(o.entries))
	for i := range o.entries {
		keys[i] = o.entries[i].Key
	}
	return keys
}

// Values returns a new slice of values in insertion order.
func (o *Object[V]) Values() []V {
	values := make([]V, len(o.entries))
	for i := range o.entries {
		values[i] = o.entries[i].Value
	}
	return values
}

// Entries returns a copy of o's entries in insertion order.
func (o *Object[V]) Entries() []Entry[V] {
	return slices.Clone(o.entries)
}

// ForEach calls fn for each entry in insertion order.
func (o *Object[V]) ForEach(fn func(key string, value V)) {
	for i := range o.entries {
		fn(o.entries[i].Key, o.entries[i].Value)
	}
}

// Clone returns a shallow copy of o.
func (o *Object[V]) Clone() *Object[V] {
	return &Object[V]{entries: slices.Clone(o.entries)}
}

// MarshalJSON returns the JSON encoding of o.
func (o *Object[V]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	enc := jsontext.NewEncoder(&buf)
	if err := o.MarshalJSONTo(enc); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJSONTo writes the JSON encoding of o to enc.
// Nested map values are encoded with deterministic key order.
func (o *Object[V]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	for i := range o.entries {
		entry := &o.entries[i]
		if err := enc.WriteToken(jsontext.String(entry.Key)); err != nil {
			return err
		}

		if marshaler, ok := any(entry.Value).(OrderedMarshaler); ok {
			if err := marshaler.MarshalJSONTo(enc); err != nil {
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

// UnmarshalJSON decodes a JSON object into o.
func (o *Object[V]) UnmarshalJSON(data []byte) error {
	dec := jsontext.NewDecoder(bytes.NewReader(data))
	if err := o.UnmarshalJSONFrom(dec); err != nil {
		return err
	}

	_, err := dec.ReadToken()
	switch {
	case errors.Is(err, io.EOF):
		return nil
	case err != nil:
		return err
	default:
		return io.ErrUnexpectedEOF
	}
}

func readObjectKey(dec *jsontext.Decoder) (string, error) {
	kind := dec.PeekKind()
	if kind == 0 {
		_, err := dec.ReadToken()
		return "", err
	}
	if kind != jsontext.KindString {
		return "", fmt.Errorf("expected string key, got %v: %w", kind, ErrExpectedStringKey)
	}

	tok, err := dec.ReadToken()
	if err != nil {
		return "", err
	}
	return tok.String(), nil
}

// UnmarshalJSONFrom decodes a JSON object from dec into o.
// UnmarshalJSONFrom replaces any existing entries in o.
func (o *Object[V]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	clear(o.entries)
	o.entries = o.entries[:0]

	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != jsontext.KindBeginObject {
		return fmt.Errorf("expected object start '{', got %v: %w", tok.Kind(), ErrExpectedObjectStart)
	}

	for dec.PeekKind() != jsontext.KindEndObject {
		key, err := readObjectKey(dec)
		if err != nil {
			return err
		}

		var value V
		if err := json.UnmarshalDecode(dec, &value); err != nil {
			return err
		}

		o.entries = append(o.entries, Entry[V]{Key: key, Value: value})
	}

	_, err = dec.ReadToken()
	return err
}

// ToMap returns a new map containing o's entries.
// The returned map does not preserve insertion order.
func (o *Object[V]) ToMap() map[string]V {
	m := make(map[string]V, len(o.entries))
	for i := range o.entries {
		m[o.entries[i].Key] = o.entries[i].Value
	}
	return m
}

// ToJSON returns the JSON encoding of o.
func (o *Object[V]) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}
