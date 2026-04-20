# Overview Specs

## Overview

`orderedobject` defines a small Go abstraction for JSON objects whose member order matters. Its scope is preserving top-level object order across mutation, iteration, and JSON encoding without turning the package into a general-purpose JSON DOM.

## Scope

- The package exposes one ordered container, `Object[V any]`, plus `Entry[V]` and `OrderedMarshaler`.
- The package is optimized for small-to-medium JSON objects where explicit ordering matters more than constant-time lookup.
- Examples and README content are explanatory only; normative design rules live in `SPECS/`.

> **Why:** Most callers need a small ordered-object type for payload shaping, config-like documents, and stable JSON output. Keeping the scope narrow preserves readability and avoids inventing a second collection framework.
>
> **Rejected:** A full JSON AST, package layering for hypothetical future growth, and performance features that only pay off for very large objects.

## Non-Goals

- Provide map-like constant-time lookups through a shadow index.
- Preserve order when converting through `map[string]V`.
- Model every JSON value kind as its own package-level abstraction.

## Forbidden

- Do not add extra indexing layers or configurability to optimize hypothetical large-object workloads.
  Use the existing slice-backed model until profiling shows a real need.
- Do not treat `FromMap` or `ToMap` as order-preserving APIs.
  Use `Object` methods or JSON encoding when order matters.
- Do not duplicate normative design rules in `README.md` or `CLAUDE.md`.
  Keep `SPECS/` as the single source of truth.

## Acceptance Criteria

- The package remains centered on ordered JSON object semantics.
- Non-goals stay explicit when new features are proposed.
- Normative design rules are maintained in `SPECS/`, not scattered across the repo.
