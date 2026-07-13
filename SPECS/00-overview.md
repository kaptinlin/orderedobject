# Overview Specs

## Overview

`orderedobject` defines a small Go abstraction for JSON objects whose member order matters. Its scope is preserving top-level object order across mutation, iteration, and JSON encoding without turning the package into a general-purpose JSON DOM.

## Scope

- The package exposes one ordered container, `Object[V any]`, plus `Entry[V]`.
- The package is optimized for small-to-medium JSON objects where explicit ordering matters more than constant-time lookup.
- Map boundaries are explicit: sorted imports invent deterministic lexical order, while exports drop ordering semantics.
- Examples and README content are explanatory only; normative design rules live in `SPECS/`.

> **Why:** Most callers need a small ordered-object type for payload shaping, config-like documents, and stable JSON output. Keeping the scope narrow preserves readability and avoids inventing a second collection framework.
>
> **Rejected:** A full JSON AST, package layering for hypothetical future growth, compatibility aliases, and performance features that only pay off for very large objects.

## Non-Goals

- Provide map-like constant-time lookups through a shadow index.
- Preserve order when converting through `map[string]V`.
- Model every JSON value kind as its own package-level abstraction.
- Maintain duplicate public paths for the same operation.

## Forbidden

- Do not add extra indexing layers or configurability to optimize hypothetical large-object workloads.
  Use the existing slice-backed model until profiling shows a real need.
- Do not treat unordered map conversion as order-preserving.
  Use `Object` methods, `FromEntries`, `FromSortedMap`, or JSON encoding when order matters.
- Do not add compatibility aliases for removed public names.
  Keep one clear public path for each operation.
- Do not duplicate normative design rules in `README.md` or `CLAUDE.md`.
  Keep `SPECS/` as the single source of truth.

## Acceptance Criteria

- The package remains centered on ordered JSON object semantics.
- Non-goals stay explicit when new features are proposed.
- Normative design rules are maintained in `SPECS/`, not scattered across the repo.
