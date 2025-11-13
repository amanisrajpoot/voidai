# Security Platform Python Agent

Security observability agent for Python applications with RASP-lite capabilities.

## Installation

```bash
pip install security-platform-agent
```

## Quick Start

### FastAPI

```python
from fastapi import FastAPI
from security_platform import init_agent, SecurityMiddleware, AgentConfig

# Initialize agent
config = AgentConfig(
    control_plane_url="https://api.example.com",
    auth_key="your-auth-key",
    service_name="my-api",
)
agent = init_agent(config)

# Create app
app = FastAPI()

# Add middleware
app.add_middleware(SecurityMiddleware)

@app.get("/")
def read_root():
    return {"message": "Hello World"}
```

### Django

```python
# settings.py
MIDDLEWARE = [
    # ... other middleware
    'security_platform.middleware.DjangoSecurityMiddleware',
]

# Configure agent
SECURITY_PLATFORM_CONTROL_PLANE_URL = "https://api.example.com"
SECURITY_PLATFORM_AUTH_KEY = "your-auth-key"
```

### Manual Usage

```python
from security_platform import init_agent, AgentConfig

config = AgentConfig.from_yaml("agent.yaml")
agent = init_agent(config)

# Capture custom events
agent.capture_event("user_action", {
    "user_id": "123",
    "action": "login",
})

# Get tracer for custom spans
tracer = agent.get_tracer("my-service")
with tracer.start_as_current_span("operation"):
    # Your code here
    pass
```

## Configuration

Create `agent.yaml`:

```yaml
control_plane_url: https://api.example.com
auth_key: your-auth-key
service_name: my-service
env: production
redaction_rules:
  - password
  - card
  - ssn
local_policy:
  mode: observe  # observe or block
telemetry_batch_size: 100
otlp_endpoint: http://localhost:4318
offline_mode: false
```

Or use environment variables:

```bash
export SECURITY_PLATFORM_CONTROL_PLANE_URL=https://api.example.com
export SECURITY_PLATFORM_AUTH_KEY=your-auth-key
export SECURITY_PLATFORM_SERVICE_NAME=my-service
```

## Features

- **Automatic Request Capture**: Middleware captures all HTTP requests/responses
- **Security Event Detection**: Detects suspicious patterns (SQL injection, path traversal, etc.)
- **PII Redaction**: Automatically redacts sensitive data
- **OpenTelemetry Integration**: Full tracing support
- **Offline Mode**: Cache events locally when control plane is unreachable
