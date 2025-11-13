package com.observability.agent

import io.opentelemetry.api.OpenTelemetry
import io.opentelemetry.api.trace.Tracer
import io.opentelemetry.exporter.otlp.http.trace.OtlpHttpSpanExporter
import io.opentelemetry.sdk.OpenTelemetrySdk
import io.opentelemetry.sdk.resources.Resource
import io.opentelemetry.sdk.trace.SdkTracerProvider
import io.opentelemetry.sdk.trace.export.BatchSpanProcessor
import io.opentelemetry.semconv.resource.attributes.ResourceAttributes

class ObservabilityAgent private constructor(private val config: AgentConfig) {
    private var tracer: Tracer? = null
    
    companion object {
        @JvmStatic
        fun initialize(config: AgentConfig): ObservabilityAgent {
            val agent = ObservabilityAgent(config)
            agent.init()
            return agent
        }
    }
    
    private fun init() {
        val resource = Resource.getDefault()
            .merge(
                Resource.builder()
                    .put(ResourceAttributes.SERVICE_NAME, config.serviceName)
                    .put(ResourceAttributes.SERVICE_VERSION, config.serviceVersion)
                    .put(ResourceAttributes.DEPLOYMENT_ENVIRONMENT, config.environment)
                    .build()
            )
        
        val spanExporter = OtlpHttpSpanExporter.builder()
            .setEndpoint(config.otlpEndpoint)
            .addHeader("Authorization", "Bearer ${config.authKey}")
            .build()
        
        val spanProcessor = BatchSpanProcessor.builder(spanExporter).build()
        
        val tracerProvider = SdkTracerProvider.builder()
            .setResource(resource)
            .addSpanProcessor(spanProcessor)
            .build()
        
        val openTelemetry = OpenTelemetrySdk.builder()
            .setTracerProvider(tracerProvider)
            .build()
        
        OpenTelemetry.setGlobalOpenTelemetry(openTelemetry)
        this.tracer = openTelemetry.getTracer("observability-agent")
    }
    
    fun trackEvent(name: String, attributes: Map<String, String> = emptyMap()) {
        val span = tracer?.spanBuilder(name)?.startSpan()
        attributes.forEach { (key, value) ->
            span?.setAttribute(key, value)
        }
        span?.end()
    }
    
    fun shutdown() {
        // Flush pending spans
    }
}

data class AgentConfig(
    val controlPlaneUrl: String,
    val authKey: String,
    val serviceName: String,
    val serviceVersion: String = "1.0.0",
    val environment: String = "production",
    val otlpEndpoint: String? = null
) {
    val otlpEndpoint: String
        get() = otlpEndpoint ?: "$controlPlaneUrl/v1/traces"
}
