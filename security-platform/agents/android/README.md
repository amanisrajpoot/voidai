# Security Platform Android SDK

Simple, automatic security observability for Android applications.

## Quick Start

### Gradle

Add to your `build.gradle.kts`:

```kotlin
dependencies {
    implementation("com.securityplatform:agent:1.0.0")
}
```

### Maven

Add to your `pom.xml`:

```xml
<dependency>
    <groupId>com.securityplatform</groupId>
    <artifactId>agent</artifactId>
    <version>1.0.0</version>
</dependency>
```

## Usage

### Basic Integration

```kotlin
import com.securityplatform.agent.Agent
import com.securityplatform.agent.AgentConfig

class MainActivity : AppCompatActivity() {
    private lateinit var agent: Agent
    
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        
        // Initialize agent
        val config = AgentConfig.fromEnvironment(this)
        config.serviceName = "my-android-app"
        
        agent = Agent.create(config)
        agent.start()
        
        // Your application code here...
    }
    
    override fun onDestroy() {
        super.onDestroy()
        agent.stop()
    }
}
```

### Application Class Integration

```kotlin
import android.app.Application
import com.securityplatform.agent.Agent
import com.securityplatform.agent.AgentConfig

class MyApplication : Application() {
    private lateinit var agent: Agent
    
    override fun onCreate() {
        super.onCreate()
        
        val config = AgentConfig.fromEnvironment(this)
        config.serviceName = "my-app"
        
        agent = Agent.create(config)
        agent.start()
    }
    
    override fun onTerminate() {
        super.onTerminate()
        agent.stop()
    }
}
```

### Environment Variables

Set via build config or system properties:

```kotlin
buildConfigField("String", "SECURITY_PLATFORM_CONTROL_PLANE_URL", "\"https://api.securityplatform.com\"")
buildConfigField("String", "SECURITY_PLATFORM_AUTH_KEY", "\"your-token-here\"")
```

## Features

- ✅ Automatic OpenTelemetry instrumentation
- ✅ PII redaction
- ✅ Policy enforcement
- ✅ Simple configuration
- ✅ Zero manual instrumentation required

## Documentation

See [docs/DEPLOYMENT.md](../../docs/DEPLOYMENT.md) for full documentation.
