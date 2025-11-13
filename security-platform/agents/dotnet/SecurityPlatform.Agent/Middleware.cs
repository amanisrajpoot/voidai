using Microsoft.AspNetCore.Http;
using System.Diagnostics;

namespace SecurityPlatform.Agent;

/// <summary>
/// ASP.NET Core Middleware for Security Platform Agent
/// </summary>
public class SecurityPlatformMiddleware
{
    private readonly RequestDelegate _next;
    private readonly Agent _agent;
    private readonly Redaction _redaction;
    private readonly PolicyEngine _policyEngine;

    public SecurityPlatformMiddleware(RequestDelegate next, Agent agent)
    {
        _next = next;
        _agent = agent;
        _redaction = new Redaction(agent.GetConfig().RedactionRules);
        _policyEngine = new PolicyEngine(agent.GetConfig());
    }

    public async Task InvokeAsync(HttpContext context)
    {
        using var activity = _agent.CreateActivitySource("security-platform-agent")
            .StartActivity("http.request");

        if (activity != null)
        {
            activity.SetTag("http.method", context.Request.Method);
            activity.SetTag("http.url", context.Request.Path + context.Request.QueryString);
        }

        try
        {
            // Capture request data
            var requestData = new Dictionary<string, object>
            {
                ["method"] = context.Request.Method,
                ["uri"] = context.Request.Path.ToString(),
                ["query"] = context.Request.QueryString.ToString()
            };

            var headers = context.Request.Headers.ToDictionary(
                h => h.Key,
                h => h.Value.ToString());
            requestData["headers"] = _redaction.Redact(headers.ToDictionary(k => k.Key, k => (object)k.Value));

            // Check policy
            var policyResult = _policyEngine.CheckPolicy("http_request", requestData);

            if ("block".Equals(_agent.GetConfig().LocalPolicy, StringComparison.OrdinalIgnoreCase) && !policyResult.IsAllowed)
            {
                context.Response.StatusCode = 403;
                await context.Response.WriteAsync("Request blocked by security policy");
                activity?.SetTag("security.blocked", true);
                activity?.SetTag("security.reason", policyResult.Reason);
                return;
            }

            await _next(context);

            activity?.SetTag("http.status_code", context.Response.StatusCode);
            activity?.SetTag("security.policy.allowed", policyResult.IsAllowed);

            if (!policyResult.IsAllowed)
            {
                activity?.SetTag("security.policy.reason", policyResult.Reason);
            }
        }
        catch (Exception ex)
        {
            activity?.RecordException(ex);
            activity?.SetTag("error", true);
            throw;
        }
    }
}

public static class SecurityPlatformMiddlewareExtensions
{
    public static IApplicationBuilder UseSecurityPlatform(this IApplicationBuilder builder, Agent agent)
    {
        return builder.UseMiddleware<SecurityPlatformMiddleware>(agent);
    }
}
