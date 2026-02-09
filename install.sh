#!/bin/bash
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

set -e

OS_TYPE=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH_RAW=$(uname -m)

case "$ARCH_RAW" in
    x86_64|amd64)  ARCH_TYPE="amd64" ;;
    aarch64|arm64) ARCH_TYPE="arm64" ;;
    *) 
        echo -e "${RED}Unsupported architecture: $ARCH_RAW${NC}"
        exit 1 
        ;;
esac

BINARY_NAME="sw_${OS_TYPE}_${ARCH_TYPE}"
DOWNLOAD_URL="https://github.com/zenith-sw/aws-role-switcher/releases/latest/download/${BINARY_NAME}"

echo -e "${CYAN}System detected: $OS_TYPE ($ARCH_TYPE)${NC}"

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo -e "${CYAN}Downloading $BINARY_NAME...${NC}"

if ! curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/sw_temp"; then
    echo -e "${RED}Download failed!${NC}"
    echo "URL: $DOWNLOAD_URL"
    exit 1
fi

echo -e "${GREEN}Download complete.${NC}"

chmod +x "$TMP_DIR/sw_temp"
INSTALL_DEST="/usr/local/bin/sw"
echo "Installing to $INSTALL_DEST..."

if [ -w "/usr/local/bin" ]; then
    mv "$TMP_DIR/sw_temp" "$INSTALL_DEST"
else
    echo -e "${YELLOW}Requesting sudo permission to install...${NC}"
    sudo mv "$TMP_DIR/sw_temp" "$INSTALL_DEST"
fi

echo "Initializing configuration..."
"$INSTALL_DEST" init

echo -e "\n${CYAN}-----${NC}"
echo -e "${GREEN}Installation Complete!${NC}"
echo "Register your first role using 'sw add'"
echo -e "${CYAN}-----${NC}"