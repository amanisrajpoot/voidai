# Quick Start Guide

This guide will help you get started with the Security Platform in minutes.

## Prerequisites

- Docker and Docker Compose (for local development)
- Node.js 18+ (for frontend SDK)
- Python 3.8+ (for Python agent)
- Go 1.21+ (for building agents)

## Local Development Setup

### 1. Start Local Control Plane

```bash
cd security-platform
docker-compose up -d
```

This starts:
- Control Plane API (port 8080)
- OpenTelemetry Collector (port 4318)
- OpenSearch (port 9200)
- MinIO for session storage (port 9000)

### 2. Initialize an Agent Project

```bash
# Build the builder CLI
cd builder
go build -o builder .

# Initialize a Node.js agent
./builder init my-node-agent --lang=node

# Initialize a Python agent
./builder init my-python-agent --lang=python
```

### 3. Frontend SDK Integration

#### Install the SDK

```bash
npm install @security-platform/frontend-sdk
```

#### Basic Usage

```javascript
import { initSecuritySDK } from '@security-platform/frontend-sdk';

// Initialize the SDK
const sdk = initSecuritySDK({
  controlPlaneUrl: 'http://localhost:8080',
  authKey: 'your-auth-key',
  serviceName: 'my-web-app',
  environment: 'development',
  redactionRules: ['password', 'card', 'ssn'],
  enableSessionReplay: true,
  enableTelemetry: true,
});

// Capture custom events
sdk.captureEvent('user_action', {
  action: 'button_click',
  page: '/dashboard',
});

// Mark a field as safe (won't be redacted)
sdk.markFieldSafe('public_id');
```

#### CDN Usage

```html
<script src="https://cdn.securityplatform.com/sdk/v0.1.0/index.umd.js"></script>
<script>
  const sdk = SecurityPlatform.initSecuritySDK({
    controlPlaneUrl: 'https://console.securityplatform.com',
    authKey: 'your-auth-key',
  });
</script>
```

### 4. Backend Agent Integration

#### Node.js / Express

```bash
npm install @security-platform/node-agent
```

```javascript
const express = require('express');
const { securityMiddleware } = require('@security-platform/node-agent');

const app = express();

// Add security middleware
app.use(securityMiddleware({
  controlPlaneUrl: 'http://localhost:8080',
  authKey: 'your-auth-key',
  serviceName: 'api-server',
  environment: 'production',
  localPolicy: 'observe', // or 'block'
}));

app.get('/api/users', (req, res) => {
  // Access security context
  const { requestId, traceId } = req.securityContext;
  res.json({ users: [] });
});

app.listen(3000);
```

#### Python / FastAPI

```bash
pip install security-platform-python-agent
```

```python
from fastapi import FastAPI
from security_platform import SecurityMiddleware, AgentConfig, create_agent

app = FastAPI()

# Create agent
config = AgentConfig(
    control_plane_url="http://localhost:8080",
    auth_key="your-auth-key",
    service_name="api-server",
    environment="production",
    local_policy="observe",
)

agent = create_agent(config)

# Add middleware
app.add_middleware(SecurityMiddleware, agent=agent)

@app.get("/api/users")
async def get_users(request: Request):
    # Access security context
    context = request.state.security_context
    return {"users": []}
```

### 5. Building and Packaging

#### Build an Agent

```bash
# Build Node.js agent
cd my-node-agent
./builder build agent --lang=node --version=1.0.0

# Build Python agent
cd my-python-agent
./builder build agent --lang=python --version=1.0.0
```

#### Create Installers

```bash
# Create DEB package
./builder package --target=deb --version=1.0.0 --out=./dist

# Create RPM package
./builder package --target=rpm --version=1.0.0 --out=./dist

# Create multiple formats
./builder package --target=deb,rpm,helm --version=1.0.0 --out=./dist
```

#### Sign Artifacts

```bash
# Sign DEB package with GPG
./builder sign --artifact ./dist/agent.deb --method=gpg --key=your-key-id

# Sign macOS DMG
./builder sign --artifact ./dist/agent.dmg --method=codesign --key="Developer ID"
```

#### Upload to Repositories

```bash
# Upload to npm
./builder upload --artifact ./dist/*.tgz --repo=npmjs.com --type=npm

# Upload to PyPI
./builder upload --artifact ./dist/*.whl --repo=pypi.org --type=pypi

# Upload to S3
./builder upload --artifact ./dist/*.deb --repo=s3://my-bucket/packages --type=s3
```

## Configuration

### Agent Configuration File

Create `agent-config.yaml`:

```yaml
control_plane:
  url: "https://console.securityplatform.com"

auth:
  bootstrap_token: "${SECURITY_PLATFORM_BOOTSTRAP_TOKEN}"

service:
  name: "${SERVICE_NAME:-my-service}"
  environment: "${ENV:-production}"

telemetry:
  batch_size: 50
  flush_interval_ms: 5000
  offline_mode: true

redaction:
  default_rules:
    - "password"
    - "card"
    - "ssn"
  hash_payloads: true

policy:
  mode: "observe"
  cache_ttl: 300
```

### Environment Variables

```bash
export SECURITY_PLATFORM_CONTROL_PLANE_URL="https://console.securityplatform.com"
export SECURITY_PLATFORM_BOOTSTRAP_TOKEN="your-token"
export SERVICE_NAME="my-service"
export ENV="production"
```

## Next Steps

- Read the [Architecture Guide](./ARCHITECTURE.md)
- Check out [Examples](./examples/)
- Review [Security Best Practices](./SECURITY.md)
- Explore [API Documentation](./API.md)
