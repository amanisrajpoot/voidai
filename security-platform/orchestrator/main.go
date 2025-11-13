package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Orchestrator handles automated security responses
type Orchestrator struct {
	router *mux.Router
	port   string
}

// Playbook represents an automated response playbook
type Playbook struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Triggers    []Trigger              `json:"triggers"`
	Actions     []Action               `json:"actions"`
	Enabled     bool                   `json:"enabled"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type Trigger struct {
	Type     string                 `json:"type"`
	Condition map[string]interface{} `json:"condition"`
}

type Action struct {
	Type    string                 `json:"type"`
	Config  map[string]interface{} `json:"config"`
	Timeout int                    `json:"timeout"`
}

func NewOrchestrator(port string) *Orchestrator {
	o := &Orchestrator{
		router: mux.NewRouter(),
		port:   port,
	}
	o.setupRoutes()
	return o
}

func (o *Orchestrator) setupRoutes() {
	api := o.router.PathPrefix("/api/v1").Subrouter()

	// Webhook endpoint for receiving incidents
	api.HandleFunc("/webhooks/incidents", o.handleIncidentWebhook).Methods("POST")

	// Playbooks management
	api.HandleFunc("/playbooks", o.handlePlaybooks).Methods("GET", "POST")
	api.HandleFunc("/playbooks/{id}", o.handlePlaybook).Methods("GET", "PUT", "DELETE")

	// Execute playbook manually
	api.HandleFunc("/playbooks/{id}/execute", o.handleExecutePlaybook).Methods("POST")

	// Health check
	o.router.HandleFunc("/health", o.handleHealth).Methods("GET")
}

func (o *Orchestrator) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (o *Orchestrator) handleIncidentWebhook(w http.ResponseWriter, r *http.Request) {
	var incident map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&incident); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received incident webhook: %v", incident)

	// Find matching playbooks
	playbooks := o.getMatchingPlaybooks(incident)

	// Execute playbooks
	for _, playbook := range playbooks {
		if err := o.executePlaybook(playbook, incident); err != nil {
			log.Printf("Failed to execute playbook %s: %v", playbook.ID, err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "processed"})
}

func (o *Orchestrator) getMatchingPlaybooks(incident map[string]interface{}) []Playbook {
	// In production, query database for matching playbooks
	// For now, return empty list
	return []Playbook{}
}

func (o *Orchestrator) executePlaybook(playbook Playbook, incident map[string]interface{}) error {
	log.Printf("Executing playbook %s for incident", playbook.ID)

	for _, action := range playbook.Actions {
		if err := o.executeAction(action, incident); err != nil {
			return fmt.Errorf("failed to execute action %s: %w", action.Type, err)
		}
	}

	return nil
}

func (o *Orchestrator) executeAction(action Action, incident map[string]interface{}) error {
	switch action.Type {
	case "block_ip_kong":
		return o.blockIPKong(action.Config, incident)
	case "block_ip_nginx":
		return o.blockIPNginx(action.Config, incident)
	case "revoke_token":
		return o.revokeToken(action.Config, incident)
	case "slack_notify":
		return o.slackNotify(action.Config, incident)
	case "jira_create":
		return o.jiraCreate(action.Config, incident)
	case "pagerduty_alert":
		return o.pagerdutyAlert(action.Config, incident)
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

func (o *Orchestrator) blockIPKong(config map[string]interface{}, incident map[string]interface{}) error {
	ip := incident["remote_addr"].(string)
	kongAdminURL := config["kong_admin_url"].(string)

	// Call Kong Admin API to block IP
	log.Printf("Blocking IP %s via Kong at %s", ip, kongAdminURL)
	// Implementation: HTTP POST to Kong Admin API
	return nil
}

func (o *Orchestrator) blockIPNginx(config map[string]interface{}, incident map[string]interface{}) error {
	ip := incident["remote_addr"].(string)
	nginxConfigPath := config["config_path"].(string)

	log.Printf("Blocking IP %s via Nginx config at %s", ip, nginxConfigPath)
	// Implementation: Update Nginx config and reload
	return nil
}

func (o *Orchestrator) revokeToken(config map[string]interface{}, incident map[string]interface{}) error {
	token := incident["token"].(string)
	iamURL := config["iam_url"].(string)

	log.Printf("Revoking token via IAM at %s", iamURL)
	// Implementation: Call IAM API to revoke token
	return nil
}

func (o *Orchestrator) slackNotify(config map[string]interface{}, incident map[string]interface{}) error {
	webhookURL := config["webhook_url"].(string)

	log.Printf("Sending Slack notification to %s", webhookURL)
	// Implementation: POST to Slack webhook
	return nil
}

func (o *Orchestrator) jiraCreate(config map[string]interface{}, incident map[string]interface{}) error {
	jiraURL := config["jira_url"].(string)

	log.Printf("Creating Jira ticket at %s", jiraURL)
	// Implementation: Create Jira issue via API
	return nil
}

func (o *Orchestrator) pagerdutyAlert(config map[string]interface{}, incident map[string]interface{}) error {
	apiKey := config["api_key"].(string)

	log.Printf("Sending PagerDuty alert")
	// Implementation: POST to PagerDuty Events API
	return nil
}

func (o *Orchestrator) handlePlaybooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var playbook Playbook
		if err := json.NewDecoder(r.Body).Decode(&playbook); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		playbook.ID = fmt.Sprintf("playbook-%d", time.Now().Unix())
		log.Printf("Creating playbook %s", playbook.ID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(playbook)
		return
	}

	// GET: List playbooks
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]Playbook{})
}

func (o *Orchestrator) handlePlaybook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playbookID := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Playbook{
		ID:   playbookID,
		Name: "Example Playbook",
	})
}

func (o *Orchestrator) handleExecutePlaybook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playbookID := vars["id"]

	var incident map[string]interface{}
	json.NewDecoder(r.Body).Decode(&incident)

	playbook := Playbook{ID: playbookID}
	if err := o.executePlaybook(playbook, incident); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "executed"})
}

func main() {
	port := "8081"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	orchestrator := NewOrchestrator(port)
	log.Printf("Starting orchestrator on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, orchestrator.router))
}
