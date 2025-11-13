package com.securityplatform.agent

import java.util.regex.Pattern

class Redactor(private val rules: List<RedactionRule>) {
    companion object {
        private const val DEFAULT_REPLACEMENT = "***REDACTED***"
    }

    fun redact(data: Any?): Any? {
        if (data == null) return null

        return when (data) {
            is Map<*, *> -> redactMap(data as Map<String, Any>)
            is List<*> -> data.map { redact(it) }
            is String -> redactString(data)
            else -> data
        }
    }

    private fun redactMap(map: Map<String, Any>): Map<String, Any> {
        val result = mutableMapOf<String, Any>()
        for ((key, value) in map) {
            val shouldRedact = rules.any { rule ->
                (rule.field != null && rule.field.equals(key, ignoreCase = true)) ||
                (rule.pattern != null && Pattern.compile(rule.pattern, Pattern.CASE_INSENSITIVE)
                    .matcher(key).find())
            }

            result[key] = if (shouldRedact) {
                rules.firstOrNull { it.field == key || it.pattern != null }?.replacement
                    ?: DEFAULT_REPLACEMENT
            } else {
                redact(value) ?: value
            }
        }
        return result
    }

    private fun redactString(str: String): String {
        var result = str
        for (rule in rules) {
            if (rule.pattern != null) {
                val replacement = rule.replacement ?: DEFAULT_REPLACEMENT
                result = Pattern.compile(rule.pattern, Pattern.CASE_INSENSITIVE)
                    .matcher(result)
                    .replaceAll(replacement)
            }
        }
        return result
    }
}
