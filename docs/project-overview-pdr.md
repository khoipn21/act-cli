# Project Overview PDR

## Overview
`act-cli` is a Go command-line tool that bootstraps, manages, and migrates ACT runtime assets for developer environments.

Core outcomes:
- Initialize runtime scaffolding from `act-kit` into project or global scope.
- Migrate Claude-style runtime layout to Codex-compatible layout.
- Provide lifecycle commands for diagnostics, config, versioning, and self-update.

## Product Goals
- Fast setup with minimal required flags.
- Safe defaults for project scope and non-destructive file writes unless forced.
- Cross-platform support for Linux, macOS, and Windows.
- Predictable migration path from `.claude` to `.codex` assets.

## Primary Users
- Developers onboarding new repos with ACT runtime.
- Teams migrating existing runtime automation to Codex.
- Maintainers managing local CLI versions and runtime configuration.

## Functional Requirements
1. `act init` installs runtime from local `act-kit` or GitHub fallback.
2. `act migrate --to codex` transforms runtime assets and hook config.
3. `act update` checks/releases binary updates from GitHub releases.
4. `act doctor` validates executable and kit resolution health.
5. `act config` supports key-value list/get/set persistence.
6. `act skills` and `act agents` enumerate catalog artifacts from kit.
7. `act plans validate` checks plan phase-file references.

## Non-Functional Requirements
- Platform compatibility: Windows/macOS/Linux.
- Runtime language: Go 1.22+.
- Simple command contract, low dependency footprint.
- Deterministic outputs and error messages for CLI automation.

## In Scope
- Runtime installation and migration orchestration.
- Local/GitHub kit source selection.
- Release asset download and binary replacement.

## Out of Scope
- Server runtime, daemon mode, or remote orchestration backend.
- Rich stateful config schemas beyond key-value map.
- Plugin marketplace or package dependency solver.

## Success Metrics
- New users can run setup (`act init` or setup scripts) without manual file copying.
- Migration completes with valid `.codex` artifacts and migrated hooks.
- CI passes on all target operating systems.

## Constraints
- Depends on `act-kit` structure (`claude/`, `docs/`, `plans/templates`).
- GitHub token required for private or token-gated kit/release access.
- Update flow on Windows uses deferred replacement script.