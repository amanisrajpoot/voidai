package com.securityplatform.agent

class PolicyEngine(private val config: PolicyConfig) {
    fun check(action: String, context: Map<String, Any>): PolicyResult {
        val mode = config.mode.lowercase()

        if (mode == "block") {
            val rules = config.rules
            if (rules != null) {
                for (rule in rules) {
                    val ruleAction = rule["action"] as? String
                    if (ruleAction != null && ruleAction.equals(action, ignoreCase = true)) {
                        val effect = rule["effect"] as? String
                        if (effect != null && effect.lowercase() == "deny") {
                            val ruleId = rule["id"] as? String ?: "unknown"
                            return PolicyResult(
                                allowed = false,
                                reason = "Blocked by policy rule: $ruleId"
                            )
                        }
                    }
                }
            }
            return PolicyResult(allowed = true)
        } else {
            return PolicyResult(
                allowed = true,
                reason = "Observe mode - action allowed"
            )
        }
    }
}
