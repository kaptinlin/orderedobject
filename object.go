// Package orderedobject provides an ordered JSON object that preserves key
// insertion order during JSON marshaling and unmarshaling.
package orderedobject

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"maps"
	"reflect"
	"slices"
)

var (
	// ErrDuplicateKey is returned when JSON or entry input contains a repeated key.
	ErrDuplicateKey = errors.New("duplicate key")
	errNilObject    = errors.New("nil object")
)

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

// New returns an empty Object.
func New[V any]() *Object[V] {
	return NewCap[V](0)
}

// NewCap returns an empty Object with capacity for n entries.
func NewCap[V any](n int) *Object[V] {
	n = max(n, 0)
	return &Object[V]{
		entries: make([]Entry[V], 0, n),
	}
}

// FromEntries returns an Object containing entries in their current order.
func FromEntries[V any](entries []Entry[V]) (*Object[V], error) {
	obj := NewCap[V](len(entries))
	for _, entry := range entries {
		if findEntryIndex(obj.entries, entry.Key) >= 0 {
			return nil, duplicateKeyError(entry.Key)
		}
		obj.entries = append(obj.entries, entry)
	}
	return obj, nil
}

// FromSortedMap returns an Object containing entries from m in lexical key order.
func FromSortedMap[V any](m map[string]V) *Object[V] {
	obj := NewCap[V](len(m))
	for _, key := range slices.Sorted(maps.Keys(m)) {
		obj.entries = append(obj.entries, Entry[V]{Key: key, Value: m[key]})
	}
	return obj
}

// FromJSON decodes data into an Object, preserving key order from the input.
// FromJSON returns an error if data is not valid JSON or does not encode an object.
func FromJSON[V any](data []byte) (*Object[V], error) {
	obj := New[V]()
	if err := obj.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return obj, nil
}

func (o *Object[V]) findKeyIndex(key string) int {
	if o == nil {
		return -1
	}
	return findEntryIndex(o.entries, key)
}

// Set stores value under key and returns o.
// If key already exists, Set updates its value without changing its position.
func (o *Object[V]) Set(key string, value V) *Object[V] {
	if o == nil {
		return nil
	}
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
	if o == nil {
		return nil
	}
	if idx := o.findKeyIndex(key); idx >= 0 {
		o.entries = slices.Delete(o.entries, idx, idx+1)
	}
	return o
}

// Len returns the number of entries in o.
func (o *Object[V]) Len() int {
	if o == nil {
		return 0
	}
	return len(o.entries)
}

// Keys returns a new slice of keys in insertion order.
func (o *Object[V]) Keys() []string {
	if o == nil {
		return []string{}
	}
	keys := make([]string, len(o.entries))
	for i := range o.entries {
		keys[i] = o.entries[i].Key
	}
	return keys
}

// Values returns a new slice of values in insertion order.
func (o *Object[V]) Values() []V {
	if o == nil {
		return []V{}
	}
	values := make([]V, len(o.entries))
	for i := range o.entries {
		values[i] = o.entries[i].Value
	}
	return values
}

// Entries returns a copy of o's entries in insertion order.
func (o *Object[V]) Entries() []Entry[V] {
	if o == nil {
		return []Entry[V]{}
	}
	return slices.Clone(o.entries)
}

// All returns an iterator over entries in insertion order.
// Mutating o during iteration is unsupported; use Entries for a snapshot.
func (o *Object[V]) All() iter.Seq2[string, V] {
	return func(yield func(string, V) bool) {
		if o == nil {
			return
		}
		for i := range o.entries {
			if !yield(o.entries[i].Key, o.entries[i].Value) {
				return
			}
		}
	}
}

// Clone returns a shallow copy of o.
func (o *Object[V]) Clone() *Object[V] {
	if o == nil {
		return New[V]()
	}
	return &Object[V]{entries: slices.Clone(o.entries)}
}

// MarshalJSON returns the JSON encoding of o.
func (o *Object[V]) MarshalJSON() ([]byte, error) {
	if o == nil {
		return []byte("null"), nil
	}

	var buf bytes.Buffer
	buf.WriteByte('{')
	for i := range o.entries {
		if i > 0 {
			buf.WriteByte(',')
		}

		entry := &o.entries[i]
		key, err := json.Marshal(entry.Key)
		if err != nil {
			return nil, fmt.Errorf("marshal key %q: %w", entry.Key, err)
		}
		buf.Write(key)
		buf.WriteByte(':')

		value, err := json.Marshal(entry.Value)
		if err != nil {
			return nil, fmt.Errorf("marshal value for key %q: %w", entry.Key, err)
		}
		buf.Write(value)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON decodes a JSON object into o.
func (o *Object[V]) UnmarshalJSON(data []byte) error {
	if o == nil {
		return errNilObject
	}

	var raw json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	entries, err := decodeEntries[V](raw)
	if err != nil {
		return err
	}
	o.replaceEntries(entries)
	return nil
}

func findEntryIndex[V any](entries []Entry[V], key string) int {
	for i := range entries {
		if entries[i].Key == key {
			return i
		}
	}
	return -1
}

func duplicateKeyError(key string) error {
	return fmt.Errorf("duplicate key %q: %w", key, ErrDuplicateKey)
}

func decodeEntries[V any](data []byte) ([]Entry[V], error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if tok != json.Delim('{') {
		return nil, &json.UnmarshalTypeError{
			Value: jsonTokenType(tok),
			Type:  reflect.TypeFor[Object[V]](),
		}
	}

	entries := make([]Entry[V], 0)
	for dec.More() {
		keyToken, err := dec.Token()
		if err != nil {
			return nil, err
		}
		key, ok := keyToken.(string)
		if !ok {
			return nil, fmt.Errorf("decode object key: expected string, got %T", keyToken)
		}

		if findEntryIndex(entries, key) >= 0 {
			return nil, duplicateKeyError(key)
		}

		var value V
		if err := dec.Decode(&value); err != nil {
			return nil, fmt.Errorf("decode value for key %q: %w", key, err)
		}

		entries = append(entries, Entry[V]{Key: key, Value: value})
	}

	if _, err = dec.Token(); err != nil {
		return nil, err
	}
	return entries, nil
}

func jsonTokenType(tok json.Token) string {
	switch tok := tok.(type) {
	case nil:
		return "null"
	case bool:
		return "bool"
	case float64, json.Number:
		return "number"
	case string:
		return "string"
	case json.Delim:
		switch tok {
		case '[':
			return "array"
		case '{':
			return "object"
		}
	}
	return fmt.Sprintf("%v", tok)
}

func (o *Object[V]) replaceEntries(entries []Entry[V]) {
	clear(o.entries)
	o.entries = entries
}

// ToUnorderedMap returns a new map containing o's entries.
// The returned map does not preserve insertion order.
func (o *Object[V]) ToUnorderedMap() map[string]V {
	if o == nil {
		return map[string]V{}
	}
	m := make(map[string]V, len(o.entries))
	for i := range o.entries {
		m[o.entries[i].Key] = o.entries[i].Value
	}
	return m
}
