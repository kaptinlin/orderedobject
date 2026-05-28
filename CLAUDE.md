# CLAUDE.md

## Project Overview

`orderedobject` is a small Go library for JSON objects whose member order matters. It preserves top-level insertion order across mutation, iteration, and JSON encoding without becoming a general-purpose JSON DOM.

For usage and runnable examples, see [README.md](README.md).

## Commands

```bash
task test  # Run go test -race ./...
task lint  # Run golangci-lint and the go mod tidy check
task fmt   # Run go fmt ./...
task vet   # Run go vet ./...
task clean # Remove build artifacts and Go caches
```

## Architecture

- `object.go` — ordered object implementation and JSON encode/decode paths
- `object_test.go` — unit, round-trip, and benchmark coverage
- `examples/` — runnable usage examples
- `SPECS/` — canonical design rules
- `Taskfile.yml` — developer commands
- `.github/workflows/ci.yml` — CI gates for test, lint, and security checks

Runtime JSON behavior depends on `github.com/go-json-experiment/json` and `jsontext`.

## Agent Workflow

### Design Phase — Read SPECS First

Before changing behavior, tests, or docs, read the relevant `SPECS/` files first. `SPECS/` is the canonical source for ordering rules, JSON semantics, and documentation boundaries.

**Workflow**:

1. Identify the relevant specs from the index below.
2. Read those specs completely.
3. Keep `README.md` usage-oriented and `CLAUDE.md` agent-oriented.
4. If a requested change conflicts with a spec, ask the user before proceeding.

## SPECS Index

| File | Purpose |
| --- | --- |
| `SPECS/00-overview.md` | Scope, non-goals, and repository boundaries |
| `SPECS/10-domain-specs.md` | Public types, mutation invariants, map conversion rules, and JSON semantics |
| `SPECS/40-architecture-specs.md` | Storage model and JSON pipeline |
| `SPECS/50-coding-standards.md` | Testing, linting, and documentation rules |

## Design Philosophy

- **KISS**: Keep one core collection type and one extension point. Avoid shadow indexes, caches, and helper layers that hide ordered behavior.
- **YAGNI**: This package is not a general-purpose JSON DOM. Do not add extra abstractions for hypothetical large-object workloads.
- **Precision over cleverness**: Make order-preserving and order-dropping operations explicit. `Set` keeps position, while `ToMap` and `FromMap` do not promise stable order.
- **APIs as language**: The common path should read like shaping a JSON object: `NewObject().Set(...).Set(...)`.
- **Never:** accidental complexity, feature gravity, abstraction theater, configurability cope.

## API Design Principles

- **Progressive Disclosure**: `NewObject`, `FromMap`, and `FromJSON` cover the common entry points, while `MarshalJSONTo` and `UnmarshalJSONFrom` remain available for streaming control.

## Coding Rules

### Must Follow

- Use the Go version declared in `go.mod`; use modern standard-library helpers when they simplify code.
- Read the relevant `SPECS/` documents before changing behavior or docs.
- Use `t.Parallel()` when a test case is safe to run concurrently.
- Use focused subtests and `cmp.Diff` for structural comparisons.
- Keep documentation role-separated: `README.md` for usage, `CLAUDE.md` for agent workflow, `SPECS/` for normative rules.

### Forbidden

- No `panic` in production code; return errors instead.
- No documentation masquerading as code — do not encode spec prose or unused rules into runtime structures.
- No working around dependency bugs — if a dependency is the source of the problem, create `reports/<dependency-name>.md` instead of reimplementing it.
- No hidden ordering layers — keep ordered behavior slice-backed and avoid unordered map round-trips in ordered JSON paths.

## Testing

- `task test` is the main gate and runs `go test -race ./...`.
- `task lint` runs golangci-lint plus a tidy check.
- Benchmarks use `testing.B.Loop()`.

## Dependency Issue Reporting

When you hit a bug or limitation in a dependency:

1. Do not work around it by reimplementing the dependency.
2. Do not skip the dependency and write a replacement inline.
3. Create `reports/<dependency-name>.md`.
4. Include the dependency version, trigger scenario, expected behavior, actual behavior, relevant errors, and any non-code workaround.
5. Continue with tasks that do not depend on the broken functionality.

## Agent Skills

No repo-local skills are currently defined for this library.
