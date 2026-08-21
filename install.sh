#!/usr/bin/env bash
set -e
REPO="infraflakes/sutils"
BINARY_NAME="sn"
INSTALL_DIR="$HOME/.local/bin"

mkdir -p "$INSTALL_DIR"

echo "Finding the latest version tag for $REPO..."

LATEST_TAG=$(curl -sI "https://github.com/$REPO/releases/latest" | awk -F '/' '/^[Ll]ocation:/ {gsub(/[[:space:]\r\n]/,""); print $NF}')

if [ -z "$LATEST_TAG" ]; then
    echo "Error: Could not resolve the latest release tag." >&2
    exit 1
fi

VERSION="${LATEST_TAG#v}" 
ASSET_NAME="sutils_${VERSION}_linux_amd64.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$ASSET_NAME"

echo "Downloading $ASSET_NAME ($LATEST_TAG)..."

if ! curl -sSL "$DOWNLOAD_URL" | tar -x -z -f - -C "$INSTALL_DIR" "$BINARY_NAME"; then
    echo "Error: Failed to download or extract the binary." >&2
    exit 1
fi

chmod +x "$INSTALL_DIR/$BINARY_NAME"
echo "Successfully installed $BINARY_NAME to $INSTALL_DIR!"

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo " Warning: $INSTALL_DIR is not in your PATH."
    echo "Add this to your ~/.bashrc or ~/.zshrc:"
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
fi
