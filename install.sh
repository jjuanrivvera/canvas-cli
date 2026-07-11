#!/bin/sh
# canvas installer for macOS and Linux.
#
#   curl -fsSL https://raw.githubusercontent.com/jjuanrivvera/canvas-cli/main/install.sh | sh
#
# Downloads the release archive matching your OS/arch, verifies its SHA-256 against the
# release checksums.txt, and installs the binary. Configure via env vars:
#   CANVAS_VERSION=v1.2.3   pin a version (default: the latest release)
#   INSTALL_DIR=/usr/local/bin        install location (default; falls back to ~/.local/bin)
#
# Windows: use Scoop (see the README). This installer is for macOS and Linux.
set -eu

# --- per-tool configuration (the only lines cliwright templates per CLI) ---
REPO="jjuanrivvera/canvas-cli"
BINARY="canvas"
VERSION="${CANVAS_VERSION:-}"

die() { printf 'error: %s\n' "$1" >&2; exit 1; }

command -v curl >/dev/null 2>&1 || die "curl is required"
command -v tar  >/dev/null 2>&1 || die "tar is required"
if command -v sha256sum >/dev/null 2>&1; then shacmd="sha256sum"
elif command -v shasum  >/dev/null 2>&1; then shacmd="shasum -a 256"
else die "need sha256sum or shasum to verify the download"; fi

# --- detect platform (goreleaser naming: lowercase os, amd64/arm64) ---
os="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$os" in
  linux | darwin) ;;
  *) die "unsupported OS: $os (this installer covers Linux and macOS; use 'go install' otherwise)" ;;
esac
arch="$(uname -m)"
case "$arch" in
  x86_64 | amd64) arch="amd64"; archre="(amd64|x86_64)" ;;
  arm64 | aarch64) arch="arm64"; archre="(arm64|aarch64)" ;;
  *) die "unsupported architecture: $arch" ;;
esac

# --- resolve the version to install ---
if [ -z "$VERSION" ]; then
  VERSION="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' | head -1 | sed -E 's/.*"tag_name" *: *"([^"]+)".*/\1/')"
  [ -n "$VERSION" ] || die "could not determine the latest release; set CANVAS_VERSION"
fi

base="https://github.com/${REPO}/releases/download/${VERSION}"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT INT TERM

# --- discover the exact archive from checksums.txt (naming-agnostic) ---
curl -fsSL "${base}/checksums.txt" -o "${tmp}/checksums.txt" \
  || die "could not download checksums.txt for ${VERSION}"
archive="$(awk '{print $2}' "${tmp}/checksums.txt" | grep -iE "_${os}_${archre}\.tar\.gz$" | head -1)"
[ -n "$archive" ] || die "no ${os}/${arch} archive in ${VERSION}"

# --- download + verify the SHA-256 ---
curl -fsSL "${base}/${archive}" -o "${tmp}/${archive}" || die "download failed: ${archive}"
want="$(awk -v f="$archive" '$2 == f {print $1}' "${tmp}/checksums.txt")"
got="$(cd "$tmp" && $shacmd "$archive" | awk '{print $1}')"
[ -n "$want" ] || die "no checksum recorded for ${archive}"
[ "$want" = "$got" ] || die "checksum mismatch for ${archive}"

# --- extract + install ---
tar -xzf "${tmp}/${archive}" -C "$tmp"
[ -f "${tmp}/${BINARY}" ] || die "binary '${BINARY}' not found in the archive"
chmod +x "${tmp}/${BINARY}"

dir="${INSTALL_DIR:-/usr/local/bin}"
if mkdir -p "$dir" 2>/dev/null && [ -w "$dir" ]; then
  mv "${tmp}/${BINARY}" "${dir}/${BINARY}"
elif command -v sudo >/dev/null 2>&1 && [ -t 0 ]; then
  printf 'installing to %s (needs sudo)\n' "$dir" >&2
  sudo mv "${tmp}/${BINARY}" "${dir}/${BINARY}"
else
  dir="${HOME}/.local/bin"
  mkdir -p "$dir"
  mv "${tmp}/${BINARY}" "${dir}/${BINARY}"
fi

printf '%s %s installed to %s/%s\n' "$BINARY" "$VERSION" "$dir" "$BINARY"
# shellcheck disable=SC2016  # the PATH hint is a literal string for the user to copy, not an expansion
case ":${PATH}:" in
  *":${dir}:"*) ;;
  *) printf 'note: %s is not on your PATH — add:\n  export PATH="%s:$PATH"\n' "$dir" "$dir" >&2 ;;
esac
