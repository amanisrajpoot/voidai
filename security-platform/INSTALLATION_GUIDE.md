# Security Platform - Complete Installation Guide

This guide covers installation and integration for all platforms and languages.

## Quick Install (One-Liner)

### Linux/macOS

```bash
curl -fsSL https://install.securityplatform.com | bash
```

### Windows (PowerShell)

```powershell
iex (New-Object Net.WebClient).DownloadString('https://install.securityplatform.com/install.ps1')
```

## Platform-Specific Installation

### Desktop Agents

#### Linux (DEB)

```bash
wget https://releases.securityplatform.com/security-platform-agent_1.0.0_amd64.deb
sudo dpkg -i security-platform-agent_1.0.0_amd64.deb
sudo systemctl enable security-platform-agent
sudo systemctl start security-platform-agent
```

#### Linux (RPM)

```bash
wget https://releases.securityplatform.com/security-platform-agent-1.0.0-1.x86_64.rpm
sudo rpm -ivh security-platform-agent-1.0.0-1.x86_64.rpm
sudo systemctl enable security-platform-agent
sudo systemctl start security-platform-agent
```

#### macOS (Homebrew)

```bash
brew tap security-platform/agent
brew install security-platform-agent
```

#### macOS (DMG/PKG)

```bash
# Download and install from releases page
# Or use installer:
curl -L https://releases.securityplatform.com/security-platform-agent-1.0.0.pkg -o /tmp/agent.pkg
sudo installer -pkg /tmp/agent.pkg -target /
```

#### Windows (Chocolatey)

```powershell
choco install security-platform-agent
```

#### Windows (MSI)

```powershell
# Download and run installer from releases page
```

## Language-Specific Agents

### Node.js

```bash
npm install @security-platform/node-agent
```

```javascript
const { Agent } = require('@security-platform/node-agent');

const agent = new Agent({
  controlPlaneUrl: process.env.SECURITY_PLATFORM_CONTROL_PLANE_URL,
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
    "control_plane_url": os.getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL"),
    "auth_key": os.getenv("SECURITY_PLATFORM_AUTH_KEY"),
    "service_name": "my-service"
})

agent.start()
```

### Java

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
    ControlPlaneURL: os.Getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL"),
    AuthKey:         os.Getenv("SECURITY_PLATFORM_AUTH_KEY"),
    ServiceName:     "my-service",
}

ag := agent.NewAgent(config)
ag.Start()
```

## Mobile SDKs

### iOS

**Swift Package Manager:**

Add to `Package.swift` or Xcode:

```
https://github.com/security-platform/agent-ios
```

**CocoaPods:**

```ruby
pod 'SecurityPlatform', '~> 1.0.0'
```

```swift
import SecurityPlatform

let config = AgentConfig()
config.controlPlaneURL = "https://api.securityplatform.com"
config.authKey = ProcessInfo.processInfo.environment["SECURITY_PLATFORM_AUTH_KEY"]
config.serviceName = "my-ios-app"

let agent = SecurityPlatformAgent(config: config)
try? agent.start()
```

### Android

**Gradle:**

```gradle
implementation 'com.securityplatform:agent:1.0.0'
```

**Maven:**

```xml
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>agent</artifactId>
    <version>1.0.0</version>
</dependency>
```

```kotlin
import com.securityplatform.agent.Agent
import com.securityplatform.agent.AgentConfig

val config = AgentConfig.fromEnvironment(this)
config.serviceName = "my-android-app"

val agent = Agent.create(config)
agent.start()
```

## Frontend SDK

### NPM

```bash
npm install @security-platform/frontend-sdk
```

### CDN

```html
<script src="https://cdn.securityplatform.com/sdk/v1/frontend-sdk.min.js"
        data-auto-init="true"
        data-control-plane-url="https://api.securityplatform.com"
        data-service-name="my-app">
</script>
```

### Manual

```typescript
import { init } from '@security-platform/frontend-sdk';

init({
  controlPlaneUrl: 'https://api.securityplatform.com',
  serviceName: 'my-app',
  enableSessionRecording: true,
});
```

## Kubernetes Operator

### Install Operator

```bash
# Install CRDs
kubectl apply -f https://raw.githubusercontent.com/security-platform/operator/main/config/crd/bases/securityplatform.io_agents.yaml

# Install operator
kubectl apply -f https://raw.githubusercontent.com/security-platform/operator/main/config/manager/manager.yaml
```

### Deploy Agent

```yaml
apiVersion: securityplatform.io/v1alpha1
kind: Agent
metadata:
  name: my-agent
spec:
  controlPlaneURL: "https://api.securityplatform.com"
  authKeySecretRef:
    name: security-platform-secret
    key: auth-key
  serviceName: "my-service"
  environment: "production"
  replicas: 2
```

## Air-Gap Installation

### Generate Bundle

```bash
builder package --target=airgap --out=./airgap-bundle
```

### Install

1. Transfer bundle to air-gapped system
2. Extract: `tar -xzf airgap-bundle.tar.gz`
3. Run installer: `sudo ./install.sh`
4. Configure: Edit `/etc/security-platform/agent-config.yaml`
5. Start service: `sudo systemctl start security-platform-agent`

## Configuration

All agents support configuration via:

1. **Environment Variables** (recommended):
   ```bash
   export SECURITY_PLATFORM_CONTROL_PLANE_URL=https://api.securityplatform.com
   export SECURITY_PLATFORM_AUTH_KEY=your-token-here
   ```

2. **Configuration File**:
   ```yaml
   control_plane_url: "https://api.securityplatform.com"
   auth_key: "your-token-here"
   service_name: "my-service"
   environment: "production"
   ```

3. **Code Configuration**:
   See language-specific examples above

## Verification

### Check Agent Status

**Linux:**
```bash
sudo systemctl status security-platform-agent
```

**macOS:**
```bash
sudo launchctl list | grep security-platform
```

**Windows:**
```powershell
Get-Service SecurityPlatformAgent
```

### View Logs

**Linux:**
```bash
sudo journalctl -u security-platform-agent -f
```

**macOS:**
```bash
tail -f /var/log/security-platform-agent.log
```

**Windows:**
```powershell
Get-Content C:\ProgramData\SecurityPlatform\logs\agent.log -Tail 50
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

3. Check agent logs (see above)

### Telemetry Not Appearing

1. Verify OTLP endpoint is accessible
2. Check collector logs
3. Verify storage backends are running

## Support

- Documentation: https://docs.securityplatform.com
- Support: support@securityplatform.com
- GitHub: https://github.com/security-platform
