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
    public void setAllowed(boolean allowed) { this.allowed = allowed; }

    public String getReason() { return reason; }
    public void setReason(String reason) { this.reason = reason; }
}
