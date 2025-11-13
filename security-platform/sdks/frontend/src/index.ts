import { init as rrwebInit, record } from 'rrweb';
import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { OTLPTraceExporter } from '@opentelemetry/exporter-otlp-http';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import * as opentelemetry from '@opentelemetry/api';

export interface SecuritySDKConfig {
  controlPlaneUrl: string;
  authKey: string;
  serviceName?: string;
  environment?: string;
  redactionRules?: string[];
  enableSessionReplay?: boolean;
  enableTelemetry?: boolean;
  batchSize?: number;
  flushInterval?: number;
}

export interface RedactionRule {
  pattern: RegExp | string;
  replacement?: string;
}

class SecuritySDK {
  private config: Required<SecuritySDKConfig>;
  private tracer: opentelemetry.Tracer;
  private sessionId: string;
  private events: any[] = [];
  private redactionRules: RedactionRule[] = [];
  private stopRecording?: () => void;

  constructor(config: SecuritySDKConfig) {
    this.config = {
      serviceName: config.serviceName || 'web-app',
      environment: config.environment || 'production',
      redactionRules: config.redactionRules || ['password', 'card', 'ssn', 'pin', 'auth.*'],
      enableSessionReplay: config.enableSessionReplay !== false,
      enableTelemetry: config.enableTelemetry !== false,
      batchSize: config.batchSize || 50,
      flushInterval: config.flushInterval || 5000,
      ...config,
    };

    this.sessionId = this.generateSessionId();
    this.setupRedactionRules();
    
    if (this.config.enableTelemetry) {
      this.setupTelemetry();
    }

    if (this.config.enableSessionReplay) {
      this.setupSessionReplay();
    }
  }

  private generateSessionId(): string {
    return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  }

  private setupRedactionRules(): void {
    this.redactionRules = this.config.redactionRules.map(rule => {
      if (typeof rule === 'string') {
        return {
          pattern: new RegExp(rule, 'gi'),
          replacement: '[REDACTED]',
        };
      }
      return rule as RedactionRule;
    });
  }

  private setupTelemetry(): void {
    const resource = new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]: this.config.serviceName,
      [SemanticResourceAttributes.SERVICE_VERSION]: '0.1.0',
      [SemanticResourceAttributes.DEPLOYMENT_ENVIRONMENT]: this.config.environment,
    });

    const provider = new WebTracerProvider({ resource });
    const exporter = new OTLPTraceExporter({
      url: `${this.config.controlPlaneUrl}/v1/traces`,
      headers: {
        'Authorization': `Bearer ${this.config.authKey}`,
      },
    });

    provider.addSpanProcessor(new BatchSpanProcessor(exporter, {
      maxQueueSize: this.config.batchSize,
      scheduledDelayMillis: this.config.flushInterval,
    }));

    provider.register();
    this.tracer = opentelemetry.trace.getTracer(this.config.serviceName);
  }

  private setupSessionReplay(): void {
    this.stopRecording = rrwebInit({
      emit(event) {
        // Store events in memory, will be batched and sent
        this.events.push(event);
        
        // Auto-flush when batch size reached
        if (this.events.length >= this.config.batchSize) {
          this.flushEvents();
        }
      }.bind(this),
      maskAllInputs: true,
      maskAllText: false,
      maskTextSelector: this.getMaskSelectors(),
      recordCanvas: false, // Disable canvas recording for privacy
      recordCrossOriginIframes: false,
    });

    // Periodic flush
    setInterval(() => {
      this.flushEvents();
    }, this.config.flushInterval);

    // Flush on page unload
    window.addEventListener('beforeunload', () => {
      this.flushEvents();
    });
  }

  private getMaskSelectors(): string[] {
    // Common selectors for sensitive fields
    return [
      'input[type="password"]',
      'input[name*="password"]',
      'input[name*="card"]',
      'input[name*="ssn"]',
      'input[name*="pin"]',
      '[data-sensitive="true"]',
    ];
  }

  private redactData(data: any): any {
    if (typeof data === 'string') {
      let redacted = data;
      for (const rule of this.redactionRules) {
        if (rule.pattern instanceof RegExp) {
          redacted = redacted.replace(rule.pattern, rule.replacement || '[REDACTED]');
        } else {
          redacted = redacted.replace(new RegExp(rule.pattern, 'gi'), rule.replacement || '[REDACTED]');
        }
      }
      return redacted;
    }

    if (typeof data === 'object' && data !== null) {
      const redacted: any = Array.isArray(data) ? [] : {};
      for (const [key, value] of Object.entries(data)) {
        // Check if key matches redaction rules
        const shouldRedact = this.config.redactionRules.some(rule => {
          const pattern = typeof rule === 'string' ? new RegExp(rule, 'gi') : rule;
          return pattern.test(key);
        });

        if (shouldRedact) {
          redacted[key] = '[REDACTED]';
        } else {
          redacted[key] = this.redactData(value);
        }
      }
      return redacted;
    }

    return data;
  }

  private async flushEvents(): Promise<void> {
    if (this.events.length === 0) return;

    const eventsToSend = [...this.events];
    this.events = [];

    // Redact sensitive data
    const redactedEvents = eventsToSend.map(event => this.redactData(event));

    try {
      await fetch(`${this.config.controlPlaneUrl}/api/v1/sessions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.config.authKey}`,
        },
        body: JSON.stringify({
          sessionId: this.sessionId,
          events: redactedEvents,
          metadata: {
            url: window.location.href,
            userAgent: navigator.userAgent,
            timestamp: Date.now(),
          },
        }),
      });
    } catch (error) {
      console.error('Failed to send session events:', error);
      // Re-queue events on failure (with limit)
      this.events.unshift(...eventsToSend);
      if (this.events.length > 1000) {
        this.events = this.events.slice(0, 1000);
      }
    }
  }

  public captureEvent(name: string, data?: Record<string, any>): void {
    if (!this.config.enableTelemetry) return;

    const span = this.tracer.startSpan(name);
    if (data) {
      const redactedData = this.redactData(data);
      Object.entries(redactedData).forEach(([key, value]) => {
        span.setAttribute(key, String(value));
      });
    }
    span.end();
  }

  public markFieldSafe(fieldName: string): void {
    // Remove from redaction rules if present
    this.config.redactionRules = this.config.redactionRules.filter(
      rule => typeof rule === 'string' && !rule.includes(fieldName)
    );
  }

  public destroy(): void {
    if (this.stopRecording) {
      this.stopRecording();
    }
    this.flushEvents();
  }
}

export function initSecuritySDK(config: SecuritySDKConfig): SecuritySDK {
  return new SecuritySDK(config);
}

export default SecuritySDK;
