param(
  [ValidateSet("project","global")]
  [string]$Scope = "project",
  [string]$KitPath = "",
  [string]$InitTarget = ".",
  [switch]$Force
)

$ErrorActionPreference = 'Stop'

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $repoRoot

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
  throw "Go is required but not found. Install Go 1.22+ first: https://go.dev/dl/"
}

$isWindows = ($env:OS -eq "Windows_NT")
$isMacOS = [System.Runtime.InteropServices.RuntimeInformation]::IsOSPlatform([System.Runtime.InteropServices.OSPlatform]::OSX)
$isLinux = [System.Runtime.InteropServices.RuntimeInformation]::IsOSPlatform([System.Runtime.InteropServices.OSPlatform]::Linux)

$osName = if ($isWindows) { "windows" } elseif ($isMacOS) { "darwin" } elseif ($isLinux) { "linux" } else { "unknown" }
$arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLower()
if ($arch -eq "x64") { $arch = "amd64" }
if ($arch -eq "arm64") { $arch = "arm64" }
Write-Host "[setup] Detected OS: $osName, Arch: $arch"

Write-Host "[setup] Running tests"
go test ./...

New-Item -ItemType Directory -Force -Path "build" | Out-Null
$binName = if ($isWindows) { "act.exe" } else { "act" }

Write-Host "[setup] Building binary"
go build -o (Join-Path "build" $binName) ./cmd/act

$installDir = if ($isWindows) {
  Join-Path $env:USERPROFILE "bin"
} elseif ($isMacOS) {
  Join-Path $env:HOME "bin"
} else {
  Join-Path $env:HOME ".local/bin"
}

New-Item -ItemType Directory -Force -Path $installDir | Out-Null
Copy-Item -Force (Join-Path "build" $binName) (Join-Path $installDir $binName)
Write-Host "[setup] Installed: $(Join-Path $installDir $binName)"

if ($isWindows) {
  $currentUserPath = [Environment]::GetEnvironmentVariable("Path", "User")
  if ([string]::IsNullOrWhiteSpace($currentUserPath)) {
    $updated = $installDir
  } else {
    $parts = $currentUserPath -split ';' | Where-Object { $_ -ne "" }
    if ($parts -contains $installDir) {
      $updated = $currentUserPath
    } else {
      $updated = "$currentUserPath;$installDir"
    }
  }
  [Environment]::SetEnvironmentVariable("Path", $updated, "User")
  if (-not (($env:Path -split ';') -contains $installDir)) {
    $env:Path = "$env:Path;$installDir"
  }
  $legacyNoExt = Join-Path $env:USERPROFILE ".local\\bin\\act"
  $legacyExe = Join-Path $env:USERPROFILE ".local\\bin\\act.exe"
  if ((Test-Path $legacyNoExt) -and ($legacyNoExt -ne (Join-Path $installDir $binName))) {
    Remove-Item -Force $legacyNoExt -ErrorAction SilentlyContinue
  }
  if ((Test-Path $legacyExe) -and ($legacyExe -ne (Join-Path $installDir $binName))) {
    Remove-Item -Force $legacyExe -ErrorAction SilentlyContinue
  }
  Write-Host "[setup] PATH updated for current user: $installDir"
} else {
  Write-Host "[setup] Ensure install dir is in PATH: $installDir"
}

$actBin = Join-Path $installDir $binName
$args = @("init", $InitTarget, "--scope", $Scope, "--non-interactive")
if ($KitPath -ne "") { $args += @("--kit", $KitPath) }
if ($Force) { $args += "--force" }

Write-Host "[setup] Running: $actBin $($args -join ' ')"
& $actBin @args
if ($LASTEXITCODE -ne 0) {
  throw "act init failed with exit code $LASTEXITCODE"
}

Write-Host "[setup] Done"
