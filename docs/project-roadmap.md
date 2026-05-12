# Project Roadmap

## Status Summary
- Phase 1: Core CLI bootstrap and command routing - Complete
- Phase 2: Runtime install/init flows - Complete
- Phase 3: Codex migration flows - Complete
- Phase 4: Release and self-update pipeline - Complete
- Phase 5: Documentation baseline - In Progress

## Completed Milestones
1. Implemented command dispatcher and help catalog.
2. Added runtime install commands (`new`, `init`) with scope handling.
3. Added migration command for Codex target with hook/agent conversion.
4. Added setup/install scripts for Bash and PowerShell.
5. Added CI and release automation workflows.
6. Added test coverage for parsing and migration critical paths.

## Near-Term Work
1. Expand migration validation tests for complex hook matcher cases.
2. Add smoke-test script for end-to-end init + migrate + update checks.
3. Harden network retry behavior for GitHub fetch operations.
4. Improve doctor diagnostics beyond executable/kit checks.

## Medium-Term Work
1. Add richer `act plans` validation (schema + lint-like checks).
2. Add optional checksum/signature verification for downloaded binaries.
3. Add structured machine-readable output mode for automation.

## Risks and Dependencies
- Dependency on `act-kit` repository structure stability.
- GitHub API availability and token permission constraints.
- Cross-platform process replacement differences on Windows.

## Definition of Done for Documentation Baseline
- `docs/` contains overview, standards, summary, architecture, roadmap, deployment, and design files.
- Content reflects current command surface and workflow behavior.
- Future updates can run as incremental `act:docs update` operations.
