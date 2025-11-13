# Security Platform Agent - iOS

iOS SDK for Security Platform with OpenTelemetry integration.

## Installation

### CocoaPods

Add to your `Podfile`:

```ruby
pod 'SecurityPlatform', '~> 1.0'
```

Then run:

```bash
pod install
```

### Swift Package Manager

Add to your `Package.swift`:

```swift
dependencies: [
    .package(url: "https://github.com/security-platform/agents", from: "1.0.0")
]
```

## Quick Start

```swift
import SecurityPlatform

// Create configuration
let config = AgentConfig(
    controlPlaneURL: "https://api.securityplatform.com",
    authKey: ProcessInfo.processInfo.environment["SECURITY_PLATFORM_AUTH_KEY"],
    serviceName: "my-ios-app",
    environment: "production"
)

// Initialize and start agent
let agent = SecurityPlatformAgent(config: config)
try? agent.start()

// Use tracer
let span = agent.getTracer()?.spanBuilder(spanName: "my-operation")
    .startSpan()
defer { span?.end() }

// Stop on app termination
defer { agent.stop() }
```

## Configuration

See `agent-config.yaml` for full configuration options.
