using System;
using System.Collections.Generic;
using System.Diagnostics;
using Microsoft.Extensions.Logging;
using OpenTelemetry;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using OpenTelemetry.Exporter;

namespace SecurityPlatform.Agent
{
    /// <summary>
    /// Security Platform Agent for .NET applications.
    /// </summary>
    public class Agent : IDisposable
    {
        private readonly AgentConfig _config;
        private readonly ILogger<Agent> _logger;
        private readonly Redactor _redactor;
        private readonly PolicyEngine _policyEngine;
        private TracerProvider? _tracerProvider;
        private bool _disposed = false;

        public Agent(AgentConfig config, ILogger<Agent>? logger = null)
        {
            _config = config ?? throw new ArgumentNullException(nameof(config));
            _logger = logger ?? new LoggerFactory().CreateLogger<Agent>();
            _redactor = new Redactor(_config.RedactionRules ?? new List<RedactionRule>());
            _policyEngine = new PolicyEngine(_config.Policy ?? new PolicyConfig());
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
                    .AddService(_config.ServiceName ?? "dotnet-service")
                    .AddAttributes(new Dictionary<string, object>
                    {
                        ["deployment.environment"] = _config.Environment ?? "production"
                    });

                var otlpEndpoint = _config.OtlpEndpoint;
                if (string.IsNullOrEmpty(otlpEndpoint))
                {
                    otlpEndpoint = $"{_config.ControlPlaneUrl}/v1/traces";
                }

                _tracerProvider = Sdk.CreateTracerProviderBuilder()
                    .SetResourceBuilder(resourceBuilder)
                    .AddAspNetCoreInstrumentation()
                    .AddHttpClientInstrumentation()
                    .AddOtlpExporter(options =>
                    {
                        options.Endpoint = new Uri(otlpEndpoint);
                        if (!string.IsNullOrEmpty(_config.AuthKey))
                        {
                            options.Headers = $"Authorization=Bearer {_config.AuthKey}";
                        }
                    })
                    .Build();

                _logger.LogInformation("Security Platform Agent started for service: {ServiceName}", _config.ServiceName);
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
            if (_tracerProvider == null)
            {
                return;
            }

            try
            {
                _tracerProvider?.Dispose();
                _tracerProvider = null;
                _logger.LogInformation("Security Platform Agent stopped");
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Error stopping agent");
            }
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
        public PolicyResult CheckPolicy(string action, Dictionary<string, object> context)
        {
            return _policyEngine.Check(action, context);
        }

        public void Dispose()
        {
            if (!_disposed)
            {
                Stop();
                _disposed = true;
            }
        }
    }
}
