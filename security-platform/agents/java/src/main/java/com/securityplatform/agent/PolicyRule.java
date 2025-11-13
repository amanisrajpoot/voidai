package com.securityplatform.agent;

/**
 * Policy rule definition.
 */
public class PolicyRule {
    private String id;
    private String action;
    private String condition;
    private String effect; // "allow" or "deny"

    public PolicyRule() {}

    public String getId() { return id; }
    public void setId(String id) { this.id = id; }

    public String getAction() { return action; }
    public void setAction(String action) { this.action = action; }

    public String getCondition() { return condition; }
    public void setCondition(String condition) { this.condition = condition; }

    public String getEffect() { return effect; }
    public void setEffect(String effect) { this.effect = effect; }
}
