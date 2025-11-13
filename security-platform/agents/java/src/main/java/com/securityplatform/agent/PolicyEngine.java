package com.securityplatform.agent;

import java.util.Map;

/**
 * Policy Engine for local policy evaluation
 */
public class PolicyEngine {
    private AgentConfig config;

    public PolicyEngine(AgentConfig config) {
        this.config = config;
    }

    public PolicyResult checkPolicy(String policyName, Map<String, Object> context) {
        // In observe mode, always allow but record
        if ("observe".equals(config.getLocalPolicy()) || config.getLocalPolicy() == null) {
            return new PolicyResult(true, "observe_mode");
        }

        // Simple policy checks (can be extended)
        // For now, always allow in block mode unless explicitly configured otherwise
        return new PolicyResult(true, "allowed");
    }

    public static class PolicyResult {
        private boolean allowed;
        private String reason;

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
