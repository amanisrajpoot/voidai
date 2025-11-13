namespace SecurityPlatform.Agent
{
    /// <summary>
    /// Policy rule definition.
    /// </summary>
    public class PolicyRule
    {
        public string? Id { get; set; }
        public string? Action { get; set; }
        public string? Condition { get; set; }
        public string? Effect { get; set; } // "allow" or "deny"
    }
}
