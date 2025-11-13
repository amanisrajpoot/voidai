package com.securityplatform.agent;

import java.util.List;
import java.util.Map;

/**
 * Policy engine for local enforcement.
 */
public class PolicyEngine {
    private AgentConfig.PolicyConfig config;

    public PolicyEngine(AgentConfig.PolicyConfig config) {
        this.config = config != null ? config : new AgentConfig.PolicyConfig();
    }

    public PolicyResult check(String action, Map<String, Object> context) {
        String mode = config.getMode();
        
        if ("block".equals(mode)) {
            // Check rules
            List<Map<String, Object>> rules = config.getRules();
            if (rules != null) {
                for (Map<String, Object> rule : rules) {
                    String ruleAction = (String) rule.get("action");
                    String effect = (String) rule.get("effect");
                    
                    if (ruleAction != null && ruleAction.equals(action)) {
                        if ("deny".equals(effect)) {
                            return new PolicyResult(false, "Blocked by policy rule: " + rule.get("id"));
                        }
                    }
                }
            }
            return new PolicyResult(true, null);
        } else {
            // Observe mode - always allow but log
            return new PolicyResult(true, "Observe mode - action allowed");
        }
    }
}
