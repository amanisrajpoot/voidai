# Security Platform Agent - Installation Guide

Complete installation guide for all platforms and deployment methods.

## Quick Install (One-Liner)

### Linux (Debian/Ubuntu)
```bash
curl -fsSL https://install.securityplatform.com | bash
```

### Linux (RHEL/CentOS)
```bash
curl -fsSL https://install.securityplatform.com | bash
```

### macOS
```bash
brew install security-platform/agent/security-platform-agent
```

### Windows (PowerShell)
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://install.securityplatform.com/install.ps1'))
```

## Platform-Specific Installation

### Frontend (Web) SDK

**NPM:**
```bash
npm install @security-platform/frontend-sdk
```

**CDN:**
```html
<script src="https://cdn.securityplatform.com/sdk/v1/frontend-sdk.min.js"
        data-auto-init="true"
        data-control-plane-url="https://api.securityplatform.com"
        data-service-name="my-app">
</script>
```

**Manual Integration:**
```typescript
import { init } from '@security-platform/frontend-sdk';

init({
  controlPlaneUrl: 'https://api.securityplatform.com',
  serviceName: 'my-app',
  enableSessionRecording: true,
});
```

### Node.js Backend

```bash
npm install @security-platform/node-agent
```

```javascript
const { Agent, createMiddleware } = require('@security-platform/node-agent');
const express = require('express');

const agent = new Agent({
  controlPlaneUrl: process.env.SECURITY_PLATFORM_CONTROL_PLANE_URL,
  authKey: process.env.SECURITY_PLATFORM_AUTH_KEY,
});

agent.start();

const app = express();
app.use(createMiddleware(agent));
```

### Python Backend

```bash
pip install security-platform-agent
```

```python
from security_platform_agent import Agent, create_middleware
from fastapi import FastAPI

agent = Agent({
    "control_plane_url": os.getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL"),
    "auth_key": os.getenv("SECURITY_PLATFORM_AUTH_KEY"),
})

agent.start()

app = FastAPI()
app.add_middleware(create_middleware(agent))
```

### Java Backend

**Maven:**
```xml
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>agent</artifactId>
    <version>1.0.0</version>
</dependency>
```

**Gradle:**
```gradle
implementation 'com.securityplatform:agent:1.0.0'
```

**Usage:**
```java
import com.securityplatform.agent.Agent;
import com.securityplatform.agent.AgentConfig;
import com.securityplatform.agent.Filter;

