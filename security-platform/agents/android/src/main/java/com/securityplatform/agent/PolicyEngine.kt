package com.securityplatform.agent

/**
 * Policy engine for local rule evaluation.
 */
class PolicyEngine(private val config: AgentConfig.PolicyConfig) {
    fun check(action: String, context: Map<String, Any>): PolicyResult {
        // In observe mode, always allow but log
        if (config.mode == "observe") {
            return PolicyResult(allowed = true, reason = "observe_mode")
        }

        // In block mode, evaluate rules
        if (config.mode == "block") {
            for (rule in config.rules) {
                if (rule.action == action) {
                    return evaluateRule(rule, context)
                }
            }
            return PolicyResult(allowed = true, reason = "no_matching_rule")
        }

        return PolicyResult(allowed = true, reason = "default_allow")
    }

    private fun evaluateRule(rule: PolicyRule, context: Map<String, Any>): PolicyResult {
        if (rule.effect == "deny") {
            return PolicyResult(allowed = false, reason = "rule_denied: ${rule.id}")
        }
        return PolicyResult(allowed = true, reason = "rule_allowed: ${rule.id}")
    }
}

/**
 * Policy rule definition.
 */
data class PolicyRule(
    val id: String,
    val action: String,
    val condition: String,
    val effect: String // "allow" or "deny"
)

/**
 * Result of a policy check.
 */
data class PolicyResult(
    val allowed: Boolean,
    val reason: String
)
