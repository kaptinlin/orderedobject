# Architecture Specs

## Overview

`orderedobject` stays a single-package library with one core implementation file for ordered-object behavior. The architecture favors explicit slice operations and token-stream JSON encoding over caches, layers, or hidden state.

## Package Topology

- Core ordered-object behavior lives in `object.go`.
- Example programs may depend on the package, but the package must not depend on examples or tooling code.
- Runtime JSON behavior depends on `github.com/go-json-experiment/json` and `jsontext`.

> **Why:** The package solves one narrow problem. Keeping the implementation local and visible makes the behavior easy to audit and reduces maintenance overhead.
>
> **Rejected:** Multi-package decomposition, generated wrappers, and indirection whose main benefit would be architectural symmetry instead of clarity.

## Storage Model

- `Object[V]` stores data as `[]Entry[V]`.
- Key lookup uses a linear scan over that slice.
- Mutating operations update the slice directly.
- Decode paths clear old entries before re-slicing so replaced values can be released for garbage collection.

## JSON Pipeline

- Marshal paths write JSON through `jsontext.Encoder`.
- The encoder writes object delimiters and string keys as tokens, then delegates value encoding.
- Nested `OrderedMarshaler` values bypass the generic value encoder so they can preserve their own order.
- Non-ordered map values are encoded with `json.Deterministic(true)` to stabilize nested map output.
- Unmarshal paths read from `jsontext.Decoder` and decode each member value directly into `V`.
- The package must not round-trip through an intermediate `map[string]V` during ordered JSON encoding or decoding.

> **Why:** Token-stream encoding preserves order directly and avoids rebuilding objects through unordered intermediates.
>
> **Rejected:** Map-based round-trips, parallel index structures, and architecture that optimizes hypothetical throughput at the cost of a harder-to-reason-about implementation.

## Forbidden

- Do not add a shadow `map[string]int` or similar index alongside the entries slice without measured need.
- Do not route ordered JSON operations through unordered maps.
- Do not split the ordered-object core across multiple packages unless a new domain boundary appears.

## Acceptance Criteria

- The canonical storage model remains slice-backed.
- Ordered JSON paths remain token-stream based.
- Architecture changes keep ordered behavior auditable from one primary implementation area.

**Origin:** Migrated from `CLAUDE.md` on 2026-04-21.
