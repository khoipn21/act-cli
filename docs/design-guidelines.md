# Design Guidelines

## CLI UX Principles
- Keep primary commands short and task-oriented.
- Use sensible defaults (`project` scope, interactive mode when terminal input exists).
- Offer non-interactive flags for automation.

## Output Style
- Human-readable summaries for long-running operations.
- Deterministic lines for scripts (`act commands`, `act config list`).
- Error messages should identify the failing flag/path/action.

## Interactive Wizards
- Use clear section headers: target, scope, options.
- Display resolved defaults before confirmation.
- Keep prompts binary/simple (`yes/no`, bounded values).

## Cross-Platform Design
- Ensure parity across Bash and PowerShell helper scripts.
- Avoid OS-specific assumptions in command behavior.
- Use platform-aware binary naming and install paths.

## Backward Compatibility
- Preserve existing command names and aliases where possible.
- Introduce new flags without breaking current defaults.
- Keep migration flow additive and explicit (`--to codex`).

## Documentation Design
- README should stay task-first (install, setup, migrate, update).
- Deep technical details belong in `docs/` files.
- Keep examples copy-paste ready for both shell families.