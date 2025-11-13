# Security Platform - Quick Reference

## Builder CLI Commands

```bash
# Initialize agent project
builder init <project-name> --lang=<language>

# Build agent
builder build agent --lang=<lang> --version=<version>

# Create packages
builder package --target=<targets> --out=<dir>

# Sign artifacts
builder sign --artifact <file> --key <key-file>

# Upload to repository
builder upload --artifact <file> --repo <repo-url> --type <type>
```

**Languages**: `node`, `python`, `java`, `dotnet`, `go`, `frontend`, `ios`, `android`, `desktop`

**Package Targets**: `deb`, `rpm`, `msi`, `dmg`, `helm`, `airgap`

## Agent Installation

### Frontend (NPM)
```bash
npm install @security-platform/frontend-sdk
```

### Node.js Backend
```bash
npm install @security-platform/node-agent
```

### Python Backend
```bash
pip install security-platform-agent
```

### Java Backend
```xml
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>agent</artifactId>
    <version>1.0.0</version>
</dependency>
```

### .NET Backend
```bash
dotnet add package SecurityPlatform.Agent
```

### Go Backend
```bash
go get github.com/security-platform/go-agent
```

## Configuration

### Environment Variables
```bash
SECURITY_PLATFORM_CONTROL_PLANE_URL=https://api.securityplatform.com
SECURITY_PLATFORM_AUTH_KEY=your-token
SECURITY_PLATFORM_SERVICE_NAME=my-service
SECURITY_PLATFORM_ENVIRONMENT=production
SECURITY_PLATFORM_OTLP_ENDPOINT=http://localhost:4318
```

### Config File (YAML)
```yaml
control_plane_url: "https://api.securityplatform.com"
auth_key: "your-token"
service_name: "my-service"
environment: "production"
policy:
  mode: "observe"  # or "block"
redaction:
  enabled: true
  default_blocklist:
    - "password"
    - "token"
```

## Quick Start Examples

### Frontend (HTML)
```html
<script src="https://cdn.securityplatform.com/sdk/v1/frontend-sdk.min.js"
        data-auto-init="true"
        data-control-plane-url="https://api.securityplatform.com"
        data-service-name="my-app">
</script>
```

### Frontend (TypeScript)
```typescript
import { init } from '@security-platform/frontend-sdk';

init({
  controlPlaneUrl: 'https://api.securityplatform.com',
  serviceName: 'my-app',
  enableSessionRecording: true,
});
```

### Node.js (Express)
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

### Python (FastAPI)
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

## API Endpoints

### Control Plane API
```
GET  /api/v1/agents
POST /api/v1/agents
GET  /api/v1/agents/{id}
PUT  /api/v1/agents/{id}
DELETE /api/v1/agents/{id}

GET  /api/v1/rules
POST /api/v1/rules
GET  /api/v1/rules/{id}
PUT  /api/v1/rules/{id}
DELETE /api/v1/rules/{id}

POST /api/v1/playbooks/{id}/execute
POST /api/v1/webhooks/{playbook_id}
```

### OTLP Endpoints
```
POST /v1/traces   (gRPC: :4317, HTTP: :4318)
POST /v1/metrics
POST /v1/logs
```

## Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

**Services**:
- Control Plane: `http://localhost:8080`
- Collector: `http://localhost:4318`
- OpenSearch: `http://localhost:9200`
- ClickHouse: `http://localhost:8123`
- MinIO: `http://localhost:9000`
- Orchestrator: `http://localhost:8081`

## Makefile Targets

```bash
make build-builder          # Build builder CLI
make build-agents          # Build all agents
make build-frontend        # Build frontend SDK
make build-node           # Build Node.js agent
make build-python         # Build Python agent
make test                 # Run all tests
make package              # Create packages
make docker-build         # Build Docker images
make clean                # Clean build artifacts
make install              # Install builder CLI
```

## File Locations

- **Builder CLI**: `builder/`
- **Frontend SDK**: `agents/frontend/`
- **Node.js Agent**: `agents/node/`
- **Python Agent**: `agents/python/`
- **Control Plane**: `control-plane/api/`
- **Orchestrator**: `orchestrator/`
- **Config Schema**: `config/agent-config.yaml`
- **Examples**: `examples/`
- **Documentation**: `docs/`

## Troubleshooting

### Agent Not Connecting
```bash
# Check connectivity
curl https://api.securityplatform.com/health

# Verify auth key
echo $SECURITY_PLATFORM_AUTH_KEY

# Check logs (Linux)
journalctl -u security-platform-agent -f
```

### Telemetry Not Appearing
```bash
# Check OTLP endpoint
curl http://localhost:4318/v1/traces

# Check collector logs
docker logs otel-collector

# Check OpenSearch
curl http://localhost:9200/_cat/indices
```

## Package Formats

- **Linux**: `.deb` (Debian/Ubuntu), `.rpm` (RHEL/CentOS)
- **macOS**: `.dmg` (signed, notarized)
- **Windows**: `.msi` (WiX installer)
- **Kubernetes**: Helm charts
- **Airgap**: Complete offline bundles

## Redaction Rules

Default blocklist includes:
- `password`, `passwd`, `pwd`
- `secret`, `token`, `api_key`
- `credit_card`, `card_number`, `cvv`
- `ssn`, `social_security`
- `pin`, `passcode`

Custom rules via regex patterns in config.

## Policy Modes

- **observe**: Log events, don't block (default)
- **block**: Enforce policies and block violations

Start in `observe` mode for 48 hours, then switch to `block`.

## Support

- **Docs**: `docs/`
- **Examples**: `examples/`
- **Architecture**: `docs/ARCHITECTURE.md`
- **Deployment**: `docs/DEPLOYMENT.md`
- **API**: `docs/API.md`
