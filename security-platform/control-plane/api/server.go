package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server represents the control plane API server
type Server struct {
	router *mux.Router
	port   string
}

// NewServer creates a new API server
func NewServer(port string) *Server {
	s := &Server{
		router: mux.NewRouter(),
		port:   port,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Sessions endpoint
	api.HandleFunc("/sessions", s.handleSessions).Methods("POST", "GET")
	api.HandleFunc("/sessions/{id}", s.handleSession).Methods("GET", "DELETE")

	// Policy endpoint
	api.HandleFunc("/policy", s.handlePolicy).Methods("GET", "PUT")
	api.HandleFunc("/policy/rules", s.handlePolicyRules).Methods("GET", "POST", "DELETE")

	// Telemetry endpoint (OTLP)
	api.HandleFunc("/v1/traces", s.handleTraces).Methods("POST")

	// Agents endpoint
	api.HandleFunc("/agents", s.handleAgents).Methods("GET")
	api.HandleFunc("/agents/{id}", s.handleAgent).Methods("GET", "PUT", "DELETE")

	// Incidents endpoint
	api.HandleFunc("/incidents", s.handleIncidents).Methods("GET", "POST")
	api.HandleFunc("/incidents/{id}", s.handleIncident).Methods("GET", "PUT")

	// Playbooks endpoint
	api.HandleFunc("/playbooks", s.handlePlaybooks).Methods("GET", "POST")
	api.HandleFunc("/playbooks/{id}", s.handlePlaybook).Methods("GET", "PUT", "DELETE")

	// Health check
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req struct {
			SessionID string                 `json:"sessionId"`
			Events    []interface{}           `json:"events"`
			Metadata  map[string]interface{} `json:"metadata"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Store session events (in production, persist to storage)
		log.Printf("Received session %s with %d events", req.SessionID, len(req.Events))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "ok",
			"sessionId": req.SessionID,
		})
		return
	}

	// GET: List sessions
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]map[string]interface{}{
		{"id": "session-1", "created": time.Now()},
	})
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	if r.Method == "DELETE" {
		// Delete session (GDPR compliance)
		log.Printf("Deleting session %s", sessionID)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// GET: Retrieve session
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      sessionID,
		"events":  []interface{}{},
		"created": time.Now(),
	})
}

func (s *Server) handlePolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var policy map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Update policy (in production, persist to storage)
		log.Printf("Updating policy")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}

	// GET: Retrieve policy
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mode": "observe",
		"rules": map[string]string{
			"sql_injection": "block",
			"xss":           "observe",
		},
	})
}

func (s *Server) handlePolicyRules(w http.ResponseWriter, r *http.Request) {
	// Policy rules CRUD
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]map[string]interface{}{})
}

func (s *Server) handleTraces(w http.ResponseWriter, r *http.Request) {
	// OTLP traces endpoint (simplified)
	var data interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received trace data")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleAgents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]map[string]interface{}{
		{"id": "agent-1", "status": "active", "lastSeen": time.Now()},
	})
}

func (s *Server) handleAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      agentID,
		"status":  "active",
		"lastSeen": time.Now(),
	})
}

func (s *Server) handleIncidents(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var incident map[string]interface{}
		json.NewDecoder(r.Body).Decode(&incident)
		log.Printf("Creating incident")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":     "incident-1",
			"status": "open",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]map[string]interface{}{})
}

func (s *Server) handleIncident(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	incidentID := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     incidentID,
		"status": "open",
	})
}

func (s *Server) handlePlaybooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var playbook map[string]interface{}
		json.NewDecoder(r.Body).Decode(&playbook)
		log.Printf("Creating playbook")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   "playbook-1",
			"name": playbook["name"],
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]map[string]interface{}{})
}

func (s *Server) handlePlaybook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playbookID := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   playbookID,
		"name": "Example Playbook",
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.port)
	log.Printf("Starting control plane API server on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
