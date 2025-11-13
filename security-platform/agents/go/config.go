package agent

// Config represents the agent configuration.
type Config struct {
	ControlPlaneURL string          `yaml:"control_plane_url" json:"control_plane_url"`
	AuthKey         string          `yaml:"auth_key" json:"auth_key"`
	ServiceName     string          `yaml:"service_name" json:"service_name"`
	Environment     string          `yaml:"environment" json:"environment"`
	Namespace       string          `yaml:"namespace" json:"namespace"`
	OtlpEndpoint    string          `yaml:"otlp_endpoint" json:"otlp_endpoint"`
	Telemetry       TelemetryConfig `yaml:"telemetry" json:"telemetry"`
	RedactionRules  []RedactionRule `yaml:"redaction_rules" json:"redaction_rules"`
	Policy          *PolicyConfig   `yaml:"policy" json:"policy"`
	Security        *SecurityConfig `yaml:"security" json:"security"`
}

// TelemetryConfig configures telemetry export.
type TelemetryConfig struct {
	BatchSize     int    `yaml:"batch_size" json:"batch_size"`
	BatchTimeout  string `yaml:"batch_timeout" json:"batch_timeout"`
	ExportTimeout string `yaml:"export_timeout" json:"export_timeout"`
	MaxQueueSize  int    `yaml:"max_queue_size" json:"max_queue_size"`
}

// PolicyConfig configures policy enforcement.
type PolicyConfig struct {
	Mode               string                   `yaml:"mode" json:"mode"` // "observe" or "block"
	AutoEnableBlocking bool                     `yaml:"auto_enable_blocking" json:"auto_enable_blocking"`
	ObservePeriodHours int                      `yaml:"observe_period_hours" json:"observe_period_hours"`
	Rules              []map[string]interface{} `yaml:"rules" json:"rules"`
}

// SecurityConfig configures security settings.
type SecurityConfig struct {
	MtlsEnabled   bool   `yaml:"mtls_enabled" json:"mtls_enabled"`
	CertificatePath string `yaml:"certificate_path" json:"certificate_path"`
	KeyPath        string `yaml:"key_path" json:"key_path"`
	CaBundlePath   string `yaml:"ca_bundle_path" json:"ca_bundle_path"`
}

// RedactionRule defines a rule for redacting PII.
type RedactionRule struct {
	Pattern     string `yaml:"pattern" json:"pattern"`
	Replacement string `yaml:"replacement" json:"replacement"`
	Field       string `yaml:"field" json:"field"`
}

// PolicyResult represents the result of a policy check.
type PolicyResult struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}