AgentConfig config = new AgentConfig();
config.setControlPlaneUrl(System.getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL"));
config.setAuthKey(System.getenv("SECURITY_PLATFORM_AUTH_KEY"));

Agent agent = new Agent(config);
agent.start();

// Add filter to servlet context
Filter filter = new Filter(agent);
```

### .NET Backend

```bash
dotnet add package SecurityPlatform.Agent
```

```csharp
using SecurityPlatform.Agent;

var config = new AgentConfig
{
    ControlPlaneUrl = Environment.GetEnvironmentVariable("SECURITY_PLATFORM_CONTROL_PLANE_URL"),
    AuthKey = Environment.GetEnvironmentVariable("SECURITY_PLATFORM_AUTH_KEY")
};

var agent = new Agent(config);
agent.Start();

// In Startup.cs or Program.cs
app.UseSecurityPlatform(agent);
```

### Go Backend

```bash
go get github.com/security-platform/go-agent
```

```go
import (
    "github.com/security-platform/go-agent/agent"
    "context"
)

config := &agent.Config{
    ServiceName:     "my-service",
    OtlpEndpoint:    "http://localhost:4318/v1/traces",
    ControlPlaneURL: os.Getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL"),
    AuthKey:         os.Getenv("SECURITY_PLATFORM_AUTH_KEY"),
}

ag := agent.NewAgent(config)
ag.Start(context.Background())
defer ag.Stop(context.Background())

// Use middleware
m := agent.NewMiddleware(ag)
http.Handle("/", m.Handler(yourHandler))
```

### iOS SDK

**CocoaPods:**
```ruby
pod 'SecurityPlatform', '~> 1.0'
```

**Swift Package Manager:**
```swift
dependencies: [
    .package(url: "https://github.com/security-platform/ios-sdk.git", from: "1.0.0")
]
```

**Usage:**
```swift
import SecurityPlatform

let config = AgentConfig(
    controlPlaneUrl: "https://api.securityplatform.com",
    authKey: "your-key",
    serviceName: "my-ios-app"
)

let agent = Agent(config: config)
agent.start()
```

### Android SDK

**Gradle:**
```gradle
dependencies {
    implementation("com.securityplatform:agent:1.0.0")
}
```

**Usage:**
```kotlin
import com.securityplatform.agent.Agent
import com.securityplatform.agent.AgentConfig

val config = AgentConfig(
    controlPlaneUrl = "https://api.securityplatform.com",
    authKey = "your-key",
    serviceName = "my-android-app"
)

val agent = Agent(config)
agent.start()
```

### Desktop Agent (Linux)

**DEB Package:**
```bash
wget https://releases.securityplatform.com/agent.deb
sudo dpkg -i agent.deb
sudo systemctl enable security-platform-agent
sudo systemctl start security-platform-agent
```

**RPM Package:**
```bash
wget https://releases.securityplatform.com/agent.rpm
sudo rpm -ivh agent.rpm
sudo systemctl enable security-platform-agent
sudo systemctl start security-platform-agent
```

### Desktop Agent (macOS)

**Homebrew:**
```bash
brew tap security-platform/agent
brew install security-platform-agent
```

**DMG:**
```bash
# Download and install from releases page
open security-platform-agent.dmg
```

### Desktop Agent (Windows)

**Chocolatey:**
```powershell
choco install security-platform-agent
```

**MSI:**
```powershell
# Download and run installer from releases page
Start-Process msiexec.exe -ArgumentList "/i agent.msi /quiet"
```

### Kubernetes

**Helm:**
```bash
helm repo add security-platform https://charts.securityplatform.com
helm repo update
helm install security-platform-agent security-platform/security-platform-agent \
  --set config.controlPlaneUrl=https://api.securityplatform.com \
  --set config.authKey=your-token
```

**Kubernetes Operator:**
```bash
# Install CRDs
kubectl apply -f https://raw.githubusercontent.com/security-platform/operator/main/config/crd/bases/securityplatform.io_agents.yaml

# Deploy agent
kubectl apply -f - <<EOF
apiVersion: securityplatform.io/v1
kind: Agent
metadata:
  name: my-agent
spec:
  controlPlaneUrl: https://api.securityplatform.com
  authKey: your-token
  serviceName: my-service
EOF
```

### Docker

```bash
docker run -d \
  -e SECURITY_PLATFORM_CONTROL_PLANE_URL=https://api.securityplatform.com \
  -e SECURITY_PLATFORM_AUTH_KEY=your-token \
  securityplatform/agent:latest
```

## Air-Gapped Installation

1. **Generate Airgap Bundle:**
```bash
builder airgap --version=1.0.0 --include=deb,rpm,msi,dmg --out=./airgap-bundle
```

2. **Transfer Bundle:**
```bash
scp -r airgap-bundle/ user@target-server:/tmp/
```

3. **Install:**
```bash
cd /tmp/airgap-bundle
chmod +x install.sh
sudo ./install.sh
```

## Configuration

After installation, configure the agent:

**Linux/macOS:**
```bash
sudo nano /etc/security-platform/config.yaml
```

**Windows:**
```powershell
notepad C:\ProgramData\SecurityPlatform\config.yaml
```

**Configuration Options:**
```yaml
control_plane_url: https://api.securityplatform.com
auth_key: your-bootstrap-token
service_name: my-service
environment: production
otlp_endpoint: http://localhost:4318/v1/traces
local_policy: observe  # or "block"
redaction:
  enabled: true
  rules: []
```

## Verification

**Linux:**
```bash
systemctl status security-platform-agent
journalctl -u security-platform-agent -f
```

**macOS:**
```bash
launchctl list | grep securityplatform
tail -f /var/log/security-platform-agent.log
```

**Windows:**
```powershell
Get-Service SecurityPlatformAgent
Get-EventLog -LogName Application -Source SecurityPlatformAgent
```

## Troubleshooting

### Agent Not Connecting

1. Check network connectivity:
```bash
curl https://api.securityplatform.com/health
```

2. Verify authentication:
```bash
echo $SECURITY_PLATFORM_AUTH_KEY
```

3. Check agent logs (see Verification section above)

### Telemetry Not Appearing

1. Verify OTLP endpoint:
```bash
curl http://localhost:4318/v1/traces
```

2. Check collector logs:
```bash
docker logs otel-collector
```

## Support

- Documentation: https://docs.securityplatform.com
- Support: support@securityplatform.com
- GitHub: https://github.com/security-platform
