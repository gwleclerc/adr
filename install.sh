#!/usr/bin/env sh
# adr installer — downloads the latest (or a pinned) release binary for your platform.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/gwleclerc/adr/main/install.sh | sh
#
# Environment variables:
#   ADR_VERSION      Release tag to install (e.g. v1.2.3). Defaults to the latest release.
#   ADR_INSTALL_DIR  Directory to install into. Defaults to /usr/local/bin if writable,
#                    otherwise $HOME/.local/bin.

set -eu

REPO="gwleclerc/adr"
BINARY="adr"

log() { printf '%s\n' "$*" >&2; }
fail() { log "error: $*"; exit 1; }

need() { command -v "$1" >/dev/null 2>&1 || fail "missing required command: $1"; }

need uname
need tar

# Prefer curl, fall back to wget.
if command -v curl >/dev/null 2>&1; then
  http_get() { curl -fsSL "$1"; }
  http_download() { curl -fsSL -o "$2" "$1"; }
elif command -v wget >/dev/null 2>&1; then
  http_get() { wget -qO- "$1"; }
  http_download() { wget -qO "$2" "$1"; }
else
  fail "either curl or wget is required"
fi

# --- Detect OS ---------------------------------------------------------------
os="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$os" in
  linux)  OS="linux" ;;
  darwin) OS="darwin" ;;
  msys*|mingw*|cygwin*) fail "on Windows, download the .zip from https://github.com/$REPO/releases" ;;
  *) fail "unsupported OS: $os" ;;
esac

# --- Detect architecture -----------------------------------------------------
arch="$(uname -m)"
case "$arch" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  i386|i686) ARCH="386" ;;
  armv7*|armv6*|arm) ARCH="armv7" ;;
  *) fail "unsupported architecture: $arch" ;;
esac

# --- Resolve version ---------------------------------------------------------
VERSION="${ADR_VERSION:-}"
if [ -z "$VERSION" ]; then
  log "resolving latest release..."
  VERSION="$(http_get "https://api.github.com/repos/$REPO/releases/latest" \
    | grep '"tag_name"' | head -n1 | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"
  [ -n "$VERSION" ] || fail "could not determine the latest release tag"
fi

# Archive names mirror the `make release` naming: adr_<tag>_<os>_<arch>.tar.gz
ARCHIVE="${BINARY}_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/${VERSION}/${ARCHIVE}"

# --- Choose install directory ------------------------------------------------
INSTALL_DIR="${ADR_INSTALL_DIR:-}"
if [ -z "$INSTALL_DIR" ]; then
  if [ -w /usr/local/bin ] 2>/dev/null; then
    INSTALL_DIR="/usr/local/bin"
  else
    INSTALL_DIR="$HOME/.local/bin"
  fi
fi
mkdir -p "$INSTALL_DIR"

# --- Download & install ------------------------------------------------------
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

log "downloading $URL"
http_download "$URL" "$TMP/$ARCHIVE" || fail "download failed (does the asset exist for $OS/$ARCH?)"

tar -xzf "$TMP/$ARCHIVE" -C "$TMP" || fail "failed to extract archive"
[ -f "$TMP/$BINARY" ] || fail "binary '$BINARY' not found in archive"

chmod +x "$TMP/$BINARY"
mv "$TMP/$BINARY" "$INSTALL_DIR/$BINARY"

log "installed $BINARY $VERSION to $INSTALL_DIR/$BINARY"
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *) log "note: $INSTALL_DIR is not in your PATH — add it, e.g. 'export PATH=\"$INSTALL_DIR:\$PATH\"'" ;;
esac
