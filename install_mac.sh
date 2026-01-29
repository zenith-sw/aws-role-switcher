#!/bin/bash
set -e

OS=$(uname -s)
ARCH=$(uname -m)

if [[ $OS != "Darwin" ]]; then
  echo "Only macOS is supported."
  exit 1
fi


BINARY_NAME="sw_darwin_amd64"
if [[ $ARCH == "arm64" ]]; then
  BINARY_NAME="sw_darwin_arm64"
fi

echo "Downloading $BINARY_NAME..."

echo "Downloading $BINARY_NAME from GitHub..."
curl -fsSL "https://github.com/zenith-sw/util-aws-role-switcher/releases/latest/download/$BINARY_NAME?raw=true" -o sw

echo "Download complete."

chmod +x sw
install sw /usr/local/bin/

# delete temporary file
rm sw

echo "Starting to initiate configuration..."
sw init

echo "-----"
echo "Installation Complete!"
echo "---"
echo "Next Step:"
echo "1. Open your config file: vi ~/.sw/config.yaml"
echo "2. Register your IAM Role ARNs under 'assume_roles'."
echo "3. Run 'sw setup {profile}' to get your credentials!"
echo "---"
echo "Example: sw setup dev"