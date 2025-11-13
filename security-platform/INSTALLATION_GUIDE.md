# Security Platform - Complete Installation Guide

This guide covers all installation methods for Security Platform agents and SDKs.

## Quick Start - One-Liner Installation

### Linux/macOS

```bash
curl -fsSL https://install.securityplatform.com | bash
```

### Windows (PowerShell)

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://install.securityplatform.com/install.ps1'))
```

## Language Agents

### Node.js

```bash
npm install @security-platform/node-agent
```

```javascript
const { Agent } = require('@security-platform/node-agent');

const agent = new Agent({
  controlPlaneUrl: 'https://api.securityplatform.com',
  authKey: process.env.SECURITY_PLATFORM_AUTH_KEY,
  serviceName: 'my-service'
});

agent.start();
```

### Python

```bash
pip install security-platform-agent
```

```python
from security_platform_agent import Agent

agent = Agent({
    'control_plane_url': 'https://api.securityplatform.com',
    'auth_key': os.getenv('SECURITY_PLATFORM_AUTH_KEY'),
    'service_name': 'my-service'
})

agent.start()
```

### Java

Add to `pom.xml`:

```xml
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>security-platform-agent</artifactId>
    <version>1.0.0</version>
</dependency>
```

Or Gradle:

```groovy
implementation 'com.securityplatform:security-platform-agent:1.0.0'
```

```java
import com.securityplatform.agent.Agent;
import com.securityplatform.agent.AgentConfig;

AgentConfig config = new AgentConfig();
config.setControlPlaneUrl("https://api.securityplatform.com");
config.setAuthKey(System.getenv("SECURITY_PLATFORM_AUTH_KEY"));
config.setServiceName("my-service");

Agent agent = new Agent(config);
agent.start();
```

### .NET

```bash
dotnet add package SecurityPlatform.Agent
```

```csharp
using SecurityPlatform.Agent;

var config = new AgentConfig
{
    ControlPlaneUrl = "https://api.securityplatform.com",
    AuthKey = Environment.GetEnvironmentVariable("SECURITY_PLATFORM_AUTH_KEY"),
    ServiceName = "my-service"
};

var agent = new Agent(config);
agent.Start();
```

### Go

```bash
go get github.com/security-platform/go-agent
```

```go
import "github.com/security-platform/go-agent"

config := &agent.Config{
    ControlPlaneURL: "https://api.securityplatform.com",
    AuthKey:         os.Getenv("SECURITY_PLATFORM_AUTH_KEY"),
    ServiceName:     "my-service",
}

ag, _ := agent.NewAgent(config)
ag.Start()
```

## Mobile SDKs

### iOS

#### CocoaPods

```ruby
pod 'SecurityPlatform', '~> 1.0'
```

#### Swift Package Manager

Add to `Package.swift`:

```swift
dependencies: [
    .package(url: "https://github.com/security-platform/agents", from: "1.0.0")
]
```

```swift
import SecurityPlatform

let config = AgentConfig(
    controlPlaneURL: "https://api.securityplatform.com",
    authKey: ProcessInfo.processInfo.environment["SECURITY_PLATFORM_AUTH_KEY"],
    serviceName: "my-ios-app"
)

let agent = SecurityPlatformAgent(config: config)
try? agent.start()
```

### Android

Add to `build.gradle`:

```groovy
dependencies {
    implementation 'com.securityplatform:agent:1.0.0'
}
```

```kotlin
import com.securityplatform.agent.Agent
import com.securityplatform.agent.AgentConfig

val config = AgentConfig(
    controlPlaneUrl = "https://api.securityplatform.com",
    authKey = System.getenv("SECURITY_PLATFORM_AUTH_KEY"),
    serviceName = "my-android-app"
)

val agent = Agent(config)
agent.start()
```

## Desktop Agents

### Linux (systemd)

```bash
# Install package
sudo dpkg -i security-platform-agent.deb
# or
sudo rpm -ivh security-platform-agent.rpm

# Configure
sudo nano /etc/security-platform/agent-config.yaml

# Start service
sudo systemctl start security-platform-agent
sudo systemctl enable security-platform-agent
```

### macOS

#### Homebrew

```bash
brew install security-platform-agent
brew services start security-platform-agent
```

#### Manual Installation

```bash
sudo cp security-platform-agent /usr/local/bin/
sudo cp com.securityplatform.agent.plist /Library/LaunchDaemons/
sudo launchctl load /Library/LaunchDaemons/com.securityplatform.agent.plist
```

### Windows

#### Chocolatey

```powershell
choco install security-platform-agent -y
net start SecurityPlatformAgent
```

#### Manual Installation

```powershell
# Install as service
security-platform-agent.exe -install

