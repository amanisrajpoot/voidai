import Foundation
import OpenTelemetryApi
import OpenTelemetrySdk
import ResourceExtension

/// Security Platform Agent for iOS applications.
public class Agent {
    private let config: AgentConfig
    private let redactor: Redactor
    private let policyEngine: PolicyEngine
    private var tracerProvider: TracerProvider?
    private var isStarted: Bool = false

    public init(config: AgentConfig) {
        self.config = config
        self.redactor = Redactor(rules: config.redactionRules)
        self.policyEngine = PolicyEngine(config: config.policy)
    }

    /// Start the agent and initialize OpenTelemetry.
    public func start() throws {
        guard !isStarted else {
            print("Agent already started")
            return
        }

        let otlpEndpoint = config.otlpEndpoint ?? "\(config.controlPlaneURL)/v1/traces"

        // Create resource
        let resource = Resource.getDefault()
            .merging(other: Resource(attributes: [
                "service.name": AttributeValue.string(config.serviceName),
                "deployment.environment": AttributeValue.string(config.environment)
            ]))

        // Create OTLP exporter (simplified - in production use proper OTLP exporter)
        // For now, we'll use stdout exporter as placeholder
        let exporter = StdoutSpanExporter()
        
        let spanProcessor = BatchSpanProcessor(spanExporter: exporter)
        
        self.tracerProvider = TracerProviderBuilder()
            .with(resource: resource)
            .with(spanProcessor: spanProcessor)
            .build()

        OpenTelemetry.registerTracerProvider(tracerProvider: tracerProvider!)
        isStarted = true
        
        print("Security Platform Agent started for service: \(config.serviceName)")
    }

    /// Stop the agent.
    public func stop() {
        guard isStarted else { return }
        tracerProvider?.shutdown()
        isStarted = false
        print("Security Platform Agent stopped")
    }

    /// Redact PII from data.
    public func redact(_ data: Any) -> Any {
        return redactor.redact(data)
    }

    /// Check if an action is allowed by policy.
    public func checkPolicy(action: String, context: [String: Any]) -> PolicyResult {
        return policyEngine.check(action: action, context: context)
    }

    /// Get a tracer for instrumentation.
    public func getTracer(name: String) -> Tracer {
        return OpenTelemetry.instance.tracerProvider.get(instrumentationName: name)
    }
}
