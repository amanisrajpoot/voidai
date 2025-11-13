# Frontend SDK

Security observability frontend SDK with session replay, RUM, and telemetry collection.

## Installation

```bash
npm install @security-platform/frontend-sdk
```

## Quick Start

```typescript
import { init } from '@security-platform/frontend-sdk';

const sdk = init({
  controlPlaneUrl: 'https://api.example.com',
  authKey: 'your-auth-key',
  serviceName: 'my-web-app',
  environment: 'production',
  enableSessionReplay: true,
  enableRUM: true,
});

// Capture custom events
sdk.captureEvent('user_action', { action: 'click', element: 'button' });
```

## CDN Usage

```html
<script src="https://cdn.example.com/security-platform-sdk.js"></script>
<script>
  SecurityPlatform.init({
    controlPlaneUrl: 'https://api.example.com',
    authKey: 'your-auth-key',
  });
</script>
```

## Configuration

- `controlPlaneUrl`: Control plane API endpoint
- `authKey`: Authentication token
- `serviceName`: Service identifier
- `environment`: Deployment environment
- `redactionRules`: Custom redaction patterns (default includes password, card, ssn, etc.)
- `batchSize`: Events per batch (default: 100)
- `enableSessionReplay`: Enable rrweb session recording
- `enableRUM`: Enable OpenTelemetry RUM

## Redaction

By default, the SDK redacts:
- password
- card
- ssn
- pin
- auth.*
- token
- secret
- api[_-]?key

Mark elements to exclude from recording:
```html
<div class="no-record">Sensitive content</div>
<div data-no-record>Sensitive content</div>
```
