import { Agent } from './agent';
import { AgentConfig, MiddlewareOptions } from './types';
import { createMiddleware } from './middleware';

export { Agent, createMiddleware };
export type { AgentConfig, MiddlewareOptions };

// Auto-initialize if environment variables are set
if (process.env.SECURITY_PLATFORM_AUTO_INIT === 'true') {
  const config: AgentConfig = {
    controlPlaneUrl: process.env.SECURITY_PLATFORM_CONTROL_PLANE_URL || '',
    authKey: process.env.SECURITY_PLATFORM_AUTH_KEY,
    serviceName: process.env.SECURITY_PLATFORM_SERVICE_NAME || 'node-service',
    environment: process.env.SECURITY_PLATFORM_ENVIRONMENT || 'production',
    otlpEndpoint: process.env.SECURITY_PLATFORM_OTLP_ENDPOINT,
  };

  const agent = new Agent(config);
  agent.start();
}
