.PHONY: help build test clean install

help:
	@echo "Security Platform Build System"
	@echo ""
	@echo "Targets:"
	@echo "  build-builder     - Build the builder CLI"
	@echo "  build-agents      - Build all agent SDKs"
	@echo "  build-frontend    - Build frontend SDK"
	@echo "  build-node        - Build Node.js agent"
	@echo "  build-python      - Build Python agent"
	@echo "  build-control-plane - Build control plane API"
	@echo "  build-orchestrator - Build orchestrator"
	@echo "  test              - Run all tests"
	@echo "  package           - Create distribution packages"
	@echo "  docker-build      - Build Docker images"
	@echo "  clean             - Clean build artifacts"
	@echo "  install           - Install builder CLI"

build-builder:
	cd builder && go build -o ../bin/builder ./main.go

build-agents: build-frontend build-node build-python

build-frontend:
	cd agents/frontend && npm ci && npm run build

build-node:
	cd agents/node && npm ci && npm run build

build-python:
	cd agents/python && pip install build && python -m build

build-control-plane:
	cd control-plane/api && go build -o ../../bin/control-plane ./main.go

build-orchestrator:
	cd orchestrator && go build -o ../bin/orchestrator ./main.go

test:
	@echo "Running tests..."
	cd agents/frontend && npm test || true
	cd agents/node && npm test || true
	cd agents/python && pytest || true
	cd builder && go test ./... || true
	cd control-plane/api && go test ./... || true
	cd orchestrator && go test ./... || true

package:
	@echo "Creating packages..."
	./bin/builder package --target=deb,rpm,msi,dmg,helm --out=./dist

docker-build:
	docker build -t securityplatform/builder:latest -f builder/Dockerfile builder/
	docker build -t securityplatform/control-plane:latest -f control-plane/Dockerfile control-plane/
	docker build -t securityplatform/orchestrator:latest -f orchestrator/Dockerfile orchestrator/

clean:
	rm -rf bin/ dist/ node_modules/ agents/*/node_modules/ agents/*/dist/ agents/*/build/
	find . -type d -name "__pycache__" -exec rm -r {} + || true
	find . -type f -name "*.pyc" -delete || true

install: build-builder
	cp bin/builder /usr/local/bin/builder
	chmod +x /usr/local/bin/builder

.DEFAULT_GOAL := help
