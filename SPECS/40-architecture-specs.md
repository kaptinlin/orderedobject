# Architecture Specs

## Overview

`orderedobject` stays a single-package library with one core implementation file for ordered-object behavior. The architecture favors explicit slice operations and direct ordered JSON byte framing over caches, layers, or hidden state.

## Package Topology

- Core ordered-object behavior lives in `object.go`.
- Executable package examples live in `example_test.go` and add no runtime dependency.
- Runtime JSON behavior depends only on the standard `encoding/json` package.

> **Why:** The package solves one narrow problem. Keeping the implementation local and visible makes the behavior easy to audit and reduces maintenance overhead.
>
> **Rejected:** Multi-package decomposition, generated wrappers, compatibility shim packages, and indirection whose main benefit would be architectural symmetry instead of clarity.

## Storage Model

- `Object[V]` stores data as `[]Entry[V]`.
- Key lookup uses a linear scan over that slice.
- Mutating operations update the slice directly.
- Ordered imports copy into a fresh entries slice.
- Decode paths build local entries first and replace the receiver only after the full replacement content is accepted.
- Successful replacement clears old entries before assigning the new sequence so replaced values can be released for garbage collection.

## JSON Pipeline

- `MarshalJSON` writes object delimiters in `entries` order and delegates every key and value to `encoding/json`.
- Nested ordered objects use the standard `json.Marshaler` contract; plain maps inherit deterministic standard-library key ordering.
- `UnmarshalJSON` validates one complete JSON value, then uses `json.Decoder` tokens to read top-level members in source order.
- Decode builds local entries and commits only after the complete object succeeds.
- The package must not round-trip through an intermediate `map[string]V` during ordered JSON encoding or decoding.

> **Why:** Direct byte framing preserves encode order without an intermediate map, while decoder tokens recover top-level source order without building a general JSON DOM.
>
> **Rejected:** Map-based round-trips, parallel index structures, and architecture that optimizes hypothetical throughput at the cost of a harder-to-reason-about implementation.

## Forbidden

- Do not add a shadow `map[string]int` or similar index alongside the entries slice without measured need.
- Do not route ordered JSON operations through unordered maps.
- Do not split the ordered-object core across multiple packages unless a new domain boundary appears.
- Do not add reader, writer, or append JSON APIs without a concrete caller that needs that exact boundary.

## Acceptance Criteria

- The canonical storage model remains slice-backed.
- Ordered JSON paths remain direct and avoid unordered map round-trips.
- Architecture changes keep ordered behavior auditable from one primary implementation area.
