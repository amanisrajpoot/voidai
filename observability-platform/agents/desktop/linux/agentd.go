package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var (
	configFile = flag.String("config", "/etc/observability-agent/config.yaml", "Configuration file path")
	version    = flag.Bool("version", false, "Print version and exit")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Println("observability-agent 1.0.0")
		os.Exit(0)
	}

	// Load configuration
	cfg, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize agent
	agent, err := NewAgent(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize agent: %v", err)
	}

	// Start agent
	if err := agent.Start(); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}
	defer agent.Shutdown()

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}

type Agent struct {
	config      *Config
	tracer      *trace.TracerProvider
	ctx         context.Context
	cancel      context.CancelFunc
}

type Config struct {
	ControlPlaneURL string        `yaml:"control_plane_url"`
	AuthKey         string        `yaml:"auth_key"`
	ServiceName     string        `yaml:"service_name"`
	Environment     string        `yaml:"environment"`
	OTLPEndpoint    string        `yaml:"otlp_endpoint"`
	BatchSize       int           `yaml:"telemetry_batch_size"`
	ExportInterval  time.Duration `yaml:"telemetry_export_interval"`
}

func NewAgent(cfg *Config) (*Agent, error) {
	ctx, cancel := context.WithCancel(context.Background())

	otlpEndpoint := cfg.OTLPEndpoint
	if otlpEndpoint == "" {
		otlpEndpoint = cfg.ControlPlaneURL + "/v1/traces"
	}

	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization": "Bearer " + cfg.AuthKey,
		}),
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		),
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return &Agent{
		config: cfg,
		tracer: tp,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (a *Agent) Start() error {
	log.Println("Starting observability agent...")
	
	// Start background tasks
	go a.collectSystemMetrics()
	go a.collectNetworkMetrics()
	
	return nil
}

func (a *Agent) Shutdown() error {
	log.Println("Shutting down observability agent...")
	a.cancel()
	
	if a.tracer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return a.tracer.Shutdown(ctx)
	}
	
	return nil
}

func (a *Agent) collectSystemMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			// Collect system metrics (CPU, memory, disk, etc.)
			// Send to OTEL collector
		}
	}
}

func (a *Agent) collectNetworkMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			// Collect network metrics (optional: pcap)
			// Send to OTEL collector
		}
	}
}

func loadConfig(path string) (*Config, error) {
	// Simplified - would use gopkg.in/yaml.v3 in real implementation
	return &Config{
		ControlPlaneURL: os.Getenv("CONTROL_PLANE_URL"),
		AuthKey:         os.Getenv("AUTH_KEY"),
		ServiceName:     "desktop-agent",
		Environment:     "production",
	}, nil
}
