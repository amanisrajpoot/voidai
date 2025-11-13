import Foundation

/// Policy engine for local rule evaluation.
public class PolicyEngine {
    private let config: AgentConfig.PolicyConfig

    public init(config: AgentConfig.PolicyConfig) {
        self.config = config
    }

    public func check(action: String, context: [String: Any]) -> PolicyResult {
        // In observe mode, always allow but log
        if config.mode == "observe" {
            return PolicyResult(allowed: true, reason: "observe_mode")
        }

        // In block mode, evaluate rules
        if config.mode == "block" {
            for rule in config.rules {
                if rule.action == action {
                    return evaluateRule(rule, context: context)
                }
            }
            return PolicyResult(allowed: true, reason: "no_matching_rule")
        }

        return PolicyResult(allowed: true, reason: "default_allow")
    }

    private func evaluateRule(_ rule: PolicyRule, context: [String: Any]) -> PolicyResult {
        if rule.effect == "deny" {
            return PolicyResult(allowed: false, reason: "rule_denied: \(rule.id)")
        }
        return PolicyResult(allowed: true, reason: "rule_allowed: \(rule.id)")
    }
}

/// Policy rule definition.
public struct PolicyRule {
    public var id: String
    public var action: String
    public var condition: String
    public var effect: String // "allow" or "deny"

    public init(id: String, action: String, condition: String, effect: String) {
        self.id = id
        self.action = action
        self.condition = condition
        self.effect = effect
    }
}

/// Result of a policy check.
public struct PolicyResult {
    public let allowed: Bool
    public let reason: String

    public init(allowed: Bool, reason: String) {
        self.allowed = allowed
        self.reason = reason
    }
}
