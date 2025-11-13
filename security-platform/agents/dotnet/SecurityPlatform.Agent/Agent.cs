using System;
using System.Collections.Generic;
using System.Linq;
using Microsoft.Extensions.Logging;
using OpenTelemetry;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using OpenTelemetry.Exporter;

namespace SecurityPlatform.Agent
{
    /// <summary>
    /// Security Platform Agent for .NET applications.
    /// Provides automatic instrumentation and telemetry collection.
    /// </summary>
    public class Agent : IDisposable
    {
        private readonly AgentConfig _config;
        private readonly ILogger<Agent> _logger;
        private readonly Redactor _redactor;
        private readonly PolicyEngine _policyEngine;
        private TracerProvider? _tracerProvider;
        private bool _disposed = false;

        /// <summary>
        /// Create a new agent instance with the given configuration.
        /// </summary>
        public Agent(AgentConfig config, ILogger<Agent>? logger = null)
        {
            _config = config ?? throw new ArgumentNullException(nameof(config));
            _logger = logger ?? new LoggerFactory().CreateLogger<Agent>();
            _redactor = new Redactor(_config.RedactionRules);
            _policyEngine = new PolicyEngine(_config.Policy);
        }

        /// <summary>
        /// Start the agent and initialize OpenTelemetry.
        /// </summary>
        public void Start()
        {
            if (_tracerProvider != null)
            {
                _logger.LogWarning("Agent already started");
                return;
            }

            try
            {
                var resourceBuilder = ResourceBuilder.CreateDefault()
                    .AddService(
                        serviceName: _config.ServiceName ?? "dotnet-service",
                        serviceVersion: _config.ServiceVersion);

                var otlpEndpoint = _config.OtlpEndpoint;
                if (string.IsNullOrEmpty(otlpEndpoint))
                {
                    otlpEndpoint = $"{_config.ControlPlaneUrl}/v1/traces";
                }

                var headers = new Dictionary<string, string>();
                var authKey = _config.AuthKey ?? Environment.GetEnvironmentVariable("SECURITY_PLATFORM_AUTH_KEY");
                if (!string.IsNullOrEmpty(authKey))
                {
                    headers["Authorization"] = $"Bearer {authKey}";
                }

                _tracerProvider = Sdk.CreateTracerProviderBuilder()
                    .SetResourceBuilder(resourceBuilder)
                    .AddAspNetCoreInstrumentation()
                    .AddHttpClientInstrumentation()
                    .AddOtlpExporter(options =>
                    {
                        options.Endpoint = new Uri(otlpEndpoint);
                        options.Headers = string.Join(",", headers.Select(kvp => $"{kvp.Key}={kvp.Value}"));
                    })
                    .Build();

                _logger.LogInformation("Security Platform Agent started successfully");
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Failed to start agent");
                throw;
            }
        }

        /// <summary>
        /// Stop the agent and shutdown OpenTelemetry.
        /// </summary>
        public void Stop()
        {
            Dispose();
        }

        /// <summary>
        /// Redact PII from data.
        /// </summary>
        public object? Redact(object? data)
        {
            return _redactor.Redact(data);
        }

        /// <summary>
        /// Check if an action is allowed by policy.
        /// </summary>
        public PolicyResult CheckPolicy(string action, Dictionary<string, object>? context = null)
        {
            return _policyEngine.Check(action, context);
        }

        /// <summary>
        /// Dispose the agent.
        /// </summary>
        public void Dispose()
        {
            if (_disposed)
            {
                return;
            }

            _tracerProvider?.Dispose();
            _tracerProvider = null;
            _disposed = true;
            _logger.LogInformation("Security Platform Agent stopped");
        }
    }
}
