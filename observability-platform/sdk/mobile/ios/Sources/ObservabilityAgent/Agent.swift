import Foundation
import OpenTelemetryApi
import OpenTelemetrySdk

public class ObservabilityAgent {
    private var tracer: Tracer?
    private var config: AgentConfig
    
    public init(config: AgentConfig) {
        self.config = config
    }
    
    public func initialize() {
        // Initialize OpenTelemetry
        let resource = Resource(attributes: [
            "service.name": AttributeValue.string(config.serviceName),
            "service.version": AttributeValue.string(config.serviceVersion),
            "deployment.environment": AttributeValue.string(config.environment)
        ])
        
        // Configure OTLP exporter
        let exporter = OtlpHttpTraceExporter(
            endpoint: config.otlpEndpoint,
            headers: ["Authorization": "Bearer \(config.authKey)"]
        )
        
        let spanProcessor = BatchSpanProcessor(spanExporter: exporter)
        let provider = TracerProviderBuilder()
            .with(resource: resource)
            .with(spanProcessor: spanProcessor)
            .build()
        
        OpenTelemetry.registerTracerProvider(tracerProvider: provider)
        self.tracer = OpenTelemetry.instance.tracerProvider.get(instrumentationName: "observability-agent")
    }
    
    public func trackEvent(name: String, attributes: [String: String] = [:]) {
        let span = tracer?.spanBuilder(spanName: name)
            .setStartTime(time: Date())
            .startSpan()
        
        for (key, value) in attributes {
            span?.setAttribute(key: key, value: value)
        }
        
        span?.end()
    }
    
    public func shutdown() {
        // Flush pending spans
    }
}

public struct AgentConfig {
    public let controlPlaneUrl: String
    public let authKey: String
    public let serviceName: String
    public let serviceVersion: String
    public let environment: String
    public let otlpEndpoint: String
    
    public init(
        controlPlaneUrl: String,
        authKey: String,
        serviceName: String,
        serviceVersion: String = "1.0.0",
        environment: String = "production",
        otlpEndpoint: String? = nil
    ) {
        self.controlPlaneUrl = controlPlaneUrl
        self.authKey = authKey
        self.serviceName = serviceName
        self.serviceVersion = serviceVersion
        self.environment = environment
        self.otlpEndpoint = otlpEndpoint ?? "\(controlPlaneUrl)/v1/traces"
    }
}
