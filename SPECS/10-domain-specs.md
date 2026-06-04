# Domain Specs

## Overview

The ordered-object domain is defined by stable insertion order, explicit mutation rules, and object-focused JSON semantics. Callers should be able to predict how keys move, how data is copied, and which operations preserve, invent, or drop ordering.

## Core Types

- `Object[V any]` stores entries in insertion order.
- `Entry[V]` represents one key-value pair.
- `OrderedMarshaler` allows a value to control its own ordered JSON encoding when nested inside an `Object`.

> **Why:** The package needs only one collection type and one extension point. That keeps the public model small while still allowing nested ordered encoders.
>
> **Rejected:** Separate mutable and immutable object types, hidden iterator objects, multiple nested-order extension mechanisms, and compatibility aliases for duplicate API paths.

## Constructors and Ordered Imports

- `New` creates an empty object.
- `NewCap` creates an empty object with an explicit capacity hint.
- Non-positive `NewCap` capacity creates an empty object without failing.
- `FromEntries` copies entries in caller-provided order.
- `FromEntries` rejects duplicate keys with an error that wraps `ErrDuplicateKey`.
- `FromSortedMap` copies all map entries in lexical key order.
- `FromUnorderedMap` copies all map entries using Go map iteration order and does not promise stable order.
- `FromJSON` decodes a JSON object while preserving member order from the input.

## Ordering and Mutation Invariants

- `Set` appends a new key at the end.
- `Set` updates an existing key in place without changing its position.
- `Delete` removes the key and closes the gap in the ordered sequence.
- `Keys`, `Values`, `Entries`, and `ForEach` observe the current insertion order.
- `Keys`, `Values`, and `Entries` return new slices; mutating those returned slices must not mutate the object.
- `Clone` is shallow: it copies the entries slice, not the values stored inside it.
- Empty-string keys are valid.
- The zero value of `Object[V]` behaves as an empty object.
- Nil read-only receivers return natural empty results where possible.
- Nil `*Object[V]` marshals as JSON `null`.

## Map Conversion Rules

- `FromSortedMap` creates deterministic lexical key order.
- `FromUnorderedMap` inherits Go's map iteration order and is therefore non-deterministic.
- `ToUnorderedMap` returns a new `map[string]V` and drops ordering semantics.

## JSON Rules

- `FromJSON` and `UnmarshalJSONFrom` accept only JSON objects.
- Non-object input must fail with an error that wraps `ErrExpectedObjectStart`.
- Non-string member names must fail with an error that wraps `ErrExpectedStringKey`.
- Duplicate member names must fail with an error that wraps `ErrDuplicateKey`.
- Trailing tokens after `UnmarshalJSON` must fail with an error that wraps `ErrTrailingToken`.
- Nil JSON encoders and decoders must fail with errors that wrap `ErrNilJSONEncoder` and `ErrNilJSONDecoder`.
- JSON decode replacement is transactional: failed decode leaves existing entries unchanged.
- Successful JSON decode replaces the receiver.
- `MarshalJSON` returns compact JSON bytes without a trailing stream newline.
- `MarshalJSONTo` emits members in insertion order and follows the provided encoder's output behavior.
- Plain map values are encoded with deterministic key ordering.
- Values implementing `OrderedMarshaler` are responsible for their own nested ordered encoding.

> **Why:** Ordered behavior must stay explicit at the API boundary. Callers should know exactly which operations preserve order, which ones discard it, and which errors represent structural JSON violations.
>
> **Rejected:** Silent best-effort coercion of non-object JSON, preserving map iteration order as if it were stable, mutating a valid receiver on failed replacement, and deep-copy semantics hidden behind `Clone`.

## Forbidden

- Do not change key order as a side effect of updating an existing key.
  Keep replacement semantics position-stable.
- Do not alias the internal entries slice through exported APIs.
  Return copies for snapshot-style accessors.
- Do not silently accept non-object JSON, duplicate keys, or trailing input.
  Return structural errors instead.
- Do not add duplicate public names for the same JSON byte operation.
  `MarshalJSON` is the direct byte API.

## Acceptance Criteria

- Public APIs preserve the stated ordering and copy semantics.
- JSON decoding errors continue to distinguish object-start, string-key, duplicate-key, trailing-token, and nil-stream violations.
- Order-dropping operations are visible at the call site.
