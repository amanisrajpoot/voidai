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

### One-Liner Installation

**Linux/macOS:**
```bash
curl -fsSL https://install.securityplatform.com | bash
```

**Windows (PowerShell):**
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://install.securityplatform.com/install.ps1'))
```

### Installing Language Agents

```bash
# Node.js
npm install @security-platform/node-agent

# Python
pip install security-platform-agent

# Java (Maven)
# Add to pom.xml: <dependency><groupId>com.securityplatform</groupId><artifactId>security-platform-agent</artifactId><version>1.0.0</version></dependency>

# .NET
dotnet add package SecurityPlatform.Agent

# Go
go get github.com/security-platform/go-agent
```

### Installing Mobile SDKs

```bash
# iOS (CocoaPods)
pod 'SecurityPlatform', '~> 1.0'

# Android (Gradle)
implementation 'com.securityplatform:agent:1.0.0'
```

### Package Managers

```bash
# Homebrew (macOS)
brew install security-platform-agent

# Chocolatey (Windows)
choco install security-platform-agent -y

# APT (Debian/Ubuntu)
curl -fsSL https://packages.securityplatform.com/install.sh | sudo bash

# YUM/DNF (RHEL/CentOS/Fedora)
sudo yum install https://packages.securityplatform.com/rpm/security-platform-agent.rpm
```

### Kubernetes

```bash
# Helm
helm install my-agent security-platform/security-platform-agent

# Kubernetes Operator
kubectl apply -f https://raw.githubusercontent.com/security-platform/agents/main/packaging/kubernetes-operator/config/crd/bases/securityplatform.io_agents.yaml
```

For complete installation instructions, see [INSTALLATION_GUIDE.md](INSTALLATION_GUIDE.md).

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
