namespace SecurityPlatform.Agent;

/// <summary>
/// Policy Engine for local policy evaluation
/// </summary>
public class PolicyEngine
{
    private readonly AgentConfig _config;

    public PolicyEngine(AgentConfig config)
    {
        _config = config;
    }

    public PolicyResult CheckPolicy(string policyName, Dictionary<string, object> context)
    {
        // In observe mode, always allow but record
        if ("observe".Equals(_config.LocalPolicy, StringComparison.OrdinalIgnoreCase) || 
            string.IsNullOrEmpty(_config.LocalPolicy))
        {
            return new PolicyResult(true, "observe_mode");
        }

        // Simple policy checks (can be extended)
        return new PolicyResult(true, "allowed");
    }

    public class PolicyResult
    {
        public bool IsAllowed { get; }
        public string Reason { get; }

        public PolicyResult(bool isAllowed, string reason)
        {
            IsAllowed = isAllowed;
            Reason = reason;
        }
    }
}
