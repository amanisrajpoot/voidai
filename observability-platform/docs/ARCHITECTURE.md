# Architecture Overview

## System Components

### 1. Control Plane / Console
- **Purpose**: Centralized management UI and API
- **Tech**: Go (Gin) REST API + React frontend (future)
- **Features**:
  - Incident management
  - Session replay viewer
  - Trace analysis
  - Rule management
  - Playbook orchestration
- **Deployment**: Docker container or Kubernetes deployment

### 2. Telemetry Bus
- **OTEL Collector**: Receives OTLP traces/metrics/logs from agents
- **Storage**: OpenSearch for structured data, MinIO for session blobs
- **Processing**: Batch processing, redaction, enrichment

### 3. Agents & SDKs

#### Backend Agents (RASP-lite)
- **Node.js**: Express/Koa middleware, OpenTelemetry SDK
- **Python**: FastAPI/Django middleware, OpenTelemetry SDK
- **Java**: Spring Boot agent, OpenTelemetry Java agent
- **.NET**: ASP.NET Core middleware, OpenTelemetry .NET SDK
- **Go**: HTTP middleware, OpenTelemetry Go SDK
- **Ruby**: Rails middleware, OpenTelemetry Ruby SDK

#### Frontend SDK
- **TypeScript/JavaScript**: Browser SDK with rrweb for session replay
- **Features**: RUM, session capture, automatic redaction

#### Mobile SDKs
- **iOS**: Swift SDK with CocoaPods/SPM
- **Android**: Kotlin/Java SDK with Gradle

#### Desktop Agents
- **Windows**: C# Windows Service
- **macOS**: Swift daemon with launchd
- **Linux**: Go daemon with systemd

### 4. Edge Modules
- **WAF**: nginx + ModSecurity + OWASP CRS
- **Kong Plugin**: Lua plugin for Kong API Gateway
- **Envoy**: HTTP filter for Envoy proxy

### 5. Orchestrator
- **Purpose**: SOAR-lite automation engine
- **Features**: Webhook-driven playbooks, integrations with Kong, Kubernetes, IAM

### 6. Builder CLI
- **Purpose**: Single tool to build, package, sign, and distribute all components
- **Tech**: Go with Cobra CLI framework
- **Features**:
  - Agent scaffolding
  - Multi-platform builds
  - Package generation (deb/rpm/msi/dmg/helm/ova)
  - Code signing
  - Repository uploads
  - Airgap bundle generation

## Data Flow

```
[Agents/SDKs] 
    ↓ (OTLP)
[OTEL Collector]
    ↓
[OpenSearch] ← [Control Plane] → [MinIO (Sessions)]
    ↓
[Dashboards/UI]
```

## Security Model

1. **mTLS**: All agent-to-control-plane communication uses mTLS
2. **Signing**: All packages signed (GPG for Linux, codesign for macOS, signtool for Windows)
3. **Redaction**: Default PII redaction rules, configurable per agent
4. **RBAC**: Role-based access control in control plane
5. **Audit Logs**: Immutable audit trail for all operations

## Deployment Modes

1. **Cloud SaaS**: Fully managed cloud deployment
2. **Hybrid**: Control plane in cloud, telemetry on-prem
3. **Fully On-Prem**: Complete on-premises installation (airgap support)

## Compliance

- **Data Models**: SaaS vs On-Prem data residency
- **Encryption**: AES-256 at rest, TLS 1.3 in transit
- **GDPR/CCPA**: Data deletion API, TTL policies, export tools
- **SBOM**: Software Bill of Materials for supply chain security
