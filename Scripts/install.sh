#!/bin/bash
# install.sh
# Linux installation script for Endmi.
#
# Steps:
#   - Detect system architecture (amd64, arm64, arm)
#   - Download target binary from GitHub Releases
#   - Install binary to /opt/endmi
#   - Create a symbolic link in /usr/local/bin for PATH access

set -e

REPO="dlcuy22/StellePlayer"
INSTALL_DIR="/opt/StellePlayer"
BIN_NAME="StellePlayer"

echo "Installing StellePlayer..."

# 1. Detect Architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64)  TARGET_ARCH="amd64" ;;
    aarch64) TARGET_ARCH="arm64" ;;
    armv7l)  TARGET_ARCH="arm" ;;
    *)       echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# 2. Get latest version from GitHub API
VERSION=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo "Error: Could not retrieve latest version."
    exit 1
fi

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/StellePlayer-linux-$TARGET_ARCH"

echo "Detected: Linux $TARGET_ARCH"
echo "Version: $VERSION"

# 3. Create Install Directory
sudo mkdir -p $INSTALL_DIR

# 4. Download Binary
echo "Downloading $DOWNLOAD_URL..."
sudo curl -L $DOWNLOAD_URL -o $INSTALL_DIR/$BIN_NAME
sudo chmod +x $INSTALL_DIR/$BIN_NAME

# 5. Add to PATH via Symlink
echo "Creating symlinks in /usr/local/bin..."
sudo ln -sf $INSTALL_DIR/$BIN_NAME /usr/local/bin/$BIN_NAME
sudo ln -sf $INSTALL_DIR/$BIN_NAME /usr/local/bin/Splayer
sudo ln -sf $INSTALL_DIR/$BIN_NAME /usr/local/bin/Splay

echo "Successfully installed StellePlayer to $INSTALL_DIR/$BIN_NAME"
echo "You can now run 'StellePlayer', 'Splayer', or 'Splay' from your terminal."