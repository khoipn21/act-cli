# Codebase Summary

## Repository Purpose
`act-cli` provides a cross-platform CLI for ACT runtime bootstrap, migration to Codex layout, diagnostics, and self-update.

## Top-Level Structure
- `cmd/act`: process entrypoint.
- `internal/app`: command router and update notice logic.
- `internal/commands`: user-facing command implementations and TUI helpers.
- `internal/kit`: runtime install/migrate and GitHub kit retrieval.
- `internal/config`: local user config persistence (`~/.act/config.json`).
- `internal/fsutil`: tree/file copy utilities.
- `internal/version`: version and repository constants.
- `scripts`: setup/install scripts for Bash and PowerShell.
- `.github/workflows`: CI and release automation.

## Current Size Snapshot (approx)
- `internal/commands`: 22 files, ~1690 LOC
- `internal/kit`: 10 files, ~1012 LOC
- `scripts`: 4 files, ~343 LOC
- `internal/app`: 3 files, ~227 LOC

## Command Surface
- `new`, `init`, `migrate`, `doctor`, `config`, `skills`, `agents`, `commands`, `plans`, `versions`, `update`.

## Key Flows
1. Init flow
- Parse args and mode.
- Optionally run interactive wizard.
- Load env values (local `.env` first, then user profile candidates).
- Resolve kit from local path or GitHub fallback.
- Install runtime files into target.

2. Migration flow
- Parse target/scope/options.
- Ensure source `.claude` exists (or install from kit).
- Copy skills/rules/hooks to `.codex` layout.
- Convert agents to TOML and migrate hooks to `hooks.json` + wrappers.
- Merge managed agent config into `.codex/config.toml`.

3. Update flow
- Fetch latest release metadata from GitHub.
- Compare semver with current version.
- Download matching binary asset.
- Replace executable (deferred script on Windows).

## Test Coverage Focus
- Command arg parsing and flag behavior.
- App routing and unknown command handling.
- Migration conversion behavior.
- Init env parsing and precedence.
- Config store load/save behavior.

## External Dependencies
- Standard library only for core implementation.
- GitHub APIs for kit/release retrieval.
- `gh` optional in scripts for authenticated release downloads.