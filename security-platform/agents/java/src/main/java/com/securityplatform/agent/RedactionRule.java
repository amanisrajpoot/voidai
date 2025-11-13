package com.securityplatform.agent;

/**
 * Redaction rule configuration.
 */
public class RedactionRule {
    private String pattern;
    private String replacement;
    private String field;

    public RedactionRule() {}

    public RedactionRule(String pattern) {
        this.pattern = pattern;
    }

    public String getPattern() { return pattern; }
    public void setPattern(String pattern) { this.pattern = pattern; }

    public String getReplacement() { return replacement; }
    public void setReplacement(String replacement) { this.replacement = replacement; }

    public String getField() { return field; }
    public void setField(String field) { this.field = field; }
}
