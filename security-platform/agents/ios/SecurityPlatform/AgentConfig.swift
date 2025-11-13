import Foundation

/// Configuration for the Security Platform Agent.
public struct AgentConfig {
    public var controlPlaneURL: String
    public var authKey: String?
    public var serviceName: String?
    public var environment: String?
    public var namespace: String?
    public var otlpEndpoint: String?
    public var telemetry: TelemetryConfig?
    public var redactionRules: [RedactionRule]?
    public var policy: PolicyConfig?
    public var security: SecurityConfig?

    public init(controlPlaneURL: String,
                authKey: String? = nil,
                serviceName: String? = nil,
                environment: String? = nil) {
        self.controlPlaneURL = controlPlaneURL
        self.authKey = authKey
        self.serviceName = serviceName
        self.environment = environment
    }
}

public struct TelemetryConfig {
    public var batchSize: Int
    public var batchTimeout: String
    public var exportTimeout: String
    public var maxQueueSize: Int

    public init(batchSize: Int = 100,
                batchTimeout: String = "5s",
                exportTimeout: String = "30s",
                maxQueueSize: Int = 2048) {
        self.batchSize = batchSize
        self.batchSize = batchTimeout
        self.exportTimeout = exportTimeout
        self.maxQueueSize = maxQueueSize
    }
}

public struct PolicyConfig {
    public var mode: String // "observe" or "block"
    public var autoEnableBlocking: Bool
    public var observePeriodHours: Int
    public var rules: [[String: Any]]?

    public init(mode: String = "observe",
                autoEnableBlocking: Bool = false,
                observePeriodHours: Int = 48) {
        self.mode = mode
        self.autoEnableBlocking = autoEnableBlocking
        self.observePeriodHours = observePeriodHours
    }
}

public struct SecurityConfig {
    public var mtlsEnabled: Bool
    public var certificatePath: String?
    public var keyPath: String?
    public var caBundlePath: String?

    public init(mtlsEnabled: Bool = true) {
        self.mtlsEnabled = mtlsEnabled
    }
}

public struct RedactionRule {
    public var pattern: String?
    public var replacement: String?
    public var field: String?

    public init(pattern: String? = nil,
                replacement: String? = nil,
                field: String? = nil) {
        self.pattern = pattern
        self.replacement = replacement
        self.field = field
    }
}

public struct PolicyResult {
    public var allowed: Bool
    public var reason: String?

    public init(allowed: Bool, reason: String? = nil) {
        self.allowed = allowed
        self.reason = reason
    }
}
