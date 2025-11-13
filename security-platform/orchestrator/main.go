package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type Playbook struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Triggers    []Trigger              `json:"triggers"`
	Actions     []Action               `json:"actions"`
	Enabled     bool                   `json:"enabled"`
}

type Trigger struct {
	Type      string                 `json:"type"`
	Condition map[string]interface{} `json:"condition"`
}

type Action struct {
	Type    string                 `json:"type"`
	Config  map[string]interface{} `json:"config"`
	Timeout int                    `json:"timeout"`
}

type Execution struct {
	PlaybookID string                 `json:"playbook_id"`
	Trigger    map[string]interface{} `json:"trigger"`
	Status     string                 `json:"status"`
	StartedAt  time.Time              `json:"started_at"`
	CompletedAt *time.Time            `json:"completed_at,omitempty"`
	Results    []ActionResult         `json:"results"`
}

type ActionResult struct {
	Action  Action `json:"action"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Output  interface{} `json:"output,omitempty"`
}

var playbooks = make(map[string]Playbook)
var executions = make([]Execution, 0)

func main() {
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/playbooks", handleListPlaybooks).Methods("GET")
	api.HandleFunc("/playbooks", handleCreatePlaybook).Methods("POST")
	api.HandleFunc("/playbooks/{id}", handleGetPlaybook).Methods("GET")
	api.HandleFunc("/playbooks/{id}", handleUpdatePlaybook).Methods("PUT")
	api.HandleFunc("/playbooks/{id}", handleDeletePlaybook).Methods("DELETE")
	api.HandleFunc("/playbooks/{id}/execute", handleExecutePlaybook).Methods("POST")
	api.HandleFunc("/webhooks/{playbook_id}", handleWebhook).Methods("POST")
	api.HandleFunc("/executions", handleListExecutions).Methods("GET")
	api.HandleFunc("/executions/{id}", handleGetExecution).Methods("GET")

	// Health check
	router.HandleFunc("/health", handleHealth).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting orchestrator on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down orchestrator")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func handleListPlaybooks(w http.ResponseWriter, r *http.Request) {
	list := make([]Playbook, 0, len(playbooks))
	for _, pb := range playbooks {
		list = append(list, pb)
	}
	json.NewEncoder(w).Encode(list)
}

func handleCreatePlaybook(w http.ResponseWriter, r *http.Request) {
	var pb Playbook
	if err := json.NewDecoder(r.Body).Decode(&pb); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playbooks[pb.ID] = pb
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pb)
}

func handleGetPlaybook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	pb, ok := playbooks[id]
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(pb)
}

func handleUpdatePlaybook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var pb Playbook
	if err := json.NewDecoder(r.Body).Decode(&pb); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pb.ID = id
	playbooks[id] = pb
	json.NewEncoder(w).Encode(pb)
}

func handleDeletePlaybook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	delete(playbooks, id)
	w.WriteHeader(http.StatusNoContent)
}

func handleExecutePlaybook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	pb, ok := playbooks[id]
	if !ok {
		http.NotFound(w, r)
		return
	}

	var triggerData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&triggerData); err != nil {
		triggerData = make(map[string]interface{})
	}

	exec := executePlaybook(pb, triggerData)
	executions = append(executions, exec)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exec)
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playbookID := vars["playbook_id"]
	pb, ok := playbooks[playbookID]
	if !ok {
		http.NotFound(w, r)
		return
	}

	var webhookData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&webhookData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exec := executePlaybook(pb, webhookData)
	executions = append(executions, exec)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exec)
}

func handleListExecutions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(executions)
}

func handleGetExecution(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// In real implementation, store executions with IDs
	json.NewEncoder(w).Encode(executions)
}

func executePlaybook(pb Playbook, triggerData map[string]interface{}) Execution {
	exec := Execution{
		PlaybookID: pb.ID,
		Trigger:    triggerData,
		Status:     "running",
		StartedAt:  time.Now(),
		Results:    make([]ActionResult, 0),
	}

	for _, action := range pb.Actions {
		result := executeAction(action, triggerData)
		exec.Results = append(exec.Results, result)
	}

	now := time.Now()
	exec.CompletedAt = &now
	exec.Status = "completed"

	return exec
}

func executeAction(action Action, context map[string]interface{}) ActionResult {
	result := ActionResult{
		Action: action,
		Success: false,
	}

	switch action.Type {
	case "block_ip":
		ip := action.Config["ip"].(string)
		// Call Kong Admin API or WAF to block IP
		result.Success = blockIP(ip)
		result.Output = map[string]string{"ip": ip, "status": "blocked"}

	case "revoke_token":
		userID := action.Config["user_id"].(string)
		// Call IAM to revoke tokens
		result.Success = revokeToken(userID)
		result.Output = map[string]string{"user_id": userID, "status": "revoked"}

	case "send_alert":
		channel := action.Config["channel"].(string)
		message := action.Config["message"].(string)
		// Send to Slack/PagerDuty/etc
		result.Success = sendAlert(channel, message)
		result.Output = map[string]string{"channel": channel, "status": "sent"}

	case "webhook":
		url := action.Config["url"].(string)
		method := action.Config["method"].(string)
		// Call external webhook
		result.Success = callWebhook(url, method, context)
		result.Output = map[string]string{"url": url, "status": "called"}

	default:
		result.Error = fmt.Sprintf("unknown action type: %s", action.Type)
	}

	return result
}

func blockIP(ip string) bool {
	// In real implementation, call Kong Admin API or WAF
	log.Printf("Blocking IP: %s", ip)
	return true
}

func revokeToken(userID string) bool {
	// In real implementation, call IAM API
	log.Printf("Revoking tokens for user: %s", userID)
	return true
}

func sendAlert(channel, message string) bool {
	// In real implementation, send to Slack/PagerDuty/etc
	log.Printf("Sending alert to %s: %s", channel, message)
	return true
}

func callWebhook(url, method string, data map[string]interface{}) bool {
	// In real implementation, make HTTP request
	log.Printf("Calling webhook %s %s", method, url)
	return true
}
