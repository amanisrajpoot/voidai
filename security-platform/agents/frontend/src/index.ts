/**
 * Security Platform Frontend SDK
 * 
 * Provides RUM (Real User Monitoring), session recording, and security observability
 * for web applications.
 */

import { init } from './sdk';
import { SDKConfig, Event, RedactionRule } from './types';
import { redact } from './redaction';

export { init, redact };
export type { SDKConfig, Event, RedactionRule };

// Auto-initialize if script tag has data attributes
if (typeof window !== 'undefined') {
  const script = document.currentScript as HTMLScriptElement;
  if (script?.dataset.autoInit === 'true') {
    const config: SDKConfig = {
      controlPlaneUrl: script.dataset.controlPlaneUrl || '',
      serviceName: script.dataset.serviceName || 'web-app',
      environment: script.dataset.environment || 'production',
      enableSessionRecording: script.dataset.enableSessionRecording === 'true',
      enableRUM: script.dataset.enableRum !== 'false',
    };
    init(config);
  }
}
