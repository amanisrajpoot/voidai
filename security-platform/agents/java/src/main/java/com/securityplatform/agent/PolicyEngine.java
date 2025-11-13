package com.securityplatform.agent;

import java.util.Map;

/**
 * Policy engine for local policy enforcement.
 */
public class PolicyEngine {
    private final AgentConfig.PolicyConfig policy;

    public PolicyEngine(AgentConfig.PolicyConfig policy) {
        this.policy = policy != null ? policy : new AgentConfig.PolicyConfig();
    }

    public PolicyResult check(String action, Map<String, Object> context) {
        // Simple policy check - can be extended with rule evaluation
        if ("block".equalsIgnoreCase(policy.getMode())) {
            // In block mode, check rules and potentially block
            return new PolicyResult(true, null);
        }
        // In observe mode, always allow but log
        return new PolicyResult(true, null);
    }

    public static class PolicyResult {
        private final boolean allowed;
        private final String reason;

        public PolicyResult(boolean allowed, String reason) {
            this.allowed = allowed;
            this.reason = reason;
        }

        public boolean isAllowed() {
            return allowed;
        }

        public String getReason() {
            return reason;
        }
    }
}
