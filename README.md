# OrderedObject

An ordered JSON object implementation that preserves insertion order, designed to work with [github.com/go-json-experiment/json](https://github.com/go-json-experiment/json) (Go's experimental encoding/json v2).

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage Examples](#usage-examples)
  - [Basic Operations](#basic-operations)
  - [JSON Operations](#json-operations)
  - [Map Operations](#map-operations)
  - [Type Safety](#type-safety)
- [API Reference](#api-reference)
- [FAQ](#faq)
- [Contributing](#contributing)
- [License](#license)

## Features

- Preserves key insertion order when marshalling to JSON
- Generic implementation that can store any type of value
- Simple API with method chaining
- Deterministic map key ordering in JSON output
- Nested ordered object support via `OrderedMarshaler` interface
- Compatible with [go-json-experiment/json](https://github.com/go-json-experiment/json)

## Installation

```bash
go get github.com/kaptinlin/orderedobject
```

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/kaptinlin/orderedobject"
	json "github.com/go-json-experiment/json"
)

func main() {
	// Create and populate an ordered object
	person := orderedobject.NewObject[any]().
		Set("name", "John").
		Set("age", 30).
		Set("city", "New York")
	
	// Convert to JSON (order preserved)
	data, _ := person.ToJSON()
	fmt.Println(string(data))
	// Output: {"name":"John","age":30,"city":"New York"}
}
```

## Usage Examples

### Basic Operations

```go
// Create and manipulate objects
obj := orderedobject.NewObject[any]().
	Set("a", 1).
	Set("b", 2).
	Set("c", 3)

// Get values
if value, found := obj.Get("a"); found {
	fmt.Println("Value:", value)
}

// Check existence
if obj.Has("b") {
	fmt.Println("Key 'b' exists")
}

// Delete keys
obj.Delete("c")

// Iterate through entries
obj.ForEach(func(key string, value any) {
	fmt.Printf("%s: %v\n", key, value)
})

// Get all entries
entries := obj.Entries()
for _, entry := range entries {
	fmt.Printf("%s: %v\n", entry.Key, entry.Value)
}

// Clone object
clone := obj.Clone()
```

### JSON Operations

```go
// Create from JSON
jsonData := []byte(`{"name":"John","age":30}`)
obj, err := orderedobject.FromJSON[any](jsonData)

// Convert to JSON (three ways)
// 1. Using ToJSON (recommended)
data1, _ := obj.ToJSON()

// 2. Using json.Marshal
data2, _ := json.Marshal(obj)

// 3. Via map (order not preserved)
m := obj.ToMap()
data3, _ := json.Marshal(m)

// Nested objects
address := orderedobject.NewObject[any]().
	Set("street", "123 Main St").
	Set("zipcode", "10001")

person := orderedobject.NewObject[any]().
	Set("name", "Alice").
	Set("address", address)

data, _ := person.ToJSON()
// Output: {"name":"Alice","address":{"street":"123 Main St","zipcode":"10001"}}
```

### Map Operations

```go
// Create from map
m := map[string]any{
	"name": "John",
	"age":  30,
}
obj := orderedobject.FromMap(m)

// Convert to map (order not preserved)
mapData := obj.ToMap()
```

### Type Safety

```go
// Using concrete types
type User struct {
	Name string
	Age  int
}

// Create type-safe object
users := orderedobject.NewObject[User]().
	Set("user1", User{Name: "John", Age: 30}).
	Set("user2", User{Name: "Alice", Age: 25})

// Type-safe value retrieval
if user, found := users.Get("user1"); found {
	fmt.Printf("User: %s, Age: %d\n", user.Name, user.Age)
}
```

## API Reference

### Types

- `Entry[V any]` — A key-value pair with `Key string` and `Value V` fields
- `Object[V any]` — An ordered collection of key-value pairs that preserves insertion order
- `OrderedMarshaler` — Interface for types that marshal themselves to JSON while preserving key order

### Sentinel Errors

- `ErrExpectedObjectStart` — Returned when the JSON token is not `{`
- `ErrExpectedStringKey` — Returned when the JSON token is not a string key

### Constructors

- `NewObject[V any](capacity ...int) *Object[V]` — Creates a new ordered object with optional pre-allocated capacity
- `FromMap[V any](m map[string]V) *Object[V]` — Creates an ordered object from a map (iteration order is non-deterministic)
- `FromJSON[V any](data []byte) (*Object[V], error)` — Parses a JSON byte slice, preserving key order

### Methods

- `Set(key string, value V) *Object[V]` — Sets a key-value pair; updates in place if the key exists
- `Get(key string) (V, bool)` — Returns the value for a key and whether the key exists
- `Has(key string) bool` — Reports whether the key exists
- `Delete(key string) *Object[V]` — Removes a key-value pair; no-op if key does not exist
- `Len() int` — Returns the number of key-value pairs
- `Keys() []string` — Returns all keys in insertion order
- `Values() []V` — Returns all values in insertion order
- `Entries() []Entry[V]` — Returns a copy of all key-value pairs in insertion order
- `ForEach(fn func(key string, value V))` — Calls fn for each key-value pair in insertion order
- `Clone() *Object[V]` — Returns a shallow copy of the ordered object
- `ToMap() map[string]V` — Converts to a standard Go map (order not preserved)
- `ToJSON() ([]byte, error)` — Converts to JSON
- `MarshalJSON() ([]byte, error)` — Implements `json.Marshaler`
- `MarshalJSONTo(enc *jsontext.Encoder) error` — Encodes to a JSON encoder (implements `OrderedMarshaler`)
- `UnmarshalJSON(data []byte) error` — Implements `json.Unmarshaler`
- `UnmarshalJSONFrom(dec *jsontext.Decoder) error` — Decodes from a JSON decoder

## FAQ

### Why use go-json-experiment/json instead of encoding/json?

go-json-experiment/json provides token-level streaming APIs (`jsontext.Encoder`/`jsontext.Decoder`) that enable efficient ordered marshalling without intermediate allocations. The standard library does not expose equivalent APIs.

### Does it support custom JSON tags?

Yes. Values are marshalled using `go-json-experiment/json`, which supports standard struct tags like `json:"field_name"` and `json:",omitempty"`.

### How does it handle nested ordered objects?

When a value implements the `OrderedMarshaler` interface, key order is preserved recursively during marshalling. When parsing JSON via `FromJSON[any]`, nested objects are decoded as `map[string]any` (order not preserved at nested levels).

### Are duplicate keys allowed?

No. `FromJSON` rejects JSON input containing duplicate keys and returns an error.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create your feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

Please ensure your code:
- Includes appropriate tests
- Follows Go code conventions
- Updates relevant documentation

## License

MIT 