#!/bin/bash
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

set -e

OS=$(uname -s)
ARCH=$(uname -m)

if [[ $OS != "Darwin" ]]; then
  echo "\n${YELLOW}Only macOS is supported.${NC}"
  exit 1
fi


BINARY_NAME="sw_darwin_amd64"
if [[ $ARCH == "arm64" ]]; then
  BINARY_NAME="sw_darwin_arm64"
fi

echo -e "${CYAN}Downloading $BINARY_NAME...${NC}\n"
curl -fsSL "https://github.com/zenith-sw/aws-role-switcher/releases/latest/download/$BINARY_NAME?raw=true" -o sw

echo -e "\n${GREEN}Download complete.${NC}\n"

chmod +x sw
install sw /usr/local/bin/

# delete temporary file
rm sw

echo "Starting to initiate configuration..."
sw init

echo -e "\n${CYAN}-----${NC}"
echo -e "${GREEN}Installation Complete!${NC}"
echo "Register your first role using 'sw add'"
echo -e "${CYAN}-----${NC}"