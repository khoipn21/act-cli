#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

SCOPE="project"
KIT_PATH=""
INIT_TARGET="."
FORCE_FLAG=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --scope)
      SCOPE="$2"; shift 2 ;;
    --kit-path)
      KIT_PATH="$2"; shift 2 ;;
    --init-target)
      INIT_TARGET="$2"; shift 2 ;;
    --force)
      FORCE_FLAG="--force"; shift ;;
    *)
      echo "Unknown argument: $1" >&2
      echo "Usage: scripts/setup.sh [--scope project|global] [--kit-path <path>] [--init-target <dir>] [--force]" >&2
      exit 1 ;;
  esac
done

if ! command -v go >/dev/null 2>&1; then
  echo "Go is required but not found. Install Go 1.22+ first: https://go.dev/dl/" >&2
  exit 1
fi

OS="$(uname -s)"
ARCH="$(uname -m)"
case "$OS" in
  Linux*)  OS_NAME="linux" ;;
  Darwin*) OS_NAME="darwin" ;;
  MINGW*|MSYS*|CYGWIN*) OS_NAME="windows" ;;
  *) OS_NAME="unknown" ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH_NAME="amd64" ;;
  arm64|aarch64) ARCH_NAME="arm64" ;;
  *) ARCH_NAME="$ARCH" ;;
esac

echo "[setup] Detected OS: ${OS_NAME}, Arch: ${ARCH_NAME}"
echo "[setup] Running tests"
go test ./...

mkdir -p build
BIN_NAME="act"
if [[ "$OS_NAME" == "windows" ]]; then
  BIN_NAME="act.exe"
fi

echo "[setup] Building binary"
go build -o "build/${BIN_NAME}" ./cmd/act

INSTALL_DIR="$HOME/.local/bin"
if [[ "$OS_NAME" == "darwin" ]]; then
  INSTALL_DIR="$HOME/bin"
fi
if [[ "$OS_NAME" == "windows" ]]; then
  INSTALL_DIR="${USERPROFILE:-$HOME}/bin"
fi
mkdir -p "$INSTALL_DIR"
cp "build/${BIN_NAME}" "$INSTALL_DIR/${BIN_NAME}"

echo "[setup] Installed: $INSTALL_DIR/${BIN_NAME}"
echo "[setup] Ensure install dir is in PATH: $INSTALL_DIR"

ACT_BIN="$INSTALL_DIR/$BIN_NAME"
INIT_ARGS=("init" "$INIT_TARGET" "--scope" "$SCOPE" "--non-interactive")
if [[ -n "$KIT_PATH" ]]; then
  INIT_ARGS+=("--kit" "$KIT_PATH")
fi
if [[ -n "$FORCE_FLAG" ]]; then
  INIT_ARGS+=("$FORCE_FLAG")
fi

echo "[setup] Running: ${ACT_BIN} ${INIT_ARGS[*]}"
"$ACT_BIN" "${INIT_ARGS[@]}"

echo "[setup] Done"