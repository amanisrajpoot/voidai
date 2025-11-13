package com.securityplatform.agent

/**
 * Redacts PII from data structures.
 */
class Redactor(private val blocklist: List<String>) {
    companion object {
        private const val REDACTED = "[REDACTED]"
        private val DEFAULT_BLOCKLIST = listOf(
            "password", "passwd", "pwd", "secret", "token",
            "api_key", "apikey", "credit_card", "card_number",
            "cvv", "ssn", "social_security", "pin", "passcode"
        )
    }

    constructor() : this(DEFAULT_BLOCKLIST)

    fun redact(data: Any?): Any? {
        if (data == null) {
            return null
        }

        return when (data) {
            is Map<*, *> -> redactMap(data as Map<String, Any?>)
            is List<*> -> data.map { redact(it) }
            is String -> redactString(data)
            else -> data
        }
    }

    private fun redactMap(map: Map<String, Any?>): Map<String, Any?> {
        val redacted = mutableMapOf<String, Any?>()
        for ((key, value) in map) {
            if (shouldRedact(key)) {
                redacted[key] = REDACTED
            } else {
                redacted[key] = redact(value)
            }
        }
        return redacted
    }

    private fun redactString(str: String): String {
        // Simple string redaction - could be enhanced with regex patterns
        return str
    }

    private fun shouldRedact(key: String?): Boolean {
        if (key == null) {
            return false
        }

        val lowerKey = key.lowercase()
        return blocklist.any { pattern ->
            if (pattern.contains("*")) {
                val regex = "^${pattern.replace("*", ".*")}$".toRegex(RegexOption.IGNORE_CASE)
                regex.matches(lowerKey)
            } else {
                lowerKey.contains(pattern.lowercase())
            }
        }
    }
}
