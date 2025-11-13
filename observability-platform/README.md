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
11. **Orchestrator**: Webhook-driven runner for automated responses

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

### Building Agents

```bash
# Initialize a new agent project
./builder init my-agent --lang=node

# Build an agent
./builder build agent --lang=node --version=1.2.0

# Package for distribution
./builder package --target=deb,rpm,msi,dmg --out=./dist

# Sign artifacts
./builder sign --artifact ./dist/agent.deb --key ./keys/private.pem

# Upload to repository
./builder upload --artifact ./dist/*.deb --repo artifactory
```

### Local Development

```bash
# Start local control plane emulator
docker-compose -f docker-compose.dev.yml up

# Run builder CLI
go run ./builder/main.go init test-agent --lang=python
```

## Deployment Modes

1. **Cloud SaaS**: Fully managed cloud deployment
2. **Hybrid**: Control plane cloud, telemetry on-prem
3. **Fully On-Prem**: Complete on-premises installation

## License

See LICENSE file for details. Community edition available under open-source license.
