param(
  [string]$Repo = "khoipn21/act-cli",
  [string]$Version = "",
  [string]$InstallDir = ""
)

$ErrorActionPreference = "Stop"

$isWindows = ($env:OS -eq "Windows_NT")
$isMacOS = [System.Runtime.InteropServices.RuntimeInformation]::IsOSPlatform([System.Runtime.InteropServices.OSPlatform]::OSX)
$isLinux = [System.Runtime.InteropServices.RuntimeInformation]::IsOSPlatform([System.Runtime.InteropServices.OSPlatform]::Linux)

$osName = if ($isWindows) { "windows" } elseif ($isMacOS) { "darwin" } elseif ($isLinux) { "linux" } else { throw "Unsupported OS" }
$arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLower()
$archName = switch ($arch) {
  "x64" { "amd64" }
  "amd64" { "amd64" }
  "arm64" { "arm64" }
  default { throw "Unsupported architecture: $arch" }
}

$ext = if ($isWindows) { ".exe" } else { "" }
if ([string]::IsNullOrWhiteSpace($InstallDir)) {
  if ($isWindows) {
    $InstallDir = Join-Path $env:USERPROFILE "bin"
  } elseif ($isMacOS -and (Test-Path "/usr/local/bin")) {
    $InstallDir = "/usr/local/bin"
  } else {
    $InstallDir = Join-Path $env:HOME ".local/bin"
  }
}

$asset = "act-$osName-$archName$ext"
$url = if ([string]::IsNullOrWhiteSpace($Version)) {
  "https://github.com/$Repo/releases/latest/download/$asset"
} else {
  "https://github.com/$Repo/releases/download/$Version/$asset"
}

Write-Host "[install] Repo: $Repo"
Write-Host "[install] OS/Arch: $osName/$archName"
Write-Host "[install] Asset: $asset"
Write-Host "[install] URL: $url"

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
$tmpFile = Join-Path ([System.IO.Path]::GetTempPath()) ("act-install-" + [System.Guid]::NewGuid().ToString() + $ext)
$tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) ("act-install-dir-" + [System.Guid]::NewGuid().ToString())
New-Item -ItemType Directory -Force -Path $tmpDir | Out-Null

if (Get-Command gh -ErrorAction SilentlyContinue) {
  try {
    gh auth status | Out-Null
    Write-Host "[install] Using authenticated gh release download"
    if ([string]::IsNullOrWhiteSpace($Version)) {
      gh release download -R $Repo -p $asset -D $tmpDir --clobber
    } else {
      gh release download $Version -R $Repo -p $asset -D $tmpDir --clobber
    }
    Copy-Item -Force (Join-Path $tmpDir $asset) $tmpFile
  } catch {
    Invoke-WebRequest -Uri $url -OutFile $tmpFile
  }
} else {
  Invoke-WebRequest -Uri $url -OutFile $tmpFile
}

$dest = Join-Path $InstallDir ("act" + $ext)
Move-Item -Force $tmpFile $dest

if (-not $isWindows) {
  & chmod +x $dest
}

function Add-ToUserPath([string]$PathToAdd) {
  if (-not $isWindows) { return }
  $currentUserPath = [Environment]::GetEnvironmentVariable("Path", "User")
  if ([string]::IsNullOrWhiteSpace($currentUserPath)) {
    $updated = $PathToAdd
  } else {
    $parts = $currentUserPath -split ';' | Where-Object { $_ -ne "" }
    if ($parts -contains $PathToAdd) {
      $updated = $currentUserPath
    } else {
      $updated = "$currentUserPath;$PathToAdd"
    }
  }
  [Environment]::SetEnvironmentVariable("Path", $updated, "User")
  if (-not (($env:Path -split ';') -contains $PathToAdd)) {
    $env:Path = "$env:Path;$PathToAdd"
  }
}

Add-ToUserPath $InstallDir

if ($isWindows) {
  $legacyNoExt = Join-Path $env:USERPROFILE ".local\\bin\\act"
  $legacyExe = Join-Path $env:USERPROFILE ".local\\bin\\act.exe"
  if ((Test-Path $legacyNoExt) -and ($legacyNoExt -ne $dest)) {
    Remove-Item -Force $legacyNoExt -ErrorAction SilentlyContinue
  }
  if ((Test-Path $legacyExe) -and ($legacyExe -ne $dest)) {
    Remove-Item -Force $legacyExe -ErrorAction SilentlyContinue
  }
}

Write-Host "[install] Installed to: $dest"
Write-Host "[install] PATH updated for current user: $InstallDir"
Write-Host "[install] Verify with: act commands"

Remove-Item -Recurse -Force $tmpDir -ErrorAction SilentlyContinue
