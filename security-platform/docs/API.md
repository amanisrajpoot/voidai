# Security Platform API Reference

## Base URL

- **Cloud SaaS**: `https://api.securityplatform.com`
- **On-Prem**: `https://your-control-plane.example.com`

## Authentication

All API requests require authentication via Bearer token:

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" https://api.securityplatform.com/api/v1/agents
```

## Endpoints

### Agents

#### List Agents
```http
GET /api/v1/agents
```

**Response:**
```json
[
  {
    "id": "agent-123",
    "service_name": "my-service",
    "environment": "production",
    "status": "active",
    "last_seen": "2024-01-15T10:30:00Z"
  }
]
```

#### Get Agent
```http
GET /api/v1/agents/{id}
```

#### Create Agent
```http
POST /api/v1/agents
Content-Type: application/json

{
  "service_name": "my-service",
  "environment": "production"
}
```

#### Update Agent
```http
PUT /api/v1/agents/{id}
Content-Type: application/json

{
  "policy": {
    "mode": "block"
  }
}
```

#### Delete Agent
```http
DELETE /api/v1/agents/{id}
```

### Rules

#### List Rules
```http
GET /api/v1/rules
```

#### Create Rule
```http
POST /api/v1/rules
Content-Type: application/json

{
  "name": "Block Suspicious IPs",
  "condition": "ip == '192.168.1.100'",
  "action": "block",
  "enabled": true
}
```

### Playbooks

#### List Playbooks
```http
GET /api/v1/playbooks
```

#### Execute Playbook
```http
POST /api/v1/playbooks/{id}/execute
Content-Type: application/json

{
  "trigger": {
    "type": "security_incident",
    "data": {
      "severity": "high",
      "ip": "192.168.1.100"
    }
  }
}
```

### Sessions

#### Upload Session Replay
```http
POST /api/v1/sessions
Content-Type: application/json

{
  "session_id": "session-123",
  "events": [...],
  "metadata": {
    "url": "https://example.com",
    "user_agent": "..."
  }
}
```

### Events

#### Send Event
```http
POST /api/v1/events
Content-Type: application/json

{
  "type": "http_request",
  "data": {
    "method": "POST",
    "path": "/api/login",
    "ip": "192.168.1.100"
  }
}
```

### OTLP Endpoints

#### Traces
```http
POST /v1/traces
Content-Type: application/x-protobuf
```

#### Metrics
```http
POST /v1/metrics
Content-Type: application/x-protobuf
```

#### Logs
```http
POST /v1/logs
Content-Type: application/x-protobuf
```

## Webhooks

### Playbook Webhook

```http
POST /api/v1/webhooks/{playbook_id}
Content-Type: application/json

{
  "event": "security_incident",
  "data": {
    "severity": "high",
    "ip": "192.168.1.100"
  }
}
```

## Rate Limits

- **API**: 1000 requests per minute per API key
- **OTLP**: 10,000 spans per second per agent

## Error Responses

```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Invalid request parameters",
    "details": {}
  }
}
```

## SDKs

Official SDKs available for:
- JavaScript/TypeScript
- Python
- Go
- Java
- .NET

See [SDK Documentation](./SDK.md) for details.
