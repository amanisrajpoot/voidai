package com.securityplatform.agent

import android.content.Context
import io.opentelemetry.api.OpenTelemetry
import io.opentelemetry.api.common.Attributes
import io.opentelemetry.api.trace.Tracer
import io.opentelemetry.exporter.otlp.http.trace.OtlpHttpSpanExporter
import io.opentelemetry.sdk.OpenTelemetrySdk
import io.opentelemetry.sdk.resources.Resource
import io.opentelemetry.sdk.trace.SdkTracerProvider
import io.opentelemetry.sdk.trace.export.BatchSpanProcessor
import io.opentelemetry.semconv.resource.attributes.ResourceAttributes
import java.util.concurrent.TimeUnit

/**
 * Security Platform Agent for Android applications.
 */
class Agent(private val config: AgentConfig) {
    private var tracerProvider: SdkTracerProvider? = null
    private var tracer: Tracer? = null
    private val redactor: Redactor = Redactor(config.redactionRules ?: emptyList())
    private val policyEngine: PolicyEngine = PolicyEngine(config.policy ?: PolicyConfig())
    private var started: Boolean = false

    /**
     * Start the agent and initialize OpenTelemetry.
     */
    fun start() {
        if (started) {
            android.util.Log.w(TAG, "Agent already started")
            return
        }

        try {
            val serviceName = config.serviceName ?: "android-app"
            val environment = config.environment ?: "production"

            val resource = Resource.getDefault()
                .merge(Resource.create(Attributes.of(
                    ResourceAttributes.SERVICE_NAME, serviceName,
                    ResourceAttributes.DEPLOYMENT_ENVIRONMENT, environment
                )))

            val otlpEndpoint = config.otlpEndpoint
                ?: "${config.controlPlaneUrl}/v1/traces"

            val headers = mutableMapOf<String, String>()
            val authKey = config.authKey
                ?: System.getenv("SECURITY_PLATFORM_AUTH_KEY")
            if (authKey != null) {
                headers["Authorization"] = "Bearer $authKey"
            }

            val spanExporter = OtlpHttpSpanExporter.builder()
                .setEndpoint(otlpEndpoint)
                .setTimeout(30, TimeUnit.SECONDS)
                .build()

            val batchSize = config.telemetry?.batchSize ?: 100
            val maxQueueSize = config.telemetry?.maxQueueSize ?: 2048

            tracerProvider = SdkTracerProvider.builder()
                .setResource(resource)
                .addSpanProcessor(
                    BatchSpanProcessor.builder(spanExporter)
                        .setMaxExportBatchSize(batchSize)
                        .setExportTimeout(30, TimeUnit.SECONDS)
                        .setScheduleDelay(5, TimeUnit.SECONDS)
                        .setMaxQueueSize(maxQueueSize)
                        .build()
                )
                .build()

            val openTelemetry = OpenTelemetrySdk.builder()
                .setTracerProvider(tracerProvider)
                .build()

            OpenTelemetry.setGlobalOpenTelemetry(openTelemetry)
            tracer = openTelemetry.getTracer("security-platform-agent", "1.0.0")
            started = true

            android.util.Log.i(TAG, "Security Platform Agent started successfully")
        } catch (e: Exception) {
            android.util.Log.e(TAG, "Failed to start agent", e)
            throw RuntimeException("Failed to start Security Platform Agent", e)
        }
    }

    /**
     * Stop the agent and shutdown OpenTelemetry.
     */
    fun stop() {
        if (!started) return

        try {
            tracerProvider?.shutdown()
            tracerProvider = null
            started = false
            android.util.Log.i(TAG, "Security Platform Agent stopped")
        } catch (e: Exception) {
            android.util.Log.e(TAG, "Error stopping agent", e)
        }
    }

    /**
     * Get the OpenTelemetry tracer.
     */
    fun getTracer(): Tracer? = tracer

    /**
     * Redact PII from data.
     */
    fun redact(data: Any?): Any? = redactor.redact(data)

    /**
     * Check if an action is allowed by policy.
     */
    fun checkPolicy(action: String, context: Map<String, Any>): PolicyResult =
        policyEngine.check(action, context)

    companion object {
        private const val TAG = "SecurityPlatformAgent"
    }
}
