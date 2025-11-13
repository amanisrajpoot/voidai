using System;
using System.Collections.Generic;

namespace SecurityPlatform.Agent
{
    /// <summary>
    /// Configuration for the Security Platform Agent.
    /// </summary>
    public class AgentConfig
    {
        /// <summary>
        /// Control plane URL (defaults to environment variable or default URL).
        /// </summary>
        public string ControlPlaneUrl { get; set; } = 
            Environment.GetEnvironmentVariable("SECURITY_PLATFORM_CONTROL_PLANE_URL") 
            ?? "https://api.securityplatform.com";

        /// <summary>
        /// Authentication key (defaults to environment variable).
        /// </summary>
        public string? AuthKey { get; set; } = 
            Environment.GetEnvironmentVariable("SECURITY_PLATFORM_AUTH_KEY");

        /// <summary>
        /// Service name.
        /// </summary>
        public string ServiceName { get; set; } = "dotnet-service";

        /// <summary>
        /// Service version.
        /// </summary>
        public string? ServiceVersion { get; set; }

        /// <summary>
        /// Deployment environment.
        /// </summary>
        public string Environment { get; set; } = "production";

        /// <summary>
        /// OTLP endpoint (optional, defaults to control plane URL + /v1/traces).
        /// </summary>
        public string? OtlpEndpoint { get; set; }

        /// <summary>
        /// Telemetry configuration.
        /// </summary>
        public TelemetryConfig Telemetry { get; set; } = new TelemetryConfig();

        /// <summary>
        /// Redaction rules for PII protection.
        /// </summary>
        public List<string> RedactionRules { get; set; } = new List<string>();

        /// <summary>
        /// Policy configuration.
        /// </summary>
        public PolicyConfig Policy { get; set; } = new PolicyConfig();
    }

    /// <summary>
    /// Telemetry configuration.
    /// </summary>
    public class TelemetryConfig
    {
        public int BatchSize { get; set; } = 100;
        public int BatchTimeoutSeconds { get; set; } = 5;
        public int ExportTimeoutSeconds { get; set; } = 30;
        public int MaxQueueSize { get; set; } = 2048;
    }

    /// <summary>
    /// Policy configuration.
    /// </summary>
    public class PolicyConfig
    {
        public string Mode { get; set; } = "observe"; // "observe" or "block"
    }
}
