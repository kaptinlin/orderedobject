# orderedobject

A generic ordered JSON object for Go that preserves top-level insertion order during mutation, iteration, and JSON encoding

## Features

- **Ordered keys**: Preserve top-level insertion order when encoding JSON.
- **Generic values**: Store `any` or concrete types with the same API.
- **Stable updates**: Update existing keys without moving their positions.
- **Explicit map bridges**: Choose sorted map import or unordered map conversion at the call site.
- **Streaming hooks**: Use `MarshalJSONTo` and `UnmarshalJSONFrom` for token-based JSON work.
- **Runnable examples**: Explore focused programs under [`examples/`](examples/README.md).

## Installation

Requires Go 1.26+.

```bash
go get github.com/kaptinlin/orderedobject
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/kaptinlin/orderedobject"
)

func main() {
	person := orderedobject.New[any]().
		Set("name", "Alice").
		Set("age", 30).
		Set("city", "New York")

	data, err := person.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}
```

Output:

```text
{"name":"Alice","age":30,"city":"New York"}
```

## API Overview

| API | Purpose |
| --- | --- |
| `New[V]()` | Create an empty ordered object. |
| `NewCap[V](n)` | Create an empty ordered object with a capacity hint. |
| `FromEntries[V](entries)` | Build from ordered entries and reject duplicate keys. |
| `FromSortedMap[V](m)` | Build from a map in lexical key order. |
| `FromUnorderedMap[V](m)` | Build from a map using Go map iteration order. |
| `FromJSON[V](data)` | Decode a JSON object while preserving member order. |
| `Set`, `Get`, `Has`, `Delete` | Mutate and query keys. |
| `Keys`, `Values`, `Entries`, `ForEach` | Iterate in insertion order. |
| `Clone` | Copy the entries slice without deep-copying stored values. |
| `ToUnorderedMap` | Convert to a plain map and drop ordering semantics. |
| `MarshalJSON`, `MarshalJSONTo` | Encode ordered content to JSON. |
| `UnmarshalJSON`, `UnmarshalJSONFrom` | Decode ordered content from JSON. |

## Examples

### Preserve order while updating values

```go
obj := orderedobject.New[int]().
	Set("first", 1).
	Set("second", 2)

obj.Set("first", 99)

fmt.Println(obj.Keys())
// [first second]
```

### Use concrete value types

```go
type User struct {
	Name string
	Age  int
}

users := orderedobject.New[User]().
	Set("alice", User{Name: "Alice", Age: 30}).
	Set("bob", User{Name: "Bob", Age: 25})

user, _ := users.Get("alice")
fmt.Println(user.Name)
// Alice
```

### Import maps explicitly

```go
settings := map[string]int{"z": 26, "a": 1}

sorted := orderedobject.FromSortedMap(settings)
fmt.Println(sorted.Keys())
// [a z]

plain := sorted.ToUnorderedMap()
fmt.Println(plain["a"])
// 1
```

### Run the bundled examples

```bash
go run ./examples/basic
go run ./examples/json
go run ./examples/map
```

See [`examples/README.md`](examples/README.md) for the full example list.

## Development

```bash
task test
task lint
```

For development guidelines, see [AGENTS.md](AGENTS.md).

## Contributing

Run `task fmt`, `task test`, and `task lint` before sending changes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
