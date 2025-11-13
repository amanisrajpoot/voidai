namespace SecurityPlatform.Agent
{
    /// <summary>
    /// Redaction rule configuration.
    /// </summary>
    public class RedactionRule
    {
        public string? Pattern { get; set; }
        public string? Replacement { get; set; }
        public string? Field { get; set; }
    }
}
