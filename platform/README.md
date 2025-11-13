# Security Observability Platform

A comprehensive, build-ready platform for installable security observability agents/SDKs across web frontends, backends, desktop, mobile, and network appliances.

## 🚀 Quick Start

### 1. Start Local Development Environment

```bash
docker-compose -f docker-compose.dev.yml up -d
```

This starts all services:
- Control Plane API: http://localhost:3000
- OpenSearch Dashboards: http://localhost:5601
- MinIO Console: http://localhost:9001
- Kong Admin: http://localhost:8001

### 2. Build Builder CLI

```bash
cd builder
go mod download
go build -o ../bin/builder ./main.go
```

### 3. Initialize an Agent

```bash
./bin/builder init my-agent --lang=python
cd my-agent
# Edit agent.yaml with your config
```

### 4. Test Python Agent

```bash
cd examples/python-fastapi
pip install -r requirements.txt
python main.py
```

See [docs/GETTING_STARTED.md](./docs/GETTING_STARTED.md) for detailed instructions.

## 📦 What's Included

### ✅ Phase 0 MVP (Complete)

- **Builder CLI** (Go) - Scaffold, build, package, sign, upload
- **Control Plane API** (Node.js/TypeScript) - Event ingestion, rules, actions
- **Frontend SDK** (TypeScript) - Session replay, RUM, PII redaction
- **Python Agent** (FastAPI) - RASP-lite, telemetry, security detection
- **Edge Components** - ModSecurity WAF, Kong integration
- **Local Dev Environment** - Docker Compose with all services
- **Documentation** - Architecture, deployment, API reference
- **Examples** - Frontend and Python quickstarts
- **CI/CD** - GitHub Actions workflows

### ⏳ Phase 1 (Planned)

- Additional language agents (Node.js, Go, Java, .NET)
- Mobile SDKs (iOS, Android)
- Desktop agents (Windows, macOS, Linux)
- Kubernetes Operator
- Orchestrator (SOAR-lite)
- Air-gap bundle generator

## 🏗️ Architecture

```
[Agents: Frontend/Backend/Mobile/Desktop]
    ↓ (OTLP/Events)
[OpenTelemetry Collector]
    ↓
[OpenSearch/ClickHouse] ← [Control Plane] → [MinIO/S3]
    ↓
[Dashboard/Console]
    ↓
[Rules Engine] → [Orchestrator] → [Actions]
```

See [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md) for details.

## 📁 Project Structure

```
platform/
├── builder/              # Go CLI for building/packaging
├── control-plane/        # Node.js API and dashboard
├── agents/
│   ├── frontend/        # TypeScript SDK (rrweb + OTLP)
│   └── backend/
│       └── python/      # Python agent (FastAPI middleware)
├── edge/
│   └── waf/             # ModSecurity + nginx container
├── config/              # Configuration files
├── docs/                # Documentation
├── examples/            # Quickstart examples
└── docker-compose.dev.yml
```

## 🛠️ Tech Stack

- **Builder**: Go 1.21+ with Cobra CLI
- **Control Plane**: Node.js 20+, Express, TypeScript
- **Telemetry**: OpenTelemetry (OTLP)
- **Storage**: OpenSearch, MinIO (S3-compatible)
- **Frontend SDK**: TypeScript, rrweb, OpenTelemetry JS
- **Backend Agent**: Python 3.8+, FastAPI, OpenTelemetry Python
- **Packaging**: deb/rpm/msi/dmg/helm/ova
- **CI/CD**: GitHub Actions

## 📚 Documentation

- [Getting Started](./docs/GETTING_STARTED.md) - Quick start guide
- [Architecture](./docs/ARCHITECTURE.md) - System design and components
- [Deployment](./docs/DEPLOYMENT.md) - Production deployment guide
- [API Reference](./docs/API.md) - Control plane API documentation
- [Project Summary](./PROJECT_SUMMARY.md) - Complete feature list

## 🎯 Key Features

- **Multi-Platform**: Web, backend, mobile, desktop, network appliances
- **Universal Builder**: Single CLI to build and package everything
- **Session Replay**: Full session recording with PII redaction
- **RASP-Lite**: Runtime Application Self-Protection capabilities
- **OpenTelemetry**: Standards-based telemetry collection
- **Security First**: Automatic PII redaction, signing, mTLS
- **Flexible Deployment**: Cloud SaaS, Hybrid, or Fully On-Prem
- **Enterprise Ready**: RBAC, audit logs, air-gap support

## 🔒 Security

- JWT authentication for agents
- mTLS for agent → control plane communication
- Automatic PII redaction (configurable rules)
- Signed packages (GPG, codesign, Authenticode)
- RBAC for control plane operations
- Audit logs and immutable event trails

## 📦 Packaging Support

- **Linux**: deb, rpm, Docker, Helm charts
- **macOS**: dmg, Homebrew tap
- **Windows**: msi, Chocolatey
- **Mobile**: CocoaPods, SPM, Maven, Gradle
- **Appliance**: OVA, VM images, PXE/Ansible

## 🚢 Deployment Modes

1. **Cloud SaaS**: Full cloud deployment
2. **Hybrid**: Control plane cloud, telemetry on-prem
3. **Fully On-Prem**: Complete on-premises installation with air-gap support

## 📝 License

MIT License - see [LICENSE](./LICENSE) file.

## 🤝 Contributing

See [docs/GETTING_STARTED.md](./docs/GETTING_STARTED.md) for development setup.

## 📞 Support

- Documentation: [docs/](./docs/)
- Examples: [examples/](./examples/)
- Issues: GitHub Issues
