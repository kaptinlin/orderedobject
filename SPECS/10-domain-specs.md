# Domain Specs

## Overview

The ordered-object domain is defined by stable insertion order, explicit mutation rules, and object-focused JSON semantics. Callers should be able to predict how keys move, how data is copied, and which operations preserve ordering.

## Core Types

- `Object[V any]` stores entries in insertion order.
- `Entry[V]` represents one key-value pair.
- `OrderedMarshaler` allows a value to control its own ordered JSON encoding when nested inside an `Object`.

> **Why:** The package needs only one collection type and one extension point. That keeps the public model small while still allowing nested ordered encoders.
>
> **Rejected:** Separate mutable and immutable object types, hidden iterator objects, and multiple nested-order extension mechanisms.

## Ordering and Mutation Invariants

- `Set` appends a new key at the end.
- `Set` updates an existing key in place without changing its position.
- `Delete` removes the key and closes the gap in the ordered sequence.
- `Keys`, `Values`, `Entries`, and `ForEach` observe the current insertion order.
- `Keys`, `Values`, and `Entries` return new slices; mutating those returned slices must not mutate the object.
- `Clone` is shallow: it copies the entries slice, not the values stored inside it.
- Empty-string keys are valid.

## Map Conversion Rules

- `FromMap` copies all entries from the input map.
- `FromMap` inherits Go's map iteration order and is therefore non-deterministic.
- `ToMap` returns a new `map[string]V` and drops ordering semantics.

## JSON Rules

- `FromJSON` and `UnmarshalJSONFrom` accept only JSON objects.
- Non-object input must fail with an error that wraps `ErrExpectedObjectStart`.
- Non-string member names must fail with an error that wraps `ErrExpectedStringKey`.
- `UnmarshalJSONFrom` clears existing entries before decoding replacement content.
- `UnmarshalJSON` rejects trailing tokens after the top-level object.
- Duplicate keys are rejected during decoding by the underlying JSON decoder.
- `MarshalJSONTo` emits members in insertion order.
- Plain map values are encoded with deterministic key ordering.
- Values implementing `OrderedMarshaler` are responsible for their own nested ordered encoding.

> **Why:** Ordered behavior must stay explicit at the API boundary. Callers should know exactly which operations preserve order, which ones discard it, and which errors represent structural JSON violations.
>
> **Rejected:** Silent best-effort coercion of non-object JSON, preserving map iteration order as if it were stable, and deep-copy semantics hidden behind `Clone`.

## Forbidden

- Do not change key order as a side effect of updating an existing key.
  Keep replacement semantics position-stable.
- Do not alias the internal entries slice through exported APIs.
  Return copies for snapshot-style accessors.
- Do not silently accept non-object JSON or trailing input.
  Return structural errors instead.

## Acceptance Criteria

- Public APIs preserve the stated ordering and copy semantics.
- JSON decoding errors continue to distinguish object-start and string-key violations.
- Order-dropping operations are documented as such and remain explicit.
