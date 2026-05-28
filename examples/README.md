# Examples

Runnable programs that show the main `orderedobject` usage paths.

Run them from the repository root:

```bash
go run ./examples/basic
go run ./examples/json
go run ./examples/map
go run ./examples/type
go run ./examples/nested
go run ./examples/array
go run ./examples/error
```

| Example | Demonstrates |
| --- | --- |
| `basic` | Creating, updating, deleting, iterating, and cloning ordered objects. |
| `json` | Encoding and decoding JSON while preserving top-level member order. |
| `map` | Bridging to and from `map[string]V` when order is not required. |
| `type` | Using concrete generic value types. |
| `nested` | Nesting ordered objects for structured JSON output. |
| `array` | Encoding slices of ordered objects. |
| `error` | Handling missing keys, type assertions, and JSON decode errors. |
