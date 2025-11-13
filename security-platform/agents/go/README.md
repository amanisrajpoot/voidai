# Security Platform Agent - Go

Go agent for Security Platform with OpenTelemetry integration.

## Installation

```bash
go get github.com/security-platform/go-agent
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "net/http"
    
    "github.com/security-platform/go-agent"
)

func main() {
    // Create configuration
    config := &agent.Config{
        ControlPlaneURL: "https://api.securityplatform.com",
        AuthKey:         os.Getenv("SECURITY_PLATFORM_AUTH_KEY"),
        ServiceName:     "my-go-service",
        Environment:     "production",
    }

    // Create and start agent
    ag, err := agent.NewAgent(config)
    if err != nil {
        log.Fatal(err)
    }

    if err := ag.Start(); err != nil {
        log.Fatal(err)
    }
    defer ag.Stop()

    // Use tracer
    ctx := context.Background()
    _, span := ag.Tracer().Start(ctx, "my-operation")
    defer span.End()

    // HTTP server with middleware
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello"))
    })

    handler := ag.HTTPMiddleware(mux)
    log.Fatal(http.ListenAndServe(":8080", handler))
}
```

## Configuration

See `agent-config.yaml` for full configuration options.
