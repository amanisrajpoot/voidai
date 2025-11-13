#!/bin/bash
# Universal installation script for Security Platform Agent
set -e

VERSION="${VERSION:-latest}"
PLATFORM="${PLATFORM:-auto}"

# Detect platform if auto
if [ "$PLATFORM" = "auto" ]; then
    case "$(uname -s)" in
        Linux*)     PLATFORM="linux";;
        Darwin*)    PLATFORM="darwin";;
        CYGWIN*)    PLATFORM="windows";;
        MINGW*)     PLATFORM="windows";;
        *)          echo "Unsupported platform"; exit 1;;
    esac
fi

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64) ARCH="amd64";;
    arm64|aarch64) ARCH="arm64";;
    *) echo "Unsupported architecture: $ARCH"; exit 1;;
esac

echo "Installing Security Platform Agent"
echo "Platform: $PLATFORM"
echo "Architecture: $ARCH"
echo "Version: $VERSION"

# Download and install based on platform
case "$PLATFORM" in
    linux)
        if command -v apt-get &> /dev/null; then
            # Debian/Ubuntu
            curl -fsSL https://packages.securityplatform.com/install.sh | sudo bash
            sudo apt-get update
            sudo apt-get install -y security-platform-agent
        elif command -v yum &> /dev/null; then
            # RHEL/CentOS
            sudo yum install -y https://packages.securityplatform.com/rpm/security-platform-agent.rpm
        elif command -v dnf &> /dev/null; then
            # Fedora
            sudo dnf install -y https://packages.securityplatform.com/rpm/security-platform-agent.rpm
        else
            echo "Unsupported Linux distribution"
            exit 1
        fi
        ;;
    darwin)
        if command -v brew &> /dev/null; then
            brew install security-platform-agent
        else
            echo "Homebrew is required for macOS installation"
            echo "Install Homebrew: https://brew.sh"
            exit 1
        fi
        ;;
    windows)
        if command -v choco &> /dev/null; then
            choco install security-platform-agent -y
        else
            echo "Chocolatey is required for Windows installation"
            echo "Install Chocolatey: https://chocolatey.org"
            exit 1
        fi
        ;;
esac

echo "Installation complete!"
echo ""
echo "Next steps:"
echo "1. Configure the agent: /etc/security-platform/agent-config.yaml"
echo "2. Set your auth key: export SECURITY_PLATFORM_AUTH_KEY=your-key"
echo "3. Start the service:"
case "$PLATFORM" in
    linux)
        echo "   sudo systemctl start security-platform-agent"
        echo "   sudo systemctl enable security-platform-agent"
        ;;
    darwin)
        echo "   brew services start security-platform-agent"
        ;;
    windows)
        echo "   net start SecurityPlatformAgent"
        ;;
esac
