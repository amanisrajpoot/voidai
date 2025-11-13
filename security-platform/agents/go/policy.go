package agent

import (
	"fmt"
	"strings"
)

// PolicyEngine enforces local policies.
type PolicyEngine struct {
	config *PolicyConfig
}

// NewPolicyEngine creates a new policy engine.
func NewPolicyEngine(config *PolicyConfig) *PolicyEngine {
	if config == nil {
		config = &PolicyConfig{Mode: "observe"}
	}
	return &PolicyEngine{config: config}
}

// Check checks if an action is allowed by policy.
func (p *PolicyEngine) Check(action string, context map[string]interface{}) *PolicyResult {
	mode := p.config.Mode
	if mode == "" {
		mode = "observe"
	}

	if strings.ToLower(mode) == "block" {
		rules := p.config.Rules
		if rules != nil {
			for _, rule := range rules {
				ruleAction, ok := rule["action"].(string)
				if !ok {
					continue
				}

				if strings.EqualFold(ruleAction, action) {
					effect, ok := rule["effect"].(string)
					if ok && strings.ToLower(effect) == "deny" {
						ruleID := "unknown"
						if id, ok := rule["id"].(string); ok {
							ruleID = id
						}
						return &PolicyResult{
							Allowed: false,
							Reason:  fmt.Sprintf("Blocked by policy rule: %s", ruleID),
						}
					}
				}
			}
		}
		return &PolicyResult{Allowed: true}
	}

	return &PolicyResult{
		Allowed: true,
		Reason:  "Observe mode - action allowed",
	}
}
