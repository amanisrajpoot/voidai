# Security Platform Architecture

## Overview

The Security Platform is a comprehensive security observability and response system designed to work across multiple deployment models (SaaS, Hybrid, On-Prem) and support various platforms (web, backend, mobile, desktop, network appliances).

## Core Components

### 1. Control Plane

The control plane is the central management interface and API for the platform.

**Components:**
- REST API (Go) for agent management, rule configuration, playbook execution
- Web Dashboard (React/TypeScript) for visualization and management
- Authentication & Authorization (RBAC)
- Policy Engine for rule evaluation

**Deployment:**
- Cloud SaaS: Hosted multi-tenant service
- On-Prem: Docker/Kubernetes deployment

### 2. Telemetry Bus

The telemetry bus collects, processes, and stores observability data.

**Components:**
- OpenTelemetry Collector for trace/metric/log collection
- OpenSearch for log and trace storage
- ClickHouse for metrics storage
- MinIO/S3 for session replay blob storage

**Data Flow:**
```
Agents → OTLP → Collector → [OpenSearch/ClickHouse/MinIO]
```

### 3. Agents & SDKs

Language-specific agents that instrument applications and send telemetry.

**Supported Platforms:**
- **Frontend**: TypeScript SDK with rrweb for session recording
- **Backend**: Node.js, Python, Java, .NET, Go agents
- **Mobile**: iOS (Swift), Android (Kotlin) SDKs
- **Desktop**: Windows service, macOS daemon, Linux daemon
- **Network**: WAF modules, API gateway plugins

**Features:**
- Automatic instrumentation via OpenTelemetry
- PII redaction before transmission
- Local policy enforcement (observe/block modes)
- Offline caching and batch upload

### 4. Edge Modules

Edge components for network-level protection.

**Components:**
- NGINX + ModSecurity WAF
- Kong API Gateway plugin
- Envoy proxy integration

**Capabilities:**
- Request/response inspection
- IP blocking
- Rate limiting
- Dynamic rule updates from control plane

### 5. Orchestrator (SOAR-lite)

Webhook-driven automation engine for response actions.

**Capabilities:**
- Playbook execution
- Integration with Kong, WAF, IAM systems
- Alert routing (Slack, PagerDuty, Jira)
- Automated response actions

## Data Flow

```
┌─────────────┐
│   Agents    │
│  (SDKs)     │
└──────┬──────┘
       │ OTLP/HTTP
       ▼
┌─────────────┐
│  Collector  │
│  (OTEL)     │
└──────┬──────┘
       │
       ├──► OpenSearch (Traces/Logs)
       ├──► ClickHouse (Metrics)
       └──► MinIO/S3 (Sessions)
              │
              ▼
       ┌─────────────┐
       │Control Plane│
       │  (API/UI)   │
       └──────┬──────┘
              │
              ▼
       ┌─────────────┐
       │Orchestrator │
       │ (Playbooks) │
       └──────┬──────┘
              │
              ├──► Kong Admin API
              ├──► WAF Rules
              ├──► IAM (Revoke)
              └──► Alerts (Slack/PagerDuty)
```

## Security & Privacy

### PII Protection

- **Redaction**: Default blocklist for common PII fields (password, credit card, SSN, etc.)
- **Hashing**: SHA256 hashing of sensitive payloads
- **Configurable Rules**: Custom regex patterns for redaction
- **Safe Fields**: Developer opt-in for fields that are safe to transmit

### Authentication & Authorization

- **Bootstrap Tokens**: Short-lived tokens for initial agent registration
- **PKI Support**: mTLS with client certificates for enterprise deployments
- **RBAC**: Role-based access control in control plane

### Data Encryption

- **In Transit**: TLS 1.3 for all communications
- **At Rest**: AES-256 encryption for stored data
- **Key Management**: Integration with OS keystores and secure enclaves

## Deployment Models

### 1. Cloud SaaS

- Control plane hosted in cloud
- Telemetry sent to cloud storage
- Agents connect to cloud endpoints

### 2. Hybrid

- Control plane in cloud
- Telemetry stored on-premises
- Agents send to on-prem collector, metadata to cloud

### 3. Fully On-Prem

- All components deployed on-premises
- Air-gapped installation support
- Complete data sovereignty

## Scalability

- **Horizontal Scaling**: All components stateless and horizontally scalable
- **Load Balancing**: Control plane and collectors behind load balancers
- **Caching**: Redis for session data and policy caching
- **Queueing**: Message queues for async event processing

## High Availability

- **Multi-Region**: Control plane deployed across regions
- **Replication**: Database replication for telemetry storage
- **Failover**: Automatic failover for collectors and control plane
- **Backup**: Regular backups of configuration and telemetry data
