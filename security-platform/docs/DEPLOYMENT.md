# Deployment Guide

## Quick Start

### 1. Install Builder CLI

```bash
# Using Homebrew (macOS)
brew install security-platform/builder

# Using npm
npm install -g @security-platform/builder-cli

# Download binary
curl -L https://github.com/security-platform/builder/releases/latest/download/builder-linux-amd64 -o /usr/local/bin/builder
chmod +x /usr/local/bin/builder
```

### 2. Initialize Agent Project

```bash
builder init my-agent --lang=node
cd my-agent
```

### 3. Configure Agent

Edit `config/agent-config.yaml`:

```yaml
control_plane_url: "https://api.securityplatform.com"
auth_key: "your-bootstrap-token"
service_name: "my-service"
environment: "production"
```

### 4. Build and Package

```bash
# Build agent
builder build agent --lang=node --version=1.0.0

# Create packages
builder package --target=deb,rpm,msi,dmg --out=./dist

# Sign packages
builder sign --artifact ./dist/agent.deb --key ./keys/private.pem
```

## Platform-Specific Installation

### Frontend (Web)

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

**Manual:**
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

### .NET Backend

```bash
dotnet add package SecurityPlatform.Agent
```

### Go Backend

```bash
go get github.com/security-platform/go-agent
```

### Linux (DEB/RPM)

```bash
# Debian/Ubuntu
wget https://releases.securityplatform.com/agent.deb
sudo dpkg -i agent.deb
sudo systemctl enable security-platform-agent
sudo systemctl start security-platform-agent

# RHEL/CentOS
wget https://releases.securityplatform.com/agent.rpm
sudo rpm -ivh agent.rpm
sudo systemctl enable security-platform-agent
sudo systemctl start security-platform-agent
```

### macOS

```bash
# Homebrew
brew install security-platform/agent

# DMG
# Download and install from releases page
```

### Windows

```powershell
# Chocolatey
choco install security-platform-agent

# MSI
# Download and run installer from releases page
```

### Kubernetes

```bash
# Add Helm repository
helm repo add security-platform https://charts.securityplatform.com
helm repo update

# Install agent
helm install security-platform-agent security-platform/security-platform-agent \
  --set config.controlPlaneUrl=https://api.securityplatform.com \
  --set config.authKey=your-token
```

### Docker

```bash
docker run -d \
  -e SECURITY_PLATFORM_CONTROL_PLANE_URL=https://api.securityplatform.com \
  -e SECURITY_PLATFORM_AUTH_KEY=your-token \
  securityplatform/agent:latest
```

## On-Premises Deployment

### Air-Gapped Installation

1. **Generate Airgap Bundle:**
```bash
builder package --target=airgap --out=./airgap-bundle
```

2. **Transfer Bundle:**
```bash
# Copy airgap-bundle/ to target system
scp -r airgap-bundle/ user@target-server:/tmp/
```

3. **Install:**
```bash
cd /tmp/airgap-bundle
chmod +x install.sh
sudo ./install.sh
```

### Docker Compose (Local Dev)

```yaml
version: '3.8'
services:
  control-plane:
    image: securityplatform/control-plane:latest
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://...
  
  collector:
    image: otel/opentelemetry-collector:latest
    volumes:
      - ./otel-collector.yaml:/etc/otel-collector-config.yaml
    command: ["--config=/etc/otel-collector-config.yaml"]
  
  opensearch:
    image: opensearchproject/opensearch:latest
    ports:
      - "9200:9200"
  
  minio:
    image: minio/minio:latest
    ports:
      - "9000:9000"
    command: server /data
```

## Configuration

See `config/agent-config.yaml` for full configuration options.

Key settings:
- `control_plane_url`: Control plane endpoint
- `auth_key`: Bootstrap token or PKI certificate path
- `policy.mode`: `observe` or `block`
- `redaction.enabled`: Enable PII redaction
- `offline.enabled`: Enable offline caching

## Troubleshooting

### Agent Not Connecting

1. Check network connectivity:
```bash
curl https://api.securityplatform.com/health
```

2. Verify authentication:
```bash
# Check auth key is set
echo $SECURITY_PLATFORM_AUTH_KEY
```

3. Check agent logs:
```bash
# Linux
journalctl -u security-platform-agent -f

# Docker
docker logs security-platform-agent
```

### Telemetry Not Appearing

1. Verify OTLP endpoint:
```bash
curl http://localhost:4318/v1/traces
```

2. Check collector logs:
```bash
docker logs otel-collector
```

3. Verify storage:
```bash
# OpenSearch
curl http://localhost:9200/_cat/indices
```

## Support

- Documentation: https://docs.securityplatform.com
- Support: support@securityplatform.com
- GitHub: https://github.com/security-platform
