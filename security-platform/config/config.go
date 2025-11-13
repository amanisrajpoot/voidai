package config

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// AgentConfig represents the agent configuration
type AgentConfig struct {
	ControlPlane ControlPlaneConfig `yaml:"control_plane" json:"control_plane"`
	Auth         AuthConfig         `yaml:"auth" json:"auth"`
	Service      ServiceConfig      `yaml:"service" json:"service"`
	Telemetry    TelemetryConfig    `yaml:"telemetry" json:"telemetry"`
	Redaction    RedactionConfig    `yaml:"redaction" json:"redaction"`
	Policy       PolicyConfig       `yaml:"policy" json:"policy"`
	Update       UpdateConfig       `yaml:"update" json:"update"`
	Logging      LoggingConfig      `yaml:"logging" json:"logging"`
	Security     SecurityConfig     `yaml:"security" json:"security"`
	Compliance   ComplianceConfig   `yaml:"compliance" json:"compliance"`
}

type ControlPlaneConfig struct {
	URL string `yaml:"url" json:"url"`
}

type AuthConfig struct {
	BootstrapToken string    `yaml:"bootstrap_token" json:"bootstrap_token"`
	PKI            *PKIConfig `yaml:"pki,omitempty" json:"pki,omitempty"`
}

type PKIConfig struct {
	CertPath string `yaml:"cert_path" json:"cert_path"`
	KeyPath  string `yaml:"key_path" json:"key_path"`
	CAPath   string `yaml:"ca_path" json:"ca_path"`
}

type ServiceConfig struct {
	Name        string `yaml:"name" json:"name"`
	Version     string `yaml:"version" json:"version"`
	Environment string `yaml:"environment" json:"environment"`
}

type TelemetryConfig struct {
	OTLPEndpoint     string `yaml:"otlp_endpoint" json:"otlp_endpoint"`
	BatchSize        int    `yaml:"batch_size" json:"batch_size"`
	FlushIntervalMs  int    `yaml:"flush_interval_ms" json:"flush_interval_ms"`
	OfflineMode      bool   `yaml:"offline_mode" json:"offline_mode"`
	OfflineCachePath string `yaml:"offline_cache_path" json:"offline_cache_path"`
}

type RedactionConfig struct {
	DefaultRules []string `yaml:"default_rules" json:"default_rules"`
	CustomRules  []string `yaml:"custom_rules" json:"custom_rules"`
	HashPayloads bool    `yaml:"hash_payloads" json:"hash_payloads"`
	SafeFields   []string `yaml:"safe_fields" json:"safe_fields"`
}

type PolicyConfig struct {
	Mode          string `yaml:"mode" json:"mode"`
	CacheTTL      int    `yaml:"cache_ttl" json:"cache_ttl"`
	LocalPolicyPath string `yaml:"local_policy_path" json:"local_policy_path"`
}

type UpdateConfig struct {
	Enabled          bool   `yaml:"enabled" json:"enabled"`
	CheckIntervalHours int  `yaml:"check_interval_hours" json:"check_interval_hours"`
	Channel          string `yaml:"channel" json:"channel"`
	VerifySignatures bool   `yaml:"verify_signatures" json:"verify_signatures"`
}

type LoggingConfig struct {
	Level  string `yaml:"level" json:"level"`
	Format string `yaml:"format" json:"format"`
	Output string `yaml:"output" json:"output"`
}

type SecurityConfig struct {
	MTLSEnabled   bool              `yaml:"mtls_enabled" json:"mtls_enabled"`
	SecureStorage SecureStorageConfig `yaml:"secure_storage" json:"secure_storage"`
	AuditLogEnabled bool            `yaml:"audit_log_enabled" json:"audit_log_enabled"`
	AuditLogPath    string          `yaml:"audit_log_path" json:"audit_log_path"`
}

type SecureStorageConfig struct {
	Type string `yaml:"type" json:"type"`
	Path string `yaml:"path" json:"path"`
}

type ComplianceConfig struct {
	SessionDataTTLDays int  `yaml:"session_data_ttl_days" json:"session_data_ttl_days"`
	ExportEnabled      bool `yaml:"export_enabled" json:"export_enabled"`
	DeletionEnabled    bool `yaml:"deletion_enabled" json:"deletion_enabled"`
}

// LoadConfig loads configuration from file
func LoadConfig(path string) (*AgentConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AgentConfig
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	} else if strings.HasSuffix(path, ".json") {
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported config format (use .yaml or .json)")
	}

	// Expand environment variables
	expandEnvVars(&config)

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// expandEnvVars expands ${VAR} and ${VAR:-default} patterns in strings
func expandEnvVars(config *AgentConfig) {
	config.ControlPlane.URL = expandString(config.ControlPlane.URL)
	config.Auth.BootstrapToken = expandString(config.Auth.BootstrapToken)
	config.Service.Name = expandString(config.Service.Name)
	config.Service.Environment = expandString(config.Service.Environment)
	config.Telemetry.OTLPEndpoint = expandString(config.Telemetry.OTLPEndpoint)
}

func expandString(s string) string {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		inner := match[2 : len(match)-1]
		parts := strings.SplitN(inner, ":-", 2)
		key := parts[0]
		defaultVal := ""
		if len(parts) > 1 {
			defaultVal = parts[1]
		}
		if val := os.Getenv(key); val != "" {
			return val
		}
		return defaultVal
	})
}

// Validate validates the configuration
func (c *AgentConfig) Validate() error {
	if c.ControlPlane.URL == "" {
		return fmt.Errorf("control_plane.url is required")
	}
	if c.Auth.BootstrapToken == "" && c.Auth.PKI == nil {
		return fmt.Errorf("either auth.bootstrap_token or auth.pki must be set")
	}
	if c.Service.Name == "" {
		return fmt.Errorf("service.name is required")
	}
	return nil
}

// GetOTLPEndpoint returns the OTLP endpoint URL
func (c *AgentConfig) GetOTLPEndpoint() string {
	if c.Telemetry.OTLPEndpoint != "" {
		return c.Telemetry.OTLPEndpoint
	}
	return fmt.Sprintf("%s/v1/traces", c.ControlPlane.URL)
}
