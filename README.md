# orderedobject

A generic ordered JSON object for Go that preserves insertion order during mutation, iteration, and JSON encoding

## Features

- **Ordered keys**: Preserve top-level insertion order when encoding JSON.
- **Generic values**: Store `any` or concrete types with the same API.
- **Stable updates**: Updating an existing key keeps its original position.
- **Streaming hooks**: Use `MarshalJSONTo` and `UnmarshalJSONFrom` for token-based JSON work.
- **Map bridges**: Convert to and from `map[string]V` when order does not matter.
- **Runnable examples**: Explore focused programs under [`examples/`](examples/README.md).

## Installation

Requires the Go version declared in `go.mod`.

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
	person := orderedobject.NewObject[any]().
		Set("name", "Alice").
		Set("age", 30).
		Set("city", "New York")

	data, err := person.ToJSON()
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
| `NewObject[V](capacity...)` | Create a new ordered object with optional capacity. |
| `FromMap[V](m)` | Build an object from a map; resulting order follows Go's map iteration order. |
| `FromJSON[V](data)` | Decode a JSON object while preserving member order. |
| `Set`, `Get`, `Has`, `Delete` | Mutate and query keys. |
| `Keys`, `Values`, `Entries`, `ForEach` | Iterate in insertion order. |
| `Clone` | Copy the entries slice without deep-copying stored values. |
| `ToMap` | Convert to a plain map and drop ordering semantics. |
| `ToJSON`, `MarshalJSON`, `MarshalJSONTo` | Encode ordered content to JSON. |
| `UnmarshalJSON`, `UnmarshalJSONFrom` | Decode ordered content from JSON. |

## Examples

### Preserve order while updating values

```go
obj := orderedobject.NewObject[int]().
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

users := orderedobject.NewObject[User]().
	Set("alice", User{Name: "Alice", Age: 30}).
	Set("bob", User{Name: "Bob", Age: 25})

user, _ := users.Get("alice")
fmt.Println(user.Name)
// Alice
```

### Run the bundled examples

```bash
go run ./examples/basic
go run ./examples/json
go run ./examples/nested
```

See [`examples/README.md`](examples/README.md) for the full example list.

## Development

```bash
task test
task lint
```

For development guidelines, see [CLAUDE.md](CLAUDE.md).

## Contributing

Run `task fmt`, `task test`, and `task lint` before sending changes.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
