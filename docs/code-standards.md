# Code Standards

## Language and Tooling
- Language: Go (`go 1.22` in `go.mod`).
- Build target entrypoint: `cmd/act/main.go`.
- Internal packages grouped by domain under `internal/`.
- CI executes `go test ./...` and `go build ./cmd/act`.

## Architecture Conventions
- Command handlers live in `internal/commands` as `RunX` functions.
- Runtime orchestration logic lives in `internal/kit`.
- Shared filesystem helpers live in `internal/fsutil`.
- App-level command routing is centralized in `internal/app`.
- Global config persistence is isolated to `internal/config`.

## Error Handling
- Return `error` from command handlers, avoid process exits outside `main`.
- Prefer explicit, actionable error messages (`unknown option`, `invalid scope`, etc).
- Wrap errors with context where crossing module boundaries.

## CLI Behavior
- Default scope is `project`; explicit global aliases supported.
- Preserve non-interactive path for automation (`--non-interactive`, `--yes`).
- Keep command usage and help catalog aligned.

## Filesystem Safety
- Copy operations should preserve tree structure and create directories as needed.
- Skip overwriting by default; require `--force` or overwrite flag path.
- Validate source existence before migration/install actions.

## Testing Standards
- Unit tests for parser logic, command arg handling, and migration behavior.
- Keep tests in same package with `_test.go` suffix.
- Test both default and edge inputs (invalid flags/scopes, missing args).

## CI and Release Standards
- CI matrix: Ubuntu, Windows, macOS.
- Release artifacts: per-OS/per-arch binaries uploaded on `v*` tags.
- Version injection via linker flags in release workflow.

## Documentation Standards
- Keep README usage examples in sync with implemented command flags.
- Keep docs in `docs/` as source of truth for architecture and roadmap.
- Update roadmap/changelog context when significant capabilities change.