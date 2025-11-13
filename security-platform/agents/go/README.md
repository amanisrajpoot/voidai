# Security Platform Go Agent

Security observability agent for Go applications with OpenTelemetry integration.

## Installation

```bash
go get github.com/securityplatform/go-agent
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "github.com/securityplatform/go-agent"
)

func main() {
    // Create agent from default config (uses environment variables)
    ag, err := agent.NewAgentFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    // Start agent
    ctx := context.Background()
    if err := ag.Start(ctx); err != nil {
        log.Fatal(err)
    }
    defer ag.Stop(ctx)

    // Redact PII
    redacted := ag.Redact(sensitiveData)

    // Check policy
    result := ag.CheckPolicy("action", map[string]interface{}{
        "user": "alice",
        "resource": "data",
    })
    if !result.Allowed {
        // Handle blocked action
    }
}
```

### HTTP Server Integration

```go
package main

import (
    "net/http"
    "github.com/securityplatform/go-agent"
)

func main() {
    ag, _ := agent.NewAgentFromEnv()
    ctx := context.Background()
    ag.Start(ctx)
    defer ag.Stop(ctx)

    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello"))
    })

    // Use agent middleware
    handler := ag.Middleware(mux)
    http.ListenAndServe(":8080", handler)
}
```

### Gin Integration

```go
import (
    "github.com/gin-gonic/gin"
    "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
    ag, _ := agent.NewAgentFromEnv()
    ctx := context.Background()
    ag.Start(ctx)
    defer ag.Stop(ctx)

    r := gin.Default()
    r.Use(otelgin.Middleware("my-service"))
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Hello"})
    })
    r.Run()
}
```

## Configuration

Create a `config.yaml` file:

```yaml
control_plane_url: "https://api.securityplatform.com"
auth_key: "${SECURITY_PLATFORM_AUTH_KEY}"
service_name: "my-go-service"
environment: "production"
otlp_endpoint: "http://localhost:4318"
telemetry:
  batch_size: 100
  batch_timeout: "5s"
  export_timeout: "30s"
  max_queue_size: 2048
policy:
  mode: "observe"  # or "block"
  auto_enable_blocking: false
redaction_rules:
  - pattern: ".*password.*"
  - pattern: ".*token.*"
```

Load from file:

```go
ag, err := agent.NewAgentFromFile("config.yaml")
```

## Environment Variables

- `SECURITY_PLATFORM_URL` - Control plane URL
- `SECURITY_PLATFORM_AUTH_KEY` - Authentication key
- `SERVICE_NAME` - Service name
- `ENVIRONMENT` - Environment (production, staging, development)

## Features

- ✅ OpenTelemetry integration
- ✅ HTTP middleware
- ✅ Gin framework support
- ✅ PII redaction
- ✅ Policy engine (observe/block modes)
- ✅ YAML configuration support
