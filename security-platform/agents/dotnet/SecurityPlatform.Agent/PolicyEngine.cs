using System.Collections.Generic;

namespace SecurityPlatform.Agent
{
    /// <summary>
    /// Policy engine for local rule evaluation.
    /// </summary>
    public class PolicyEngine
    {
        private readonly AgentConfig.PolicyConfig _policyConfig;

        public PolicyEngine(AgentConfig.PolicyConfig policyConfig)
        {
            _policyConfig = policyConfig ?? new AgentConfig.PolicyConfig();
        }

        public PolicyResult Check(string action, Dictionary<string, object> context)
        {
            // In observe mode, always allow but log
            if (_policyConfig.Mode == "observe")
            {
                return new PolicyResult(true, "observe_mode");
            }

            // In block mode, evaluate rules
            if (_policyConfig.Mode == "block")
            {
                // Simple implementation - can be enhanced with rule evaluation
                foreach (var rule in _policyConfig.Rules)
                {
                    if (rule.Action == action)
                    {
                        return EvaluateRule(rule, context);
                    }
                }
                // Default allow if no matching rule
                return new PolicyResult(true, "no_matching_rule");
            }

            return new PolicyResult(true, "default_allow");
        }

        private PolicyResult EvaluateRule(PolicyRule rule, Dictionary<string, object> context)
        {
            // Simple evaluation - can be enhanced with expression engine
            if (rule.Effect == "deny")
            {
                return new PolicyResult(false, $"rule_denied: {rule.Id}");
            }
            return new PolicyResult(true, $"rule_allowed: {rule.Id}");
        }
    }
}
