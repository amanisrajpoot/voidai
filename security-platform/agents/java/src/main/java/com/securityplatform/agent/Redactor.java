package com.securityplatform.agent;

import java.util.List;
import java.util.Map;
import java.util.regex.Pattern;

/**
 * Redacts PII from data structures.
 */
public class Redactor {
    private final List<String> blocklist;
    private static final String REDACTED = "[REDACTED]";

    public Redactor(List<String> blocklist) {
        this.blocklist = blocklist != null ? blocklist : getDefaultBlocklist();
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

    @SuppressWarnings("unchecked")
    private Map<String, Object> redactMap(Map<String, Object> map) {
        Map<String, Object> redacted = new java.util.HashMap<>();
        for (Map.Entry<String, Object> entry : map.entrySet()) {
            String key = entry.getKey();
            Object value = entry.getValue();

            if (shouldRedact(key)) {
                redacted.put(key, REDACTED);
            } else {
                redacted.put(key, redact(value));
            }
        }
        return redacted;
    }

    private List<Object> redactList(List<Object> list) {
        List<Object> redacted = new java.util.ArrayList<>();
        for (Object item : list) {
            redacted.add(redact(item));
        }
        return redacted;
    }

    private String redactString(String str) {
        // Simple string redaction - could be enhanced with regex patterns
        return str;
    }

    private boolean shouldRedact(String key) {
        if (key == null) {
            return false;
        }

        String lowerKey = key.toLowerCase();
        for (String pattern : blocklist) {
            if (pattern.contains("*")) {
                String regex = pattern.replace("*", ".*");
                if (Pattern.matches(regex, lowerKey)) {
                    return true;
                }
            } else if (lowerKey.contains(pattern.toLowerCase())) {
                return true;
            }
        }
        return false;
    }

    private static List<String> getDefaultBlocklist() {
        return List.of(
            "password", "passwd", "pwd", "secret", "token",
            "api_key", "apikey", "credit_card", "card_number",
            "cvv", "ssn", "social_security", "pin", "passcode"
        );
    }
}
