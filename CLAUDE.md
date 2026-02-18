# CLAUDE.md

## Project Overview

`orderedobject` is a generic ordered JSON object for Go that preserves insertion order of keys. It uses [go-json-experiment/json](https://github.com/go-json-experiment/json) (experimental encoding/json v2) for token-level streaming marshal/unmarshal.

## Build Commands

```bash
task test       # Run all tests
task lint       # Run golangci-lint and go mod tidy check
make fmt        # Format code
make vet        # Run go vet
task verify     # Run all: deps, fmt, vet, lint, test
task clean      # Remove build artifacts and caches
```

## Architecture

Single-file library (`object.go`) with one core type:

- `Object[V any]` — ordered key-value store backed by `[]Entry[V]`
- Key lookup is linear scan (`findKeyIndex`); suitable for small-to-medium objects
- JSON marshalling uses `jsontext.Encoder` streaming API for zero-intermediate-allocation output
- Nested order preservation via `OrderedMarshaler` interface (checked with type assertion)
- Map values use `json.Deterministic(true)` for sorted key output

## Key Design Decisions

- **No map index**: Linear scan over entries slice, not a parallel `map[string]int`. Keeps the implementation simple and correct for typical JSON object sizes.
- **Shallow clone**: `Clone()` copies the entries slice but not the values themselves.
- **Duplicate key rejection**: `FromJSON` relies on `go-json-experiment/json` default behavior which rejects duplicate keys.
- **Clear on unmarshal**: `UnmarshalJSONFrom` clears existing entries before decoding, using `clear()` for proper GC of old references.

## Go Version

Go 1.26. Uses `slices.Clone`, `slices.Delete`, `clear()`, `for range N`, `testing.B.Loop()`.

## Dependencies

- `github.com/go-json-experiment/json` — JSON v2 experimental library (encoder/decoder/marshal)

## Testing

- All tests use `t.Parallel()`
- Standard library assertions only (no testify)
- Benchmarks use `b.Loop()` (Go 1.24+)
- Run with race detector: `go test -race ./...`

## Lint

golangci-lint version managed by `.golangci.version` file (currently v2.9.0 installed). Config in `.golangci.yml` if present, otherwise defaults.
