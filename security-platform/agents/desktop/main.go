package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/security-platform/go-agent/agent"
	"gopkg.in/yaml.v3"
)

var (
	configPath = flag.String("config", "", "Path to configuration file")
	serviceMode = flag.Bool("service", false, "Run as a service/daemon")
)

type Config struct {
	ControlPlaneURL string `yaml:"control_plane_url"`
	AuthKey         string `yaml:"auth_key"`
	ServiceName     string `yaml:"service_name"`
	Version         string `yaml:"version"`
	Environment     string `yaml:"environment"`
	OtlpEndpoint    string `yaml:"otlp_endpoint"`
	LocalPolicy     string `yaml:"local_policy"`
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	// Load configuration
	config := loadConfig()

	// Create agent
	agentConfig := &agent.Config{
		ServiceName:     config.ServiceName,
		Version:         config.Version,
		Environment:     config.Environment,
		OtlpEndpoint:    config.OtlpEndpoint,
		ControlPlaneURL: config.ControlPlaneURL,
		AuthKey:         config.AuthKey,
		LocalPolicy:     config.LocalPolicy,
	}

	ag := agent.NewAgent(agentConfig)

	// Start agent
	if err := ag.Start(ctx); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}

	log.Println("Security Platform Desktop Agent started")

	// Wait for context cancellation
	<-ctx.Done()

	// Stop agent
	if err := ag.Stop(ctx); err != nil {
		log.Printf("Error stopping agent: %v", err)
	}

	log.Println("Security Platform Desktop Agent stopped")
}

func loadConfig() *Config {
	config := &Config{
		ServiceName:  "desktop-agent",
		Version:      "1.0.0",
		Environment:  "production",
		OtlpEndpoint: "http://localhost:4318/v1/traces",
		LocalPolicy:  "observe",
	}

	// Try to find config file
	configFile := *configPath
	if configFile == "" {
		// Try common locations
		homeDir, _ := os.UserHomeDir()
		possiblePaths := []string{
			"./config.yaml",
			"./agent-config.yaml",
			filepath.Join(homeDir, ".security-platform", "config.yaml"),
			"/etc/security-platform/config.yaml",
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				configFile = path
				break
			}
		}
	}

	if configFile != "" {
		data, err := os.ReadFile(configFile)
		if err == nil {
			if err := yaml.Unmarshal(data, config); err != nil {
				log.Printf("Warning: Failed to parse config file: %v", err)
			}
		}
	}

	// Override with environment variables
	if url := os.Getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL"); url != "" {
		config.ControlPlaneURL = url
	}
	if key := os.Getenv("SECURITY_PLATFORM_AUTH_KEY"); key != "" {
		config.AuthKey = key
	}
	if name := os.Getenv("SECURITY_PLATFORM_SERVICE_NAME"); name != "" {
		config.ServiceName = name
	}

	return config
}
