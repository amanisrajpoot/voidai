import Foundation

/// Configuration for the Security Platform Agent.
public struct AgentConfig {
    public var controlPlaneURL: String
    public var authKey: String?
    public var serviceName: String
    public var environment: String
    public var namespace: String
    public var otlpEndpoint: String?
    public var telemetry: TelemetryConfig
    public var redactionRules: [RedactionRule]
    public var policy: PolicyConfig
    public var security: SecurityConfig

    public init(
        controlPlaneURL: String = "https://api.securityplatform.com",
        authKey: String? = nil,
        serviceName: String = "ios-service",
        environment: String = "production",
        namespace: String = "default",
        otlpEndpoint: String? = nil,
        telemetry: TelemetryConfig = TelemetryConfig(),
        redactionRules: [RedactionRule] = [],
        policy: PolicyConfig = PolicyConfig(),
        security: SecurityConfig = SecurityConfig()
    ) {
        self.controlPlaneURL = controlPlaneURL
        self.authKey = authKey
        self.serviceName = serviceName
        self.environment = environment
        self.namespace = namespace
        self.otlpEndpoint = otlpEndpoint
        self.telemetry = telemetry
        self.redactionRules = redactionRules
        self.policy = policy
        self.security = security
    }

    public struct TelemetryConfig {
        public var batchSize: Int
        public var batchTimeout: String
        public var exportTimeout: String
        public var maxQueueSize: Int

        public init(
            batchSize: Int = 100,
            batchTimeout: String = "5s",
            exportTimeout: String = "30s",
            maxQueueSize: Int = 2048
        ) {
            self.batchSize = batchSize
            self.batchTimeout = batchTimeout
            self.exportTimeout = exportTimeout
            self.maxQueueSize = maxQueueSize
        }
    }

    public struct PolicyConfig {
        public var mode: String // "observe" or "block"
        public var autoEnableBlocking: Bool
        public var observePeriodHours: Int
        public var rules: [PolicyRule]

        public init(
            mode: String = "observe",
            autoEnableBlocking: Bool = false,
            observePeriodHours: Int = 48,
            rules: [PolicyRule] = []
        ) {
            self.mode = mode
            self.autoEnableBlocking = autoEnableBlocking
            self.observePeriodHours = observePeriodHours
            self.rules = rules
        }
    }

    public struct SecurityConfig {
        public var mtlsEnabled: Bool
        public var certificatePath: String?
        public var keyPath: String?
        public var caBundlePath: String?

        public init(
            mtlsEnabled: Bool = false,
            certificatePath: String? = nil,
            keyPath: String? = nil,
            caBundlePath: String? = nil
        ) {
            self.mtlsEnabled = mtlsEnabled
            self.certificatePath = certificatePath
            self.keyPath = keyPath
            self.caBundlePath = caBundlePath
        }
    }
}
