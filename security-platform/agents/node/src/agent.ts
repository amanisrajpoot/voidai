import { NodeSDK } from '@opentelemetry/sdk-trace-node';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { OTLPTraceExporter } from '@opentelemetry/exporter-otlp-http';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';
import { AgentConfig } from './types';
import { Redactor } from './redaction';
import { PolicyEngine } from './policy';

export class Agent {
  private sdk: NodeSDK | null = null;
  private config: AgentConfig;
  private redactor: Redactor;
  private policyEngine: PolicyEngine;

  constructor(config: AgentConfig) {
    this.config = config;
    this.redactor = new Redactor(config.redactionRules || []);
    this.policyEngine = new PolicyEngine(config.policy || { mode: 'observe' });
  }

  start(): void {
    const resource = new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]: this.config.serviceName || 'node-service',
      [SemanticResourceAttributes.DEPLOYMENT_ENVIRONMENT]: this.config.environment || 'production',
    });

    const exporter = new OTLPTraceExporter({
      url: this.config.otlpEndpoint || `${this.config.controlPlaneUrl}/v1/traces`,
      headers: this.config.authKey ? { Authorization: `Bearer ${this.config.authKey}` } : {},
    });

    this.sdk = new NodeSDK({
      resource,
      traceExporter: exporter,
      instrumentations: [getNodeAutoInstrumentations()],
      spanProcessor: new BatchSpanProcessor(exporter, {
        maxQueueSize: this.config.telemetry?.maxQueueSize || 2048,
        maxExportBatchSize: this.config.telemetry?.batchSize || 100,
        scheduledDelayMillis: this.parseDuration(this.config.telemetry?.batchTimeout || '5s'),
        exportTimeoutMillis: this.parseDuration(this.config.telemetry?.exportTimeout || '30s'),
      }),
    });

    this.sdk.start();
  }

  stop(): void {
    this.sdk?.shutdown();
  }

  redact(data: any): any {
    return this.redactor.redact(data);
  }

  checkPolicy(action: string, context: Record<string, any>): { allowed: boolean; reason?: string } {
    return this.policyEngine.check(action, context);
  }

  private parseDuration(duration: string): number {
    const match = duration.match(/^(\d+)([smh])$/);
    if (!match) return 5000; // default 5s

    const value = parseInt(match[1], 10);
    const unit = match[2];

    switch (unit) {
      case 's':
        return value * 1000;
      case 'm':
        return value * 60 * 1000;
      case 'h':
        return value * 60 * 60 * 1000;
      default:
        return 5000;
    }
  }
}
