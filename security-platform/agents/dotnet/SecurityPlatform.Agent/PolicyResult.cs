namespace SecurityPlatform.Agent
{
    /// <summary>
    /// Result of a policy check.
    /// </summary>
    public class PolicyResult
    {
        public bool Allowed { get; }
        public string Reason { get; }

        public PolicyResult(bool allowed, string reason)
        {
            Allowed = allowed;
            Reason = reason;
        }
    }
}
