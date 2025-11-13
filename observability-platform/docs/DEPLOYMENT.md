# Deployment Guide

## Quick Start (Local Development)

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for building)
- Node.js 20+ (for frontend SDK)

### Start Control Plane Stack

```bash
cd control-plane
docker-compose up -d
```

This starts:
- OpenSearch (port 9200)
- OpenSearch Dashboards (port 5601)
- OTEL Collector (ports 4317, 4318)
- MinIO (ports 9000, 9001)
- Control Plane API (port 8080)

### Install Agent (Node.js Example)

```bash
# Initialize agent project
./builder init my-agent --lang=node

# Build agent
./builder build agent --lang=node --version=1.0.0

# Install in your Node.js app
cd my-agent
npm install
npm run build
npm link  # or npm publish
```

### Configure Agent

Create `config.yaml`:

```yaml
control_plane_url: http://localhost:8080
auth_key: your-auth-key-here
service_name: my-service
environment: development
otlp_endpoint: http://localhost:4318
redaction_rules:
  - password
  - token
  - api_key
```

## Production Deployment

### Kubernetes Deployment

```bash
# Install Helm chart
helm install observability-platform ./helm/observability-platform \
  --set controlPlane.url=https://control.example.com \
  --set opensearch.replicas=3 \
  --set minio.persistence.enabled=true
```

### Airgap Installation

```bash
# Generate airgap bundle
./airgap/generate.sh 1.0.0

# Transfer bundle to target system
scp airgap-bundle-1.0.0.tar.gz user@target-server:/tmp/

# On target system
tar -xzf airgap-bundle-1.0.0.tar.gz
cd airgap-bundle-1.0.0
./scripts/install.sh
```

### Package Installation

#### Debian/Ubuntu
```bash
sudo dpkg -i observability-agent_1.0.0_amd64.deb
sudo systemctl enable observability-agent
sudo systemctl start observability-agent
```

#### RedHat/CentOS
```bash
sudo rpm -i observability-agent-1.0.0-1.x86_64.rpm
sudo systemctl enable observability-agent
sudo systemctl start observability-agent
```

#### macOS
```bash
# Install via Homebrew
brew install observability-platform/observability-agent

# Or install DMG
open observability-agent-1.0.0.dmg
```

#### Windows
```bash
# Install via Chocolatey
choco install observability-agent

# Or install MSI
msiexec /i observability-agent-1.0.0.msi
```

## Configuration

### Agent Configuration

All agents support YAML or JSON configuration:

```yaml
control_plane_url: https://control.example.com
auth_key: ${AUTH_KEY}  # Supports env var substitution
service_name: my-service
environment: production
otlp_endpoint: https://otel.example.com:4318
telemetry_batch_size: 100
redaction_rules:
  - password
  - passwd
  - token
  - api_key
  - credit_card
local_policy:
  mode: observe  # observe | block | alert
  rules: []
```

### Control Plane Configuration

Environment variables:
- `OPENSEARCH_URL`: OpenSearch endpoint
- `MINIO_ENDPOINT`: MinIO endpoint
- `JWT_SECRET`: JWT signing secret
- `PORT`: API server port (default: 8080)

## Upgrades

### Automatic Updates

Agents support automatic updates via signed manifests:

```yaml
auto_update:
  enabled: true
  check_interval: 1h
  channel: stable  # stable | beta | alpha
```

### Manual Updates

```bash
# Download new version
wget https://releases.example.com/observability-agent-1.1.0.deb

# Verify signature
dpkg-sig --verify observability-agent-1.1.0.deb

# Upgrade
sudo dpkg -i observability-agent-1.1.0.deb
sudo systemctl restart observability-agent
```

## Monitoring

### Health Checks

- Control Plane: `GET /api/v1/health`
- Agent: `GET http://localhost:9090/health` (if health endpoint enabled)

### Metrics

Prometheus metrics available at:
- Control Plane: `http://localhost:8080/metrics`
- OTEL Collector: `http://localhost:8888/metrics`

## Troubleshooting

### Agent Not Sending Data

1. Check agent configuration
2. Verify network connectivity to OTEL collector
3. Check agent logs: `journalctl -u observability-agent` (Linux)

### Control Plane Not Receiving Data

1. Verify OTEL collector is running
2. Check OpenSearch connectivity
3. Review collector logs: `docker logs otel-collector`
