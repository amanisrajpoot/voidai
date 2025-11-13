package com.securityplatform.agent;

import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.api.common.Attributes;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.api.trace.propagation.W3CTraceContextPropagator;
import io.opentelemetry.context.propagation.ContextPropagators;
import io.opentelemetry.exporter.otlp.trace.OtlpSpanExporter;
import io.opentelemetry.sdk.OpenTelemetrySdk;
import io.opentelemetry.sdk.resources.Resource;
import io.opentelemetry.sdk.trace.SdkTracerProvider;
import io.opentelemetry.sdk.trace.export.BatchSpanProcessor;
import io.opentelemetry.semconv.ResourceAttributes;

import java.util.Map;
import java.util.concurrent.TimeUnit;

/**
 * Security Platform Java Agent
 */
public class Agent {
    private OpenTelemetry openTelemetry;
    private AgentConfig config;
    private Tracer tracer;
    private boolean started = false;

    public Agent(AgentConfig config) {
        this.config = config;
    }

    public void start() {
        if (started) {
            return;
        }

        Resource resource = Resource.getDefault()
            .merge(Resource.builder()
                .put(ResourceAttributes.SERVICE_NAME, config.getServiceName() != null ? config.getServiceName() : "unknown-service")
                .put(ResourceAttributes.SERVICE_VERSION, config.getVersion() != null ? config.getVersion() : "1.0.0")
                .put(ResourceAttributes.DEPLOYMENT_ENVIRONMENT, config.getEnvironment() != null ? config.getEnvironment() : "production")
                .build());

        OtlpSpanExporter spanExporter = OtlpSpanExporter.builder()
            .setEndpoint(config.getOtlpEndpoint() != null ? config.getOtlpEndpoint() : "http://localhost:4318/v1/traces")
            .setTimeout(30, TimeUnit.SECONDS)
            .build();

        SdkTracerProvider tracerProvider = SdkTracerProvider.builder()
            .setResource(resource)
            .addSpanProcessor(BatchSpanProcessor.builder(spanExporter).build())
            .build();

        this.openTelemetry = OpenTelemetrySdk.builder()
            .setTracerProvider(tracerProvider)
            .setPropagators(ContextPropagators.create(W3CTraceContextPropagator.getInstance()))
            .buildAndRegisterGlobal();

        this.tracer = openTelemetry.getTracer("security-platform-agent", config.getVersion());
        this.started = true;
    }

    public void stop() {
        if (!started) {
            return;
        }

        if (openTelemetry instanceof OpenTelemetrySdk) {
            OpenTelemetrySdk sdk = (OpenTelemetrySdk) openTelemetry;
            sdk.getSdkTracerProvider().shutdown();
        }
        started = false;
    }

    public Tracer getTracer() {
        return tracer;
    }

    public OpenTelemetry getOpenTelemetry() {
        return openTelemetry;
    }

    public AgentConfig getConfig() {
        return config;
    }

    public boolean isStarted() {
        return started;
    }
}
