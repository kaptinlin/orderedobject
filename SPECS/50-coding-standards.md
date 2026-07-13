# Coding Standards

## Overview

This repository favors small, explicit Go code and lightweight tooling. Standards here cover how to validate behavior, how to keep docs authoritative, and which supporting tools are acceptable.

## Language and Dependency Rules

- The package targets the Go version declared in `go.mod`.
- Prefer standard-library building blocks unless a dependency directly supports the package's ordered-JSON contract.
- `encoding/json` is the runtime JSON dependency.
- `github.com/google/go-cmp/cmp` is acceptable in tests for structural diffs.

> **Why:** The package is small enough that each dependency should justify itself clearly. The standard library covers the ordered byte boundary, and `go-cmp` keeps test failures readable without introducing an assertion framework.
>
> **Rejected:** Heavy test DSLs, convenience dependencies that do not strengthen the ordered-object contract, and duplicated design guidance outside `SPECS/`.

## Testing Rules

- Tests should call `t.Parallel()` when the case is safe to run concurrently.
- Use focused named subtests when they clarify observable behavior.
- Prefer direct standard-library checks plus `cmp.Diff` for comparisons.
- Do not add assertion frameworks.
- Test user-visible behavior and public contracts, not copies of spec prose.
- Do not add spec mirror tests when stronger behavior tests already prove the invariant.
- Model fuzzers must assert their complete final state for every input, including
  inputs that contain no complete operation.
- When benchmarks are added, use `testing.B.Loop()`.
- `task test` runs `go test -race -count=1 ./...`.
- `task fuzz` runs bounded operation-model and JSON-transaction fuzz targets.

## Lint and Documentation Rules

- `task lint` must pass before shipping changes.
- `task lint` includes golangci-lint, a `go mod tidy` diff check, and a check-only
  Go formatting pass.
- `task verify` is the final check-only gate and runs vet, lint, tests,
  vulnerability scanning, and bounded fuzzing without rewriting source files.
- `SPECS/**` is canonical design documentation and must stay markdownlint-clean.
- `README.md` is usage-oriented, and `CLAUDE.md` is agent-oriented.
- Do not create policy-only gate scripts whose only job is to restate `SPECS/`, `README.md`, or `CLAUDE.md`.

## Forbidden

- Do not add testify-style assertion layers or similar test frameworks.
- Do not exclude `SPECS/**` from markdown linting.
- Do not place normative design rules only in code comments, `README.md`, or `CLAUDE.md`.
- Do not enforce documentation policy with custom scripts unless a program consumes the rule as product behavior.

## Acceptance Criteria

- `task verify` passes without rewriting tracked source files.
- New tests follow the concurrency and comparison rules above.
- Normative design guidance remains centralized in `SPECS/`.
