# CLAUDE.md

## Project Overview

`orderedobject` is a generic ordered JSON object for Go that preserves insertion order of keys.

## Commands

```bash
task test         # Run all tests with race detection
task lint         # Run golangci-lint and the tidy check
task fmt          # Format Go code
task vet          # Run go vet
task verify       # Run deps, fmt, vet, lint, test, and vuln
task markdownlint # Lint Markdown files
task clean        # Remove build artifacts and caches
```

## SPECS Index

- `SPECS/00-overview.md` — package scope and non-goals
- `SPECS/10-domain-specs.md` — core types, invariants, and JSON rules
- `SPECS/40-architecture-specs.md` — storage model and JSON pipeline
- `SPECS/50-coding-standards.md` — testing, lint, and documentation rules

## Working Rules

- Treat `SPECS/` as the canonical home for normative design rules.
- Keep `README.md` focused on usage and examples.
- Run `task lint` and `task test` before finishing a change.
