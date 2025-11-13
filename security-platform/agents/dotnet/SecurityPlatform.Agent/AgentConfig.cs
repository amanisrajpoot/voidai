using System.Collections.Generic;

namespace SecurityPlatform.Agent
{
    /// <summary>
    /// Configuration for the Security Platform Agent.
    /// </summary>
    public class AgentConfig
    {
        public string ControlPlaneUrl { get; set; } = "https://api.securityplatform.com";
        public string? AuthKey { get; set; }
        public string? ServiceName { get; set; }
        public string? Environment { get; set; }
        public string? Namespace { get; set; }
        public string? OtlpEndpoint { get; set; }
        public TelemetryConfig? Telemetry { get; set; }
        public List<RedactionRule>? RedactionRules { get; set; }
        public PolicyConfig? Policy { get; set; }
        public SecurityConfig? Security { get; set; }

        public class TelemetryConfig
        {
            public int BatchSize { get; set; } = 100;
            public string BatchTimeout { get; set; } = "5s";
            public string ExportTimeout { get; set; } = "30s";
            public int MaxQueueSize { get; set; } = 2048;
        }

        public class PolicyConfig
        {
            public string Mode { get; set; } = "observe"; // "observe" or "block"
            public bool AutoEnableBlocking { get; set; } = false;
            public int ObservePeriodHours { get; set; } = 48;
            public List<Dictionary<string, object>>? Rules { get; set; }
        }

        public class SecurityConfig
        {
            public bool MtlsEnabled { get; set; } = true;
            public string? CertificatePath { get; set; }
            public string? KeyPath { get; set; }
            public string? CaBundlePath { get; set; }
        }
    }

    public class RedactionRule
    {
        public string? Pattern { get; set; }
        public string? Replacement { get; set; }
        public string? Field { get; set; }
    }

    public class PolicyResult
    {
        public bool Allowed { get; set; }
        public string? Reason { get; set; }

        public PolicyResult(bool allowed, string? reason = null)
        {
            Allowed = allowed;
            Reason = reason;
        }
    }
}
