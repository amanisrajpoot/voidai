# Security Platform - Installation Guide

Complete installation guide for all agents, SDKs, and components.

## Quick Start

### One-Line Installation

**macOS (Homebrew):**
```bash
brew tap securityplatform/agent && brew install security-platform-agent
```

**Windows (Chocolatey):**
```powershell
choco install security-platform-agent
```

**Linux (DEB):**
```bash
curl -fsSL https://install.securityplatform.com/linux | sudo bash
```

**Linux (RPM):**
```bash
curl -fsSL https://install.securityplatform.com/linux-rpm | sudo bash
```

## Language Agents

### Java

**Maven:**
```xml
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>security-platform-agent</artifactId>
    <version>1.0.0</version>
</dependency>
```

**Gradle:**
```gradle
implementation 'com.securityplatform:security-platform-agent:1.0.0'
```

### .NET

**NuGet Package Manager:**
```powershell
Install-Package SecurityPlatform.Agent
```

**PackageReference:**
```xml
<PackageReference Include="SecurityPlatform.Agent" Version="1.0.0" />
```

### Go

```bash
go get github.com/securityplatform/go-agent
```

### Node.js

```bash
npm install @security-platform/node-agent
```

### Python

```bash
pip install security-platform-agent
```

## Mobile SDKs

### iOS

**Swift Package Manager:**
Add to `Package.swift`:
```swift
dependencies: [
    .package(url: "https://github.com/securityplatform/ios-sdk", from: "1.0.0")
]
```

**CocoaPods:**
```ruby
pod 'SecurityPlatform', '~> 1.0'
```

### Android

**Gradle:**
```gradle
implementation 'com.securityplatform:security-platform-agent:1.0.0'
```

**Maven:**
```xml
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>security-platform-agent</artifactId>
    <version>1.0.0</version>
</dependency>
```

## Desktop Agents

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
```

## Kubernetes Operator

### Using Helm

```bash
helm install security-platform-operator ./helm/security-platform-operator
```

### Using kubectl

```bash
# Install CRDs
kubectl apply -f deploy/crds/

# Install operator
kubectl apply -f deploy/operator.yaml
```

## Airgap Installation

For air-gapped environments:

```bash
# Generate airgap bundle
./builder airgap --version 1.0.0 --output security-platform-airgap-1.0.0.tar.gz

# Transfer bundle to air-gapped system
# Extract and run install script
tar -xzf security-platform-airgap-1.0.0.tar.gz
cd security-platform-airgap-1.0.0
./install.sh
```

## Configuration

All agents use a common configuration format. Create `config.yaml`:

```yaml
control_plane_url: "https://api.securityplatform.com"
auth_key: "${SECURITY_PLATFORM_AUTH_KEY}"
service_name: "my-service"
environment: "production"
telemetry:
  otlp_endpoint: "http://localhost:4318"
  batch_size: 100
  batch_timeout: "5s"
policy:
  mode: "observe"  # or "block"
```

## Environment Variables

- `SECURITY_PLATFORM_URL` - Control plane URL
- `SECURITY_PLATFORM_AUTH_KEY` - Authentication key
- `SERVICE_NAME` - Service name
- `ENVIRONMENT` - Environment (production, staging, development)

## Verification

After installation, verify the agent is running:

```bash
# Check status (Linux/macOS)
systemctl status security-platform-agent
# or
launchctl list | grep security-platform

# Check status (Windows)
Get-Service -Name SecurityPlatformAgent
```

## Troubleshooting

### Agent not starting

1. Check configuration file syntax
2. Verify authentication key is set
3. Check network connectivity to control plane
4. Review logs:
   - Linux: `journalctl -u security-platform-agent`
   - macOS: `tail -f /var/log/security-platform-agent.log`
   - Windows: Event Viewer → Applications

### Connection issues

1. Verify firewall rules allow outbound HTTPS
2. Check proxy settings if behind corporate proxy
3. Verify control plane URL is correct
4. Test connectivity: `curl https://api.securityplatform.com/health`

## Next Steps

After installation:
1. Configure your application to use the agent
2. Set up authentication keys
3. Configure policies and redaction rules
4. Monitor telemetry in the control plane

For detailed integration guides, see:
- [Java Integration Guide](agents/java/README.md)
- [.NET Integration Guide](agents/dotnet/README.md)
- [Go Integration Guide](agents/go/README.md)
- [Node.js Integration Guide](agents/node/README.md)
- [Python Integration Guide](agents/python/README.md)
