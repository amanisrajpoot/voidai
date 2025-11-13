# Security Platform iOS SDK

Security observability SDK for iOS applications with OpenTelemetry integration.

## Installation

### Swift Package Manager

Add to your `Package.swift`:

```swift
dependencies: [
    .package(url: "https://github.com/securityplatform/ios-sdk", from: "1.0.0")
]
```

Or in Xcode:
1. File → Add Packages...
2. Enter: `https://github.com/securityplatform/ios-sdk`
3. Select version: `1.0.0`

### CocoaPods

```ruby
pod 'SecurityPlatform', '~> 1.0'
```

## Quick Start

```swift
import SecurityPlatform

// Create agent configuration
let config = AgentConfig(
    controlPlaneURL: "https://api.securityplatform.com",
    authKey: ProcessInfo.processInfo.environment["SECURITY_PLATFORM_AUTH_KEY"],
    serviceName: "my-ios-app",
    environment: "production"
)

// Initialize and start agent
let agent = Agent(config: config)
try agent.start()

// Redact PII
let redacted = agent.redact(sensitiveData)

// Check policy
let result = agent.checkPolicy(action: "api_call", context: ["user": "alice"])
if !result.allowed {
    // Handle blocked action
}

// Get tracer for custom spans
let tracer = agent.getTracer(name: "my-app")
let span = tracer.spanBuilder(spanName: "api-request").startSpan()
// ... use span
span.end()
```

## Configuration

```swift
let config = AgentConfig(
    controlPlaneURL: "https://api.securityplatform.com",
    authKey: "your-auth-key",
    serviceName: "my-app",
    environment: "production",
    otlpEndpoint: "http://localhost:4318",
    telemetry: AgentConfig.TelemetryConfig(
        batchSize: 100,
        batchTimeout: "5s",
        exportTimeout: "30s",
        maxQueueSize: 2048
    ),
    policy: AgentConfig.PolicyConfig(
        mode: "observe",
        autoEnableBlocking: false
    ),
    redactionRules: [
        RedactionRule(pattern: ".*password.*"),
        RedactionRule(pattern: ".*token.*")
    ]
)
```

## Features

- ✅ OpenTelemetry integration
- ✅ PII redaction
- ✅ Policy engine (observe/block modes)
- ✅ Swift Package Manager support
- ✅ CocoaPods support
- ✅ iOS 13+ support
