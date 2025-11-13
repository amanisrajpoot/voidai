# Security Platform .NET Agent

Simple, automatic security observability for .NET applications.

## Quick Start

### NuGet Package

```bash
dotnet add package SecurityPlatform.Agent
```

Or add to your `.csproj`:

```xml
<PackageReference Include="SecurityPlatform.Agent" Version="1.0.0" />
```

## Usage

### Basic Integration

```csharp
using SecurityPlatform.Agent;

// Initialize agent
var config = new AgentConfig
{
    ControlPlaneUrl = "https://api.securityplatform.com",
    AuthKey = Environment.GetEnvironmentVariable("SECURITY_PLATFORM_AUTH_KEY"),
    ServiceName = "my-dotnet-service"
};

var agent = new Agent(config);
agent.Start();

// Your application code here...

// Shutdown on application exit
agent.Dispose();
```

### ASP.NET Core Integration

```csharp
using SecurityPlatform.Agent;
using Microsoft.Extensions.DependencyInjection;

public class Startup
{
    public void ConfigureServices(IServiceCollection services)
    {
        services.AddControllers();
        
        // Add Security Platform Agent
        var agentConfig = new AgentConfig
        {
            ControlPlaneUrl = Configuration["SecurityPlatform:ControlPlaneUrl"],
            AuthKey = Configuration["SecurityPlatform:AuthKey"],
            ServiceName = "my-aspnet-app"
        };
        
        var agent = new Agent(agentConfig);
        agent.Start();
        
        services.AddSingleton(agent);
    }
}
```

### appsettings.json

```json
{
  "SecurityPlatform": {
    "ControlPlaneUrl": "https://api.securityplatform.com",
    "AuthKey": "your-token-here"
  }
}
```

### Environment Variables

```bash
export SECURITY_PLATFORM_CONTROL_PLANE_URL=https://api.securityplatform.com
export SECURITY_PLATFORM_AUTH_KEY=your-token-here
```

## Features

- ✅ Automatic OpenTelemetry instrumentation
- ✅ ASP.NET Core middleware support
- ✅ PII redaction
- ✅ Policy enforcement
- ✅ Simple configuration
- ✅ Zero manual instrumentation required

## Documentation

See [docs/DEPLOYMENT.md](../../docs/DEPLOYMENT.md) for full documentation.
