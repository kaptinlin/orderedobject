# Coding Standards

## Overview

This repository favors small, explicit Go code and lightweight tooling. Standards here cover how to validate behavior, how to keep docs authoritative, and which supporting tools are acceptable.

## Language and Dependency Rules

- The package targets the Go version declared in `go.mod`.
- Prefer standard-library building blocks unless a dependency directly supports the package's ordered-JSON contract.
- `github.com/go-json-experiment/json` is the runtime JSON dependency.
- `github.com/google/go-cmp/cmp` is acceptable in tests for structural diffs.

> **Why:** The package is small enough that each dependency should justify itself clearly. `go-json-experiment/json` provides the token-level APIs the package is built around, and `go-cmp` keeps test failures readable without introducing an assertion framework.
>
> **Rejected:** Heavy test DSLs, convenience dependencies that do not strengthen the ordered-object contract, and duplicated design guidance outside `SPECS/`.

## Testing Rules

- Tests should call `t.Parallel()` when the case is safe to run concurrently.
- Use focused named subtests when they clarify observable behavior.
- Prefer direct standard-library checks plus `cmp.Diff` for comparisons.
- Do not add assertion frameworks.
- When benchmarks are added, use `testing.B.Loop()`.
- The package gate is `task test`, which runs `go test -race ./...`.

## Lint and Documentation Rules

- `task lint` must pass before shipping changes.
- `task lint` includes golangci-lint and a `go mod tidy` diff check.
- `SPECS/**` is canonical design documentation and must stay markdownlint-clean.
- `README.md` is usage-oriented, and `CLAUDE.md` is agent-oriented.

## Forbidden

- Do not add testify-style assertion layers or similar test frameworks.
- Do not exclude `SPECS/**` from markdown linting.
- Do not place normative design rules only in code comments, `README.md`, or `CLAUDE.md`.

## Acceptance Criteria

- `task lint` and `task test` pass.
- New tests follow the concurrency and comparison rules above.
- Normative design guidance remains centralized in `SPECS/`.
