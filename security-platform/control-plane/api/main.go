package main

import (
	"log"
	"os"

	"github.com/security-platform/control-plane/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := api.NewServer(port)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
