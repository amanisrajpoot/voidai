package com.securityplatform.agent;

import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.api.common.Attributes;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.exporter.otlp.http.trace.OtlpHttpSpanExporter;
import io.opentelemetry.sdk.OpenTelemetrySdk;
import io.opentelemetry.sdk.resources.Resource;
import io.opentelemetry.sdk.trace.SdkTracerProvider;
import io.opentelemetry.sdk.trace.export.BatchSpanProcessor;
import io.opentelemetry.semconv.ResourceAttributes;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;
import java.util.concurrent.TimeUnit;

/**
 * Security Platform Agent for Java applications.
 * Provides automatic instrumentation and telemetry collection.
 */
public class Agent {
    private static final Logger logger = LoggerFactory.getLogger(Agent.class);
    
    private final AgentConfig config;
    private OpenTelemetry openTelemetry;
    private Redactor redactor;
    private PolicyEngine policyEngine;
    private boolean started = false;

    /**
     * Create a new agent instance with the given configuration.
     */
    public Agent(AgentConfig config) {
        this.config = config;
        this.redactor = new Redactor(config.getRedactionRules());
        this.policyEngine = new PolicyEngine(config.getPolicy());
    }

    /**
     * Start the agent and initialize OpenTelemetry.
     */
    public void start() {
        if (started) {
            logger.warn("Agent already started");
            return;
        }

        try {
            Resource resource = Resource.getDefault()
                .merge(Resource.create(Attributes.of(
                    ResourceAttributes.SERVICE_NAME, config.getServiceName(),
                    ResourceAttributes.DEPLOYMENT_ENVIRONMENT, config.getEnvironment()
                )));

            String otlpEndpoint = config.getOtlpEndpoint();
            if (otlpEndpoint == null || otlpEndpoint.isEmpty()) {
                otlpEndpoint = config.getControlPlaneUrl() + "/v1/traces";
            }

            OtlpHttpSpanExporter spanExporter = OtlpHttpSpanExporter.builder()
                .setEndpoint(otlpEndpoint)
                .setHeaders(() -> {
                    Map<String, String> headers = new java.util.HashMap<>();
                    String authKey = config.getAuthKey();
                    if (authKey != null && !authKey.isEmpty()) {
                        headers.put("Authorization", "Bearer " + authKey);
                    }
                    return headers;
                })
                .build();

            SdkTracerProvider tracerProvider = SdkTracerProvider.builder()
                .addSpanProcessor(BatchSpanProcessor.builder(spanExporter)
                    .setMaxQueueSize(config.getTelemetry().getMaxQueueSize())
                    .setMaxExportBatchSize(config.getTelemetry().getBatchSize())
                    .setExportTimeout(config.getTelemetry().getExportTimeout(), TimeUnit.SECONDS)
                    .setScheduleDelay(config.getTelemetry().getBatchTimeout(), TimeUnit.SECONDS)
                    .build())
                .setResource(resource)
                .build();

            this.openTelemetry = OpenTelemetrySdk.builder()
                .setTracerProvider(tracerProvider)
                .build();

            started = true;
            logger.info("Security Platform Agent started successfully");
        } catch (Exception e) {
            logger.error("Failed to start agent", e);
            throw new RuntimeException("Failed to start Security Platform Agent", e);
        }
    }

    /**
     * Stop the agent and shutdown OpenTelemetry.
     */
    public void stop() {
        if (!started) {
            return;
        }

        try {
            if (openTelemetry instanceof OpenTelemetrySdk) {
                ((OpenTelemetrySdk) openTelemetry).getSdkTracerProvider().shutdown();
            }
            started = false;
            logger.info("Security Platform Agent stopped");
        } catch (Exception e) {
            logger.error("Error stopping agent", e);
        }
    }

    /**
     * Get the OpenTelemetry instance.
     */
    public OpenTelemetry getOpenTelemetry() {
        return openTelemetry;
    }

    /**
     * Get a tracer for manual instrumentation.
     */
    public Tracer getTracer(String instrumentationName) {
        return openTelemetry.getTracer(instrumentationName);
    }

    /**
     * Redact PII from data.
     */
    public Object redact(Object data) {
        return redactor.redact(data);
    }

    /**
     * Check if an action is allowed by policy.
     */
    public PolicyResult checkPolicy(String action, Map<String, Object> context) {
        return policyEngine.check(action, context);
    }

    /**
     * Check if the agent is started.
     */
    public boolean isStarted() {
        return started;
    }
}
