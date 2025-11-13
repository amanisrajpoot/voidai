# Architecture Overview

## System Components

### 1. Control Plane
- **Purpose**: Central management console and API
- **Tech**: Node.js/Express, TypeScript
- **Features**:
  - Event ingestion and storage
  - Rule management
  - Action orchestration
  - Dashboard API

### 2. Telemetry Bus
- **OpenTelemetry Collector**: Receives OTLP traces/metrics/logs
- **Storage**: OpenSearch/ClickHouse for events, MinIO/S3 for session blobs
- **Protocol**: OTLP (HTTP/gRPC)

### 3. Agents & SDKs

#### Frontend SDK
- **Tech**: TypeScript, rrweb, OpenTelemetry JS
- **Features**: Session replay, RUM, automatic PII redaction

#### Backend Agents
- **Python**: FastAPI/Starlette middleware, OTLP exporter
- **Node.js**: Express/Koa middleware
- **Java**: Spring Boot agent, Java agent
- **.NET**: ASP.NET Core middleware
- **Go**: HTTP middleware, OTLP exporter

#### Mobile SDKs
- **iOS**: Swift package, CocoaPods/SPM
- **Android**: Kotlin/Java, Gradle/Maven

#### Desktop Agents
- **Windows**: C# service, MSI installer
- **macOS**: Swift daemon, DMG installer
- **Linux**: Go binary, systemd service, deb/rpm

### 4. Edge Modules
- **WAF**: ModSecurity + nginx, OWASP CRS rules
- **Kong Plugin**: Admin API integration
- **Envoy**: WASM filter or HTTP filter

### 5. Builder CLI
- **Tech**: Go, Cobra CLI framework
- **Features**:
  - Agent scaffolding
  - Multi-language builds
  - Cross-platform packaging
  - Code signing
  - Air-gap bundle generation

## Data Flow

```
[Frontend/Backend/Mobile] 
    ↓ (OTLP/Events)
[OTEL Collector]
    ↓
[OpenSearch/ClickHouse] ← [Control Plane] → [MinIO/S3]
    ↓
[Dashboard/Console]
    ↓
[Rules Engine] → [Orchestrator] → [Actions: Block IP, Revoke Token, Alert]
```

## Security Model

1. **Authentication**: JWT tokens, mTLS for agent → control plane
2. **Redaction**: Configurable rules, automatic PII detection
3. **Signing**: All packages signed (GPG/codesign/signtool)
4. **RBAC**: Role-based access control in control plane
5. **Audit**: Immutable event trails

## Deployment Modes

### Cloud SaaS
- Control plane: Cloud-hosted
- Telemetry: Cloud storage
- Agents: Connect to cloud endpoint

### Hybrid
- Control plane: Cloud-hosted
- Telemetry: On-prem storage
- Agents: Connect to cloud, store locally

### Fully On-Prem
- Control plane: On-prem deployment
- Telemetry: On-prem storage
- Agents: Connect to on-prem endpoint
- Air-gap support: Offline bundle installation

## Packaging & Distribution

### Package Formats
- **Linux**: deb (Debian/Ubuntu), rpm (RHEL/CentOS), Docker, Helm
- **macOS**: dmg, Homebrew tap
- **Windows**: msi, Chocolatey
- **Mobile**: CocoaPods, SPM, Maven, Gradle
- **Appliance**: OVA, VM image, PXE/Ansible

### Signing
- **Linux**: GPG signatures
- **macOS**: codesign + notarization
- **Windows**: Authenticode signing

## Update Mechanism

1. Agent checks control plane for updates
2. Downloads signed manifest
3. Verifies signature
4. Downloads delta or full update
5. Applies update with rollback capability
