.PHONY: build builder agents frontend control-plane test clean

# Build builder CLI
builder:
	cd builder && go build -o ../bin/builder ./main.go

# Build all agents
agents:
	@echo "Building agents..."
	@for lang in node python go java dotnet; do \
		echo "Building $$lang agent..."; \
		./bin/builder build agent --lang=$$lang --version=1.0.0 || true; \
	done

# Build frontend SDK
frontend:
	cd sdk/frontend && npm install && npm run build

# Build control plane
control-plane:
	cd control-plane/api && go build -o ../../bin/control-plane ./cmd/api

# Build orchestrator
orchestrator:
	cd orchestrator && go build -o ../bin/orchestrator .

# Run tests
test:
	cd builder && go test ./...
	cd control-plane/api && go test ./...
	cd orchestrator && go test ./...

# Clean build artifacts
clean:
	rm -rf bin/ dist/ node_modules/ */node_modules/ */dist/ */target/ */build/

# Start local development environment
dev:
	cd control-plane && docker-compose up -d

# Stop local development environment
dev-stop:
	cd control-plane && docker-compose down

# Generate airgap bundle
airgap:
	./airgap/generate.sh 1.0.0

# Install builder
install-builder: builder
	sudo cp bin/builder /usr/local/bin/observability-builder

# All-in-one build
all: builder agents frontend control-plane orchestrator
