package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/security-platform/go-agent"
)

func main() {
	// Create configuration
	config := &agent.Config{
		ControlPlaneURL: os.Getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL"),
		AuthKey:         os.Getenv("SECURITY_PLATFORM_AUTH_KEY"),
		ServiceName:     "security-platform-agent",
		Environment:     "production",
	}

	// Create and start agent
	ag, err := agent.NewAgent(config)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	if err := ag.Start(); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}
	defer ag.Stop()

	log.Println("Security Platform Agent started")

	// Create systemd service file
	servicePath := "/etc/systemd/system/security-platform-agent.service"
	createSystemdService(servicePath)

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}

func createSystemdService(path string) {
	service := `[Unit]
Description=Security Platform Agent
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/security-platform-agent
Restart=always
RestartSec=10
Environment="SECURITY_PLATFORM_CONTROL_PLANE_URL=https://api.securityplatform.com"
Environment="SECURITY_PLATFORM_AUTH_KEY=%i"

[Install]
WantedBy=multi-user.target
`

	if err := os.WriteFile(path, []byte(service), 0644); err != nil {
		log.Printf("Warning: Failed to create systemd service: %v", err)
	}
}
