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
- Compatible with go-json-experiment/json
- Rich set of constructors and utility methods

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

- `Entry[V any]`: Represents a key-value pair
- `Object[V any]`: An ordered collection of key-value pairs

### Functions

- `NewObject[V any](capacity ...int) *Object[V]`: Creates a new ordered object
- `FromMap[V any](m map[string]V) *Object[V]`: Creates an ordered object from a map
- `FromJSON[V any](data []byte) (*Object[V], error)`: Creates an ordered object from JSON

### Methods

- `Set(key string, value V) *Object[V]`: Sets a key-value pair
- `Get(key string) (V, bool)`: Gets a value by key
- `Has(key string) bool`: Checks if a key exists
- `Delete(key string) *Object[V]`: Removes a key-value pair
- `Length() int`: Returns the number of key-value pairs
- `ForEach(fn func(key string, value V))`: Iterates through key-value pairs
- `Clone() *Object[V]`: Creates a deep copy of the object
- `Entries() []Entry[V]`: Returns all key-value pairs
- `ToMap() map[string]V`: Converts to a standard Go map
- `ToJSON() ([]byte, error)`: Converts to JSON
- `MarshalJSON() ([]byte, error)`: Implements json.Marshaler

## FAQ

### Q: Why choose go-json-experiment/json over the standard library?
A: go-json-experiment/json provides better performance and more features while maintaining compatibility with the standard library.

### Q: Does it support custom JSON tags?
A: Yes, it supports standard struct tags like `json:"field_name"`.

### Q: How does it handle key order?
A: The library preserves the insertion order of keys when marshalling to JSON, unlike standard Go maps.

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