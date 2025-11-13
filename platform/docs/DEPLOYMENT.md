# Deployment Guide

## Local Development

### Prerequisites
- Docker and Docker Compose
- Node.js 20+
- Go 1.21+
- Python 3.8+ (for Python agent)

### Quick Start

```bash
# Start all services
docker-compose -f docker-compose.dev.yml up -d

# Check services
docker-compose -f docker-compose.dev.yml ps

# View logs
docker-compose -f docker-compose.dev.yml logs -f
```

### Services
- Control Plane API: http://localhost:3000
- OpenSearch: http://localhost:9200
- OpenSearch Dashboards: http://localhost:5601
- MinIO Console: http://localhost:9001 (minioadmin/minioadmin)
- Kong Admin: http://localhost:8001
- WAF: http://localhost:8080

## Production Deployment

### Control Plane (Kubernetes)

```bash
# Install Helm chart
helm install security-platform ./packaging/helm/security-platform-chart \
  --set controlPlane.url=https://api.example.com \
  --set storage.opensearch.enabled=true \
  --set storage.minio.enabled=true

# Or use the operator
kubectl apply -f ./packaging/kubernetes/operator.yaml
```

### On-Premises Installation

#### Air-Gap Bundle

```bash
# Generate air-gap bundle
./builder airgap --version=1.0.0 --out=./airgap-bundle

# Transfer bundle to air-gapped environment
# Extract and run installer
cd airgap-bundle
./install.sh
```

#### Docker Compose (On-Prem)

```bash
# Copy docker-compose.yml to server
scp docker-compose.prod.yml user@server:/opt/security-platform/

# SSH into server
ssh user@server
cd /opt/security-platform
docker-compose -f docker-compose.prod.yml up -d
```

### Agent Installation

#### Linux (deb)
```bash
wget https://releases.example.com/security-agent_1.0.0_amd64.deb
sudo dpkg -i security-agent_1.0.0_amd64.deb
sudo systemctl enable security-agent
sudo systemctl start security-agent
```

#### Linux (rpm)
```bash
wget https://releases.example.com/security-agent-1.0.0-1.x86_64.rpm
sudo rpm -i security-agent-1.0.0-1.x86_64.rpm
sudo systemctl enable security-agent
sudo systemctl start security-agent
```

#### macOS
```bash
# Homebrew
brew tap security-platform/tap
brew install security-agent

# Or download DMG
open https://releases.example.com/security-agent-1.0.0.dmg
```

#### Windows
```powershell
# Chocolatey
choco install security-agent

# Or download MSI
# Run installer GUI
```

#### Python Agent
```bash
pip install security-platform-agent
```

#### Node.js Agent
```bash
npm install @security-platform/node-agent
```

## Configuration

### Agent Configuration

Create `agent.yaml`:

```yaml
control_plane_url: https://api.example.com
auth_key: your-bootstrap-token
service_name: my-service
env: production
redaction_rules:
  - password
  - card
  - ssn
local_policy:
  mode: observe
telemetry_batch_size: 100
otlp_endpoint: http://localhost:4318
offline_mode: false
```

### Environment Variables

```bash
export SECURITY_PLATFORM_CONTROL_PLANE_URL=https://api.example.com
export SECURITY_PLATFORM_AUTH_KEY=your-key
export SECURITY_PLATFORM_SERVICE_NAME=my-service
export SECURITY_PLATFORM_ENV=production
```

## Upgrades

### Agent Auto-Update

Agents check for updates every 24 hours by default. To force update:

```bash
# Linux
sudo security-agent update

# Or restart service to pick up new version
sudo systemctl restart security-agent
```

### Control Plane Upgrade

```bash
# Kubernetes
helm upgrade security-platform ./packaging/helm/security-platform-chart \
  --set image.tag=v1.1.0

# Docker Compose
docker-compose pull
docker-compose up -d
```

## Monitoring

### Health Checks

```bash
# Control plane
curl http://localhost:3000/health

# Agent
curl http://localhost:9090/health
```

### Metrics

- Prometheus metrics: http://localhost:8888/metrics
- OpenSearch metrics: http://localhost:9200/_prometheus/metrics

## Troubleshooting

### Agent Not Connecting

1. Check network connectivity:
   ```bash
   curl https://api.example.com/health
   ```

2. Verify auth key:
   ```bash
   security-agent config show
   ```

3. Check logs:
   ```bash
   journalctl -u security-agent -f  # Linux
   log show --predicate 'process == "security-agent"' --last 1h  # macOS
   ```

### Control Plane Issues

1. Check service status:
   ```bash
   docker-compose ps
   ```

2. View logs:
   ```bash
   docker-compose logs control-plane
   ```

3. Check OpenSearch:
   ```bash
   curl http://localhost:9200/_cluster/health
   ```
