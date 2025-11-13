# Security Platform .NET Agent

Security observability agent for .NET applications with OpenTelemetry integration.

## Installation

### NuGet Package Manager

```powershell
Install-Package SecurityPlatform.Agent
```

### .NET CLI

```bash
dotnet add package SecurityPlatform.Agent
```

### PackageReference

```xml
<PackageReference Include="SecurityPlatform.Agent" Version="1.0.0" />
```

## Quick Start

### Basic Usage

```csharp
using SecurityPlatform.Agent;

// Create agent from default config (uses environment variables)
var agent = AgentFactory.CreateDefault();
agent.Start();

// Or load from YAML config file
var agent = AgentFactory.FromYamlFile("/path/to/config.yaml");
agent.Start();

// Redact PII
var redacted = agent.Redact(sensitiveData);

// Check policy
var result = agent.CheckPolicy("action", context);
if (!result.Allowed)
{
    // Handle blocked action
}

// Stop agent on shutdown
agent.Stop();
```

### ASP.NET Core Integration

```csharp
using SecurityPlatform.Agent;
using SecurityPlatform.Agent.Extensions;

// In Program.cs or Startup.cs
var builder = WebApplication.CreateBuilder(args);

// Add Security Platform Agent
builder.Services.AddSecurityPlatformAgent(builder.Configuration);

var app = builder.Build();

// Use Security Platform middleware
app.UseSecurityPlatform();

app.Run();
```

### Configuration (appsettings.yaml)

```yaml
SecurityPlatform:
  control_plane_url: "https://api.securityplatform.com"
  auth_key: "${SECURITY_PLATFORM_AUTH_KEY}"
  service_name: "my-dotnet-service"
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

## Environment Variables

- `SECURITY_PLATFORM_URL` - Control plane URL
- `SECURITY_PLATFORM_AUTH_KEY` - Authentication key
- `SERVICE_NAME` - Service name
- `ENVIRONMENT` - Environment (production, staging, development)

## Features

- ✅ OpenTelemetry integration
- ✅ Automatic instrumentation (ASP.NET Core, HttpClient)
- ✅ PII redaction
- ✅ Policy engine (observe/block modes)
- ✅ ASP.NET Core middleware
- ✅ Dependency injection support
