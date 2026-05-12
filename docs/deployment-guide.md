# Deployment Guide

## Prerequisites
- Go 1.22+
- GitHub CLI (`gh`) optional but recommended for private repo workflows.
- GitHub token env for private access when needed:
  - `ACT_GITHUB_TOKEN` (preferred)
  - `GITHUB_TOKEN`
  - `GH_TOKEN`

## Local Build
```bash
go test ./...
go build -o build/act ./cmd/act
```

## One-Run Setup
- Bash:
```bash
./scripts/setup.sh --kit-path ../act-kit --scope project --init-target . --force
```
- PowerShell:
```powershell
./scripts/setup.ps1 -KitPath ..\act-kit -Scope project -InitTarget . -Force
```

## Install From Release
- Bash installer:
```bash
curl -fsSL https://raw.githubusercontent.com/khoipn21/act-cli/main/scripts/install.sh | bash
```
- PowerShell installer:
```powershell
irm https://raw.githubusercontent.com/khoipn21/act-cli/main/scripts/install.ps1 | iex
```

## CI
- Workflow: `.github/workflows/ci.yml`
- Trigger: pushes to `main`, pull requests.
- Jobs: `go test ./...`, `go build ./cmd/act` on Ubuntu/Windows/macOS.

## Release
- Workflow: `.github/workflows/release.yml`
- Trigger: tags matching `v*`.
- Builds matrix artifacts:
  - linux amd64/arm64
  - darwin amd64/arm64
  - windows amd64
- Publishes binaries to GitHub Releases.

## Post-Release Verification
1. `act update --check` reports latest tag.
2. Fresh install works on each OS.
3. `act versions` prints injected version.
4. `act init` and `act migrate --to codex` smoke tests pass.

## Rollback
- Reinstall a known tag using installer `--version <tag>`.
- Validate command surface with `act commands`.