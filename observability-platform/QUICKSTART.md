# Quick Start Guide

## Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for builder CLI)
- Node.js 20+ (for frontend SDK and Node.js agents)
- Python 3.8+ (for Python agents)

## 1. Build the Builder CLI

```bash
cd observability-platform/builder
go mod download
go build -o ../../bin/builder ./main.go
```

Or use Make:

```bash
make builder
```

## 2. Start Local Control Plane

```bash
cd observability-platform/control-plane
docker-compose up -d
```

This starts:
- OpenSearch on http://localhost:9200
- OpenSearch Dashboards on http://localhost:5601
- OTEL Collector on ports 4317 (gRPC) and 4318 (HTTP)
- MinIO on http://localhost:9000 (console: http://localhost:9001)
- Control Plane API on http://localhost:8080

Verify it's running:

```bash
curl http://localhost:8080/api/v1/health
```

## 3. Create and Build an Agent

### Node.js Agent Example

```bash
# Initialize agent project
./bin/builder init my-node-agent --lang=node

# Navigate to project
cd my-node-agent

# Install dependencies
npm install

# Build
npm run build

# Use in your Express app
```

In your Express app:

```javascript
const { ObservabilityAgent } = require('observability-agent-node');
const { observabilityMiddleware } = require('observability-agent-node/middleware');

const agent = new ObservabilityAgent();
agent.init({
  controlPlaneUrl: 'http://localhost:8080',
  authKey: 'your-auth-key',
  serviceName: 'my-service',
  environment: 'development',
});

app.use(observabilityMiddleware);
```

### Python Agent Example

```bash
# Initialize agent project
./bin/builder init my-python-agent --lang=python

# Install
cd my-python-agent
pip install -e .

# Use in FastAPI app
```

In your FastAPI app:

```python
from fastapi import FastAPI
from observability_agent import ObservabilityAgent, observability_middleware

app = FastAPI()
agent = ObservabilityAgent({
    'control_plane_url': 'http://localhost:8080',
    'auth_key': 'your-auth-key',
    'service_name': 'my-service',
    'environment': 'development',
})
agent.init()
app.add_middleware(observability_middleware)
```

## 4. Build Frontend SDK

```bash
cd observability-platform/sdk/frontend
npm install
npm run build
npm pack
```

Use in your frontend:

```typescript
import { FrontendSDK } from '@observability-platform/frontend-sdk';

const sdk = new FrontendSDK();
sdk.init({
  controlPlaneUrl: 'http://localhost:8080',
  authKey: 'your-auth-key',
  serviceName: 'my-web-app',
  environment: 'development',
  enableSessionReplay: true,
});
```

## 5. Package Agents for Distribution

```bash
# Build agent
./bin/builder build agent --lang=node --version=1.0.0

# Package as DEB
./bin/builder package --target=deb --out=./dist --version=1.0.0

# Sign package (requires GPG key)
./bin/builder sign --artifact ./dist/observability-agent_1.0.0_amd64.deb --key ./keys/private.pem

# Upload to repository
./bin/builder upload --artifact ./dist/*.deb --repo https://artifactory.example.com --repo-type=artifactory --token <token>
```

## 6. Generate Airgap Bundle

```bash
# Build all components first
make all

# Generate airgap bundle
./airgap/generate.sh 1.0.0

# Transfer to target system and install
scp airgap-bundle-1.0.0.tar.gz user@target:/tmp/
ssh user@target "cd /tmp && tar -xzf airgap-bundle-1.0.0.tar.gz && cd airgap-bundle-1.0.0 && ./scripts/install.sh"
```

## 7. Deploy to Kubernetes

```bash
# Install Helm chart
helm install observability-platform ./helm/observability-platform \
  --set global.controlPlaneUrl=https://control.example.com \
  --set opensearch.replicas=3 \
  --set minio.persistence.enabled=true
```

## Testing the Setup

1. **Send test trace from agent:**
   ```bash
   curl -X POST http://localhost:4318/v1/traces \
     -H "Content-Type: application/json" \
     -d '{"resourceSpans": [{"resource": {"attributes": [{"key": "service.name", "value": {"stringValue": "test-service"}}]}}]}'
   ```

2. **Check OpenSearch:**
   ```bash
   curl http://localhost:9200/observability-traces/_search?pretty
   ```

3. **View in OpenSearch Dashboards:**
   Open http://localhost:5601 and create index pattern `observability-*`

## Next Steps

- Read [ARCHITECTURE.md](docs/ARCHITECTURE.md) for system design
- Read [DEPLOYMENT.md](docs/DEPLOYMENT.md) for production deployment
- Check [IMPLEMENTATION_STATUS.md](IMPLEMENTATION_STATUS.md) for what's implemented
