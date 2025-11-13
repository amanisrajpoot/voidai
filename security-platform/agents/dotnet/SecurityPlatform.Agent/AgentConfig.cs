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
        public string? ServiceName { get; set; } = "dotnet-service";
        public string? Environment { get; set; } = "production";
        public string? Namespace { get; set; } = "default";
        public string? OtlpEndpoint { get; set; }
        public TelemetryConfig? Telemetry { get; set; } = new TelemetryConfig();
        public List<RedactionRule>? RedactionRules { get; set; } = new List<RedactionRule>();
        public PolicyConfig? Policy { get; set; } = new PolicyConfig();
        public SecurityConfig? Security { get; set; } = new SecurityConfig();

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
            public List<PolicyRule> Rules { get; set; } = new List<PolicyRule>();
        }

        public class SecurityConfig
        {
            public bool MtlsEnabled { get; set; } = false;
            public string? CertificatePath { get; set; }
            public string? KeyPath { get; set; }
            public string? CaBundlePath { get; set; }
        }
    }
}
