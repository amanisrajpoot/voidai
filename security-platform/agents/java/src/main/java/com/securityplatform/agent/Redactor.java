package com.securityplatform.agent;

import java.util.List;
import java.util.Map;
import java.util.regex.Pattern;

/**
 * Redacts PII from data structures.
 */
public class Redactor {
    private static final String DEFAULT_REPLACEMENT = "***REDACTED***";
    private List<RedactionRule> rules;

    public Redactor(List<RedactionRule> rules) {
        this.rules = rules != null ? rules : List.of();
    }

    public Object redact(Object data) {
        if (data == null) {
            return null;
        }

        if (data instanceof Map) {
            return redactMap((Map<String, Object>) data);
        } else if (data instanceof List) {
            return redactList((List<?>) data);
        } else if (data instanceof String) {
            return redactString((String) data);
        }

        return data;
    }

    private Map<String, Object> redactMap(Map<String, Object> map) {
        Map<String, Object> result = new java.util.HashMap<>(map);
        for (Map.Entry<String, Object> entry : map.entrySet()) {
            String key = entry.getKey();
            Object value = entry.getValue();

            // Check if key matches any redaction rule
            boolean shouldRedact = false;
            for (RedactionRule rule : rules) {
                if (rule.getField() != null && rule.getField().equals(key)) {
                    shouldRedact = true;
                    break;
                }
                if (rule.getPattern() != null && Pattern.compile(rule.getPattern(), Pattern.CASE_INSENSITIVE).matcher(key).find()) {
                    shouldRedact = true;
                    break;
                }
            }

            if (shouldRedact) {
                result.put(key, DEFAULT_REPLACEMENT);
            } else {
                result.put(key, redact(value));
            }
        }
        return result;
    }

    private List<?> redactList(List<?> list) {
        return list.stream().map(this::redact).collect(java.util.stream.Collectors.toList());
    }

    private String redactString(String str) {
        for (RedactionRule rule : rules) {
            if (rule.getPattern() != null) {
                Pattern pattern = Pattern.compile(rule.getPattern(), Pattern.CASE_INSENSITIVE);
                String replacement = rule.getReplacement() != null ? rule.getReplacement() : DEFAULT_REPLACEMENT;
                str = pattern.matcher(str).replaceAll(replacement);
            }
        }
        return str;
    }
}
