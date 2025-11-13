namespace SecurityPlatform.Agent;

/// <summary>
/// Configuration for Security Platform Agent
/// </summary>
public class AgentConfig
{
    public string? ControlPlaneUrl { get; set; }
    public string? AuthKey { get; set; }
    public string? ServiceName { get; set; }
    public string? Version { get; set; }
    public string? Environment { get; set; }
    public string? OtlpEndpoint { get; set; }
    public List<string>? RedactionRules { get; set; }
    public string? LocalPolicy { get; set; } // "observe" or "block"
    public int? TelemetryBatchSize { get; set; }
    public Dictionary<string, string>? CustomAttributes { get; set; }
}
