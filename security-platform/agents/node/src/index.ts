import { NodeSDK } from '@opentelemetry/sdk-node';
import { OTLPTraceExporter } from '@opentelemetry/exporter-otlp-http';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { HttpInstrumentation } from '@opentelemetry/instrumentation-http';
import { ExpressInstrumentation } from '@opentelemetry/instrumentation-express';
import { Express, Request, Response, NextFunction } from 'express';
import * as opentelemetry from '@opentelemetry/api';

export interface AgentConfig {
  controlPlaneUrl: string;
  authKey: string;
  serviceName?: string;
  environment?: string;
  redactionRules?: string[];
  localPolicy?: 'observe' | 'block';
  telemetryBatchSize?: number;
  otlpEndpoint?: string;
}

export interface SecurityContext {
  requestId: string;
  traceId: string;
  sessionId?: string;
  userId?: string;
  taintTags?: string[];
}

class SecurityAgent {
  private config: Required<AgentConfig>;
  private sdk: NodeSDK;
  private tracer: opentelemetry.Tracer;
  private redactionRules: RegExp[] = [];
  private localPolicy: Map<string, 'block' | 'observe'> = new Map();

  constructor(config: AgentConfig) {
    this.config = {
      serviceName: config.serviceName || 'node-service',
      environment: config.environment || 'production',
      redactionRules: config.redactionRules || ['password', 'card', 'ssn', 'pin', 'auth.*'],
      localPolicy: config.localPolicy || 'observe',
      telemetryBatchSize: config.telemetryBatchSize || 50,
      otlpEndpoint: config.otlpEndpoint || `${config.controlPlaneUrl}/v1/traces`,
      ...config,
    };

    this.setupRedactionRules();
    this.setupTelemetry();
    this.loadLocalPolicy();
  }

  private setupRedactionRules(): void {
    this.redactionRules = this.config.redactionRules.map(rule => 
      new RegExp(rule, 'gi')
    );
  }

  private setupTelemetry(): void {
    const resource = new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]: this.config.serviceName,
      [SemanticResourceAttributes.SERVICE_VERSION]: '0.1.0',
      [SemanticResourceAttributes.DEPLOYMENT_ENVIRONMENT]: this.config.environment,
    });

    const exporter = new OTLPTraceExporter({
      url: this.config.otlpEndpoint,
      headers: {
        'Authorization': `Bearer ${this.config.authKey}`,
      },
    });

    this.sdk = new NodeSDK({
      resource,
      traceExporter: exporter,
      instrumentations: [
        new HttpInstrumentation(),
        new ExpressInstrumentation(),
      ],
    });

    this.sdk.start();
    this.tracer = opentelemetry.trace.getTracer(this.config.serviceName);
  }

  private async loadLocalPolicy(): Promise<void> {
    // Load policy from control plane or local cache
    try {
      const response = await fetch(`${this.config.controlPlaneUrl}/api/v1/policy`, {
        headers: {
          'Authorization': `Bearer ${this.config.authKey}`,
        },
      });

      if (response.ok) {
        const policy = await response.json();
        // Update local policy cache
        for (const [pattern, action] of Object.entries(policy.rules || {})) {
          this.localPolicy.set(pattern, action as 'block' | 'observe');
        }
      }
    } catch (error) {
      console.warn('Failed to load policy from control plane, using local cache:', error);
    }
  }

  private redactData(data: any): any {
    if (typeof data === 'string') {
      let redacted = data;
      for (const rule of this.redactionRules) {
        redacted = redacted.replace(rule, '[REDACTED]');
      }
      return redacted;
    }

    if (typeof data === 'object' && data !== null) {
      const redacted: any = Array.isArray(data) ? [] : {};
      for (const [key, value] of Object.entries(data)) {
        const shouldRedact = this.redactionRules.some(rule => rule.test(key));
        redacted[key] = shouldRedact ? '[REDACTED]' : this.redactData(value);
      }
      return redacted;
    }

    return data;
  }

  private checkPolicy(context: SecurityContext): 'block' | 'observe' {
    // Check taint tags against policy
    if (context.taintTags) {
      for (const tag of context.taintTags) {
        const action = this.localPolicy.get(tag);
        if (action === 'block') {
          return 'block';
        }
      }
    }

    return this.config.localPolicy === 'block' ? 'block' : 'observe';
  }

  public middleware() {
    return (req: Request, res: Response, next: NextFunction) => {
      const span = this.tracer.startSpan('http_request');
      const requestId = `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      const traceId = span.spanContext().traceId;

      const context: SecurityContext = {
        requestId,
        traceId,
        sessionId: req.headers['x-session-id'] as string,
        userId: req.headers['x-user-id'] as string,
      };

      // Capture request metadata
      span.setAttributes({
        'http.method': req.method,
        'http.url': req.url,
        'http.route': req.route?.path || req.path,
        'request.id': requestId,
      });

      // Redact sensitive headers
      const redactedHeaders = this.redactData(req.headers);
      span.setAttribute('http.request.headers', JSON.stringify(redactedHeaders));

      // Check body for sensitive data
      if (req.body) {
        const redactedBody = this.redactData(req.body);
        span.setAttribute('http.request.body', JSON.stringify(redactedBody));
      }

      // Check policy
      const action = this.checkPolicy(context);
      if (action === 'block') {
        span.setAttribute('security.action', 'blocked');
        span.end();
        return res.status(403).json({ error: 'Request blocked by security policy' });
      }

      // Add context to request
      (req as any).securityContext = context;

      // Capture response
      res.on('finish', () => {
        span.setAttributes({
          'http.status_code': res.statusCode,
          'security.action': action,
        });
        span.end();
      });

      next();
    };
  }

  public captureEvent(name: string, data?: Record<string, any>, taintTags?: string[]): void {
    const span = this.tracer.startSpan(name);
    
    if (taintTags) {
      span.setAttribute('security.taint_tags', taintTags.join(','));
    }

    if (data) {
      const redactedData = this.redactData(data);
      Object.entries(redactedData).forEach(([key, value]) => {
        span.setAttribute(key, String(value));
      });
    }

    span.end();
  }

  public shutdown(): Promise<void> {
    return this.sdk.shutdown();
  }
}

export function createAgent(config: AgentConfig): SecurityAgent {
  return new SecurityAgent(config);
}

export function securityMiddleware(config: AgentConfig) {
  const agent = createAgent(config);
  return agent.middleware();
}

export default SecurityAgent;
