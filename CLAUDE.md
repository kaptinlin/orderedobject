# CLAUDE.md

## Project Overview

`orderedobject` is a small Go library for JSON objects whose member order matters. It preserves top-level insertion order across mutation, iteration, and JSON encoding without becoming a general-purpose JSON DOM.

For usage and runnable examples, see [README.md](README.md). `AGENTS.md` is a symlink to this file.

## Commands

```bash
task test   # Run go test -race -count=1 ./...
task lint   # Check module tidiness, formatting, and golangci-lint
task fuzz   # Run both fuzz targets for a bounded 10 seconds each
task verify # Run the complete check-only gate
task fmt    # Run Go formatting
task vet    # Run go vet ./...
task clean  # Remove build artifacts and Go caches
```

## Architecture

- `object.go` — ordered object implementation and JSON encode/decode paths
- `object_test.go` — unit, round-trip, and benchmark coverage
- `object_fuzz_test.go` — operation-sequence model fuzzing
- `json_fuzz_test.go` — JSON transaction and round-trip fuzzing
- `example_test.go` — executable canonical usage examples
- `SPECS/` — canonical design rules
- `.references/` — reference implementations used as evidence, not conformance targets
- `Taskfile.yml` — developer commands
- `.github/workflows/ci.yml` — CI gates for test, lint, security, and fuzz checks

`MarshalJSON` frames members directly in insertion order and delegates keys and
values to `encoding/json`. `UnmarshalJSON` uses decoder tokens only to recover
top-level source order before committing a complete replacement.

## Agent Operating Rules

- Read the relevant `SPECS/` before changing behavior, tests, or docs.
- Keep changes surgical; do not refactor unrelated code while solving a narrow task.
- Prefer the simplest implementation that preserves the ordered-object contract.
- Verify behavior through public APIs and user-visible JSON output.
- Fail loudly with errors; do not hide structural JSON problems.
- Do not create policy-only gate scripts that merely restate docs or specs.
- Do not add spec mirror tests when stronger behavior tests already cover the invariant.
- Keep README usage-oriented, CLAUDE.md agent-oriented, and SPECS normative.

## Agent Workflow

### Design Phase — Read SPECS First

Before changing behavior, tests, or docs, read the relevant `SPECS/` files first. `SPECS/` is the canonical source for ordering rules, JSON semantics, and documentation boundaries.

**Workflow**:

1. Identify the relevant specs from the index below.
2. Read those specs completely.
3. Design the change within the slice-backed model and direct ordered JSON byte boundary.
4. If a requested change conflicts with a spec, ask the user before proceeding.

### Implementation Phase — Read References as Evidence

Before implementation, inspect at least two relevant projects under
`.references/`. Borrow only proof shapes or lessons that match current local
pressure; do not copy their API, storage model, dependencies, or compatibility
constraints without local evidence.

## SPECS Index

| File | Purpose |
| --- | --- |
| `SPECS/00-overview.md` | Scope, non-goals, public boundaries, and repository boundaries |
| `SPECS/10-domain-specs.md` | Public types, constructors, mutation invariants, map conversion rules, and JSON semantics |
| `SPECS/40-architecture-specs.md` | Storage model and JSON pipeline |
| `SPECS/50-coding-standards.md` | Testing, linting, and documentation rules |

## References Index

| Path | Relevant evidence |
| --- | --- |
| `.references/orderedmap/` | Explicit ordered traversal and copy behavior |
| `.references/go-ordered-map/` | Ordered JSON byte framing, examples, and bounded fuzzing |
| `.references/go-json-experiment-json/` | JSON dispatch and stream encoder semantics used for comparison only |

## Design Philosophy

- **KISS**: Keep one core collection type and one extension point. Avoid shadow indexes, caches, and helper layers that hide ordered behavior.
- **YAGNI**: This package is not a general-purpose JSON DOM. Do not add abstractions for hypothetical large-object workloads.
- **Precision over cleverness**: Order-preserving operations should be direct; order-dropping operations must be explicit at the call site.
- **APIs as language**: The common path should read like shaping JSON: `New().Set(...).Set(...)`.
- **Never:** accidental complexity, feature gravity, compatibility shims, duplicate public paths, or configurability theater.

## API Design Principles

- **Progressive Disclosure**: `New`, `NewCap`, `FromEntries`, `FromSortedMap`, and `FromJSON` cover construction without duplicate paths.
- **Explicit boundaries**: Sorted map import invents deterministic order; map export explicitly drops ordering semantics.
- **Single path**: Do not add aliases or wrappers that duplicate an existing public operation.

## Coding Rules

### Must Follow

- Use the Go version declared in `go.mod`; use modern standard-library helpers when they simplify code.
- Read the relevant `SPECS/` documents before changing behavior or docs.
- Use `t.Parallel()` when a test case is safe to run concurrently.
- Use focused subtests and `cmp.Diff` for structural comparisons.
- Frame ordered JSON directly, use decoder tokens only for top-level source order, and avoid unordered map round-trips.

### Forbidden

- No `panic` in production code; return errors instead.
- No documentation masquerading as code — do not encode spec prose or unused rules into runtime structures.
- No working around dependency bugs — if a dependency is the source of the problem, create `reports/<dependency-name>.md` instead of reimplementing it.
- No hidden ordering layers — keep ordered behavior slice-backed unless profiling proves a real need.
- No compatibility aliases or duplicate public names for the same operation.
- No policy-only gates or spec mirror tests that restate documentation without proving behavior.

## Testing

- `task test` runs `go test -race -count=1 ./...`.
- `task lint` checks module tidiness, Go formatting, and golangci-lint.
- `task fuzz` runs the operation and JSON fuzz targets for bounded durations.
- `task verify` is the final check-only gate and adds vet and vulnerability scanning.
- Benchmarks use `testing.B.Loop()`.
- Test error contracts with `errors.Is`.

## Dependency Issue Reporting

When you hit a bug or limitation in a dependency:

1. Do not work around it by reimplementing the dependency.
2. Do not skip the dependency and write a replacement inline.
3. Create `reports/<dependency-name>.md`.
4. Include the dependency version, trigger scenario, expected behavior, actual behavior, relevant errors, and any non-code workaround.
5. Continue with tasks that do not depend on the broken functionality.

## Agent Skills

Repo-local skills live under `.agents/skills/`; `.claude/skills` points to the
same directory. Use the smallest relevant skill.

| Skill | Use when |
| --- | --- |
| `tdd-implementing` | Implementing behavior through RED/GREEN/refactor |
| `library-test-covering` | Expanding behavior, error, fuzz, or regression coverage |
| `agent-md-writing` | Updating `CLAUDE.md` / `AGENTS.md` |
| `readme-writing` | Updating user-facing usage documentation |
| `spec-writing` | Updating durable contracts under `SPECS/` |
| `code-review` | Reviewing the complete implementation diff before commit |
| `committing` | Staging and creating the verified commit |
| `releasing` | Selecting, tagging, publishing, and verifying a release |
