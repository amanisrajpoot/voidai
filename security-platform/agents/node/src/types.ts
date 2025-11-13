export interface AgentConfig {
  controlPlaneUrl: string;
  authKey?: string;
  serviceName?: string;
  environment?: string;
  namespace?: string;
  otlpEndpoint?: string;
  telemetry?: {
    batchSize?: number;
    batchTimeout?: string;
    exportTimeout?: string;
    maxQueueSize?: number;
  };
  redactionRules?: RedactionRule[];
  policy?: {
    mode?: 'observe' | 'block';
    autoEnableBlocking?: boolean;
    observePeriodHours?: number;
    rules?: PolicyRule[];
  };
  security?: {
    mtlsEnabled?: boolean;
    certificatePath?: string;
    keyPath?: string;
    caBundlePath?: string;
  };
}

export interface RedactionRule {
  pattern: string;
  replacement?: string;
  field?: string;
}

export interface PolicyRule {
  id: string;
  action: string;
  condition: string; // JSONPath or expression
  effect: 'allow' | 'deny';
}

export interface MiddlewareOptions {
  mode?: 'observe' | 'block';
  captureResponseBody?: boolean;
  sendToControlPlane?: boolean;
}
