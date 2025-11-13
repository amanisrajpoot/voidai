using System;
using System.Collections.Generic;
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
            _logger = logger ?? new Microsoft.Extensions.Logging.Abstractions.NullLogger<Agent>();
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
                    .AddService(
                        serviceName: _config.ServiceName ?? "dotnet-service",
                        serviceVersion: "1.0.0")
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
                        options.Headers = GetAuthHeaders();
                    })
                    .SetBatchExportProcessorOptions(new BatchExportProcessorOptions<Activity>
                    {
                        MaxQueueSize = _config.Telemetry?.MaxQueueSize ?? 2048,
                        MaxExportBatchSize = _config.Telemetry?.BatchSize ?? 100,
                        ScheduledDelayMilliseconds = ParseDuration(_config.Telemetry?.BatchTimeout ?? "5s"),
                        ExporterTimeoutMilliseconds = ParseDuration(_config.Telemetry?.ExportTimeout ?? "30s")
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

        private Dictionary<string, string> GetAuthHeaders()
        {
            var headers = new Dictionary<string, string>();
            var authKey = _config.AuthKey ?? Environment.GetEnvironmentVariable("SECURITY_PLATFORM_AUTH_KEY");
            if (!string.IsNullOrEmpty(authKey))
            {
                headers["Authorization"] = $"Bearer {authKey}";
            }
            return headers;
        }

        private int ParseDuration(string duration)
        {
            if (string.IsNullOrEmpty(duration))
                return 5000;

            var match = System.Text.RegularExpressions.Regex.Match(duration, @"^(\d+)([smh])$");
            if (!match.Success)
                return 5000;

            var value = int.Parse(match.Groups[1].Value);
            var unit = match.Groups[2].Value;

            return unit switch
            {
                "s" => value * 1000,
                "m" => value * 60 * 1000,
                "h" => value * 60 * 60 * 1000,
                _ => 5000
            };
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
