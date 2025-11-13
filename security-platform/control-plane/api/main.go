package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/sessions", handleSessions).Methods("POST")
	api.HandleFunc("/errors", handleErrors).Methods("POST")
	api.HandleFunc("/events", handleEvents).Methods("POST")
	api.HandleFunc("/agents", handleAgents).Methods("GET", "POST")
	api.HandleFunc("/agents/{id}", handleAgent).Methods("GET", "PUT", "DELETE")
	api.HandleFunc("/rules", handleRules).Methods("GET", "POST")
	api.HandleFunc("/rules/{id}", handleRule).Methods("GET", "PUT", "DELETE")
	api.HandleFunc("/playbooks", handlePlaybooks).Methods("GET", "POST")
	api.HandleFunc("/playbooks/{id}/execute", handlePlaybookExecute).Methods("POST")

	// OTLP endpoint
	router.HandleFunc("/v1/traces", handleOTLPTraces).Methods("POST")
	router.HandleFunc("/v1/metrics", handleOTLPMetrics).Methods("POST")
	router.HandleFunc("/v1/logs", handleOTLPLogs).Methods("POST")

	// Health check
	router.HandleFunc("/health", handleHealth).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Starting server", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server shutdown failed", zap.Error(err))
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func handleSessions(w http.ResponseWriter, r *http.Request) {
	// Store session replay data
	w.WriteHeader(http.StatusAccepted)
}

func handleErrors(w http.ResponseWriter, r *http.Request) {
	// Store error events
	w.WriteHeader(http.StatusAccepted)
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	// Store general events
	w.WriteHeader(http.StatusAccepted)
}

func handleAgents(w http.ResponseWriter, r *http.Request) {
	// List or create agents
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[]`))
}

func handleAgent(w http.ResponseWriter, r *http.Request) {
	// Get, update, or delete agent
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id":"` + id + `"}`))
}

func handleRules(w http.ResponseWriter, r *http.Request) {
	// List or create rules
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[]`))
}

func handleRule(w http.ResponseWriter, r *http.Request) {
	// Get, update, or delete rule
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id":"` + id + `"}`))
}

func handlePlaybooks(w http.ResponseWriter, r *http.Request) {
	// List or create playbooks
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[]`))
}

func handlePlaybookExecute(w http.ResponseWriter, r *http.Request) {
	// Execute playbook
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"playbook_id":"` + id + `","status":"executed"}`))
}

func handleOTLPTraces(w http.ResponseWriter, r *http.Request) {
	// Handle OTLP traces
	w.WriteHeader(http.StatusOK)
}

func handleOTLPMetrics(w http.ResponseWriter, r *http.Request) {
	// Handle OTLP metrics
	w.WriteHeader(http.StatusOK)
}

func handleOTLPLogs(w http.ResponseWriter, r *http.Request) {
	// Handle OTLP logs
	w.WriteHeader(http.StatusOK)
}
