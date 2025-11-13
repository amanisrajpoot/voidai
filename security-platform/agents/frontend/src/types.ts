export interface SDKConfig {
  controlPlaneUrl: string;
  authKey?: string;
  serviceName?: string;
  environment?: string;
  otlpEndpoint?: string;
  enableRUM?: boolean;
  enableSessionRecording?: boolean;
  redactInputs?: boolean;
  maskText?: boolean;
  recordCanvas?: boolean;
  recordCrossOriginIframes?: boolean;
  batchSize?: number;
  redactionRules?: RedactionRule[];
}

export interface Event {
  type: string;
  timestamp: number;
  data: Record<string, any>;
}

export interface RedactionRule {
  pattern: string; // Regex pattern
  replacement?: string; // Default: "[REDACTED]"
  field?: string; // Specific field name, or "*" for all fields
}
