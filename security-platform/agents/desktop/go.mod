module github.com/security-platform/desktop-agent

go 1.21

require (
	github.com/security-platform/go-agent v0.0.0
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.21.0
	go.opentelemetry.io/otel/sdk v1.21.0
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/security-platform/go-agent => ../go
