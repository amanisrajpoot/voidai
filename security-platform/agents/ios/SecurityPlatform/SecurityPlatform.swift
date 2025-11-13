import Foundation
import OpenTelemetryApi
import OpenTelemetrySdk
import OpenTelemetryProtocolExporterHttp

/// Security Platform Agent for iOS applications.
public class SecurityPlatformAgent {
    private let config: AgentConfig
    private var tracerProvider: TracerProvider?
    private var tracer: Tracer?
    private let redactor: Redactor
    private let policyEngine: PolicyEngine
    private var started: Bool = false

    public init(config: AgentConfig) {
        self.config = config
        self.redactor = Redactor(rules: config.redactionRules ?? [])
        self.policyEngine = PolicyEngine(config: config.policy ?? PolicyConfig())
    }

    /// Start the agent and initialize OpenTelemetry.
    public func start() throws {
        guard !started else {
            print("Warning: Agent already started")
            return
        }

        let serviceName = config.serviceName ?? "ios-app"
        let environment = config.environment ?? "production"

        // Create resource
        let resource = Resource(attributes: [
            "service.name": AttributeValue.string(serviceName),
            "deployment.environment": AttributeValue.string(environment)
        ])

        // Determine OTLP endpoint
        let otlpEndpoint = config.otlpEndpoint ?? "\(config.controlPlaneURL)/v1/traces"

        // Create OTLP exporter
        var headers: [String: String] = [:]
        if let authKey = config.authKey ?? ProcessInfo.processInfo.environment["SECURITY_PLATFORM_AUTH_KEY"] {
            headers["Authorization"] = "Bearer \(authKey)"
        }

        let exporter = OtlpHttpTraceExporter(endpoint: otlpEndpoint, headers: headers)

        // Create tracer provider
        let batchSize = config.telemetry?.batchSize ?? 100
        let maxQueueSize = config.telemetry?.maxQueueSize ?? 2048

        let spanProcessor = BatchSpanProcessor(spanExporter: exporter,
                                               scheduleDelay: TimeInterval(5),
                                               maxQueueSize: maxQueueSize,
                                               maxExportBatchSize: batchSize)

        let provider = TracerProviderBuilder()
            .add(spanProcessor: spanProcessor)
            .with(resource: resource)
            .build()

        OpenTelemetry.registerTracerProvider(tracerProvider: provider)
        self.tracerProvider = provider
        self.tracer = provider.get(instrumentationName: "security-platform-agent", instrumentationVersion: "1.0.0")
        self.started = true

        print("Security Platform Agent started successfully")
    }

    /// Stop the agent and shutdown OpenTelemetry.
    public func stop() {
        guard started else { return }
        tracerProvider?.shutdown()
        started = false
        print("Security Platform Agent stopped")
    }

    /// Get the OpenTelemetry tracer.
    public func getTracer() -> Tracer? {
        return tracer
    }

    /// Redact PII from data.
    public func redact(_ data: Any) -> Any {
        return redactor.redact(data)
    }

    /// Check if an action is allowed by policy.
    public func checkPolicy(action: String, context: [String: Any]) -> PolicyResult {
        return policyEngine.check(action: action, context: context)
    }
}
