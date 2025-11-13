#!/bin/bash
# One-liner installer - can be piped directly from curl
# Usage: curl -fsSL https://install.securityplatform.com | bash

set -e

echo "Installing Security Platform Agent..."

# Download and run main installer
INSTALLER_URL="${INSTALLER_URL:-https://raw.githubusercontent.com/security-platform/agents/main/scripts/install.sh}"

curl -fsSL "$INSTALLER_URL" | bash
