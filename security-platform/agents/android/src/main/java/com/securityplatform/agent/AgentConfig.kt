package com.securityplatform.agent

/**
 * Configuration for the Security Platform Agent.
 */
data class AgentConfig(
    val controlPlaneUrl: String,
    val authKey: String? = null,
    val serviceName: String? = null,
    val environment: String? = null,
    val namespace: String? = null,
    val otlpEndpoint: String? = null,
    val telemetry: TelemetryConfig? = null,
    val redactionRules: List<RedactionRule>? = null,
    val policy: PolicyConfig? = null,
    val security: SecurityConfig? = null
)

data class TelemetryConfig(
    val batchSize: Int = 100,
    val batchTimeout: String = "5s",
    val exportTimeout: String = "30s",
    val maxQueueSize: Int = 2048
)

data class PolicyConfig(
    val mode: String = "observe", // "observe" or "block"
    val autoEnableBlocking: Boolean = false,
    val observePeriodHours: Int = 48,
    val rules: List<Map<String, Any>>? = null
)

data class SecurityConfig(
    val mtlsEnabled: Boolean = true,
    val certificatePath: String? = null,
    val keyPath: String? = null,
    val caBundlePath: String? = null
)

data class RedactionRule(
    val pattern: String? = null,
    val replacement: String? = null,
    val field: String? = null
)

data class PolicyResult(
    val allowed: Boolean,
    val reason: String? = null
)
