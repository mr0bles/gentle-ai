#!/bin/bash
set -e

# gentle-ai install script (build from source)
# Usage: curl -fsSL https://raw.githubusercontent.com/mr0bles/gentle-ai/main/install.sh | bash

REPO="mr0bles/gentle-ai"
INSTALL_DIR="${HOME}/.local/bin"
TEMP_DIR=$(mktemp -d)

echo "📦 Installing gentle-ai from source..."

# Clone repo
git clone --depth 1 "https://github.com/${REPO}.git" "$TEMP_DIR/gentle-ai"
cd "$TEMP_DIR/gentle-ai"

# Detect OS/ARCH
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

# Build
GOOS="$GOOS" GOARCH="$GOARCH" go build -o "${TEMP_DIR}/gentle-ai" .

# Install
mkdir -p "$INSTALL_DIR"
mv "${TEMP_DIR}/gentle-ai" "$INSTALL_DIR/"
chmod +x "${INSTALL_DIR}/gentle-ai"

# Cleanup
rm -rf "$TEMP_DIR"

echo "✅ gentle-ai installed to ${INSTALL_DIR}/gentle-ai"
echo ""
echo "Add to PATH if needed:"
echo "   export PATH=\"\${PATH}:${INSTALL_DIR}\" >> ~/.bashrc"
echo "   source ~/.bashrc"
