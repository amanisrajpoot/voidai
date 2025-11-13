import Foundation

class PolicyEngine {
    private let config: PolicyConfig

    init(config: PolicyConfig) {
        self.config = config
    }

    func check(action: String, context: [String: Any]) -> PolicyResult {
        let mode = config.mode.lowercased()

        if mode == "block" {
            if let rules = config.rules {
                for rule in rules {
                    if let ruleAction = rule["action"] as? String,
                       ruleAction.lowercased() == action.lowercased() {
                        if let effect = rule["effect"] as? String,
                           effect.lowercased() == "deny" {
                            let ruleID = rule["id"] as? String ?? "unknown"
                            return PolicyResult(allowed: false,
                                               reason: "Blocked by policy rule: \(ruleID)")
                        }
                    }
                }
            }
            return PolicyResult(allowed: true)
        } else {
            return PolicyResult(allowed: true, reason: "Observe mode - action allowed")
        }
    }
}
