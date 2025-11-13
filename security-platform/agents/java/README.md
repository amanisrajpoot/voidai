# Security Platform Java Agent

Simple, automatic security observability for Java applications.

## Quick Start

### Maven

Add to your `pom.xml`:

```xml
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>agent</artifactId>
    <version>1.0.0</version>
</dependency>
```

### Gradle

Add to your `build.gradle`:

```gradle
implementation 'com.securityplatform:agent:1.0.0'
```

## Usage

### Basic Integration

```java
import com.securityplatform.agent.Agent;
import com.securityplatform.agent.AgentConfig;

// Initialize agent
AgentConfig config = new AgentConfig();
config.setControlPlaneUrl("https://api.securityplatform.com");
config.setAuthKey(System.getenv("SECURITY_PLATFORM_AUTH_KEY"));
config.setServiceName("my-java-service");

Agent agent = new Agent(config);
agent.start();

// Your application code here...

// Shutdown on application exit
agent.stop();
```

### Spring Boot Integration

```java
import com.securityplatform.agent.Agent;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class SecurityPlatformConfig {
    
    @Value("${security.platform.control-plane-url}")
    private String controlPlaneUrl;
    
    @Value("${security.platform.auth-key}")
    private String authKey;
    
    @Bean
    public Agent securityPlatformAgent() {
        AgentConfig config = new AgentConfig();
        config.setControlPlaneUrl(controlPlaneUrl);
        config.setAuthKey(authKey);
        config.setServiceName("my-spring-app");
        
        Agent agent = new Agent(config);
        agent.start();
        return agent;
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
- ✅ PII redaction
- ✅ Policy enforcement
- ✅ Simple configuration
- ✅ Zero manual instrumentation required

## Documentation

See [docs/DEPLOYMENT.md](../../docs/DEPLOYMENT.md) for full documentation.
