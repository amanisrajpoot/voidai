package agent

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// Agent represents the Security Platform Agent.
type Agent struct {
	config       *Config
	tracerProvider *trace.TracerProvider
	redactor     *Redactor
	policyEngine *PolicyEngine
	started      bool
}

// Config holds the agent configuration.
type Config struct {
	ControlPlaneURL string
	AuthKey         string
	ServiceName     string
	Environment     string
	OtlpEndpoint    string
	Telemetry       TelemetryConfig
	RedactionRules  []string
	Policy          PolicyConfig
}

// TelemetryConfig holds telemetry configuration.
type TelemetryConfig struct {
	BatchSize      int
	BatchTimeout   time.Duration
	ExportTimeout  time.Duration
	MaxQueueSize   int
}

// PolicyConfig holds policy configuration.
type PolicyConfig struct {
	Mode string // "observe" or "block"
}

// NewAgent creates a new agent instance.
func NewAgent(config *Config) *Agent {
	if config == nil {
		config = &Config{}
	}

	// Set defaults from environment variables
	if config.ControlPlaneURL == "" {
		config.ControlPlaneURL = os.Getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL")
		if config.ControlPlaneURL == "" {
			config.ControlPlaneURL = "https://api.securityplatform.com"
		}
	}

	if config.AuthKey == "" {
		config.AuthKey = os.Getenv("SECURITY_PLATFORM_AUTH_KEY")
	}

	if config.ServiceName == "" {
		config.ServiceName = "go-service"
	}

	if config.Environment == "" {
		config.Environment = "production"
	}

	if config.Telemetry.BatchSize == 0 {
		config.Telemetry.BatchSize = 100
	}

	if config.Telemetry.BatchTimeout == 0 {
		config.Telemetry.BatchTimeout = 5 * time.Second
	}

	if config.Telemetry.ExportTimeout == 0 {
		config.Telemetry.ExportTimeout = 30 * time.Second
	}

	if config.Telemetry.MaxQueueSize == 0 {
		config.Telemetry.MaxQueueSize = 2048
	}

	if config.Policy.Mode == "" {
		config.Policy.Mode = "observe"
	}

	return &Agent{
		config:       config,
		redactor:     NewRedactor(config.RedactionRules),
		policyEngine: NewPolicyEngine(&config.Policy),
	}
}

// Start initializes and starts the agent.
func (a *Agent) Start() error {
	if a.started {
		return fmt.Errorf("agent already started")
	}

	ctx := context.Background()

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(a.config.ServiceName),
			semconv.DeploymentEnvironmentKey.String(a.config.Environment),
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
	}

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithHeaders(headers),
	)
	if err != nil {
		return fmt.Errorf("failed to create exporter: %w", err)
	}

	// Create tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			trace.WithMaxQueueSize(a.config.Telemetry.MaxQueueSize),
			trace.WithBatchTimeout(a.config.Telemetry.BatchTimeout),
			trace.WithExportTimeout(a.config.Telemetry.ExportTimeout),
		),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	a.tracerProvider = tp
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

// GetTracer returns a tracer for manual instrumentation.
func (a *Agent) GetTracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// Redact redacts PII from data.
func (a *Agent) Redact(data interface{}) interface{} {
	return a.redactor.Redact(data)
}

// CheckPolicy checks if an action is allowed by policy.
func (a *Agent) CheckPolicy(action string, context map[string]interface{}) *PolicyResult {
	return a.policyEngine.Check(action, context)
}

// IsStarted returns whether the agent is started.
func (a *Agent) IsStarted() bool {
	return a.started
}
