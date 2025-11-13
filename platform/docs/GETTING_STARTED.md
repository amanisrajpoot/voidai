# Getting Started

This guide will help you get started with the Security Observability Platform.

## Prerequisites

- Docker and Docker Compose
- Node.js 20+ (for frontend SDK and control plane)
- Go 1.21+ (for builder CLI)
- Python 3.8+ (for Python agent)

## Quick Start

### 1. Start Local Development Environment

```bash
cd platform
docker-compose -f docker-compose.dev.yml up -d
```

This starts:
- Control Plane API (port 3000)
- OpenTelemetry Collector (ports 4317, 4318)
- OpenSearch (port 9200)
- OpenSearch Dashboards (port 5601)
- MinIO (ports 9000, 9001)
- Kong API Gateway (ports 8000, 8001)
- ModSecurity WAF (port 8080)

### 2. Build the Builder CLI

```bash
cd platform/builder
go mod download
go build -o ../../bin/builder ./main.go
```

### 3. Initialize an Agent Project

```bash
./bin/builder init my-python-agent --lang=python
cd my-python-agent
```

Edit `agent.yaml` with your configuration:

```yaml
control_plane_url: http://localhost:3000
auth_key: demo-key
service_name: my-service
env: development
```

### 4. Build and Test Python Agent

```bash
# Install Python agent locally
cd ../agents/backend/python
pip install -e .

# Run example
cd ../../examples/python-fastapi
pip install -r requirements.txt
python main.py
```

### 5. Test Frontend SDK

```bash
# Build frontend SDK
cd ../../agents/frontend
npm install
npm run build

# Open example
open ../../examples/frontend-quickstart/index.html
```

## Next Steps

### Create Your First Rule

```bash
curl -X POST http://localhost:3000/api/v1/rules \
  -H "Authorization: Bearer demo-key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Block Admin Access",
    "pattern": ".*/admin.*",
    "action": "block",
    "enabled": true
  }'
```

### View Events

```bash
curl http://localhost:3000/api/v1/events \
  -H "Authorization: Bearer demo-key"
```

### View in OpenSearch Dashboards

1. Open http://localhost:5601
2. Create index pattern: `events-*`
3. Explore your events

## Production Deployment

See [DEPLOYMENT.md](./DEPLOYMENT.md) for production deployment instructions.

## Architecture

See [ARCHITECTURE.md](./ARCHITECTURE.md) for detailed architecture documentation.

## API Reference

See [API.md](./API.md) for complete API documentation.

## Support

- Documentation: [docs/](./docs/)
- Examples: [examples/](./examples/)
- Issues: GitHub Issues
