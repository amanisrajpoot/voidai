package agent

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfigFromFile loads configuration from a YAML file.
func LoadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// LoadConfigFromEnv loads configuration from environment variables.
func LoadConfigFromEnv() *Config {
	config := DefaultConfig()

	if url := os.Getenv("SECURITY_PLATFORM_URL"); url != "" {
		config.ControlPlaneURL = url
	}
	if key := os.Getenv("SECURITY_PLATFORM_AUTH_KEY"); key != "" {
		config.AuthKey = key
	}
	if service := os.Getenv("SERVICE_NAME"); service != "" {
		config.ServiceName = service
	}
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		config.Environment = env
	}

	return config
}

// NewAgentFromFile creates a new agent from a configuration file.
func NewAgentFromFile(path string) (*Agent, error) {
	config, err := LoadConfigFromFile(path)
	if err != nil {
		return nil, err
	}
	return NewAgent(config)
}

// NewAgentFromEnv creates a new agent from environment variables.
func NewAgentFromEnv() (*Agent, error) {
	config := LoadConfigFromEnv()
	return NewAgent(config)
}
