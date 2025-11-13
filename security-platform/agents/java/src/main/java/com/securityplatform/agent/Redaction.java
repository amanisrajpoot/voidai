package com.securityplatform.agent;

import java.util.*;
import java.util.regex.Pattern;

/**
 * PII Redaction utility
 */
public class Redaction {
    private static final List<Pattern> DEFAULT_PATTERNS = Arrays.asList(
        Pattern.compile("(?i)(password|passwd|pwd)\\s*[:=]\\s*([^\\s,}]+)", Pattern.CASE_INSENSITIVE),
        Pattern.compile("(?i)(api[_-]?key|apikey)\\s*[:=]\\s*([^\\s,}]+)", Pattern.CASE_INSENSITIVE),
        Pattern.compile("(?i)(token|bearer)\\s*[:=]\\s*([^\\s,}]+)", Pattern.CASE_INSENSITIVE),
        Pattern.compile("\\b\\d{4}[\\s-]?\\d{4}[\\s-]?\\d{4}[\\s-]?\\d{4}\\b"), // Credit card
        Pattern.compile("\\b\\d{3}-\\d{2}-\\d{4}\\b"), // SSN
        Pattern.compile("\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b") // Email (optional)
    );

    private List<Pattern> patterns;

    public Redaction(List<String> customRules) {
        this.patterns = new ArrayList<>(DEFAULT_PATTERNS);
        if (customRules != null) {
            for (String rule : customRules) {
                try {
                    this.patterns.add(Pattern.compile(rule, Pattern.CASE_INSENSITIVE));
                } catch (Exception e) {
                    // Invalid regex, skip
                }
            }
        }
    }

    public String redact(String input) {
        if (input == null || input.isEmpty()) {
            return input;
        }

        String result = input;
        for (Pattern pattern : patterns) {
            result = pattern.matcher(result).replaceAll("[REDACTED]");
        }
        return result;
    }

    public Map<String, Object> redact(Map<String, Object> data) {
        if (data == null) {
            return null;
        }

        Map<String, Object> redacted = new HashMap<>();
        for (Map.Entry<String, Object> entry : data.entrySet()) {
            String key = entry.getKey();
            Object value = entry.getValue();

            // Check if key itself should be redacted
            if (shouldRedactKey(key)) {
                redacted.put(key, "[REDACTED]");
            } else if (value instanceof String) {
                redacted.put(key, redact((String) value));
            } else if (value instanceof Map) {
                @SuppressWarnings("unchecked")
                Map<String, Object> mapValue = (Map<String, Object>) value;
                redacted.put(key, redact(mapValue));
            } else {
                redacted.put(key, value);
            }
        }
        return redacted;
    }

    private boolean shouldRedactKey(String key) {
        String lowerKey = key.toLowerCase();
        return lowerKey.contains("password") || 
               lowerKey.contains("passwd") || 
               lowerKey.contains("pwd") ||
               lowerKey.contains("api_key") ||
               lowerKey.contains("apikey") ||
               lowerKey.contains("token") ||
               lowerKey.contains("secret");
    }
}
