# Universal Security Observability Platform

A comprehensive security observability platform with installable agents/SDKs for web frontends, backends, desktop, mobile, and network appliances.

## Architecture Overview

### Components

1. **Control Plane / Console**: SaaS or on-prem UI for incidents, session replay, traces, rules, playbooks, and orchestration
2. **Telemetry Bus & Storage**: OTEL collector, event store (OpenSearch/ClickHouse), session storage (S3/MinIO)
3. **Edge Modules**: WAF (ModSecurity/nginx) packages, reverse-proxy integrations
4. **API Gateway Integrations**: Kong/Envoy plugins and admin connectors
5. **Language Agents / RASP SDKs**: Runtime libraries for Python, Node, Java, .NET, Go, Ruby
6. **Frontend SDK**: JS/TS SDK + rrweb integration for RUM/session capture
7. **Mobile SDKs**: iOS (Swift) and Android (Kotlin/Java) SDKs
8. **Desktop Agents**: Windows service, macOS daemon, Linux daemon
9. **Network Installer / Appliance**: Docker/VM images and Kubernetes Operator
10. **Builder / CLI**: Single CLI tool to generate all installers and packages
11. **Orchestrator**: Webhook-driven SOAR-lite runner

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

### Installing Agents

```bash
# Frontend (npm)
npm install @security-platform/frontend-sdk

# Python backend
pip install security-platform-agent

# Node.js backend
npm install @security-platform/node-agent

# Java backend
# Add to pom.xml or build.gradle

# .NET backend
dotnet add package SecurityPlatform.Agent

# Go backend
go get github.com/security-platform/go-agent
```

## Project Structure

```
security-platform/
├── builder/              # Builder CLI (Go)
├── agents/               # Agent SDKs
│   ├── frontend/        # TypeScript frontend SDK
│   ├── node/            # Node.js agent
│   ├── python/          # Python agent
│   ├── java/            # Java agent
│   ├── dotnet/          # .NET agent
│   ├── go/              # Go agent
│   ├── ios/             # iOS Swift SDK
│   ├── android/         # Android Kotlin SDK
│   ├── desktop/         # Desktop agents (Windows/macOS/Linux)
│   └── network/         # Network appliance components
├── control-plane/        # Control plane API and dashboard
├── telemetry/            # OTEL collector configs, storage
├── edge/                 # WAF, API gateway integrations
├── orchestrator/         # SOAR-lite orchestrator
├── packaging/            # Packaging scripts and templates
├── config/               # Configuration schemas
├── ci/                   # CI/CD templates
└── docs/                 # Documentation
```

## Deployment Modes

1. **Cloud SaaS**: Full cloud deployment with telemetry to cloud
2. **Hybrid**: Control plane in cloud, telemetry on-prem
3. **Fully On-Prem**: Complete on-premises deployment

## License

See LICENSE file for details.
