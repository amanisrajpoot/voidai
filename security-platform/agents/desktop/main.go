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
	"time"

	"github.com/security-platform/go-agent"
)

var (
	configPath = flag.String("config", "", "Path to agent config file")
	service    = flag.Bool("service", false, "Run as system service/daemon")
)

func main() {
	flag.Parse()

	// Load configuration
	config := loadConfig(*configPath)

	// Create agent
	ag := agent.NewAgent(config)

	// Start agent
	if err := ag.Start(); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}
	defer ag.Stop()

	if *service {
		// Run as service/daemon
		runService(ag)
	} else {
		// Run in foreground
		runForeground(ag)
	}
}

func loadConfig(configPath string) *agent.Config {
	// Default config
	config := &agent.Config{
		ServiceName: "desktop-agent",
		Environment: "production",
	}

	// Load from file if provided
	if configPath != "" {
		// TODO: Load from YAML/JSON file
		log.Printf("Loading config from %s", configPath)
	}

	// Override with environment variables
	if url := os.Getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL"); url != "" {
		config.ControlPlaneURL = url
	}
	if key := os.Getenv("SECURITY_PLATFORM_AUTH_KEY"); key != "" {
		config.AuthKey = key
	}

	return config
}

func runForeground(ag *agent.Agent) {
	log.Println("Security Platform Desktop Agent running in foreground")
	log.Println("Press Ctrl+C to stop")

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}

func runService(ag *agent.Agent) {
	log.Println("Security Platform Desktop Agent running as service")

	// Create signal channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Main service loop
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case sig := <-sigChan:
			log.Printf("Received signal: %v", sig)
			if sig == syscall.SIGTERM || sig == os.Interrupt {
				return
			}
		case <-ticker.C:
			// Periodic health check
			if !ag.IsStarted() {
				log.Println("Agent not started, attempting restart...")
				if err := ag.Start(); err != nil {
					log.Printf("Failed to restart agent: %v", err)
				}
			}
		}
	}
}

// Install service (platform-specific)
func installService() error {
	// Platform-specific installation
	// Windows: Use golang.org/x/sys/windows/svc
	// macOS: Use launchd plist
	// Linux: Use systemd service file
	return fmt.Errorf("service installation not implemented")
}

// Uninstall service
func uninstallService() error {
	return fmt.Errorf("service uninstallation not implemented")
}
