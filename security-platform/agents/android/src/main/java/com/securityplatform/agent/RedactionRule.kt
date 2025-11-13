package com.securityplatform.agent

/**
 * Redaction rule configuration.
 */
data class RedactionRule(
    val pattern: String,
    val replacement: String? = null,
    val field: String? = null
)
