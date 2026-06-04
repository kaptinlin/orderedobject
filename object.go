// Package orderedobject provides an ordered JSON object that preserves key
// insertion order during JSON marshaling and unmarshaling.
package orderedobject

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"maps"
	"reflect"
	"slices"

	json "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

var (
	// ErrExpectedObjectStart is returned when the next JSON token is not '{'.
	ErrExpectedObjectStart = errors.New("expected object start")
	// ErrExpectedStringKey is returned when the next JSON token is not a string key.
	ErrExpectedStringKey = errors.New("expected string key")
	// ErrDuplicateKey is returned when JSON or entry input contains a repeated key.
	ErrDuplicateKey = errors.New("duplicate key")
	// ErrTrailingToken is returned when JSON input contains data after the top-level object.
	ErrTrailingToken = errors.New("trailing token")
	// ErrNilJSONEncoder is returned when JSON encoding is asked to use a nil encoder.
	ErrNilJSONEncoder = errors.New("nil JSON encoder")
	// ErrNilJSONDecoder is returned when JSON decoding is asked to use a nil decoder.
	ErrNilJSONDecoder = errors.New("nil JSON decoder")

	errNilObject = errors.New("nil object")
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
		if obj.findKeyIndex(entry.Key) >= 0 {
			return nil, fmt.Errorf("duplicate key %q: %w", entry.Key, ErrDuplicateKey)
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

// FromUnorderedMap returns an Object containing entries from m in Go map iteration order.
func FromUnorderedMap[V any](m map[string]V) *Object[V] {
	obj := NewCap[V](len(m))
	for k, v := range m {
		obj.entries = append(obj.entries, Entry[V]{Key: k, Value: v})
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

// ForEach calls fn for each entry in insertion order.
func (o *Object[V]) ForEach(fn func(key string, value V)) {
	if o == nil {
		return
	}
	for _, entry := range slices.Clone(o.entries) {
		fn(entry.Key, entry.Value)
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
	var buf bytes.Buffer
	enc := jsontext.NewEncoder(&buf)
	if err := o.MarshalJSONTo(enc); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}

// MarshalJSONTo writes the JSON encoding of o to enc.
// Nested map values are encoded with deterministic key order.
func (o *Object[V]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if enc == nil {
		return ErrNilJSONEncoder
	}
	if o == nil {
		return enc.WriteToken(jsontext.Null)
	}
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	for i := range o.entries {
		entry := &o.entries[i]
		if err := enc.WriteToken(jsontext.String(entry.Key)); err != nil {
			return err
		}

		if marshaler, ok := any(entry.Value).(OrderedMarshaler); ok {
			if isNilOrderedMarshaler(marshaler) {
				if err := enc.WriteToken(jsontext.Null); err != nil {
					return err
				}
				continue
			}
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
	if o == nil {
		return errNilObject
	}
	dec := jsontext.NewDecoder(bytes.NewReader(data), jsontext.AllowDuplicateNames(true))
	entries, err := decodeEntries[V](dec)
	if err != nil {
		return err
	}

	tok, err := dec.ReadToken()
	switch {
	case errors.Is(err, io.EOF):
		o.replaceEntries(entries)
		return nil
	case err != nil:
		return fmt.Errorf("invalid trailing content: %w", errors.Join(ErrTrailingToken, err))
	default:
		return fmt.Errorf("unexpected trailing token %v: %w", tok.Kind(), ErrTrailingToken)
	}
}

func isNilOrderedMarshaler(marshaler OrderedMarshaler) bool {
	v := reflect.ValueOf(marshaler)
	return v.Kind() == reflect.Pointer && v.IsNil()
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

func findEntryIndex[V any](entries []Entry[V], key string) int {
	for i := range entries {
		if entries[i].Key == key {
			return i
		}
	}
	return -1
}

func decodeEntries[V any](dec *jsontext.Decoder) ([]Entry[V], error) {
	if dec == nil {
		return nil, ErrNilJSONDecoder
	}

	tok, err := dec.ReadToken()
	if err != nil {
		return nil, err
	}
	if tok.Kind() != jsontext.KindBeginObject {
		return nil, fmt.Errorf("expected object start '{', got %v: %w", tok.Kind(), ErrExpectedObjectStart)
	}

	entries := make([]Entry[V], 0)
	for dec.PeekKind() != jsontext.KindEndObject {
		key, err := readObjectKey(dec)
		if err != nil {
			return nil, err
		}

		if findEntryIndex(entries, key) >= 0 {
			return nil, fmt.Errorf("duplicate key %q: %w", key, ErrDuplicateKey)
		}

		var value V
		if err := json.UnmarshalDecode(dec, &value); err != nil {
			return nil, err
		}

		entries = append(entries, Entry[V]{Key: key, Value: value})
	}

	if _, err = dec.ReadToken(); err != nil {
		return nil, err
	}
	return entries, nil
}

func (o *Object[V]) replaceEntries(entries []Entry[V]) {
	clear(o.entries)
	o.entries = entries
}

// UnmarshalJSONFrom decodes a JSON object from dec into o.
// UnmarshalJSONFrom replaces any existing entries in o only after a successful decode.
func (o *Object[V]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if o == nil {
		return errNilObject
	}
	entries, err := decodeEntries[V](dec)
	if err != nil {
		return err
	}
	o.replaceEntries(entries)
	return nil
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
