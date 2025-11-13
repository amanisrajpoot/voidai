# Security Platform Go Agent

Simple, automatic security observability for Go applications.

## Quick Start

### Install

```bash
go get github.com/security-platform/go-agent
```

## Usage

### Basic Integration

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/security-platform/go-agent"
)

func main() {
    // Initialize agent
    config := &agent.Config{
        ControlPlaneURL: "https://api.securityplatform.com",
        AuthKey:         os.Getenv("SECURITY_PLATFORM_AUTH_KEY"),
        ServiceName:     "my-go-service",
        Environment:     "production",
    }
    
    ag := agent.NewAgent(config)
    if err := ag.Start(); err != nil {
        log.Fatalf("Failed to start agent: %v", err)
    }
    defer ag.Stop()
    
    // Your application code here...
}
```

### HTTP Server Integration

```go
package main

import (
    "net/http"
    
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    "github.com/security-platform/go-agent"
)

func main() {
    ag := agent.NewAgent(&agent.Config{
        ServiceName: "my-http-service",
    })
    ag.Start()
    defer ag.Stop()
    
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    http.Handle("/", otelhttp.NewHandler(handler, "http-server"))
    http.ListenAndServe(":8080", nil)
}
```

### Environment Variables

```bash
export SECURITY_PLATFORM_CONTROL_PLANE_URL=https://api.securityplatform.com
export SECURITY_PLATFORM_AUTH_KEY=your-token-here
```

## Features

- ✅ Automatic OpenTelemetry instrumentation
- ✅ PII redaction
- ✅ Policy enforcement
- ✅ Simple configuration
- ✅ Zero manual instrumentation required

## Documentation

See [docs/DEPLOYMENT.md](../../docs/DEPLOYMENT.md) for full documentation.
