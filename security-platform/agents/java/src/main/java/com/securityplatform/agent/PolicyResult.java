package com.securityplatform.agent;

/**
 * Result of a policy check.
 */
public class PolicyResult {
    private boolean allowed;
    private String reason;

    public PolicyResult(boolean allowed, String reason) {
        this.allowed = allowed;
        this.reason = reason;
    }

    public boolean isAllowed() { return allowed; }
    public String getReason() { return reason; }
}
