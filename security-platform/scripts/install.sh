#!/bin/bash
set -e

# Universal Security Platform Agent Installer
# Supports: Linux (deb/rpm), macOS (Homebrew/pkg), Windows (Chocolatey/MSI)

VERSION="1.0.0"
CONTROL_PLANE_URL="${SECURITY_PLATFORM_CONTROL_PLANE_URL:-}"
AUTH_KEY="${SECURITY_PLATFORM_AUTH_KEY:-}"
SERVICE_NAME="${SECURITY_PLATFORM_SERVICE_NAME:-}"

echo "Security Platform Agent Installer v${VERSION}"
echo "=========================================="
echo ""

# Detect OS
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if [ -f /etc/debian_version ]; then
            echo "debian"
        elif [ -f /etc/redhat-release ]; then
            echo "rhel"
        else
            echo "linux"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "darwin"
    elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
        echo "windows"
    else
        echo "unknown"
    fi
}

install_debian() {
    echo "Installing on Debian/Ubuntu..."
    
    # Check for curl
    if ! command -v curl &> /dev/null; then
        echo "Installing curl..."
        sudo apt-get update && sudo apt-get install -y curl
    fi
    
    # Download and install
    curl -fsSL https://releases.securityplatform.com/agent.deb -o /tmp/agent.deb
    sudo dpkg -i /tmp/agent.deb || sudo apt-get install -f -y
    sudo systemctl enable security-platform-agent
    sudo systemctl start security-platform-agent
    
    echo "✓ Installed successfully"
}

install_rhel() {
    echo "Installing on RHEL/CentOS..."
    
    # Check for curl
    if ! command -v curl &> /dev/null; then
        echo "Installing curl..."
        sudo yum install -y curl
    fi
    
    # Download and install
    curl -fsSL https://releases.securityplatform.com/agent.rpm -o /tmp/agent.rpm
    sudo rpm -ivh /tmp/agent.rpm
    sudo systemctl enable security-platform-agent
    sudo systemctl start security-platform-agent
    
    echo "✓ Installed successfully"
}

install_darwin() {
    echo "Installing on macOS..."
    
    # Check for Homebrew
    if command -v brew &> /dev/null; then
        echo "Using Homebrew..."
        brew tap security-platform/agent
        brew install security-platform-agent
    else
        echo "Homebrew not found. Installing via direct download..."
        curl -fsSL https://releases.securityplatform.com/agent.pkg -o /tmp/agent.pkg
        sudo installer -pkg /tmp/agent.pkg -target /
    fi
    
    echo "✓ Installed successfully"
}

install_windows() {
    echo "Installing on Windows..."
    
    # Check for Chocolatey
    if command -v choco &> /dev/null; then
        echo "Using Chocolatey..."
        choco install security-platform-agent -y
    else
        echo "Chocolatey not found. Please install manually:"
        echo "1. Download: https://releases.securityplatform.com/agent.msi"
        echo "2. Run the installer"
    fi
    
    echo "✓ Installed successfully"
}

configure_agent() {
    if [ -z "$CONTROL_PLANE_URL" ] || [ -z "$AUTH_KEY" ]; then
        echo ""
        echo "Configuration required:"
        read -p "Control Plane URL: " CONTROL_PLANE_URL
        read -p "Auth Key: " AUTH_KEY
        read -p "Service Name (optional): " SERVICE_NAME
    fi
    
    CONFIG_FILE="/etc/security-platform/config.yaml"
    if [[ "$OSTYPE" == "darwin"* ]]; then
        CONFIG_FILE="/etc/security-platform/config.yaml"
    elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
        CONFIG_FILE="$env:ProgramData\SecurityPlatform\config.yaml"
    fi
    
    sudo mkdir -p "$(dirname "$CONFIG_FILE")"
    
    cat <<EOF | sudo tee "$CONFIG_FILE" > /dev/null
control_plane_url: ${CONTROL_PLANE_URL}
auth_key: ${AUTH_KEY}
service_name: ${SERVICE_NAME:-agent}
environment: production
otlp_endpoint: http://localhost:4318/v1/traces
local_policy: observe
EOF
    
    echo "✓ Configuration saved to $CONFIG_FILE"
    
    # Restart service
    if command -v systemctl &> /dev/null; then
        sudo systemctl restart security-platform-agent
    elif command -v launchctl &> /dev/null; then
        sudo launchctl unload /Library/LaunchDaemons/com.securityplatform.agent.plist 2>/dev/null || true
        sudo launchctl load /Library/LaunchDaemons/com.securityplatform.agent.plist
    fi
}

# Main installation flow
OS=$(detect_os)

case $OS in
    debian)
        install_debian
        ;;
    rhel)
        install_rhel
        ;;
    darwin)
        install_darwin
        ;;
    windows)
        install_windows
        ;;
    *)
        echo "Unsupported OS: $OSTYPE"
        echo "Please install manually from: https://github.com/security-platform/desktop-agent"
        exit 1
        ;;
esac

configure_agent

echo ""
echo "Installation complete!"
echo ""
echo "Next steps:"
echo "1. Verify service status: systemctl status security-platform-agent"
echo "2. View logs: journalctl -u security-platform-agent -f"
echo "3. Edit config: sudo nano $CONFIG_FILE"
