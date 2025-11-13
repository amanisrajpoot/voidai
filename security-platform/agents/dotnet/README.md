# Security Platform Agent - .NET

.NET agent for Security Platform with OpenTelemetry integration.

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

### ASP.NET Core

```csharp
using SecurityPlatform.Agent;

var builder = WebApplication.CreateBuilder(args);

// Configure agent
var agentConfig = new AgentConfig
{
    ControlPlaneUrl = "https://api.securityplatform.com",
    AuthKey = Environment.GetEnvironmentVariable("SECURITY_PLATFORM_AUTH_KEY"),
    ServiceName = "my-dotnet-service",
    Environment = "production"
};

// Add agent to services
builder.Services.AddSecurityPlatformAgent(agentConfig);

var app = builder.Build();

// Use agent middleware
app.UseSecurityPlatformAgent();

app.Run();
```

### Manual Usage

```csharp
using SecurityPlatform.Agent;

var config = new AgentConfig
{
    ControlPlaneUrl = "https://api.securityplatform.com",
    AuthKey = Environment.GetEnvironmentVariable("SECURITY_PLATFORM_AUTH_KEY"),
    ServiceName = "my-service",
    Environment = "production"
};

var agent = new Agent(config);
agent.Start();

// Use OpenTelemetry
using var activity = ActivitySource.StartActivity("my-operation");
// Your code here

agent.Stop();
```

## Configuration

See `agent-config.yaml` for full configuration options.
