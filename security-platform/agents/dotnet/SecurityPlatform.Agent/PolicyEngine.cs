using System;
using System.Collections.Generic;
using System.Linq;

namespace SecurityPlatform.Agent
{
    /// <summary>
    /// Policy engine for local enforcement.
    /// </summary>
    public class PolicyEngine
    {
        private readonly AgentConfig.PolicyConfig _config;

        public PolicyEngine(AgentConfig.PolicyConfig config)
        {
            _config = config ?? new AgentConfig.PolicyConfig();
        }

        public PolicyResult Check(string action, Dictionary<string, object> context)
        {
            var mode = _config.Mode ?? "observe";

            if (mode == "block")
            {
                var rules = _config.Rules;
                if (rules != null)
                {
                    foreach (var rule in rules)
                    {
                        if (rule.TryGetValue("action", out var ruleActionObj) &&
                            ruleActionObj is string ruleAction &&
                            ruleAction.Equals(action, StringComparison.OrdinalIgnoreCase))
                        {
                            if (rule.TryGetValue("effect", out var effectObj) &&
                                effectObj is string effect &&
                                effect == "deny")
                            {
                                var ruleId = rule.TryGetValue("id", out var idObj) ? idObj.ToString() : "unknown";
                                return new PolicyResult(false, $"Blocked by policy rule: {ruleId}");
                            }
                        }
                    }
                }
                return new PolicyResult(true);
            }
            else
            {
                return new PolicyResult(true, "Observe mode - action allowed");
            }
        }
    }
}
