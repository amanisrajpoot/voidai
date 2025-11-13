#!/bin/bash
# Universal Security Platform Installation Script
# Supports: Linux (DEB/RPM), macOS, Windows (via WSL)

set -e

VERSION="1.0.0"
CONTROL_PLANE_URL="${SECURITY_PLATFORM_CONTROL_PLANE_URL:-https://api.securityplatform.com}"
AUTH_KEY="${SECURITY_PLATFORM_AUTH_KEY}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if [ -f /etc/debian_version ]; then
            OS="debian"
        elif [ -f /etc/redhat-release ]; then
            OS="rhel"
        elif [ -f /etc/arch-release ]; then
            OS="arch"
        else
            OS="linux"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="darwin"
    else
        print_error "Unsupported OS: $OSTYPE"
        exit 1
    fi
    print_info "Detected OS: $OS"
}

install_debian() {
    print_info "Installing on Debian/Ubuntu..."
    
    # Check for curl/wget
    if ! command -v curl &> /dev/null && ! command -v wget &> /dev/null; then
        print_info "Installing curl..."
        sudo apt-get update && sudo apt-get install -y curl
    fi
    
    # Download DEB package
    DOWNLOAD_URL="https://releases.securityplatform.com/security-platform-agent_${VERSION}_amd64.deb"
    TEMP_FILE=$(mktemp)
    
    if command -v curl &> /dev/null; then
        curl -L -o "$TEMP_FILE" "$DOWNLOAD_URL"
    else
        wget -O "$TEMP_FILE" "$DOWNLOAD_URL"
    fi
    
    # Install
    sudo dpkg -i "$TEMP_FILE" || sudo apt-get install -f -y
    rm "$TEMP_FILE"
}

install_rhel() {
    print_info "Installing on RHEL/CentOS/Fedora..."
    
    # Check for curl/wget
    if ! command -v curl &> /dev/null && ! command -v wget &> /dev/null; then
        print_info "Installing curl..."
        sudo yum install -y curl || sudo dnf install -y curl
    fi
    
    # Download RPM package
    DOWNLOAD_URL="https://releases.securityplatform.com/security-platform-agent-${VERSION}-1.x86_64.rpm"
    TEMP_FILE=$(mktemp)
    
    if command -v curl &> /dev/null; then
        curl -L -o "$TEMP_FILE" "$DOWNLOAD_URL"
    else
        wget -O "$TEMP_FILE" "$DOWNLOAD_URL"
    fi
    
    # Install
    sudo rpm -ivh "$TEMP_FILE" || sudo yum install -y "$TEMP_FILE" || sudo dnf install -y "$TEMP_FILE"
    rm "$TEMP_FILE"
}

install_darwin() {
    print_info "Installing on macOS..."
    
    # Check for Homebrew
    if command -v brew &> /dev/null; then
        print_info "Installing via Homebrew..."
        brew tap security-platform/agent
        brew install security-platform-agent
        return
    fi
    
    # Fallback to direct download
    DOWNLOAD_URL="https://releases.securityplatform.com/security-platform-agent-${VERSION}.pkg"
    TEMP_FILE=$(mktemp).pkg
    
    curl -L -o "$TEMP_FILE" "$DOWNLOAD_URL"
    sudo installer -pkg "$TEMP_FILE" -target /
    rm "$TEMP_FILE"
}

configure_agent() {
    CONFIG_DIR="/etc/security-platform"
    CONFIG_FILE="$CONFIG_DIR/agent-config.yaml"
    
    if [ ! -d "$CONFIG_DIR" ]; then
        sudo mkdir -p "$CONFIG_DIR"
    fi
    
    if [ ! -f "$CONFIG_FILE" ]; then
        print_info "Creating configuration file..."
        sudo tee "$CONFIG_FILE" > /dev/null <<EOF
control_plane_url: "$CONTROL_PLANE_URL"
auth_key: "${AUTH_KEY:-}"
service_name: "desktop-agent"
environment: "production"
EOF
    else
        print_warn "Configuration file already exists at $CONFIG_FILE"
    fi
    
    if [ -n "$AUTH_KEY" ]; then
        print_info "Updating auth key in configuration..."
        sudo sed -i "s/auth_key:.*/auth_key: \"$AUTH_KEY\"/" "$CONFIG_FILE"
    fi
}

start_service() {
    if [[ "$OS" == "darwin" ]]; then
        print_info "Starting service via launchd..."
        sudo launchctl load /Library/LaunchDaemons/com.securityplatform.agent.plist 2>/dev/null || true
    else
        print_info "Enabling and starting service..."
        sudo systemctl enable security-platform-agent || true
        sudo systemctl start security-platform-agent || true
        sudo systemctl status security-platform-agent --no-pager || true
    fi
}

main() {
    echo "=========================================="
    echo "Security Platform Agent Installer"
    echo "=========================================="
    echo ""
    
    detect_os
    
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
        *)
            print_error "Automatic installation not supported for $OS"
            print_info "Please visit https://docs.securityplatform.com/installation"
            exit 1
            ;;
    esac
    
    configure_agent
    start_service
    
    echo ""
    echo "=========================================="
    echo "Installation complete!"
    echo "=========================================="
    echo "Configuration: /etc/security-platform/agent-config.yaml"
    echo "Logs: /var/log/security-platform-agent.log"
    echo ""
    print_info "Next steps:"
    echo "  1. Edit configuration: sudo nano /etc/security-platform/agent-config.yaml"
    echo "  2. Set your auth key"
    echo "  3. Restart service: sudo systemctl restart security-platform-agent"
    echo ""
}

main "$@"
