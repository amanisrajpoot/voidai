# Security Platform - Quick Start Guide

Get up and running with Security Platform in 5 minutes.

## Step 1: Choose Your Platform

### Backend Application

**Java:**
```xml
<!-- Add to pom.xml -->
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>security-platform-agent</artifactId>
    <version>1.0.0</version>
</dependency>
```

```java
import com.securityplatform.agent.Agent;
import com.securityplatform.agent.AgentFactory;

Agent agent = AgentFactory.createDefault();
agent.start();
```

**.NET:**
```bash
dotnet add package SecurityPlatform.Agent
```

```csharp
using SecurityPlatform.Agent;

var agent = AgentFactory.CreateDefault();
agent.Start();
```

**Go:**
```bash
go get github.com/securityplatform/go-agent
```

```go
import "github.com/securityplatform/go-agent"

ag, _ := agent.NewAgentFromEnv()
ag.Start(context.Background())
```

**Node.js:**
```bash
npm install @security-platform/node-agent
```

```javascript
const { Agent } = require('@security-platform/node-agent');

const agent = new Agent({
  controlPlaneUrl: process.env.SECURITY_PLATFORM_URL,
  authKey: process.env.SECURITY_PLATFORM_AUTH_KEY
});
agent.start();
```

**Python:**
```bash
pip install security-platform-agent
```

```python
from securityplatform.agent import Agent, AgentFactory

agent = AgentFactory.create_default()
agent.start()
```

### Mobile Application

**iOS:**
```swift
import SecurityPlatform

let config = AgentConfig(
    controlPlaneURL: "https://api.securityplatform.com",
    authKey: ProcessInfo.processInfo.environment["SECURITY_PLATFORM_AUTH_KEY"]
)
let agent = Agent(config: config)
try agent.start()
```

**Android:**
```kotlin
import com.securityplatform.agent.Agent
import com.securityplatform.agent.AgentConfig

val config = AgentConfig(
    controlPlaneUrl = "https://api.securityplatform.com",
    authKey = BuildConfig.SECURITY_PLATFORM_AUTH_KEY
)
val agent = Agent.create(config, this)
agent.start()
```

### Desktop

**macOS:**
```bash
brew install security-platform-agent
brew services start security-platform-agent
```

**Windows:**
```powershell
choco install security-platform-agent
Start-Service -Name SecurityPlatformAgent
```

**Linux:**
```bash
# DEB
curl -fsSL https://install.securityplatform.com/linux | sudo bash

# RPM
curl -fsSL https://install.securityplatform.com/linux-rpm | sudo bash
```

## Step 2: Configure

Set environment variables:

```bash
export SECURITY_PLATFORM_URL="https://api.securityplatform.com"
export SECURITY_PLATFORM_AUTH_KEY="your-auth-key-here"
export SERVICE_NAME="my-service"
export ENVIRONMENT="production"
```

Or create `config.yaml`:

```yaml
control_plane_url: "https://api.securityplatform.com"
auth_key: "your-auth-key-here"
service_name: "my-service"
environment: "production"
```

## Step 3: Verify

Check that the agent is running:

```bash
# Check logs
tail -f /var/log/security-platform-agent.log

# Or check service status
systemctl status security-platform-agent  # Linux
launchctl list | grep security-platform    # macOS
Get-Service SecurityPlatformAgent          # Windows
```

## Step 4: Integrate

### HTTP Middleware

**Java (Spring Boot):**
```java
@Configuration
public class SecurityPlatformConfig {
    @Bean
    public FilterRegistrationBean<SecurityPlatformFilter> securityPlatformFilter(Agent agent) {
        FilterRegistrationBean<SecurityPlatformFilter> registration = 
            new FilterRegistrationBean<>();
        registration.setFilter(new SecurityPlatformFilter(agent));
        registration.addUrlPatterns("/*");
        return registration;
    }
}
```

**.NET (ASP.NET Core):**
```csharp
app.UseSecurityPlatform();
```

**Go (net/http):**
```go
handler := agent.Middleware(mux)
http.ListenAndServe(":8080", handler)
```

**Node.js (Express):**
```javascript
app.use(agent.middleware());
```

**Python (FastAPI):**
```python
app.add_middleware(SecurityPlatformMiddleware, agent=agent)
```

## Step 5: Monitor

View telemetry in the Security Platform control plane:
1. Log in to https://app.securityplatform.com
2. Navigate to Services
3. Find your service
4. View traces, metrics, and logs

## Troubleshooting

**Agent not connecting?**
- Verify `SECURITY_PLATFORM_AUTH_KEY` is set correctly
- Check network connectivity: `curl https://api.securityplatform.com/health`
- Review agent logs for errors

**No telemetry appearing?**
- Ensure agent is started: `agent.start()`
- Check OTLP endpoint configuration
- Verify service name matches in control plane

**Need help?**
- Documentation: [docs/](docs/)
- Examples: [examples/](examples/)
- Support: support@securityplatform.com

## Next Steps

- Configure [redaction rules](docs/REDACTION.md)
- Set up [policies](docs/POLICIES.md)
- Enable [session recording](docs/SESSION_RECORDING.md) (mobile/frontend)
- Integrate with [CI/CD](docs/CI_CD.md)
