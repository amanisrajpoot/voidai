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
import timber.log.Timber

/**
 * Security Platform Agent for Android applications.
 */
class Agent private constructor(
    private val config: AgentConfig,
    private val context: Context?
) {
    private val redactor: Redactor = Redactor(config.redactionRules)
    private val policyEngine: PolicyEngine = PolicyEngine(config.policy)
    private var tracerProvider: SdkTracerProvider? = null
    private var isStarted: Boolean = false

    companion object {
        @JvmStatic
        fun create(config: AgentConfig, context: Context? = null): Agent {
            return Agent(config, context)
        }

        @JvmStatic
        fun createDefault(context: Context? = null): Agent {
            val config = AgentConfig.default()
            return Agent(config, context)
        }
    }

    /**
     * Start the agent and initialize OpenTelemetry.
     */
    fun start() {
        if (isStarted) {
            Timber.w("Agent already started")
            return
        }

        try {
            val otlpEndpoint = config.otlpEndpoint
                ?: "${config.controlPlaneUrl}/v1/traces"

            val resource = Resource.getDefault()
                .merge(
                    Resource.create(
                        Attributes.of(
                            ResourceAttributes.SERVICE_NAME, config.serviceName,
                            ResourceAttributes.DEPLOYMENT_ENVIRONMENT, config.environment
                        )
                    )
                )

            val spanExporter = OtlpHttpSpanExporter.builder()
                .setEndpoint(otlpEndpoint)
                .setTimeout(java.time.Duration.parse("PT${config.telemetry.exportTimeout}"))
                .build()

            tracerProvider = SdkTracerProvider.builder()
                .setResource(resource)
                .addSpanProcessor(
                    BatchSpanProcessor.builder(spanExporter)
                        .setMaxQueueSize(config.telemetry.maxQueueSize)
                        .setExportBatchSize(config.telemetry.batchSize)
                        .setScheduleDelay(
                            java.time.Duration.parse("PT${config.telemetry.batchTimeout}")
                        )
                        .build()
                )
                .build()

            OpenTelemetrySdk.builder()
                .setTracerProvider(tracerProvider)
                .buildAndRegisterGlobal()

            isStarted = true
            Timber.i("Security Platform Agent started for service: ${config.serviceName}")
        } catch (e: Exception) {
            Timber.e(e, "Failed to start agent")
            throw RuntimeException("Failed to start agent", e)
        }
    }

    /**
     * Stop the agent and shutdown OpenTelemetry.
     */
    fun stop() {
        if (!isStarted) {
            return
        }

        try {
            tracerProvider?.shutdown()
            tracerProvider = null
            isStarted = false
            Timber.i("Security Platform Agent stopped")
        } catch (e: Exception) {
            Timber.e(e, "Error stopping agent")
        }
    }

    /**
     * Get the OpenTelemetry instance.
     */
    fun getOpenTelemetry(): OpenTelemetry {
        return io.opentelemetry.api.GlobalOpenTelemetry.get()
    }

    /**
     * Get a tracer for the given instrumentation name.
     */
    fun getTracer(instrumentationName: String): Tracer {
        return io.opentelemetry.api.GlobalOpenTelemetry.getTracer(instrumentationName)
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
    fun checkPolicy(action: String, context: Map<String, Any>): PolicyResult {
        return policyEngine.check(action, context)
    }

    fun isStarted(): Boolean = isStarted
    fun getConfig(): AgentConfig = config
}
