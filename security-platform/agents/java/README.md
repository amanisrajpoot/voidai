# Security Platform Java Agent

Security observability agent for Java applications with OpenTelemetry integration.

## Installation

### Maven

```xml
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>security-platform-agent</artifactId>
    <version>1.0.0</version>
</dependency>
```

### Gradle

```gradle
implementation 'com.securityplatform:security-platform-agent:1.0.0'
```

## Quick Start

### Basic Usage

```java
import com.securityplatform.agent.Agent;
import com.securityplatform.agent.AgentFactory;

// Create agent from default config (uses environment variables)
Agent agent = AgentFactory.createDefault();
agent.start();

// Or load from YAML config file
Agent agent = AgentFactory.fromYamlFile("/path/to/config.yaml");
agent.start();

// Get OpenTelemetry instance for instrumentation
OpenTelemetry openTelemetry = agent.getOpenTelemetry();
Tracer tracer = agent.getTracer("my-service");

// Redact PII
Object redacted = agent.redact(sensitiveData);

// Check policy
PolicyResult result = agent.checkPolicy("action", context);
if (!result.isAllowed()) {
    // Handle blocked action
}

// Stop agent on shutdown
agent.stop();
```

### Servlet Filter Integration

```java
import com.securityplatform.agent.filter.SecurityPlatformFilter;

// In web.xml or ServletContextInitializer
FilterRegistration.Dynamic filter = servletContext.addFilter(
    "securityPlatformFilter",
    new SecurityPlatformFilter(agent)
);
filter.addMappingForUrlPatterns(null, false, "/*");
```

### Spring Boot Integration

```java
@Configuration
public class SecurityPlatformConfig {
    @Bean
    public Agent securityPlatformAgent() {
        Agent agent = AgentFactory.fromYamlFile("classpath:security-platform.yaml");
        agent.start();
        return agent;
    }

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

## Configuration

Create a `security-platform.yaml` file:

```yaml
control_plane_url: "https://api.securityplatform.com"
auth_key: "${SECURITY_PLATFORM_AUTH_KEY}"
service_name: "my-java-service"
environment: "production"

telemetry:
  otlp_endpoint: "http://localhost:4318"
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
- ✅ Automatic instrumentation (servlet, Spring)
- ✅ PII redaction
- ✅ Policy engine (observe/block modes)
- ✅ Servlet filter integration
- ✅ Spring Boot support
