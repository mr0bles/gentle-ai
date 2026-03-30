#!/bin/bash
set -e

# gentle-ai install script (build from source with automatic Go installation)
# Usage: curl -fsSL https://raw.githubusercontent.com/mr0bles/gentle-ai/main/scripts/install.sh | bash

REPO="mr0bles/gentle-ai"
INSTALL_DIR="${HOME}/.local/bin"
TEMP_DIR=$(mktemp -d)

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "📦 Installing gentle-ai from source..."

# ============================================================================
# Install Go if not available
# ============================================================================
install_go() {
    echo "⬇️  Installing Go..."

    OS="$(uname -s)"
    ARCH="$(uname -m)"

    case "$OS" in
        Linux) ;;
        *) echo "❌ Unsupported OS: $OS"; exit 1 ;;
    esac

    case "$ARCH" in
        x86_64)   GOARCH="amd64" ;;
        aarch64)   GOARCH="arm64" ;;
        armv7l)    GOARCH="armv6l" ;;
        *)         echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
    esac

    # Use latest known stable Go version (1.23.x is latest as of 2026)
    GO_VERSION="1.23.5"
    
    echo "   Using Go ${GO_VERSION}"
    
    GO_URL="https://go.dev/dl/go${GO_VERSION}.linux-${GOARCH}.tar.gz"

    echo "   Downloading Go ${GO_VERSION}..."

    # Download and install Go
    curl -fsSL "$GO_URL" -o /tmp/go.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz

    # Add to PATH for this session
    export PATH="/usr/local/go/bin:${PATH}"

    echo "✅ Go ${GO_VERSION} installed"
}

# Check if Go is available
if ! command -v go &> /dev/null; then
    install_go
else
    # Check Go version (need 1.21+)
    GO_VERSION_NUM=$(go version 2>/dev/null | grep -oP 'go\K[0-9]+\.[0-9]+' | head -1)
    GO_MAJOR=$(echo $GO_VERSION_NUM | cut -d. -f1)
    GO_MINOR=$(echo $GO_VERSION_NUM | cut -d. -f2)

    if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 21 ]); then
        echo -e "${YELLOW}Go version too old ($GO_VERSION_NUM), upgrading...${NC}"
        install_go
    else
        echo "✅ Go $(go version | grep -oP 'go\K[0-9.]+') found"
    fi
fi

# Ensure PATH includes Go
export PATH="/usr/local/go/bin:${PATH}:${HOME}/.local/bin"

# ============================================================================
# Clone and build
# ============================================================================

# Detect OS/ARCH for Go
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux)  GOOS="linux" ;;
    Darwin) GOOS="darwin" ;;
    *)      echo "❌ Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64)  GOARCH="amd64" ;;
    aarch64) GOARCH="arm64" ;;
    arm64)   GOARCH="arm64" ;;
    *)       echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "   Building for ${GOOS}/${GOARCH}..."

# Clone repo
git clone --depth 1 "https://github.com/${REPO}.git" "$TEMP_DIR/gentle-ai"
cd "$TEMP_DIR/gentle-ai"

# Build
GOOS="$GOOS" GOARCH="$GOARCH" go build -o "${TEMP_DIR}/gentle-ai" ./cmd/gentle-ai

# Install
mkdir -p "$INSTALL_DIR"
mv "${TEMP_DIR}/gentle-ai" "$INSTALL_DIR/"
chmod +x "${INSTALL_DIR}/gentle-ai"

# Cleanup
rm -rf "$TEMP_DIR"

echo -e "${GREEN}✅ gentle-ai installed to ${INSTALL_DIR}/gentle-ai${NC}"
echo ""
echo "Add to PATH if needed:"
echo "   echo 'export PATH=\"\${PATH}:${INSTALL_DIR}\"' >> ~/.bashrc"
echo "   source ~/.bashrc"
