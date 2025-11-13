package com.securityplatform.agent

import java.util.regex.Pattern

/**
 * Redacts PII from data structures.
 */
class Redactor(private val rules: List<RedactionRule>) {
    companion object {
        private val DEFAULT_BLOCKLIST = listOf(
            "password", "passwd", "pwd", "secret", "token", "api_key", "apikey",
            "credit_card", "card_number", "cvv", "ssn", "social_security", "pin", "passcode"
        )
    }

    private val patterns: List<Pattern> = buildList {
        // Add default patterns
        DEFAULT_BLOCKLIST.forEach { field ->
            add(Pattern.compile("(?i).*${Pattern.quote(field)}.*"))
        }
        // Add custom rules
        rules.forEach { rule ->
            if (rule.pattern.isNotEmpty()) {
                try {
                    add(Pattern.compile("(?i)${rule.pattern}"))
                } catch (e: Exception) {
                    // Invalid pattern, skip
                }
            }
        }
    }

    fun redact(data: Any?): Any? {
        if (data == null) return null

        return when (data) {
            is Map<*, *> -> redactMap(data as Map<String, Any>)
            is List<*> -> redactList(data as List<Any>)
            is String -> redactString(data)
            else -> data
        }
    }

    private fun redactMap(map: Map<String, Any>): Map<String, Any?> {
        return map.mapValues { (key, value) ->
            if (shouldRedact(key)) {
                "***REDACTED***"
            } else {
                redact(value)
            }
        }
    }

    private fun redactList(list: List<Any>): List<Any?> {
        return list.map { redact(it) }
    }

    private fun redactString(str: String): String {
        // Simple string redaction
        return str
    }

    private fun shouldRedact(fieldName: String): Boolean {
        return patterns.any { it.matcher(fieldName).matches() }
    }
}
