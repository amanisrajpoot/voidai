# Security Platform Desktop Agents

Desktop agents for Windows, macOS, and Linux that run as system services/daemons.

## Installation

### Linux (systemd)

```bash
# Install binary
sudo cp security-platform-agent /usr/local/bin/
sudo chmod +x /usr/local/bin/security-platform-agent

# Install systemd service
sudo cp systemd/security-platform-agent.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable security-platform-agent
sudo systemctl start security-platform-agent
```

### macOS (LaunchDaemon)

```bash
# Install binary
sudo cp security-platform-agent /usr/local/bin/
sudo chmod +x /usr/local/bin/security-platform-agent

# Install LaunchDaemon
sudo cp macos/com.securityplatform.agent.plist /Library/LaunchDaemons/
sudo launchctl load /Library/LaunchDaemons/com.securityplatform.agent.plist
```

### Windows (Service)

```powershell
# Run as Administrator
.\install-service.ps1

# Or manually:
# sc.exe create SecurityPlatformAgent binPath="C:\Program Files\SecurityPlatform\security-platform-agent.exe -config=C:\ProgramData\SecurityPlatform\config.yaml -service" start=auto
# sc.exe start SecurityPlatformAgent
```

## Configuration

Create configuration file at:
- Linux: `/etc/security-platform/config.yaml`
- macOS: `/etc/security-platform/config.yaml`
- Windows: `C:\ProgramData\SecurityPlatform\config.yaml`

```yaml
control_plane_url: "https://api.securityplatform.com"
auth_key: "${SECURITY_PLATFORM_AUTH_KEY}"
service_name: "desktop-agent"
environment: "production"
```

## Features

- ✅ Cross-platform support (Windows, macOS, Linux)
- ✅ System service/daemon integration
- ✅ Automatic startup on boot
- ✅ Graceful shutdown handling
- ✅ Logging to system logs