# Start service
net start SecurityPlatformAgent
```

## Package Managers

### Homebrew (macOS/Linux)

```bash
brew tap security-platform/agents
brew install security-platform-agent
```

### Chocolatey (Windows)

```powershell
choco install security-platform-agent -y
```

### APT (Debian/Ubuntu)

```bash
curl -fsSL https://packages.securityplatform.com/install.sh | sudo bash
sudo apt-get update
sudo apt-get install security-platform-agent
```

### YUM/DNF (RHEL/CentOS/Fedora)

```bash
sudo yum install https://packages.securityplatform.com/rpm/security-platform-agent.rpm
# or
sudo dnf install https://packages.securityplatform.com/rpm/security-platform-agent.rpm
```

## Kubernetes

### Using Helm

```bash
helm repo add security-platform https://charts.securityplatform.com
helm install my-agent security-platform/security-platform-agent \
  --set controlPlaneUrl=https://api.securityplatform.com \
  --set authKeySecretRef.name=security-platform-auth \
  --set authKeySecretRef.key=auth-key
```

### Using Kubernetes Operator

1. Install CRDs:

```bash
kubectl apply -f https://raw.githubusercontent.com/security-platform/agents/main/packaging/kubernetes-operator/config/crd/bases/securityplatform.io_agents.yaml
```

2. Install Operator:

```bash
kubectl apply -f https://raw.githubusercontent.com/security-platform/agents/main/packaging/kubernetes-operator/config/manager/manager.yaml
```

3. Create Secret:

```bash
kubectl create secret generic security-platform-auth \
  --from-literal=auth-key=your-auth-key
```

4. Create Agent resource:

```yaml
apiVersion: securityplatform.io/v1
kind: Agent
metadata:
  name: my-agent
spec:
  controlPlaneUrl: "https://api.securityplatform.com"
  authKeySecretRef:
    name: security-platform-auth
    key: auth-key
  serviceName: "my-service"
  replicas: 2
```

```bash
kubectl apply -f agent.yaml
```

## Air-Gapped Installation

1. Download airgap bundle:

```bash
wget https://releases.securityplatform.com/airgap/security-platform-agent-1.0.0-airgap.tar.gz
```

2. Transfer to air-gapped system

3. Extract and install:

```bash
tar -xzf security-platform-agent-1.0.0-airgap.tar.gz
cd security-platform-agent-airgap
chmod +x install.sh
sudo ./install.sh
```

For Windows:

```powershell
# Extract ZIP file
Expand-Archive security-platform-agent-1.0.0-airgap.zip
cd security-platform-agent-airgap
.\install.bat
```

## Configuration

All agents use the same configuration format. Configuration file location:

- **Linux/macOS**: `/etc/security-platform/agent-config.yaml`
- **Windows**: `C:\ProgramData\SecurityPlatform\agent-config.yaml`
- **Kubernetes**: ConfigMap or environment variables

Example configuration:

```yaml
control_plane_url: "https://api.securityplatform.com"
auth_key: ""  # Or use SECURITY_PLATFORM_AUTH_KEY environment variable
service_name: "my-service"
environment: "production"

telemetry:
  otlp_endpoint: "http://localhost:4318"
  batch_size: 100
  batch_timeout: "5s"
  max_queue_size: 2048

policy:
  mode: "observe"  # or "block"
  auto_enable_blocking: false
  observe_period_hours: 48

security:
  mtls_enabled: true
```

## Verification

### Check Agent Status

**Linux:**
```bash
systemctl status security-platform-agent
```

**macOS:**
```bash
launchctl list | grep security-platform
```

**Windows:**
```powershell
sc query SecurityPlatformAgent
```

### Test Connection

```bash
security-platform-agent --version
security-platform-agent --test-connection
```

## Troubleshooting

### Agent Not Starting

1. Check logs:
   - Linux: `journalctl -u security-platform-agent`
   - macOS: `tail -f /var/log/security-platform-agent.log`
   - Windows: Event Viewer → Applications → SecurityPlatformAgent

2. Verify configuration:
```bash
security-platform-agent --validate-config
```

3. Check network connectivity:
```bash
curl https://api.securityplatform.com/health
```

### Authentication Issues

Ensure `SECURITY_PLATFORM_AUTH_KEY` is set:

```bash
export SECURITY_PLATFORM_AUTH_KEY=your-key-here
```

Or configure in `agent-config.yaml`:

```yaml
auth_key: "your-key-here"
```

## Next Steps

- [Architecture Documentation](docs/ARCHITECTURE.md)
- [API Reference](docs/API.md)
- [Deployment Guide](docs/DEPLOYMENT.md)
- [Examples](examples/)
