#!/usr/bin/env bash
set -euo pipefail

REPO="khoipn21/act-cli"
VERSION=""
INSTALL_DIR=""

usage() {
  cat <<'EOF'
Usage: install.sh [--repo <owner/repo>] [--version <tag>] [--install-dir <path>]

Examples:
  curl -fsSL https://raw.githubusercontent.com/khoipn21/act-cli/main/scripts/install.sh | bash
  curl -fsSL https://raw.githubusercontent.com/khoipn21/act-cli/main/scripts/install.sh | bash -s -- --version v0.1.0
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo)
      REPO="$2"; shift 2 ;;
    --version)
      VERSION="$2"; shift 2 ;;
    --install-dir)
      INSTALL_DIR="$2"; shift 2 ;;
    -h|--help)
      usage; exit 0 ;;
    *)
      echo "Unknown argument: $1" >&2
      usage
      exit 1 ;;
  esac
done

if command -v curl >/dev/null 2>&1; then
  DOWNLOADER="curl -fsSL"
elif command -v wget >/dev/null 2>&1; then
  DOWNLOADER="wget -qO-"
else
  echo "curl or wget is required." >&2
  exit 1
fi

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Linux*) OS_NAME="linux"; EXT="" ;;
  Darwin*) OS_NAME="darwin"; EXT="" ;;
  MINGW*|MSYS*|CYGWIN*) OS_NAME="windows"; EXT=".exe" ;;
  *)
    echo "Unsupported OS: $OS" >&2
    exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH_NAME="amd64" ;;
  arm64|aarch64) ARCH_NAME="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1 ;;
esac

if [[ -z "$INSTALL_DIR" ]]; then
  if [[ "$OS_NAME" == "windows" ]]; then
    INSTALL_DIR="${USERPROFILE:-$HOME}/bin"
  elif [[ "$OS_NAME" == "darwin" && -w "/usr/local/bin" ]]; then
    INSTALL_DIR="/usr/local/bin"
  else
    INSTALL_DIR="$HOME/.local/bin"
  fi
fi

ASSET="act-${OS_NAME}-${ARCH_NAME}${EXT}"
if [[ -n "$VERSION" ]]; then
  URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"
else
  URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"
fi

echo "[install] Repo: ${REPO}"
echo "[install] OS/Arch: ${OS_NAME}/${ARCH_NAME}"
echo "[install] Asset: ${ASSET}"
echo "[install] URL: ${URL}"

TMP_FILE="$(mktemp)"
TMP_DIR="$(mktemp -d)"
mkdir -p "$INSTALL_DIR"

if command -v gh >/dev/null 2>&1 && gh auth status >/dev/null 2>&1; then
  echo "[install] Using authenticated gh release download"
  if [[ -n "$VERSION" ]]; then
    gh release download "$VERSION" -R "$REPO" -p "$ASSET" -D "$TMP_DIR" --clobber
  else
    gh release download -R "$REPO" -p "$ASSET" -D "$TMP_DIR" --clobber
  fi
  cp "$TMP_DIR/$ASSET" "$TMP_FILE"
else
  if [[ "$DOWNLOADER" == "curl -fsSL" ]]; then
    curl -fsSL "$URL" -o "$TMP_FILE"
  else
    wget -q "$URL" -O "$TMP_FILE"
  fi
fi

DEST="${INSTALL_DIR}/act${EXT}"
mv "$TMP_FILE" "$DEST"
if [[ "$OS_NAME" != "windows" ]]; then
  chmod +x "$DEST"
fi

echo "[install] Installed to: $DEST"
echo "[install] If needed, add to PATH: $INSTALL_DIR"
echo "[install] Verify with: act commands"

rm -rf "$TMP_DIR"
