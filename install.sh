#!/usr/bin/env sh
# install.sh — download the latest canvas release into a local bin dir.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/jjuanrivvera/canvas-cli/main/install.sh | sh
#
# Environment:
#   CANVAS_VERSION       pin a version, e.g. "v1.10.0" (default: latest)
#   CANVAS_INSTALL_DIR   install location (default: $HOME/.local/bin)
#   NO_COLOR             disable colored output

set -eu

# --- configuration (per project) ---------------------------------------------
REPO="jjuanrivvera/canvas-cli"         # GitHub <owner>/<repo>
PROJECT="canvas-cli"                   # release asset prefix (goreleaser project_name)
BINARY="canvas"                        # installed binary (name inside the archive)
VERSION_IN_NAME=0                      # 1 if the asset filename embeds the version
ARCH_AMD64="x86_64"                    # token used for x86_64 hosts in asset names
VERSION="${CANVAS_VERSION:-latest}"
INSTALL_DIR="${CANVAS_INSTALL_DIR:-$HOME/.local/bin}"
# ------------------------------------------------------------------------------

if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
  C_RESET=$(printf '\033[0m'); C_BOLD=$(printf '\033[1m'); C_DIM=$(printf '\033[2m')
  C_CYAN=$(printf '\033[36m'); C_GREEN=$(printf '\033[32m')
  C_YELLOW=$(printf '\033[33m'); C_RED=$(printf '\033[31m')
else
  C_RESET=''; C_BOLD=''; C_DIM=''; C_CYAN=''; C_GREEN=''; C_YELLOW=''; C_RED=''
fi
step() { printf '  %s→%s %s\n' "$C_CYAN" "$C_RESET" "$*"; }
ok()   { printf '  %s✓%s %s\n' "$C_GREEN" "$C_RESET" "$*"; }
warn() { printf '  %s!%s %s\n' "$C_YELLOW" "$C_RESET" "$*"; }
err()  { printf '  %s✗%s %s\n' "$C_RED" "$C_RESET" "$*" >&2; exit 1; }
has()  { command -v "$1" >/dev/null 2>&1; }

# --- detect platform ----------------------------------------------------------
case "$(uname -s)" in
  Darwin) os=darwin ;;
  Linux)  os=linux ;;
  *) err "unsupported OS: $(uname -s) (Windows users: download from https://github.com/$REPO/releases)" ;;
esac
case "$(uname -m)" in
  x86_64|amd64)  arch="$ARCH_AMD64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) err "unsupported architecture: $(uname -m)" ;;
esac
step "Detected ${os}/${arch}"

# --- resolve version (latest by default) --------------------------------------
if [ "$VERSION" = "latest" ]; then
  VERSION=$(curl -fsSLI -o /dev/null -w '%{url_effective}' \
    "https://github.com/$REPO/releases/latest" 2>/dev/null \
    | sed -n 's|.*/tag/\(v[^/]*\).*|\1|p')
  [ -n "$VERSION" ] || err "could not resolve the latest version (check network)"
fi
case "$VERSION" in v*) ;; *) VERSION="v$VERSION" ;; esac
ver_clean="${VERSION#v}"
step "Installing ${C_BOLD}${PROJECT} ${VERSION}${C_RESET}"

# --- build asset name + base URL ----------------------------------------------
if [ "$VERSION_IN_NAME" = "1" ]; then
  asset="${PROJECT}_${ver_clean}_${os}_${arch}.tar.gz"
else
  asset="${PROJECT}_${os}_${arch}.tar.gz"
fi
base="https://github.com/$REPO/releases/download/$VERSION"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

# --- download -----------------------------------------------------------------
step "Downloading ${asset}"
curl -fsSL "$base/$asset" -o "$tmp/$asset" \
  || err "failed to download $base/$asset (no release asset for ${os}/${arch}?)"
curl -fsSL "$base/checksums.txt" -o "$tmp/checksums.txt" \
  || err "failed to download checksums.txt; refusing to install without verification"

# --- verify sha256 (required) -------------------------------------------------
expected=$(grep " ${asset}\$" "$tmp/checksums.txt" | awk '{print $1}')
[ -n "$expected" ] || err "${asset} not listed in checksums.txt; refusing to install"
if has sha256sum; then
  actual=$(sha256sum "$tmp/$asset" | awk '{print $1}')
elif has shasum; then
  actual=$(shasum -a 256 "$tmp/$asset" | awk '{print $1}')
else
  err "no sha256 tool found; install sha256sum or shasum and retry"
fi
[ "$expected" = "$actual" ] || err "checksum mismatch (expected=${expected} got=${actual})"
step "Verified checksum"

# --- verify cosign signature (optional, only if cosign is installed) ----------
if has cosign; then
  identity="^https://github.com/${REPO}/\.github/workflows/.+@refs/tags/"
  issuer="https://token.actions.githubusercontent.com"
  verify_checksums() {
    cosign verify-blob "$@" \
      --certificate-identity-regexp "$identity" --certificate-oidc-issuer "$issuer" \
      "$tmp/checksums.txt" >/dev/null 2>&1
  }
  if curl -fsSL "$base/checksums.txt.cosign.bundle" -o "$tmp/cosign.bundle" 2>/dev/null; then
    if verify_checksums --bundle "$tmp/cosign.bundle"; then
      step "Verified cosign signature"
    else
      err "cosign signature verification failed"
    fi
  elif curl -fsSL "$base/checksums.txt.pem" -o "$tmp/checksums.pem" 2>/dev/null \
    && curl -fsSL "$base/checksums.txt.sig" -o "$tmp/checksums.sig" 2>/dev/null; then
    if verify_checksums --certificate "$tmp/checksums.pem" --signature "$tmp/checksums.sig"; then
      step "Verified cosign signature"
    else
      err "cosign signature verification failed"
    fi
  else
    warn "no cosign signature published for this release; skipping (checksum verified)"
  fi
fi

# --- extract & install --------------------------------------------------------
tar -xzf "$tmp/$asset" -C "$tmp"
[ -f "$tmp/$BINARY" ] || err "archive did not contain a '$BINARY' binary"
mkdir -p "$INSTALL_DIR"
if has install; then
  install -m 0755 "$tmp/$BINARY" "$INSTALL_DIR/$BINARY"
else
  mv "$tmp/$BINARY" "$INSTALL_DIR/$BINARY"
  chmod +x "$INSTALL_DIR/$BINARY"
fi
ok "Installed ${C_BOLD}${BINARY} ${VERSION}${C_RESET} to ${INSTALL_DIR}/${BINARY}"

# --- PATH hint ----------------------------------------------------------------
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    printf '\n'
    warn "${INSTALL_DIR} is not on your PATH"
    printf '    add this to your shell profile (~/.zshrc, ~/.bashrc, ...):\n'
    # shellcheck disable=SC2016  # literal $PATH is intentional
    printf '      %sexport PATH="%s:$PATH"%s\n' "$C_DIM" "$INSTALL_DIR" "$C_RESET"
    ;;
esac

printf '\n  %sNext steps%s\n' "$C_BOLD" "$C_RESET"
printf '    %s%s --help%s\n' "$C_CYAN" "$BINARY" "$C_RESET"
printf '    %s%s completion --help%s   set up shell completion\n\n' "$C_CYAN" "$BINARY" "$C_RESET"
