package com.securityplatform.agent;

import java.util.Map;

/**
 * Policy engine for local rule evaluation.
 */
public class PolicyEngine {
    private AgentConfig.PolicyConfig policyConfig;

    public PolicyEngine(AgentConfig.PolicyConfig policyConfig) {
        this.policyConfig = policyConfig != null ? policyConfig : new AgentConfig.PolicyConfig();
    }

    public PolicyResult check(String action, Map<String, Object> context) {
        // In observe mode, always allow but log
        if ("observe".equals(policyConfig.getMode())) {
            return new PolicyResult(true, "observe_mode");
        }

        // In block mode, evaluate rules
        if ("block".equals(policyConfig.getMode())) {
            // Simple implementation - can be enhanced with rule evaluation
            for (PolicyRule rule : policyConfig.getRules()) {
                if (rule.getAction().equals(action)) {
                    return evaluateRule(rule, context);
                }
            }
            // Default allow if no matching rule
            return new PolicyResult(true, "no_matching_rule");
        }

        return new PolicyResult(true, "default_allow");
    }

    private PolicyResult evaluateRule(PolicyRule rule, Map<String, Object> context) {
        // Simple evaluation - can be enhanced with expression engine
        if ("deny".equals(rule.getEffect())) {
            return new PolicyResult(false, "rule_denied: " + rule.getId());
        }
        return new PolicyResult(true, "rule_allowed: " + rule.getId());
    }
}
