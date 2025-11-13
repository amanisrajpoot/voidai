# Security Platform Agent - Java

Java agent for Security Platform with OpenTelemetry integration.

## Installation

### Maven

Add to your `pom.xml`:

```xml
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>security-platform-agent</artifactId>
    <version>1.0.0</version>
</dependency>
```

### Gradle

Add to your `build.gradle`:

```groovy
dependencies {
    implementation 'com.securityplatform:security-platform-agent:1.0.0'
}
```

## Quick Start

```java
import com.securityplatform.agent.Agent;
import com.securityplatform.agent.AgentConfig;

// Create configuration
AgentConfig config = new AgentConfig();
config.setControlPlaneUrl("https://api.securityplatform.com");
config.setAuthKey(System.getenv("SECURITY_PLATFORM_AUTH_KEY"));
config.setServiceName("my-java-service");
config.setEnvironment("production");

// Initialize and start agent
Agent agent = new Agent(config);
agent.start();

// Use OpenTelemetry tracer
var span = agent.getTracer()
    .spanBuilder("my-operation")
    .startSpan();

try (var scope = span.makeCurrent()) {
    // Your code here
} finally {
    span.end();
}

// Stop agent on shutdown
agent.stop();
```

## Servlet Filter

For web applications, use the ServletFilter:

```java
// In web.xml or ServletContextInitializer
FilterRegistration.Dynamic filter = servletContext.addFilter("securityPlatform", 
    new ServletFilter(agent));
filter.addMappingForUrlPatterns(null, false, "/*");
```

## Configuration

See `agent-config.yaml` for full configuration options.
