# API Reference

## Control Plane API

Base URL: `https://api.example.com/api/v1`

### Authentication

All API requests require a Bearer token:

```http
Authorization: Bearer <token>
```

### Events

#### POST /events
Ingest security events from agents.

**Request Body:**
```json
{
  "events": [
    {
      "type": "http_request",
      "data": {
        "method": "GET",
        "path": "/api/users",
        "status_code": 200
      },
      "timestamp": 1234567890,
      "sessionId": "session-123"
    }
  ],
  "sessionId": "session-123"
}
```

**Response:**
```json
{
  "success": true,
  "received": 1
}
```

#### GET /events
List events with filtering.

**Query Parameters:**
- `limit` (default: 100): Number of events to return
- `offset` (default: 0): Pagination offset
- `type`: Filter by event type
- `sessionId`: Filter by session ID

**Response:**
```json
{
  "events": [...],
  "total": 1000,
  "limit": 100,
  "offset": 0
}
```

### Traces (OTLP)

#### POST /traces
OpenTelemetry trace export endpoint.

**Request:** OTLP protobuf format

**Response:** 200 OK

### Sessions

#### GET /sessions/:sessionId
Get session metadata.

**Response:**
```json
{
  "sessionId": "session-123",
  "startedAt": "2024-01-01T00:00:00Z",
  "events": [...]
}
```

#### POST /sessions/:sessionId/replay
Get session replay data (rrweb events).

**Response:**
```json
{
  "sessionId": "session-123",
  "events": [...],
  "metadata": {...}
}
```

### Rules

#### GET /rules
List all active rules.

**Response:**
```json
{
  "rules": [
    {
      "id": "rule-1",
      "name": "SQL Injection Detection",
      "pattern": ".*(union|select).*",
      "action": "alert",
      "enabled": true
    }
  ]
}
```

#### POST /rules
Create a new rule.

**Request Body:**
```json
{
  "name": "Custom Rule",
  "pattern": ".*suspicious.*",
  "action": "block",
  "enabled": true
}
```

### Actions

#### POST /actions/block-ip
Block an IP address.

**Request Body:**
```json
{
  "ip": "192.168.1.100",
  "reason": "Suspicious activity",
  "duration": 3600
}
```

#### POST /actions/block-route
Block a specific route.

**Request Body:**
```json
{
  "route": "/admin/*",
  "reason": "Unauthorized access attempt"
}
```

#### POST /actions/revoke-token
Revoke a user token.

**Request Body:**
```json
{
  "token": "token-123",
  "userId": "user-456",
  "reason": "Security incident"
}
```

#### POST /actions/alert
Send an alert.

**Request Body:**
```json
{
  "channel": "slack",
  "message": "Security alert: SQL injection detected",
  "severity": "high"
}
```

## SDK APIs

### Frontend SDK

```typescript
import { init } from '@security-platform/frontend-sdk';

const sdk = init({
  controlPlaneUrl: 'https://api.example.com',
  authKey: 'your-key',
});

// Capture event
sdk.captureEvent('custom_event', { data: 'value' });

// Destroy SDK
sdk.destroy();
```

### Python Agent

```python
from security_platform import init_agent, AgentConfig

config = AgentConfig(
    control_plane_url="https://api.example.com",
    auth_key="your-key",
)
agent = init_agent(config)

# Capture event
agent.capture_event("custom_event", {"data": "value"})

# Get tracer
tracer = agent.get_tracer("my-service")
with tracer.start_as_current_span("operation"):
    # Your code
    pass
```
