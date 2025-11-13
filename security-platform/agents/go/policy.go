package agent

// PolicyEngine evaluates security policies.
type PolicyEngine struct {
	config *PolicyConfig
}

// NewPolicyEngine creates a new PolicyEngine.
func NewPolicyEngine(config *PolicyConfig) *PolicyEngine {
	if config == nil {
		config = &PolicyConfig{
			Mode:               "observe",
			AutoEnableBlocking: false,
			ObservePeriodHours: 48,
		}
	}
	return &PolicyEngine{config: config}
}

// PolicyResult represents the result of a policy check.
type PolicyResult struct {
	Allowed bool
	Reason  string
}

// Check checks if an action is allowed by policy.
func (p *PolicyEngine) Check(action string, context map[string]interface{}) *PolicyResult {
	// In observe mode, always allow but log
	if p.config.Mode == "observe" {
		return &PolicyResult{Allowed: true, Reason: "observe_mode"}
	}

	// In block mode, evaluate rules
	if p.config.Mode == "block" {
		for _, rule := range p.config.Rules {
			if rule.Action == action {
				return p.evaluateRule(&rule, context)
			}
		}
		// Default allow if no matching rule
		return &PolicyResult{Allowed: true, Reason: "no_matching_rule"}
	}

	return &PolicyResult{Allowed: true, Reason: "default_allow"}
}

func (p *PolicyEngine) evaluateRule(rule *PolicyRule, context map[string]interface{}) *PolicyResult {
	// Simple evaluation - can be enhanced with expression engine
	if rule.Effect == "deny" {
		return &PolicyResult{Allowed: false, Reason: "rule_denied: " + rule.ID}
	}
	return &PolicyResult{Allowed: true, Reason: "rule_allowed: " + rule.ID}
}

// PolicyRule represents a policy rule definition.
type PolicyRule struct {
	ID        string `yaml:"id" json:"id"`
	Action    string `yaml:"action" json:"action"`
	Condition string `yaml:"condition" json:"condition"`
	Effect    string `yaml:"effect" json:"effect"` // "allow" or "deny"
}
