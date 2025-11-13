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
import io.opentelemetry.semconv.ResourceAttributes

/**
 * Security Platform Agent for Android applications.
 * Provides automatic instrumentation and telemetry collection.
 */
class Agent private constructor(private val config: AgentConfig) {
    private var openTelemetry: OpenTelemetry? = null
    private val redactor: Redactor = Redactor(config.redactionRules ?: emptyList())
    private val policyEngine: PolicyEngine = PolicyEngine(config.policy ?: PolicyConfig())
    private var started: Boolean = false

    companion object {
        /**
         * Create a new agent instance with the given configuration.
         */
        @JvmStatic
        fun create(config: AgentConfig): Agent {
            return Agent(config)
        }
    }

    /**
     * Start the agent and initialize OpenTelemetry.
     */
    fun start() {
        if (started) {
            android.util.Log.w(TAG, "Agent already started")
            return
        }

        try {
            val resource = Resource.getDefault()
                .merge(
                    Resource.create(
                        Attributes.of(
                            ResourceAttributes.SERVICE_NAME, config.serviceName ?: "android-app",
                            ResourceAttributes.DEPLOYMENT_ENVIRONMENT, config.environment ?: "production"
                        )
                    )
                )

            val otlpEndpoint = config.otlpEndpoint
                ?: "${config.controlPlaneUrl ?: "https://api.securityplatform.com"}/v1/traces"

            val headers = mutableMapOf<String, String>()
            config.authKey?.let {
                headers["Authorization"] = "Bearer $it"
            }

            val spanExporter = OtlpHttpSpanExporter.builder()
                .setEndpoint(otlpEndpoint)
                .setHeaders(headers)
                .build()

            val tracerProvider = SdkTracerProvider.builder()
                .addSpanProcessor(
                    BatchSpanProcessor.builder(spanExporter)
                        .setMaxQueueSize(config.telemetry?.maxQueueSize ?: 2048)
                        .setMaxExportBatchSize(config.telemetry?.batchSize ?: 100)
                        .setExportTimeoutMillis((config.telemetry?.exportTimeoutSeconds ?: 30) * 1000L)
                        .setScheduleDelayMillis((config.telemetry?.batchTimeoutSeconds ?: 5) * 1000L)
                        .build()
                )
                .setResource(resource)
                .build()

            this.openTelemetry = OpenTelemetrySdk.builder()
                .setTracerProvider(tracerProvider)
                .build()

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
        if (!started) {
            return
        }

        try {
            (openTelemetry as? OpenTelemetrySdk)?.sdkTracerProvider?.shutdown()
            started = false
            android.util.Log.i(TAG, "Security Platform Agent stopped")
        } catch (e: Exception) {
            android.util.Log.e(TAG, "Error stopping agent", e)
        }
    }

    /**
     * Get a tracer for manual instrumentation.
     */
    fun getTracer(instrumentationName: String): Tracer {
        return openTelemetry?.getTracer(instrumentationName)
            ?: throw IllegalStateException("Agent not started")
    }

    /**
     * Redact PII from data.
     */
    fun redact(data: Any?): Any? {
        return redactor.redact(data)
    }

    /**
     * Check if an action is allowed by policy.
     */
    fun checkPolicy(action: String, context: Map<String, Any>? = null): PolicyResult {
        return policyEngine.check(action, context ?: emptyMap())
    }

    private companion object {
        const val TAG = "SecurityPlatformAgent"
    }
}
