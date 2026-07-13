# orderedobject

A generic ordered JSON object for Go that preserves top-level insertion order during mutation, iteration, and JSON encoding

## Features

- **Ordered keys**: Preserve top-level insertion order when encoding JSON.
- **Generic values**: Store `any` or concrete types with the same API.
- **Stable updates**: Update existing keys without moving their positions.
- **Standard JSON**: Use the standard `encoding/json` package for nested values and integration.
- **Lazy iteration**: Traverse in insertion order and stop without collecting a snapshot.
- **Explicit map bridges**: Choose sorted map import or unordered map conversion at the call site.

## Installation

Requires Go 1.26.5+.

```bash
go get github.com/kaptinlin/orderedobject
```

## Quick Start

Create an ordered object with `Set`, then encode it through the standard
`encoding/json` package.

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kaptinlin/orderedobject"
)

func main() {
	person := orderedobject.New[any]().
		Set("name", "Alice").
		Set("age", 30).
		Set("city", "New York")

	data, err := json.Marshal(person)
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

## Common Usage

### Preserve order while updating values

Use `Set` to replace a value without moving its key.

```go
obj := orderedobject.New[int]().
	Set("first", 1).
	Set("second", 2)

obj.Set("first", 99)

fmt.Println(obj.Keys())
// [first second]
```

### Iterate with early stop

Use `All` when you want ordered traversal without allocating a snapshot.

```go
for key, value := range obj.All() {
	fmt.Println(key, value)
	if key == "first" {
		break
	}
}
```

Use `Entries` instead when the loop needs a stable snapshot while the object is
being mutated.

### Decode ordered JSON

Use `FromJSON` to preserve top-level member order from JSON input.

```go
obj, err := orderedobject.FromJSON[int]([]byte(`{"second":2,"first":1}`))
if err != nil {
	log.Fatal(err)
}

fmt.Println(obj.Keys())
// [second first]
```

### Use concrete value types

Choose a concrete type when every member has the same JSON shape.

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

Use a sorted import when the source is an unordered Go map.

```go
settings := map[string]int{"z": 26, "a": 1}

sorted := orderedobject.FromSortedMap(settings)
fmt.Println(sorted.Keys())
// [a z]

plain := sorted.ToUnorderedMap()
fmt.Println(plain["a"])
// 1
```

## API Overview

| API | Purpose |
| --- | --- |
| `New[V]()` | Create an empty ordered object. |
| `NewCap[V](n)` | Create an empty ordered object with a capacity hint. |
| `FromEntries[V](entries)` | Build from ordered entries and reject duplicate keys. |
| `FromSortedMap[V](m)` | Build from a map in lexical key order. |
| `FromJSON[V](data)` | Decode a JSON object while preserving member order. |
| `Set`, `Get`, `Has`, `Delete` | Mutate and query keys. |
| `All` | Iterate lazily in insertion order with early-stop support. |
| `Keys`, `Values`, `Entries` | Collect insertion-ordered snapshots. |
| `Clone` | Copy the entries slice without deep-copying stored values. |
| `ToUnorderedMap` | Convert to a plain map and drop ordering semantics. |
| `MarshalJSON`, `UnmarshalJSON` | Integrate ordered objects with `encoding/json`. |

## Development

```bash
task test
task lint
task fuzz
task verify
```

For development guidelines, see [AGENTS.md](AGENTS.md).

## Contributing

Run `task fmt` and `task verify` before sending changes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
