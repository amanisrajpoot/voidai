using System;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using SecurityPlatform.Agent.Middleware;

namespace SecurityPlatform.Agent.Extensions
{
    /// <summary>
    /// Extension methods for dependency injection.
    /// </summary>
    public static class ServiceCollectionExtensions
    {
        /// <summary>
        /// Add Security Platform Agent to the service collection.
        /// </summary>
        public static IServiceCollection AddSecurityPlatformAgent(
            this IServiceCollection services,
            IConfiguration configuration)
        {
            var config = new AgentConfig();
            configuration.GetSection("SecurityPlatform").Bind(config);
            
            services.AddSingleton(config);
            services.AddSingleton<Agent>(sp =>
            {
                var logger = sp.GetRequiredService<ILogger<Agent>>();
                var agent = new Agent(config, logger);
                agent.Start();
                return agent;
            });

            return services;
        }

        /// <summary>
        /// Use Security Platform middleware in the request pipeline.
        /// </summary>
        public static Microsoft.AspNetCore.Builder.IApplicationBuilder UseSecurityPlatform(
            this Microsoft.AspNetCore.Builder.IApplicationBuilder app)
        {
            return app.UseMiddleware<SecurityPlatformMiddleware>();
        }
    }
}
