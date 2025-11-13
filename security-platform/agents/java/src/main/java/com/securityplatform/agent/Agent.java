package com.securityplatform.agent;

import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.api.common.Attributes;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.exporter.otlp.http.trace.OtlpHttpSpanExporter;
import io.opentelemetry.sdk.OpenTelemetrySdk;
import io.opentelemetry.sdk.resources.Resource;
import io.opentelemetry.sdk.trace.SdkTracerProvider;
import io.opentelemetry.sdk.trace.export.BatchSpanProcessor;
import io.opentelemetry.semconv.resource.attributes.ResourceAttributes;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.Map;
import java.util.concurrent.TimeUnit;

/**
 * Security Platform Agent for Java applications.
 */
public class Agent {
    private static final Logger logger = LoggerFactory.getLogger(Agent.class);
    
    private AgentConfig config;
    private OpenTelemetry openTelemetry;
    private Tracer tracer;
    private Redactor redactor;
    private PolicyEngine policyEngine;
    private boolean started = false;

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
                    ResourceAttributes.SERVICE_NAME, config.getServiceName() != null ? config.getServiceName() : "java-service",
                    ResourceAttributes.DEPLOYMENT_ENVIRONMENT, config.getEnvironment() != null ? config.getEnvironment() : "production"
                )));

            String otlpEndpoint = config.getOtlpEndpoint();
            if (otlpEndpoint == null || otlpEndpoint.isEmpty()) {
                otlpEndpoint = config.getControlPlaneUrl() + "/v1/traces";
            }

            OtlpHttpSpanExporter spanExporter = OtlpHttpSpanExporter.builder()
                .setEndpoint(otlpEndpoint)
                .setTimeout(Duration.ofSeconds(30))
                .build();

            SdkTracerProvider tracerProvider = SdkTracerProvider.builder()
                .setResource(resource)
                .addSpanProcessor(BatchSpanProcessor.builder(spanExporter)
                    .setMaxExportBatchSize(config.getTelemetry().getBatchSize())
                    .setExportTimeout(Duration.ofSeconds(30))
                    .setScheduleDelay(Duration.ofMillis(5000))
                    .build())
                .build();

            this.openTelemetry = OpenTelemetrySdk.builder()
                .setTracerProvider(tracerProvider)
                .build();

            this.tracer = openTelemetry.getTracer("security-platform-agent", "1.0.0");
            this.started = true;

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
            this.started = false;
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
     * Get the tracer instance.
     */
    public Tracer getTracer() {
        return tracer;
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

    public boolean isStarted() {
        return started;
    }
}
