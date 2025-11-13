/**
 * Security Platform Frontend SDK
 * Provides session replay, RUM, and telemetry collection
 */

import { record } from 'rrweb';
import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { OTLPTraceExporter } from '@opentelemetry/exporter-otlp-http';
import { Resource } from '@opentelemetry/resources';
import { SEMRESATTRS_SERVICE_NAME, SEMRESATTRS_DEPLOYMENT_ENVIRONMENT } from '@opentelemetry/semantic-conventions';

export interface SDKConfig {
  controlPlaneUrl: string;
  authKey: string;
  serviceName?: string;
  environment?: string;
  redactionRules?: string[];
  batchSize?: number;
  enableSessionReplay?: boolean;
  enableRUM?: boolean;
}

export interface Event {
  type: string;
  data: any;
  timestamp: number;
  sessionId: string;
  traceId?: string;
}

const DEFAULT_REDACTION_RULES = [
  'password',
  'card',
  'ssn',
  'pin',
  'auth.*',
  'token',
  'secret',
  'api[_-]?key',
];

export class SecuritySDK {
  private config: SDKConfig;
  private sessionId: string;
  private events: Event[] = [];
  private stopRecording?: () => void;
  private tracerProvider?: WebTracerProvider;

  constructor(config: SDKConfig) {
    this.config = {
      redactionRules: DEFAULT_REDACTION_RULES,
      batchSize: 100,
      enableSessionReplay: true,
      enableRUM: true,
      ...config,
    };
    this.sessionId = this.generateSessionId();
    this.initialize();
  }

  private generateSessionId(): string {
    return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  }

  private initialize(): void {
    if (this.config.enableRUM) {
      this.initializeTelemetry();
    }

    if (this.config.enableSessionReplay) {
      this.initializeSessionReplay();
    }

    // Batch upload events
    this.startBatching();
  }

  private initializeTelemetry(): void {
    const resource = new Resource({
      [SEMRESATTRS_SERVICE_NAME]: this.config.serviceName || 'web-app',
      [SEMRESATTRS_DEPLOYMENT_ENVIRONMENT]: this.config.environment || 'production',
    });

    const exporter = new OTLPTraceExporter({
      url: `${this.config.controlPlaneUrl}/v1/traces`,
      headers: {
        'Authorization': `Bearer ${this.config.authKey}`,
      },
    });

    this.tracerProvider = new WebTracerProvider({
      resource,
    });

    this.tracerProvider.addSpanProcessor(new BatchSpanProcessor(exporter));
    this.tracerProvider.register();
  }

  private initializeSessionReplay(): void {
    this.stopRecording = record({
      emit(event) {
        // Redact sensitive data before storing
        const redactedEvent = this.redactEvent(event);
        this.events.push({
          type: 'session_replay',
          data: redactedEvent,
          timestamp: Date.now(),
          sessionId: this.sessionId,
        });
      }.bind(this),
      maskAllInputs: true,
      maskAllText: false,
      blockClass: 'no-record',
      blockSelector: '[data-no-record]',
    });
  }

  private redactEvent(event: any): any {
    if (!event || typeof event !== 'object') {
      return event;
    }

    const redacted = { ...event };

    for (const key in redacted) {
      if (this.shouldRedact(key)) {
        redacted[key] = '[REDACTED]';
      } else if (typeof redacted[key] === 'object') {
        redacted[key] = this.redactEvent(redacted[key]);
      }
    }

    return redacted;
  }

  private shouldRedact(key: string): boolean {
    if (!this.config.redactionRules) {
      return false;
    }

    const lowerKey = key.toLowerCase();
    return this.config.redactionRules.some(rule => {
      const regex = new RegExp(rule.replace('*', '.*'), 'i');
      return regex.test(lowerKey);
    });
  }

  private startBatching(): void {
    setInterval(() => {
      if (this.events.length > 0) {
        this.flushEvents();
      }
    }, 5000); // Flush every 5 seconds

    // Also flush on page unload
    window.addEventListener('beforeunload', () => {
      this.flushEvents();
    });
  }

  private async flushEvents(): Promise<void> {
    if (this.events.length === 0) {
      return;
    }

    const batch = this.events.splice(0, this.config.batchSize || 100);

    try {
      await fetch(`${this.config.controlPlaneUrl}/v1/events`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.config.authKey}`,
        },
        body: JSON.stringify({
          sessionId: this.sessionId,
          events: batch,
        }),
        keepalive: true, // Ensure request completes even if page unloads
      });
    } catch (error) {
      console.error('Failed to send events:', error);
      // Re-add events to queue for retry
      this.events.unshift(...batch);
    }
  }

  public captureEvent(type: string, data: any): void {
    this.events.push({
      type,
      data: this.redactEvent(data),
      timestamp: Date.now(),
      sessionId: this.sessionId,
    });
  }

  public destroy(): void {
    if (this.stopRecording) {
      this.stopRecording();
    }
    if (this.tracerProvider) {
      this.tracerProvider.shutdown();
    }
    this.flushEvents();
  }
}

// Default export for easy initialization
export function init(config: SDKConfig): SecuritySDK {
  return new SecuritySDK(config);
}

// Global initialization helper
if (typeof window !== 'undefined') {
  (window as any).SecurityPlatform = { init, SecuritySDK };
}
