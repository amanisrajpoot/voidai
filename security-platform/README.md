# Security Observability Platform

A comprehensive security observability platform with installable agents/SDKs for web frontends, backends, desktop, mobile, and network appliances.

## Architecture Overview

This platform provides:
- **Control Plane / Console**: SaaS or on-prem UI for incidents, session replay, traces, rules, playbooks, and orchestration
- **Telemetry Bus & Storage**: OTEL collector, event store (OpenSearch/ClickHouse), session storage (S3/MinIO)
- **Edge Modules**: WAF (ModSecurity/nginx) packages, reverse-proxy integrations
- **API Gateway Integrations**: Kong/Envoy plugins and admin connectors
- **Language Agents / RASP SDKs**: Runtime libraries for Python, Node, Java, .NET, Go, Ruby
- **Frontend SDK**: JS/TS SDK + rrweb integration for RUM/session capture
- **Mobile SDKs**: iOS (Swift) and Android (Kotlin/Java) SDKs
- **Desktop Agents**: Windows service, macOS daemon, Linux daemon
- **Network Installer / Appliance**: Docker/VM images and Kubernetes Operator
- **Builder / CLI**: Single CLI tool to generate all installers and packages
- **Orchestrator**: Webhook-driven runner for automated responses

## Project Structure

```
security-platform/
├── builder/              # Go-based CLI builder
├── control-plane/        # Web console and API
├── agents/               # Language-specific agents
│   ├── node/
│   ├── python/
│   ├── java/
│   ├── dotnet/
│   └── go/
├── sdks/                 # SDKs for various platforms
│   ├── frontend/         # TypeScript/JavaScript SDK
│   ├── ios/              # Swift SDK
│   └── android/          # Kotlin/Java SDK
├── desktop-agents/       # Desktop host agents
│   ├── windows/
│   ├── darwin/
│   └── linux/
├── edge/                 # Edge modules (WAF, API Gateway)
│   ├── waf/
│   ├── kong-plugin/
│   └── envoy-plugin/
├── telemetry/            # Telemetry infrastructure
│   ├── collector/        # OTEL collector configs
│   └── storage/          # Storage backends
├── orchestrator/         # SOAR-lite orchestrator
├── packaging/            # Packaging templates and scripts
└── docs/                 # Documentation

```

## Tech Stack

- **CLI / Builder**: Go + goreleaser
- **Telemetry**: OpenTelemetry (OTLP) + OpenSearch / ClickHouse
- **Frontend session**: rrweb + @opentelemetry/sdk-trace-web
- **Backend middleware**: Express / FastAPI packages + OpenTelemetry SDKs
- **WAF**: nginx + ModSecurity + OWASP CRS
- **API Gateway**: Kong (admin API)
- **Storage**: MinIO (S3) for session blobs
- **Packaging**: dpkg/rpm/WiX (MSI)/pkgbuild + codesign/signtool
- **CI**: GitHub Actions

## Quick Start

### Building Packages

```bash
# Initialize a new agent project
./builder init my-agent --lang=node

# Build an agent
./builder build agent --lang=node --version=1.2.0

# Package for distribution
./builder package --target=deb,rpm,msi,dmg,helm --out=./dist

# Sign artifacts
./builder sign --artifact ./dist/agent.deb --key ./keys/private.pem

# Upload to repository
./builder upload --artifact ./dist/*.deb --repo artifactory
```

### Using SDKs

#### Frontend (JavaScript/TypeScript)
```javascript
import { initSecuritySDK } from '@security-platform/frontend-sdk';

initSecuritySDK({
  controlPlaneUrl: 'https://console.example.com',
  authKey: 'your-auth-key',
  redactionRules: ['password', 'card', 'ssn']
});
```

#### Backend (Node.js)
```javascript
const { securityMiddleware } = require('@security-platform/node-agent');

app.use(securityMiddleware({
  controlPlaneUrl: 'https://console.example.com',
  authKey: 'your-auth-key'
}));
```

#### Backend (Python)
```python
from security_platform import SecurityMiddleware

app.add_middleware(SecurityMiddleware,
    control_plane_url='https://console.example.com',
    auth_key='your-auth-key'
)
```

## Deployment Modes

1. **Cloud SaaS**: Telemetry to cloud control plane
2. **Hybrid**: Control plane cloud, telemetry on-prem
3. **Fully On-Prem**: All components within customer network

## Roadmap

### Phase 0 — MVP (Current)
- Control-plane proto (dashboard + API)
- OpenTelemetry collector + OpenSearch backend
- Frontend SDK (rrweb + opentelemetry) + FastAPI agent middleware
- ModSecurity container + Kong plugin + simple orchestrator
- Builder CLI skeleton that can produce a deb and an npm package

### Phase 1 — Expand agents & packaging
- Java/.NET/Go agents, mobile SDKs, signed installers
- Homebrew/choco packages, airgap bundle generator & Helm charts
- Kubernetes Operator

### Phase 2 — Enterprise features
- Full session replay redaction, UEBA/ML models
- SOAR integrations, on-prem installer UI, managed upgrades

## License

See LICENSE file for details.
