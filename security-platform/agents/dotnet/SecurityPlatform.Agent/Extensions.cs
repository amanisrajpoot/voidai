using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;

namespace SecurityPlatform.Agent
{
    /// <summary>
    /// Extension methods for easy integration.
    /// </summary>
    public static class Extensions
    {
        /// <summary>
        /// Add Security Platform Agent to the service collection.
        /// </summary>
        public static IServiceCollection AddSecurityPlatformAgent(
            this IServiceCollection services,
            AgentConfig config)
        {
            services.AddSingleton(config);
            services.AddSingleton<Agent>();
            return services;
        }

        /// <summary>
        /// Use Security Platform Agent middleware.
        /// </summary>
        public static IApplicationBuilder UseSecurityPlatformAgent(this IApplicationBuilder app)
        {
            var agent = app.ApplicationServices.GetRequiredService<Agent>();
            agent.Start();
            return app;
        }
    }
}
