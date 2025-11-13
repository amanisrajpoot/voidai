using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.Logging;

namespace SecurityPlatform.Agent.Middleware
{
    /// <summary>
    /// ASP.NET Core middleware for Security Platform Agent.
    /// </summary>
    public class SecurityPlatformMiddleware
    {
        private readonly RequestDelegate _next;
        private readonly Agent _agent;
        private readonly ILogger<SecurityPlatformMiddleware> _logger;

        public SecurityPlatformMiddleware(
            RequestDelegate next,
            Agent agent,
            ILogger<SecurityPlatformMiddleware> logger)
        {
            _next = next;
            _agent = agent;
            _logger = logger;
        }

        public async Task InvokeAsync(HttpContext context)
        {
            // Check policy
            var policyContext = new Dictionary<string, object>
            {
                ["method"] = context.Request.Method,
                ["path"] = context.Request.Path.Value ?? "",
                ["remote_addr"] = context.Connection.RemoteIpAddress?.ToString() ?? ""
            };

            var policyResult = _agent.CheckPolicy("http_request", policyContext);

            if (!policyResult.Allowed)
            {
                _logger.LogWarning("Request blocked by policy: {Reason}", policyResult.Reason);
                context.Response.StatusCode = 403;
                await context.Response.WriteAsync("Request blocked by security policy");
                return;
            }

            // Continue with request
            await _next(context);
        }
    }
}
