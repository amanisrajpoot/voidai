# Security Platform Agent - Android

Android SDK for Security Platform with OpenTelemetry integration.

## Installation

### Gradle

Add to your `build.gradle`:

```groovy
dependencies {
    implementation 'com.securityplatform:agent:1.0.0'
}
```

### Maven

```xml
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>agent</artifactId>
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

        // Create configuration
        val config = AgentConfig(
            controlPlaneUrl = "https://api.securityplatform.com",
            authKey = System.getenv("SECURITY_PLATFORM_AUTH_KEY"),
            serviceName = "my-android-app",
            environment = "production"
        )

        // Initialize and start agent
        agent = Agent(config)
        agent.start()

        // Use tracer
        val span = agent.getTracer()?.spanBuilder("my-operation")
            ?.startSpan()
        try {
            // Your code here
        } finally {
            span?.end()
        }
    }

    override fun onDestroy() {
        super.onDestroy()
        agent.stop()
    }
}
```

## Configuration

See `agent-config.yaml` for full configuration options.
