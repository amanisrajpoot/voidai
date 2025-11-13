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

	// Create plist for launchd
	plistPath := "/Library/LaunchDaemons/com.securityplatform.agent.plist"
	createLaunchdPlist(plistPath)

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}

func createLaunchdPlist(path string) {
	plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.securityplatform.agent</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/security-platform-agent</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/var/log/security-platform-agent.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/security-platform-agent.error.log</string>
</dict>
</plist>`

	if err := os.WriteFile(path, []byte(plist), 0644); err != nil {
		log.Printf("Warning: Failed to create launchd plist: %v", err)
	}
}
