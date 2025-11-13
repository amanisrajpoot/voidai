# Security Platform Desktop Agent

Cross-platform desktop agent for Windows, macOS, and Linux.

## Quick Start

### Download Binary

```bash
# Linux
wget https://releases.securityplatform.com/desktop-agent-linux-amd64
chmod +x desktop-agent-linux-amd64
sudo mv desktop-agent-linux-amd64 /usr/local/bin/security-platform-agent

# macOS
curl -L https://releases.securityplatform.com/desktop-agent-darwin-amd64 -o /usr/local/bin/security-platform-agent
chmod +x /usr/local/bin/security-platform-agent

# Windows
# Download from releases page and run installer
```

### Run in Foreground

```bash
export SECURITY_PLATFORM_CONTROL_PLANE_URL=https://api.securityplatform.com
export SECURITY_PLATFORM_AUTH_KEY=your-token-here
security-platform-agent
```

### Install as Service/Daemon

#### Linux (systemd)

```bash
# Copy systemd service file
sudo cp linux/systemd.service /etc/systemd/system/security-platform-agent.service

# Edit configuration
sudo nano /etc/systemd/system/security-platform-agent.service

# Enable and start
sudo systemctl enable security-platform-agent
sudo systemctl start security-platform-agent
```

#### macOS (launchd)

```bash
# Copy plist
sudo cp macos/launchd.plist /Library/LaunchDaemons/com.securityplatform.agent.plist

# Load service
sudo launchctl load /Library/LaunchDaemons/com.securityplatform.agent.plist
```

#### Windows

```powershell
# Install service (requires admin)
security-platform-agent.exe -install

# Start service
net start SecurityPlatformAgent
```

## Configuration

Create `/etc/security-platform/agent.yaml`:

```yaml
control_plane_url: "https://api.securityplatform.com"
auth_key: "your-token-here"
service_name: "desktop-agent"
environment: "production"
```

## Features

- ✅ Cross-platform support (Windows, macOS, Linux)
- ✅ Automatic service/daemon installation
- ✅ Health monitoring and auto-restart
- ✅ Simple configuration
- ✅ Zero manual setup required

## Documentation

See [docs/DEPLOYMENT.md](../../docs/DEPLOYMENT.md) for full documentation.
