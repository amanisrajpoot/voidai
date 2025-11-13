# Security Platform Android SDK

Security observability SDK for Android applications with OpenTelemetry integration.

## Installation

### Gradle (Kotlin DSL)

```kotlin
dependencies {
    implementation("com.securityplatform:security-platform-agent:1.0.0")
}
```

### Gradle (Groovy)

```groovy
dependencies {
    implementation 'com.securityplatform:security-platform-agent:1.0.0'
}
```

### Maven

```xml
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>security-platform-agent</artifactId>
    <version>1.0.0</version>
</dependency>
```

## Quick Start

```kotlin
import com.securityplatform.agent.Agent
import com.securityplatform.agent.AgentConfig

class MainActivity : AppCompatActivity() {
    private lateinit var agent: Agent

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        // Create agent configuration
        val config = AgentConfig(
            controlPlaneUrl = "https://api.securityplatform.com",
            authKey = BuildConfig.SECURITY_PLATFORM_AUTH_KEY,
            serviceName = "my-android-app",
            environment = "production"
        )

        // Initialize and start agent
        agent = Agent.create(config, this)
        agent.start()

        // Redact PII
        val redacted = agent.redact(sensitiveData)

        // Check policy
        val result = agent.checkPolicy("api_call", mapOf("user" to "alice"))
        if (!result.allowed) {
            // Handle blocked action
        }

        // Get tracer for custom spans
        val tracer = agent.getTracer("my-app")
        val span = tracer.spanBuilder("api-request").startSpan()
        // ... use span
        span.end()
    }

    override fun onDestroy() {
        super.onDestroy()
        agent.stop()
    }
}
```

## Configuration

```kotlin
val config = AgentConfig(
    controlPlaneUrl = "https://api.securityplatform.com",
    authKey = "your-auth-key",
    serviceName = "my-app",
    environment = "production",
    otlpEndpoint = "http://localhost:4318",
    telemetry = AgentConfig.TelemetryConfig(
        batchSize = 100,
        batchTimeout = "5s",
        exportTimeout = "30s",
        maxQueueSize = 2048
    ),
    policy = AgentConfig.PolicyConfig(
        mode = "observe",
        autoEnableBlocking = false
    ),
    redactionRules = listOf(
        RedactionRule(pattern = ".*password.*"),
        RedactionRule(pattern = ".*token.*")
    )
)
```

## Features

- ✅ OpenTelemetry integration
- ✅ PII redaction
- ✅ Policy engine (observe/block modes)
- ✅ Android 5.0+ (API 21+) support
- ✅ Kotlin and Java support
- ✅ Maven and Gradle support
