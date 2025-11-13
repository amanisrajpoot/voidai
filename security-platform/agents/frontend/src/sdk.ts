import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { OTLPTraceExporter } from '@opentelemetry/exporter-otlp-http';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { getWebAutoInstrumentations } from '@opentelemetry/auto-instrumentations-web';
import { record, type eventWithTime } from 'rrweb';
import { SDKConfig } from './types';
import { Redactor } from './redaction';

let tracerProvider: WebTracerProvider | null = null;
let stopRecording: (() => void) | null = null;
let redactor: Redactor;

export function init(config: SDKConfig): void {
  // Initialize redaction
  redactor = new Redactor(config.redactionRules || []);

  // Initialize OpenTelemetry
  if (config.enableRUM !== false) {
    const resource = new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]: config.serviceName || 'web-app',
      [SemanticResourceAttributes.DEPLOYMENT_ENVIRONMENT]: config.environment || 'production',
    });

    const exporter = new OTLPTraceExporter({
      url: config.otlpEndpoint || `${config.controlPlaneUrl}/v1/traces`,
      headers: config.authKey ? { Authorization: `Bearer ${config.authKey}` } : {},
    });

    tracerProvider = new WebTracerProvider({
      resource,
      instrumentations: [getWebAutoInstrumentations()],
    });

    tracerProvider.addSpanProcessor(new BatchSpanProcessor(exporter));
    tracerProvider.register();
  }

  // Initialize session recording
  if (config.enableSessionRecording) {
    const events: eventWithTime[] = [];
    
    stopRecording = record({
      emit(event) {
        events.push(event);
        
        // Batch and send events
        if (events.length >= (config.batchSize || 50)) {
          sendSessionEvents(events.splice(0), config);
        }
      },
      maskAllInputs: config.redactInputs !== false,
      maskTextSelector: config.maskText !== false ? '*' : undefined,
      recordCanvas: config.recordCanvas || false,
      recordCrossOriginIframes: config.recordCrossOriginIframes || false,
    });

    // Send remaining events on page unload
    window.addEventListener('beforeunload', () => {
      if (events.length > 0) {
        sendSessionEvents(events, config);
      }
    });
  }

  // Capture errors
  window.addEventListener('error', (event) => {
    captureError(event.error, config);
  });

  window.addEventListener('unhandledrejection', (event) => {
    captureError(event.reason, config);
  });
}

function sendSessionEvents(events: eventWithTime[], config: SDKConfig): void {
  const redactedEvents = events.map(event => redactor.redact(event));
  
  fetch(`${config.controlPlaneUrl}/api/sessions`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...(config.authKey ? { Authorization: `Bearer ${config.authKey}` } : {}),
    },
    body: JSON.stringify({
      sessionId: getSessionId(),
      events: redactedEvents,
      metadata: {
        url: window.location.href,
        userAgent: navigator.userAgent,
        timestamp: Date.now(),
      },
    }),
    keepalive: true,
  }).catch(err => {
    console.error('Failed to send session events:', err);
  });
}

function captureError(error: Error, config: SDKConfig): void {
  if (!tracerProvider) return;

  const tracer = tracerProvider.getTracer('security-platform-sdk');
  const span = tracer.startSpan('error');
  
  span.setAttributes({
    'error.type': error.name,
    'error.message': error.message,
    'error.stack': error.stack || '',
  });
  
  span.recordException(error);
  span.end();

  // Also send to control plane
  fetch(`${config.controlPlaneUrl}/api/errors`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...(config.authKey ? { Authorization: `Bearer ${config.authKey}` } : {}),
    },
    body: JSON.stringify({
      error: {
        name: error.name,
        message: error.message,
        stack: error.stack,
      },
      url: window.location.href,
      timestamp: Date.now(),
    }),
  }).catch(err => {
    console.error('Failed to send error:', err);
  });
}

function getSessionId(): string {
  let sessionId = sessionStorage.getItem('sp_session_id');
  if (!sessionId) {
    sessionId = `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    sessionStorage.setItem('sp_session_id', sessionId);
  }
  return sessionId;
}

export function stop(): void {
  if (stopRecording) {
    stopRecording();
    stopRecording = null;
  }
  if (tracerProvider) {
    tracerProvider.shutdown();
    tracerProvider = null;
  }
}
