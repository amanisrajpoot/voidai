package main

import (
	"log"
	"os"

	"github.com/observability-platform/control-plane/internal/api"
	"github.com/observability-platform/control-plane/internal/config"
)

func main() {
	cfg := config.Load()

	server := api.NewServer(cfg)
	
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
