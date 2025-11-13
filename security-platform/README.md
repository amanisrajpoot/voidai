# Security Platform - Universal Security Observability Platform

Complete end-to-end security observability platform with agents and SDKs for all platforms.

## ✨ Features

- **Multi-Platform Support**: Web, Backend (Node.js, Python, Java, .NET, Go), Mobile (iOS, Android), Desktop (Windows, macOS, Linux)
- **Simple Installation**: One-line installers for all platforms
- **Easy Integration**: Minimal code changes required
- **Kubernetes Ready**: Helm charts and Operator support
- **Air-Gapped Support**: Complete offline installation bundles
- **Package Managers**: Homebrew, Chocolatey, DEB, RPM, MSI, DMG

## 🚀 Quick Start

### Install Agent (Linux/macOS)
```bash
curl -fsSL https://install.securityplatform.com | bash
```

### Install Agent (Windows)
```powershell
choco install security-platform-agent
```

### Add to Your Application

**Node.js:**
```bash
npm install @security-platform/node-agent
```
```javascript
const { Agent, createMiddleware } = require('@security-platform/node-agent');
const agent = new Agent({ controlPlaneUrl: process.env.SECURITY_PLATFORM_CONTROL_PLANE_URL });
agent.start();
app.use(createMiddleware(agent));
```

**Python:**
```bash
pip install security-platform-agent
```
```python
from security_platform_agent import Agent, create_middleware
agent = Agent({"control_plane_url": os.getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL")})
agent.start()
app.add_middleware(create_middleware(agent))
```

**Java:**
```xml
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>agent</artifactId>
    <version>1.0.0</version>
</dependency>
```

**Go:**
```bash
go get github.com/security-platform/go-agent
```

**iOS:**
```ruby
pod 'SecurityPlatform', '~> 1.0'
```

**Android:**
```gradle
implementation("com.securityplatform:agent:1.0.0")
```

## 📦 Available Agents & SDKs

### Backend Agents
- ✅ **Node.js** - Express/Fastify middleware
- ✅ **Python** - FastAPI/Django middleware  
- ✅ **Java** - Servlet Filter
- ✅ **.NET** - ASP.NET Core middleware
- ✅ **Go** - HTTP middleware

### Frontend SDKs
- ✅ **JavaScript/TypeScript** - Browser SDK with session recording

### Mobile SDKs
- ✅ **iOS** - Swift SDK with CocoaPods/SPM
- ✅ **Android** - Kotlin SDK with Gradle

### Desktop Agents
- ✅ **Linux** - systemd service (DEB/RPM)
- ✅ **macOS** - launchd daemon (Homebrew/DMG)
- ✅ **Windows** - Windows Service (Chocolatey/MSI)

### Kubernetes
- ✅ **Helm Chart** - Production-ready Helm chart
- ✅ **Operator** - Kubernetes Operator with CRDs

## 🛠️ Installation Methods

### Package Managers

**Homebrew (macOS):**
```bash
brew tap security-platform/agent
brew install security-platform-agent
```

**Chocolatey (Windows):**
```powershell
choco install security-platform-agent
```

**DEB (Debian/Ubuntu):**
```bash
wget https://releases.securityplatform.com/agent.deb
sudo dpkg -i agent.deb
```

**RPM (RHEL/CentOS):**
```bash
wget https://releases.securityplatform.com/agent.rpm
sudo rpm -ivh agent.rpm
```

### Kubernetes

**Helm:**
```bash
helm repo add security-platform https://charts.securityplatform.com
helm install agent security-platform/security-platform-agent
```

**Operator:**
```bash
kubectl apply -f https://raw.githubusercontent.com/security-platform/operator/main/config/crd/bases/securityplatform.io_agents.yaml
```

### Docker
```bash
docker run -d \
  -e SECURITY_PLATFORM_CONTROL_PLANE_URL=https://api.securityplatform.com \
  -e SECURITY_PLATFORM_AUTH_KEY=your-token \
  securityplatform/agent:latest
```

### Air-Gapped
```bash
# Generate bundle
builder airgap --version=1.0.0 --out=./bundle

# Transfer and install
scp -r bundle/ server:/tmp/
ssh server "cd /tmp/bundle && sudo ./install.sh"
```

## 📚 Documentation

- [Quick Start Guide](./QUICK_START.md) - Get started in minutes
- [Installation Guide](./INSTALLATION.md) - Detailed installation instructions
- [Architecture](./docs/ARCHITECTURE.md) - System architecture
- [API Reference](./docs/API.md) - API documentation
- [Deployment Guide](./docs/DEPLOYMENT.md) - Deployment options

## 🏗️ Building from Source

```bash
# Build builder CLI
cd builder && go build -o ../bin/builder

# Initialize agent project
./bin/builder init my-agent --lang=node

# Build agent
./bin/builder build agent --lang=node --version=1.0.0

# Package for distribution
./bin/builder package --target=deb,rpm,msi,dmg --out=./dist

# Generate airgap bundle
./bin/builder airgap --version=1.0.0 --out=./bundle
```

## 🎯 Key Features

- **Zero-Config**: Works out of the box with sensible defaults
- **Automatic Instrumentation**: No code changes needed for basic setup
- **PII Redaction**: Built-in protection for sensitive data
- **Policy Enforcement**: Local observe/block modes
- **Offline Support**: Works in air-gapped environments
- **Multi-Platform**: Same API across all platforms

## 🔧 Configuration

Configuration via environment variables or config file:

```yaml
control_plane_url: https://api.securityplatform.com
auth_key: your-token
service_name: my-service
environment: production
otlp_endpoint: http://localhost:4318/v1/traces
local_policy: observe
```

## 📊 Project Structure

```
security-platform/
├── agents/              # All agent implementations
│   ├── node/           # Node.js agent
│   ├── python/         # Python agent
│   ├── java/           # Java agent
│   ├── dotnet/         # .NET agent
│   ├── go/             # Go agent
│   ├── frontend/       # Frontend SDK
│   ├── ios/            # iOS SDK
│   ├── android/        # Android SDK
│   └── desktop/        # Desktop agents
├── builder/             # Builder CLI
├── packaging/          # Packaging scripts
│   ├── homebrew/       # Homebrew formula
│   ├── chocolatey/     # Chocolatey package
│   ├── helm/           # Helm charts
│   └── kubernetes-operator/  # K8s operator
├── scripts/            # Installation scripts
└── docs/               # Documentation
```

## 🤝 Contributing

Contributions welcome! See [HOW_TO_CONTRIBUTE.md](../HOW_TO_CONTRIBUTE.md) for guidelines.

## 📄 License

See [LICENSE](./LICENSE) file for details.

## 🆘 Support

- Documentation: https://docs.securityplatform.com
- GitHub Issues: https://github.com/security-platform/issues
- Support Email: support@securityplatform.com
