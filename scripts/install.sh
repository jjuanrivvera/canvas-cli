#!/bin/bash
# Canvas CLI Installation Script for Unix-like systems (macOS/Linux)

set -e

# Configuration
REPO="jjuanrivvera/canvas-cli"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="canvas"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        darwin)
            OS="Darwin"
            ;;
        linux)
            OS="Linux"
            ;;
        *)
            log_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac

    case "$ARCH" in
        x86_64|amd64)
            ARCH="x86_64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7l)
            ARCH="armv7"
            ;;
        *)
            log_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac

    log_info "Detected platform: $OS $ARCH"
}

# Get latest release version
get_latest_version() {
    log_info "Fetching latest release..."

    VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$VERSION" ]; then
        log_error "Failed to fetch latest version"
        exit 1
    fi

    log_info "Latest version: $VERSION"
}

# Download and install binary
install_binary() {
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/canvas_${VERSION#v}_${OS}_${ARCH}.tar.gz"
    TMP_DIR=$(mktemp -d)

    log_info "Downloading from: $DOWNLOAD_URL"

    if ! curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/canvas.tar.gz"; then
        log_error "Failed to download binary"
        rm -rf "$TMP_DIR"
        exit 1
    fi

    log_info "Extracting archive..."
    tar -xzf "$TMP_DIR/canvas.tar.gz" -C "$TMP_DIR"

    log_info "Installing to $INSTALL_DIR..."

    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/"
    else
        log_warn "Requires sudo for installation to $INSTALL_DIR"
        sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/"
    fi

    chmod +x "$INSTALL_DIR/$BINARY_NAME"

    rm -rf "$TMP_DIR"

    log_info "Installed successfully!"
}

# Verify installation
verify_installation() {
    log_info "Verifying installation..."

    if ! command -v $BINARY_NAME &> /dev/null; then
        log_error "$BINARY_NAME not found in PATH"
        log_warn "You may need to add $INSTALL_DIR to your PATH"
        log_warn "Add this to your shell profile:"
        log_warn "  export PATH=\"\$PATH:$INSTALL_DIR\""
        exit 1
    fi

    VERSION_OUTPUT=$($BINARY_NAME version)
    log_info "Verification successful!"
    echo "$VERSION_OUTPUT"
}

# Install shell completion
install_completion() {
    log_info "Installing shell completion..."

    SHELL_NAME=$(basename "$SHELL")

    case "$SHELL_NAME" in
        bash)
            COMPLETION_DIR="$HOME/.canvas-completion"
            mkdir -p "$COMPLETION_DIR"
            $BINARY_NAME completion bash > "$COMPLETION_DIR/canvas.bash"

            if ! grep -q "source $COMPLETION_DIR/canvas.bash" "$HOME/.bashrc"; then
                echo "source $COMPLETION_DIR/canvas.bash" >> "$HOME/.bashrc"
                log_info "Added completion to ~/.bashrc"
            fi
            ;;
        zsh)
            COMPLETION_DIR="${fpath[1]:-$HOME/.zsh/completion}"
            mkdir -p "$COMPLETION_DIR"
            $BINARY_NAME completion zsh > "$COMPLETION_DIR/_canvas"
            log_info "Added completion to $COMPLETION_DIR/_canvas"
            log_warn "Restart your shell or run: autoload -U compinit && compinit"
            ;;
        fish)
            COMPLETION_DIR="$HOME/.config/fish/completions"
            mkdir -p "$COMPLETION_DIR"
            $BINARY_NAME completion fish > "$COMPLETION_DIR/canvas.fish"
            log_info "Added completion to $COMPLETION_DIR/canvas.fish"
            ;;
        *)
            log_warn "Shell completion not configured for $SHELL_NAME"
            log_info "You can manually run: canvas completion $SHELL_NAME"
            ;;
    esac
}

# Main installation process
main() {
    echo "════════════════════════════════════════"
    echo "  Canvas CLI Installation Script"
    echo "════════════════════════════════════════"
    echo ""

    detect_platform
    get_latest_version
    install_binary
    verify_installation

    echo ""
    read -p "Install shell completion? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        install_completion
    fi

    echo ""
    echo "════════════════════════════════════════"
    log_info "Installation complete!"
    echo "════════════════════════════════════════"
    echo ""
    echo "Next steps:"
    echo "  1. Authenticate: canvas auth login --instance https://canvas.instructure.com"
    echo "  2. Test: canvas courses list"
    echo "  3. Get help: canvas --help"
    echo ""
    echo "Documentation: https://github.com/$REPO/tree/main/docs"
}

main
