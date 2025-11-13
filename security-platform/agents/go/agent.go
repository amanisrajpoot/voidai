package agent

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// Agent represents the Security Platform Go Agent
type Agent struct {
	tracerProvider *trace.TracerProvider
	config         *Config
	started        bool
}

// Config holds the agent configuration
type Config struct {
	ServiceName     string
	Version         string
	Environment     string
	OtlpEndpoint    string
	ControlPlaneURL string
	AuthKey         string
	LocalPolicy     string // "observe" or "block"
}

// NewAgent creates a new agent instance
func NewAgent(config *Config) *Agent {
	if config == nil {
		config = &Config{}
	}
	return &Agent{
		config:  config,
		started: false,
	}
}

// Start initializes and starts the agent
func (a *Agent) Start(ctx context.Context) error {
	if a.started {
		return nil
	}

	endpoint := a.config.OtlpEndpoint
	if endpoint == "" {
		endpoint = "http://localhost:4318/v1/traces"
	}

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
	)
	if err != nil {
		return err
	}

	serviceName := a.config.ServiceName
	if serviceName == "" {
		serviceName = "unknown-service"
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(a.config.Version),
			semconv.DeploymentEnvironment(a.config.Environment),
		),
	)
	if err != nil {
		return err
	}

	a.tracerProvider = trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(a.tracerProvider)
	a.started = true
	return nil
}

// Stop shuts down the agent
func (a *Agent) Stop(ctx context.Context) error {
	if !a.started {
		return nil
	}

	if a.tracerProvider != nil {
		if err := a.tracerProvider.Shutdown(ctx); err != nil {
			return err
		}
	}
	a.started = false
	return nil
}

// GetTracerProvider returns the tracer provider
func (a *Agent) GetTracerProvider() *trace.TracerProvider {
	return a.tracerProvider
}

// GetConfig returns the agent configuration
func (a *Agent) GetConfig() *Config {
	return a.config
}

// IsStarted returns whether the agent is started
func (a *Agent) IsStarted() bool {
	return a.started
}
