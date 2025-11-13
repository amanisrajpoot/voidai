package agent

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"
)

// Agent represents the Security Platform Agent for Go applications.
type Agent struct {
	config       *Config
	logger       *zap.Logger
	redactor     *Redactor
	policyEngine *PolicyEngine
	tracer       *trace.TracerProvider
	started      bool
}

// NewAgent creates a new Agent instance with the given configuration.
func NewAgent(config *Config) (*Agent, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	agent := &Agent{
		config:       config,
		logger:       logger,
		redactor:     NewRedactor(config.RedactionRules),
		policyEngine: NewPolicyEngine(config.Policy),
	}

	return agent, nil
}

// Start initializes and starts the agent.
func (a *Agent) Start(ctx context.Context) error {
	if a.started {
		a.logger.Warn("Agent already started")
		return nil
	}

	otlpEndpoint := a.config.OtlpEndpoint
	if otlpEndpoint == "" {
		otlpEndpoint = fmt.Sprintf("%s/v1/traces", a.config.ControlPlaneURL)
	}

	// Create OTLP exporter
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", a.config.AuthKey),
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

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

	// Create tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			trace.WithBatchTimeout(time.Duration(a.config.Telemetry.BatchTimeout)),
			trace.WithMaxExportBatchSize(a.config.Telemetry.BatchSize),
			trace.WithMaxQueueSize(a.config.Telemetry.MaxQueueSize),
		),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	a.tracer = tp
	a.started = true

	a.logger.Info("Security Platform Agent started",
		zap.String("service", a.config.ServiceName),
		zap.String("environment", a.config.Environment),
	)

	return nil
}

// Stop shuts down the agent.
func (a *Agent) Stop(ctx context.Context) error {
	if !a.started {
		return nil
	}

	if a.tracer != nil {
		if err := a.tracer.Shutdown(ctx); err != nil {
			a.logger.Error("Failed to shutdown tracer", zap.Error(err))
			return err
		}
	}

	a.started = false
	a.logger.Info("Security Platform Agent stopped")
	return nil
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

// GetConfig returns the agent configuration.
func (a *Agent) GetConfig() *Config {
	return a.config
}
