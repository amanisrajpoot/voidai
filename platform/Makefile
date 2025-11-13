.PHONY: help build test clean dev docker-build docker-up docker-down

help:
	@echo "Available targets:"
	@echo "  build          - Build all components"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  dev            - Start development environment"
	@echo "  docker-build   - Build Docker images"
	@echo "  docker-up      - Start Docker services"
	@echo "  docker-down    - Stop Docker services"

build:
	@echo "Building Builder CLI..."
	cd builder && go build -o ../bin/builder ./main.go
	@echo "Building Frontend SDK..."
	cd agents/frontend && npm install && npm run build
	@echo "Building Control Plane..."
	cd control-plane && npm install && npm run build

test:
	@echo "Testing Builder CLI..."
	cd builder && go test ./...
	@echo "Testing Frontend SDK..."
	cd agents/frontend && npm test
	@echo "Testing Control Plane..."
	cd control-plane && npm test

clean:
	rm -rf bin/
	rm -rf agents/frontend/dist/
	rm -rf agents/frontend/node_modules/
	rm -rf control-plane/dist/
	rm -rf control-plane/node_modules/
	find . -name "*.pyc" -delete
	find . -name "__pycache__" -type d -exec rm -r {} +

dev:
	docker-compose -f docker-compose.dev.yml up

docker-build:
	docker-compose -f docker-compose.dev.yml build

docker-up:
	docker-compose -f docker-compose.dev.yml up -d

docker-down:
	docker-compose -f docker-compose.dev.yml down

builder-init:
	./bin/builder init my-agent --lang=python

builder-build:
	./bin/builder build agent --lang=python --version=0.1.0

builder-package:
	./bin/builder package --target=deb,rpm --version=0.1.0 --out=./dist
