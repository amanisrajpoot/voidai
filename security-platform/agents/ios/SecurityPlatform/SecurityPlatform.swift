//
//  SecurityPlatform.swift
//  SecurityPlatform
//
//  Security Platform SDK for iOS
//

import Foundation
import OpenTelemetryApi
import OpenTelemetrySdk
import OpenTelemetryProtocolExporterHttp

/// Security Platform Agent for iOS applications.
@objc public class SecurityPlatformAgent: NSObject {
    private var config: AgentConfig
    private var tracerProvider: TracerProvider?
    private var redactor: Redactor
    private var policyEngine: PolicyEngine
    private var started: Bool = false
    
    /// Initialize the agent with configuration.
    @objc public init(config: AgentConfig) {
        self.config = config
        self.redactor = Redactor(blocklist: config.redactionRules ?? [])
        self.policyEngine = PolicyEngine(policy: config.policy ?? PolicyConfig())
        super.init()
    }
    
    /// Start the agent and initialize OpenTelemetry.
    @objc public func start() throws {
        guard !started else {
            print("Warning: Agent already started")
            return
        }
        
        let resource = Resource(attributes: [
            ResourceAttributes.serviceName.rawValue: AttributeValue.string(config.serviceName ?? "ios-app"),
            ResourceAttributes.deploymentEnvironment.rawValue: AttributeValue.string(config.environment ?? "production")
        ])
        
        let otlpEndpoint = config.otlpEndpoint ?? "\(config.controlPlaneURL ?? "https://api.securityplatform.com")/v1/traces"
        
        var headers: [String: String] = [:]
        if let authKey = config.authKey {
            headers["Authorization"] = "Bearer \(authKey)"
        }
        
        let exporter = OtlpHttpTraceExporter(endpoint: otlpEndpoint, headers: headers)
        
        let spanProcessor = BatchSpanProcessor(spanExporter: exporter)
        
        tracerProvider = TracerProviderBuilder()
            .with(resource: resource)
            .with(spanProcessor: spanProcessor)
            .build()
        
        OpenTelemetrySDK.instance.tracerProvider = tracerProvider!
        started = true
        
        print("Security Platform Agent started successfully")
    }
    
    /// Stop the agent.
    @objc public func stop() {
        guard started else { return }
        tracerProvider?.shutdown()
        started = false
        print("Security Platform Agent stopped")
    }
    
    /// Redact PII from data.
    @objc public func redact(_ data: Any) -> Any {
        return redactor.redact(data)
    }
    
    /// Check if an action is allowed by policy.
    @objc public func checkPolicy(action: String, context: [String: Any]?) -> PolicyResult {
        return policyEngine.check(action: action, context: context ?? [:])
    }
}

/// Configuration for the Security Platform Agent.
@objc public class AgentConfig: NSObject {
    @objc public var controlPlaneURL: String?
    @objc public var authKey: String?
    @objc public var serviceName: String?
    @objc public var environment: String?
    @objc public var otlpEndpoint: String?
    @objc public var redactionRules: [String]?
    @objc public var policy: PolicyConfig?
    
    @objc public override init() {
        super.init()
        controlPlaneURL = ProcessInfo.processInfo.environment["SECURITY_PLATFORM_CONTROL_PLANE_URL"] ?? "https://api.securityplatform.com"
        authKey = ProcessInfo.processInfo.environment["SECURITY_PLATFORM_AUTH_KEY"]
        serviceName = "ios-app"
        environment = "production"
    }
}

/// Policy configuration.
@objc public class PolicyConfig: NSObject {
    @objc public var mode: String = "observe"
    
    @objc public override init() {
        super.init()
    }
}

/// Redactor for PII protection.
class Redactor {
    private let blocklist: [String]
    
    init(blocklist: [String]) {
        self.blocklist = blocklist.isEmpty ? Redactor.defaultBlocklist : blocklist
    }
    
    func redact(_ data: Any) -> Any {
        if let dict = data as? [String: Any] {
            return redactDictionary(dict)
        } else if let array = data as? [Any] {
            return array.map { redact($0) }
        }
        return data
    }
    
    private func redactDictionary(_ dict: [String: Any]) -> [String: Any] {
        var redacted: [String: Any] = [:]
        for (key, value) in dict {
            if shouldRedact(key) {
                redacted[key] = "[REDACTED]"
            } else {
                redacted[key] = redact(value)
            }
        }
        return redacted
    }
    
    private func shouldRedact(_ key: String) -> Bool {
        let lowerKey = key.lowercased()
        return blocklist.contains { pattern in
            lowerKey.contains(pattern.lowercased())
        }
    }
    
    static let defaultBlocklist = [
        "password", "passwd", "pwd", "secret", "token",
        "api_key", "apikey", "credit_card", "card_number",
        "cvv", "ssn", "social_security", "pin", "passcode"
    ]
}

/// Policy engine.
class PolicyEngine {
    private let policy: PolicyConfig
    
    init(policy: PolicyConfig) {
        self.policy = policy
    }
    
    func check(action: String, context: [String: Any]) -> PolicyResult {
        // Simple policy check
        return PolicyResult(allowed: true)
    }
}

/// Policy result.
@objc public class PolicyResult: NSObject {
    @objc public let allowed: Bool
    @objc public let reason: String?
    
    @objc public init(allowed: Bool, reason: String? = nil) {
        self.allowed = allowed
        self.reason = reason
        super.init()
    }
}
