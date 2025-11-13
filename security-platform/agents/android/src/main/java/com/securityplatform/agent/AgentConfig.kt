package com.securityplatform.agent

import android.content.Context

/**
 * Configuration for the Security Platform Agent.
 */
data class AgentConfig(
    var controlPlaneUrl: String? = null,
    var authKey: String? = null,
    var serviceName: String? = "android-app",
    var environment: String? = "production",
    var otlpEndpoint: String? = null,
    var telemetry: TelemetryConfig? = null,
    var redactionRules: List<String>? = null,
    var policy: PolicyConfig? = null
) {
    companion object {
        /**
         * Create configuration from environment variables or defaults.
         */
        @JvmStatic
        fun fromEnvironment(context: Context? = null): AgentConfig {
            val controlPlaneUrl = System.getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL")
                ?: "https://api.securityplatform.com"
            val authKey = System.getenv("SECURITY_PLATFORM_AUTH_KEY")
            
            return AgentConfig(
                controlPlaneUrl = controlPlaneUrl,
                authKey = authKey,
                telemetry = TelemetryConfig(),
                policy = PolicyConfig()
            )
        }
    }
}

/**
 * Telemetry configuration.
 */
data class TelemetryConfig(
    var batchSize: Int = 100,
    var batchTimeoutSeconds: Int = 5,
    var exportTimeoutSeconds: Int = 30,
    var maxQueueSize: Int = 2048
)

/**
 * Policy configuration.
 */
data class PolicyConfig(
    var mode: String = "observe" // "observe" or "block"
)
