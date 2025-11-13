package com.securityplatform.agent

import java.util.regex.Pattern

/**
 * PII Redaction utility
 */
class Redaction(private val customRules: List<String>? = null) {
    private val patterns: List<Pattern>

    companion object {
        private val defaultPatterns = listOf(
            Pattern.compile("(?i)(password|passwd|pwd)\\s*[:=]\\s*([^\\s,}]+)", Pattern.CASE_INSENSITIVE),
            Pattern.compile("(?i)(api[_-]?key|apikey)\\s*[:=]\\s*([^\\s,}]+)", Pattern.CASE_INSENSITIVE),
            Pattern.compile("(?i)(token|bearer)\\s*[:=]\\s*([^\\s,}]+)", Pattern.CASE_INSENSITIVE),
            Pattern.compile("\\b\\d{4}[\\s-]?\\d{4}[\\s-]?\\d{4}[\\s-]?\\d{4}\\b"), // Credit card
            Pattern.compile("\\b\\d{3}-\\d{2}-\\d{4}\\b") // SSN
        )
    }

    init {
        val compiledPatterns = mutableListOf<Pattern>()
        compiledPatterns.addAll(defaultPatterns)

        customRules?.forEach { rule ->
            try {
                compiledPatterns.add(Pattern.compile(rule, Pattern.CASE_INSENSITIVE))
            } catch (e: Exception) {
                // Invalid regex, skip
            }
        }

        this.patterns = compiledPatterns
    }

    fun redact(input: String): String {
        if (input.isEmpty()) return input

        var result = input
        patterns.forEach { pattern ->
            result = pattern.matcher(result).replaceAll("[REDACTED]")
        }
        return result
    }

    fun redact(data: Map<String, Any?>): Map<String, Any?> {
        val redacted = mutableMapOf<String, Any?>()

        data.forEach { (key, value) ->
            if (shouldRedactKey(key)) {
                redacted[key] = "[REDACTED]"
            } else when (value) {
                is String -> redacted[key] = redact(value)
                is Map<*, *> -> {
                    @Suppress("UNCHECKED_CAST")
                    redacted[key] = redact(value as Map<String, Any?>)
                }
                else -> redacted[key] = value
            }
        }

        return redacted
    }

    private fun shouldRedactKey(key: String): Boolean {
        val lowerKey = key.lowercase()
        return lowerKey.contains("password") ||
               lowerKey.contains("passwd") ||
               lowerKey.contains("pwd") ||
               lowerKey.contains("api_key") ||
               lowerKey.contains("apikey") ||
               lowerKey.contains("token") ||
               lowerKey.contains("secret")
    }
}
