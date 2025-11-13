# Security Platform iOS SDK

Simple, automatic security observability for iOS applications.

## Quick Start

### Swift Package Manager

Add to your `Package.swift`:

```swift
dependencies: [
    .package(url: "https://github.com/security-platform/agent-ios", from: "1.0.0")
]
```

Or add via Xcode: File → Add Packages → Enter URL

### CocoaPods

Add to your `Podfile`:

```ruby
pod 'SecurityPlatform', '~> 1.0.0'
```

## Usage

### Basic Integration

```swift
import SecurityPlatform

// Initialize agent
let config = AgentConfig()
config.controlPlaneURL = "https://api.securityplatform.com"
config.authKey = ProcessInfo.processInfo.environment["SECURITY_PLATFORM_AUTH_KEY"]
config.serviceName = "my-ios-app"

let agent = SecurityPlatformAgent(config: config)
try? agent.start()

// Your application code here...

// Shutdown on app termination
agent.stop()
```

### AppDelegate Integration

```swift
import UIKit
import SecurityPlatform

@main
class AppDelegate: UIResponder, UIApplicationDelegate {
    var agent: SecurityPlatformAgent?
    
    func application(_ application: UIApplication, 
                    didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {
        let config = AgentConfig()
        config.serviceName = "my-app"
        agent = SecurityPlatformAgent(config: config)
        try? agent?.start()
        return true
    }
}
```

### Environment Variables

Set in your Xcode scheme or Info.plist:

```
SECURITY_PLATFORM_CONTROL_PLANE_URL=https://api.securityplatform.com
SECURITY_PLATFORM_AUTH_KEY=your-token-here
```

## Features

- ✅ Automatic OpenTelemetry instrumentation
- ✅ PII redaction
- ✅ Policy enforcement
- ✅ Simple configuration
- ✅ Zero manual instrumentation required

## Documentation

See [docs/DEPLOYMENT.md](../../docs/DEPLOYMENT.md) for full documentation.
