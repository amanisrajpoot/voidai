# Security Platform Agent - Quick Start Guide

## 🚀 One-Line Installation

### Linux
```bash
curl -fsSL https://install.securityplatform.com | bash
```

### macOS
```bash
brew install security-platform/agent/security-platform-agent
```

### Windows
```powershell
choco install security-platform-agent
```

## 📦 Language-Specific Integration

### Frontend (JavaScript/TypeScript)
```bash
npm install @security-platform/frontend-sdk
```
```javascript
import { init } from '@security-platform/frontend-sdk';
init({ controlPlaneUrl: 'https://api.securityplatform.com', serviceName: 'my-app' });
```

### Node.js Backend
```bash
npm install @security-platform/node-agent
```
```javascript
const { Agent, createMiddleware } = require('@security-platform/node-agent');
const agent = new Agent({ controlPlaneUrl: process.env.SECURITY_PLATFORM_CONTROL_PLANE_URL });
agent.start();
app.use(createMiddleware(agent));
```

### Python Backend
```bash
pip install security-platform-agent
```
```python
from security_platform_agent import Agent, create_middleware
agent = Agent({"control_plane_url": os.getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL")})
agent.start()
app.add_middleware(create_middleware(agent))
```

### Java Backend
```xml
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>agent</artifactId>
    <version>1.0.0</version>
</dependency>
```
```java
Agent agent = new Agent(config);
agent.start();
// Add Filter to servlet context
```

### .NET Backend
```bash
dotnet add package SecurityPlatform.Agent
```
```csharp
var agent = new Agent(config);
agent.Start();
app.UseSecurityPlatform(agent);
```

### Go Backend
```bash
go get github.com/security-platform/go-agent
```
```go
agent := agent.NewAgent(config)
agent.Start(ctx)
// Use middleware
```

### iOS (Swift)
```ruby
pod 'SecurityPlatform', '~> 1.0'
```
```swift
let agent = Agent(config: config)
agent.start()
```

### Android (Kotlin)
```gradle
implementation("com.securityplatform:agent:1.0.0")
```
```kotlin
val agent = Agent(config)
agent.start()
```

## 🐳 Kubernetes Installation

```bash
helm repo add security-platform https://charts.securityplatform.com
helm install agent security-platform/security-platform-agent \
  --set config.controlPlaneUrl=https://api.securityplatform.com \
  --set config.authKey=your-token
```

Or use the Operator:
```bash
kubectl apply -f - <<EOF
apiVersion: securityplatform.io/v1
kind: Agent
metadata:
  name: my-agent
spec:
  controlPlaneUrl: https://api.securityplatform.com
  authKey: your-token
EOF
```

## 🔒 Air-Gapped Installation

```bash
# Generate bundle
builder airgap --version=1.0.0 --out=./bundle

# Transfer and install
scp -r bundle/ server:/tmp/
ssh server "cd /tmp/bundle && sudo ./install.sh"
```

## ⚙️ Configuration

Set environment variables or edit config file:

**Environment Variables:**
```bash
export SECURITY_PLATFORM_CONTROL_PLANE_URL=https://api.securityplatform.com
export SECURITY_PLATFORM_AUTH_KEY=your-token
export SECURITY_PLATFORM_SERVICE_NAME=my-service
```

**Config File** (`/etc/security-platform/config.yaml`):
```yaml
control_plane_url: https://api.securityplatform.com
auth_key: your-token
service_name: my-service
environment: production
otlp_endpoint: http://localhost:4318/v1/traces
local_policy: observe
```

## ✅ Verification

```bash
# Check service status
systemctl status security-platform-agent

# View logs
journalctl -u security-platform-agent -f

# Test connectivity
curl https://api.securityplatform.com/health
```

## 🛠️ Building from Source

```bash
# Build builder CLI
cd builder && go build -o ../bin/builder

# Build agent
./bin/builder build agent --lang=node --version=1.0.0

# Package
./bin/builder package --target=deb,rpm,msi,dmg --out=./dist
```

## 📚 Next Steps

- Read [INSTALLATION.md](./INSTALLATION.md) for detailed instructions
- Check [docs/](./docs/) for architecture and API documentation
- See [examples/](./examples/) for integration examples

## 🆘 Support

- Documentation: https://docs.securityplatform.com
- GitHub: https://github.com/security-platform
- Support: support@securityplatform.com
