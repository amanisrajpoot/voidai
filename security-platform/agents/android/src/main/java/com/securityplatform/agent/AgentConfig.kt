package com.securityplatform.agent

/**
 * Configuration for the Security Platform Agent.
 */
data class AgentConfig(
    val controlPlaneUrl: String = "https://api.securityplatform.com",
    val authKey: String? = null,
    val serviceName: String = "android-service",
    val environment: String = "production",
    val namespace: String = "default",
    val otlpEndpoint: String? = null,
    val telemetry: TelemetryConfig = TelemetryConfig(),
    val redactionRules: List<RedactionRule> = emptyList(),
    val policy: PolicyConfig = PolicyConfig(),
    val security: SecurityConfig = SecurityConfig()
) {
    companion object {
        fun default(): AgentConfig {
            return AgentConfig(
                controlPlaneUrl = System.getenv("SECURITY_PLATFORM_URL")
                    ?: "https://api.securityplatform.com",
                authKey = System.getenv("SECURITY_PLATFORM_AUTH_KEY"),
                serviceName = System.getenv("SERVICE_NAME") ?: "android-service",
                environment = System.getenv("ENVIRONMENT") ?: "production"
            )
        }
    }

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
        val rules: List<PolicyRule> = emptyList()
    )

    data class SecurityConfig(
        val mtlsEnabled: Boolean = false,
        val certificatePath: String? = null,
        val keyPath: String? = null,
        val caBundlePath: String? = null
    )
}
