#!/bin/bash
set -e

echo "Installing Security Platform Desktop Agent..."

# Detect OS
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    if command -v systemctl &> /dev/null; then
        echo "Installing as systemd service..."
        sudo cp security-platform-agent /usr/local/bin/
        sudo chmod +x /usr/local/bin/security-platform-agent
        sudo mkdir -p /etc/security-platform
        sudo cp config.yaml /etc/security-platform/ 2>/dev/null || true
        sudo tee /etc/systemd/system/security-platform-agent.service > /dev/null <<EOF
[Unit]
Description=Security Platform Desktop Agent
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/security-platform-agent -service
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
        sudo systemctl daemon-reload
        sudo systemctl enable security-platform-agent
        sudo systemctl start security-platform-agent
        echo "Service installed and started"
    else
        echo "Installing as regular binary..."
        sudo cp security-platform-agent /usr/local/bin/
        sudo chmod +x /usr/local/bin/security-platform-agent
    fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    echo "Installing as launchd daemon..."
    sudo cp security-platform-agent /usr/local/bin/
    sudo chmod +x /usr/local/bin/security-platform-agent
    sudo mkdir -p /etc/security-platform
    sudo cp config.yaml /etc/security-platform/ 2>/dev/null || true
    sudo tee /Library/LaunchDaemons/com.securityplatform.agent.plist > /dev/null <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.securityplatform.agent</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/security-platform-agent</string>
        <string>-service</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
EOF
    sudo launchctl load /Library/LaunchDaemons/com.securityplatform.agent.plist
    echo "Daemon installed and started"
else
    echo "Unsupported OS: $OSTYPE"
    exit 1
fi

echo "Installation complete!"
