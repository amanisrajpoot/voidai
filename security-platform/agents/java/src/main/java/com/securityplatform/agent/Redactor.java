package com.securityplatform.agent;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.regex.Pattern;

/**
 * Redacts PII from data structures.
 */
public class Redactor {
    private static final List<String> DEFAULT_BLOCKLIST = List.of(
        "password", "passwd", "pwd", "secret", "token", "api_key", "apikey",
        "credit_card", "card_number", "cvv", "ssn", "social_security", "pin", "passcode"
    );

    private List<Pattern> patterns = new ArrayList<>();

    public Redactor(List<RedactionRule> rules) {
        // Add default patterns
        for (String field : DEFAULT_BLOCKLIST) {
            patterns.add(Pattern.compile("(?i).*" + Pattern.quote(field) + ".*"));
        }

        // Add custom rules
        if (rules != null) {
            for (RedactionRule rule : rules) {
                patterns.add(Pattern.compile(rule.getPattern(), Pattern.CASE_INSENSITIVE));
            }
        }
    }

    public Object redact(Object data) {
        if (data == null) {
            return null;
        }

        if (data instanceof Map) {
            return redactMap((Map<String, Object>) data);
        } else if (data instanceof List) {
            return redactList((List<Object>) data);
        } else if (data instanceof String) {
            return redactString((String) data);
        }

        return data;
    }

    private Map<String, Object> redactMap(Map<String, Object> map) {
        Map<String, Object> redacted = new HashMap<>();
        for (Map.Entry<String, Object> entry : map.entrySet()) {
            String key = entry.getKey();
            Object value = entry.getValue();

            if (shouldRedact(key)) {
                redacted.put(key, "***REDACTED***");
            } else {
                redacted.put(key, redact(value));
            }
        }
        return redacted;
    }

    private List<Object> redactList(List<Object> list) {
        List<Object> redacted = new ArrayList<>();
        for (Object item : list) {
            redacted.add(redact(item));
        }
        return redacted;
    }

    private String redactString(String str) {
        // Simple string redaction - could be enhanced with regex matching
        return str;
    }

    private boolean shouldRedact(String fieldName) {
        for (Pattern pattern : patterns) {
            if (pattern.matcher(fieldName).matches()) {
                return true;
            }
        }
        return false;
    }
}
