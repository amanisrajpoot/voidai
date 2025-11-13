using OpenTelemetry;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using System.Diagnostics;

namespace SecurityPlatform.Agent;

/// <summary>
/// Security Platform .NET Agent
/// </summary>
public class Agent
{
    private TracerProvider? _tracerProvider;
    private readonly AgentConfig _config;
    private bool _started = false;

    public Agent(AgentConfig config)
    {
        _config = config ?? throw new ArgumentNullException(nameof(config));
    }

    public void Start()
    {
        if (_started)
        {
            return;
        }

        var resourceBuilder = ResourceBuilder.CreateDefault()
            .AddService(
                serviceName: _config.ServiceName ?? "unknown-service",
                serviceVersion: _config.Version ?? "1.0.0")
            .AddAttributes(new Dictionary<string, object>
            {
                ["deployment.environment"] = _config.Environment ?? "production"
            });

        _tracerProvider = Sdk.CreateTracerProviderBuilder()
            .SetResourceBuilder(resourceBuilder)
            .AddOtlpExporter(options =>
            {
                options.Endpoint = new Uri(_config.OtlpEndpoint ?? "http://localhost:4318/v1/traces");
            })
            .Build();

        _started = true;
    }

    public void Stop()
    {
        if (!_started)
        {
            return;
        }

        _tracerProvider?.Dispose();
        _started = false;
    }

    public ActivitySource CreateActivitySource(string name)
    {
        return new ActivitySource(name, _config.Version ?? "1.0.0");
    }

    public bool IsStarted => _started;
    
    public AgentConfig GetConfig() => _config;
}
