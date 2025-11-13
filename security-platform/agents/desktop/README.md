# Desktop Agents

Desktop agents for Windows, macOS, and Linux that run as system services/daemons.

## Windows

### Installation

```powershell
# Install as Windows Service
security-platform-agent.exe -install

# Start service
net start SecurityPlatformAgent
```

### Configuration

Set environment variables:
- `SECURITY_PLATFORM_CONTROL_PLANE_URL`
- `SECURITY_PLATFORM_AUTH_KEY`

## macOS

### Installation

```bash
# Copy binary
sudo cp security-platform-agent /usr/local/bin/

# Install launchd plist
sudo cp com.securityplatform.agent.plist /Library/LaunchDaemons/

# Load service
sudo launchctl load /Library/LaunchDaemons/com.securityplatform.agent.plist
```

### Configuration

Set environment variables in plist or use:
```bash
launchctl setenv SECURITY_PLATFORM_CONTROL_PLANE_URL "https://api.securityplatform.com"
launchctl setenv SECURITY_PLATFORM_AUTH_KEY "your-key"
```

## Linux

### Installation

```bash
# Copy binary
sudo cp security-platform-agent /usr/local/bin/

# Install systemd service
sudo cp security-platform-agent.service /etc/systemd/system/

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable security-platform-agent
sudo systemctl start security-platform-agent
```

### Configuration

Edit `/etc/systemd/system/security-platform-agent.service` or use:
```bash
sudo systemctl edit security-platform-agent
```
