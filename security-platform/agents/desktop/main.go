package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/securityplatform/go-agent"
)

var (
	configPath = flag.String("config", "/etc/security-platform/config.yaml", "Path to configuration file")
	service    = flag.Bool("service", false, "Run as system service/daemon")
)

func main() {
	flag.Parse()

	// Load configuration
	config, err := agent.LoadConfigFromFile(*configPath)
	if err != nil {
		// Fallback to environment variables
		config = agent.LoadConfigFromEnv()
	}

	// Create agent
	ag, err := agent.NewAgent(config)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Start agent
	ctx := context.Background()
	if err := ag.Start(ctx); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}
	defer ag.Stop(ctx)

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if *service {
		// Run as service/daemon
		log.Println("Security Platform Desktop Agent running as service...")
		<-sigChan
	} else {
		// Run interactively
		log.Println("Security Platform Desktop Agent started. Press Ctrl+C to stop.")
		<-sigChan
	}

	log.Println("Shutting down...")
}
