package agent

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// Agent represents the Security Platform Agent.
type Agent struct {
	config       *Config
	tracerProvider *sdktrace.TracerProvider
	tracer       trace.Tracer
	redactor     *Redactor
	policyEngine *PolicyEngine
	started      bool
}

// NewAgent creates a new agent instance.
func NewAgent(config *Config) (*Agent, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	redactor := NewRedactor(config.RedactionRules)
	policyEngine := NewPolicyEngine(config.Policy)

	return &Agent{
		config:       config,
		redactor:     redactor,
		policyEngine: policyEngine,
	}, nil
}

// Start initializes and starts the agent.
func (a *Agent) Start() error {
	if a.started {
		return fmt.Errorf("agent already started")
	}

	serviceName := a.config.ServiceName
	if serviceName == "" {
		serviceName = "go-service"
	}

	environment := a.config.Environment
	if environment == "" {
		environment = "production"
	}

	// Create resource
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentKey.String(environment),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Determine OTLP endpoint
	otlpEndpoint := a.config.OtlpEndpoint
	if otlpEndpoint == "" {
		otlpEndpoint = fmt.Sprintf("%s/v1/traces", a.config.ControlPlaneURL)
	}

	// Create OTLP exporter
	headers := make(map[string]string)
	if a.config.AuthKey != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", a.config.AuthKey)
	} else if authKey := os.Getenv("SECURITY_PLATFORM_AUTH_KEY"); authKey != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", authKey)
	}

	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithHeaders(headers),
	)
	if err != nil {
		return fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create tracer provider
	batchSize := a.config.Telemetry.BatchSize
	if batchSize == 0 {
		batchSize = 100
	}

	maxQueueSize := a.config.Telemetry.MaxQueueSize
	if maxQueueSize == 0 {
		maxQueueSize = 2048
	}

	batchTimeout := parseDuration(a.config.Telemetry.BatchTimeout)
	if batchTimeout == 0 {
		batchTimeout = 5 * time.Second
	}

	exportTimeout := parseDuration(a.config.Telemetry.ExportTimeout)
	if exportTimeout == 0 {
		exportTimeout = 30 * time.Second
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(batchTimeout),
			sdktrace.WithExportTimeout(exportTimeout),
			sdktrace.WithMaxExportBatchSize(batchSize),
			sdktrace.WithMaxQueueSize(maxQueueSize),
		),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	a.tracerProvider = tracerProvider
	a.tracer = tracerProvider.Tracer("security-platform-agent", trace.WithInstrumentationVersion("1.0.0"))
	a.started = true

	return nil
}

// Stop shuts down the agent.
func (a *Agent) Stop() error {
	if !a.started {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.tracerProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown tracer provider: %w", err)
	}

	a.started = false
	return nil
}

// Tracer returns the OpenTelemetry tracer.
func (a *Agent) Tracer() trace.Tracer {
	return a.tracer
}

// Redact redacts PII from data.
func (a *Agent) Redact(data interface{}) interface{} {
	return a.redactor.Redact(data)
}

// CheckPolicy checks if an action is allowed by policy.
func (a *Agent) CheckPolicy(action string, context map[string]interface{}) *PolicyResult {
	return a.policyEngine.Check(action, context)
}

func parseDuration(duration string) time.Duration {
	if duration == "" {
		return 0
	}

	var value int
	var unit string
	fmt.Sscanf(duration, "%d%s", &value, &unit)

	switch unit {
	case "s":
		return time.Duration(value) * time.Second
	case "m":
		return time.Duration(value) * time.Minute
	case "h":
		return time.Duration(value) * time.Hour
	default:
		return 0
	}
}
