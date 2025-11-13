using System.Collections.Generic;

namespace SecurityPlatform.Agent
{
    /// <summary>
    /// Policy engine for local policy enforcement.
    /// </summary>
    public class PolicyEngine
    {
        private readonly PolicyConfig _policy;

        public PolicyEngine(PolicyConfig policy)
        {
            _policy = policy ?? new PolicyConfig();
        }

        public PolicyResult Check(string action, Dictionary<string, object>? context = null)
        {
            // Simple policy check - can be extended with rule evaluation
            if (_policy.Mode == "block")
            {
                // In block mode, check rules and potentially block
                return new PolicyResult(true, null);
            }
            // In observe mode, always allow but log
            return new PolicyResult(true, null);
        }
    }

    /// <summary>
    /// Result of a policy check.
    /// </summary>
    public class PolicyResult
    {
        public bool Allowed { get; }
        public string? Reason { get; }

        public PolicyResult(bool allowed, string? reason = null)
        {
            Allowed = allowed;
            Reason = reason;
        }
    }
}
