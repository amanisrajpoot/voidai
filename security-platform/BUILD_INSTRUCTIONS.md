# Build Instructions

## Prerequisites

- **Go 1.21+** - For builder CLI and Go services
- **Node.js 18+** - For frontend SDK and Node.js agents
- **Python 3.8+** - For Python agents
- **Docker & Docker Compose** - For local development
- **Build Tools** (optional, for packaging):
  - `dpkg-deb` (for DEB packages)
  - `rpmbuild` (for RPM packages)
  - WiX Toolset (for MSI packages on Windows)
  - `hdiutil` (for DMG packages on macOS)

## Building the Builder CLI

```bash
cd security-platform/builder
go mod download
go build -o builder .
```

The `builder` binary will be created in the current directory.

## Building Frontend SDK

```bash
cd security-platform/sdks/frontend
npm install
npm run build
```

Outputs:
- `dist/index.js` - CommonJS
- `dist/index.esm.js` - ES Module
- `dist/index.umd.js` - UMD bundle
- `dist/index.d.ts` - TypeScript definitions

## Building Node.js Agent

```bash
cd security-platform/agents/node
npm install
npm run build
```

## Building Python Agent

```bash
cd security-platform/agents/python
pip install -e .
# Or build distribution packages:
python setup.py sdist bdist_wheel
```

## Building Control Plane API

```bash
cd security-platform/control-plane/api
go mod download
go build -o control-plane .
```

## Building Orchestrator

```bash
cd security-platform/orchestrator
go mod download
go build -o orchestrator .
```

## Building Docker Images

### Control Plane
```bash
cd security-platform/control-plane
docker build -t security-platform/control-plane:latest .
```

### Orchestrator
```bash
cd security-platform/orchestrator
docker build -t security-platform/orchestrator:latest .
```

### WAF
```bash
cd security-platform/edge/waf
docker build -t security-platform/waf:latest .
```

## Using the Builder CLI

### Initialize a Project
```bash
./builder init my-agent --lang=node
# or
./builder init my-agent --lang=python
```

### Build an Agent
```bash
./builder build agent --lang=node --version=1.0.0 --out=./dist
./builder build agent --lang=python --version=1.0.0 --out=./dist
```

### Package for Distribution
```bash
# Single format
./builder package --target=deb --version=1.0.0 --out=./dist

# Multiple formats
./builder package --target=deb,rpm,helm --version=1.0.0 --out=./dist
```

### Sign Artifacts
```bash
# GPG signing (Linux packages)
./builder sign --artifact ./dist/agent.deb --method=gpg --key=your-key-id

# macOS codesign
./builder sign --artifact ./dist/agent.dmg --method=codesign --key="Developer ID"

# Windows signtool
./builder sign --artifact ./dist/agent.msi --method=signtool --key=./cert.pfx
```

### Upload to Repositories
```bash
# npm
./builder upload --artifact ./dist/*.tgz --repo=npmjs.com --type=npm

# PyPI
./builder upload --artifact ./dist/*.whl --repo=pypi.org --type=pypi

# S3
./builder upload --artifact ./dist/*.deb --repo=s3://bucket/packages --type=s3
```

## Local Development Setup

### Start All Services
```bash
cd security-platform
docker-compose up -d
```

This starts:
- Control Plane API: http://localhost:8080
- OpenTelemetry Collector: http://localhost:4318
- OpenSearch: http://localhost:9200
- MinIO: http://localhost:9000 (console: 9001)
- Orchestrator: http://localhost:8081
- WAF: http://localhost:80

### Verify Services
```bash
# Control Plane health
curl http://localhost:8080/health

# Orchestrator health
curl http://localhost:8081/health

# OpenSearch health
curl http://localhost:9200/_cluster/health
```

## Testing

### Frontend SDK
```bash
cd security-platform/sdks/frontend
npm test
```

### Node.js Agent
```bash
cd security-platform/agents/node
npm test
```

### Python Agent
```bash
cd security-platform/agents/python
pytest
```

## Production Builds

### Using GitHub Actions
The `.github/workflows/build.yml` workflow automatically:
1. Builds the builder CLI
2. Builds all agents
3. Creates packages
4. Signs artifacts (on release)
5. Uploads to GitHub Releases

### Manual Production Build
```bash
# Set version
export VERSION=1.0.0

# Build all components
cd security-platform/builder && go build -o ../../builder .
cd ../../

# Build agents
./builder build agent --lang=node --version=$VERSION
./builder build agent --lang=python --version=$VERSION

# Create packages
./builder package --target=deb,rpm,msi,dmg,helm --version=$VERSION

# Sign (requires signing keys)
./builder sign --artifact ./dist/*.deb --method=gpg --key=$GPG_KEY_ID
./builder sign --artifact ./dist/*.dmg --method=codesign --key="$CODESIGN_IDENTITY"

# Upload (requires credentials)
./builder upload --artifact ./dist/*.tgz --repo=npmjs.com --type=npm
./builder upload --artifact ./dist/*.whl --repo=pypi.org --type=pypi
```

## Troubleshooting

### Go Module Issues
```bash
go clean -modcache
go mod download
```

### Node.js Build Issues
```bash
rm -rf node_modules package-lock.json
npm install
```

### Python Build Issues
```bash
pip install --upgrade pip setuptools wheel
pip install -e .
```

### Docker Issues
```bash
docker-compose down -v
docker-compose up --build
```

## Cross-Platform Building

### Using goreleaser (Recommended)
Create `.goreleaser.yml`:
```yaml
builds:
  - main: ./builder/main.go
    binary: builder
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
```

Run:
```bash
goreleaser release --snapshot
```

### Manual Cross-Compilation
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o builder-linux-amd64 .

# macOS
GOOS=darwin GOARCH=amd64 go build -o builder-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o builder-darwin-arm64 .

# Windows
GOOS=windows GOARCH=amd64 go build -o builder-windows-amd64.exe .
```

## CI/CD Integration

### GitHub Actions
See `.github/workflows/build.yml` for complete example.

### GitLab CI
```yaml
build:
  stage: build
  script:
    - cd builder && go build -o builder .
    - ./builder build agent --lang=node --version=$CI_COMMIT_TAG
    - ./builder package --target=deb,rpm --version=$CI_COMMIT_TAG
```

### Jenkins
```groovy
stage('Build') {
    sh 'cd builder && go build -o builder .'
    sh './builder build agent --lang=node --version=${VERSION}'
    sh './builder package --target=deb,rpm --version=${VERSION}'
}
```

## Next Steps

1. Review `docs/QUICKSTART.md` for usage examples
2. Check `IMPLEMENTATION_PLAN.md` for roadmap
3. Explore `docs/examples/` for code samples
4. Read `SUMMARY.md` for overview
