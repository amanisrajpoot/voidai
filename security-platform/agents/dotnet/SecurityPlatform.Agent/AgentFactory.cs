using System;
using System.IO;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;

namespace SecurityPlatform.Agent
{
    /// <summary>
    /// Factory for creating Agent instances from configuration.
    /// </summary>
    public static class AgentFactory
    {
        /// <summary>
        /// Create agent from YAML configuration file.
        /// </summary>
        public static Agent FromYamlFile(string configPath, ILogger<Agent>? logger = null)
        {
            var configuration = new ConfigurationBuilder()
                .SetBasePath(Directory.GetCurrentDirectory())
                .AddYamlFile(configPath, optional: false)
                .Build();

            var config = new AgentConfig();
            configuration.Bind(config);
            return new Agent(config, logger);
        }

        /// <summary>
        /// Create agent from configuration object.
        /// </summary>
        public static Agent FromConfig(AgentConfig config, ILogger<Agent>? logger = null)
        {
            return new Agent(config, logger);
        }

        /// <summary>
        /// Create agent with default configuration from environment variables.
        /// </summary>
        public static Agent CreateDefault(ILogger<Agent>? logger = null)
        {
            var config = new AgentConfig
            {
                ControlPlaneUrl = Environment.GetEnvironmentVariable("SECURITY_PLATFORM_URL") 
                    ?? "https://api.securityplatform.com",
                AuthKey = Environment.GetEnvironmentVariable("SECURITY_PLATFORM_AUTH_KEY"),
                ServiceName = Environment.GetEnvironmentVariable("SERVICE_NAME") ?? "dotnet-service",
                Environment = Environment.GetEnvironmentVariable("ENVIRONMENT") ?? "production"
            };
            return new Agent(config, logger);
        }
    }
}
