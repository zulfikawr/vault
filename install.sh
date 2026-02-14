#!/bin/bash

set -e

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture names
case "$ARCH" in
  x86_64)
    ARCH="amd64"
    ;;
  aarch64)
    ARCH="arm64"
    ;;
  arm64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Map OS names
case "$OS" in
  linux)
    OS="linux"
    ;;
  darwin)
    OS="darwin"
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

echo "Detecting system: $OS-$ARCH"
echo ""

# Get latest release info
echo "Fetching latest release..."
LATEST_RELEASE=$(curl -s https://api.github.com/repos/zulfikawr/vault/releases/latest)
DOWNLOAD_URL=$(echo "$LATEST_RELEASE" | grep "browser_download_url" | grep "vault-$OS-$ARCH" | cut -d '"' -f 4 | head -n 1)
VERSION=$(echo "$LATEST_RELEASE" | grep '"tag_name"' | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
  echo "Could not find release for $OS-$ARCH"
  exit 1
fi

echo "Found version: $VERSION"
echo ""

# Download binary
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

cd "$TEMP_DIR"

echo "Downloading..."
curl -s -L -o vault "$DOWNLOAD_URL"
chmod +x vault

echo ""

# Install to /usr/local/bin
echo "Installing to /usr/local/bin/vault..."
if [ -w /usr/local/bin ]; then
  mv vault /usr/local/bin/
else
  sudo mv vault /usr/local/bin/
fi

echo ""

# Verify installation
if command -v vault &> /dev/null; then
  INSTALLED_VERSION=$(vault version 2>/dev/null || echo "unknown")
  echo "✓ Vault installed successfully!"
  echo "Version: $INSTALLED_VERSION"
  echo ""
  echo "Next steps:"
  echo "  vault init --email \"email@example.com\" --password \"yourpassword\" --username \"yourusername\""
  echo "  vault serve"
else
  echo "Installation failed"
  exit 1
fi
