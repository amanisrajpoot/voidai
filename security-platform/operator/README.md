# Security Platform Kubernetes Operator

Kubernetes operator for managing Security Platform Agents.

## Installation

### Using Helm

```bash
helm install security-platform-operator ./helm/security-platform-operator
```

### Using kubectl

```bash
# Install CRDs
kubectl apply -f deploy/crds/

# Install operator
kubectl apply -f deploy/operator.yaml
```

## Usage

Create an Agent resource:

```yaml
apiVersion: securityplatform.io/v1
kind: Agent
metadata:
  name: my-agent
  namespace: default
spec:
  controlPlaneUrl: "https://api.securityplatform.com"
  authKeySecret:
    name: security-platform-secret
    key: auth-key
  serviceName: "my-service"
  environment: "production"
  telemetry:
    batchSize: 100
    batchTimeout: "5s"
    exportTimeout: "30s"
    maxQueueSize: 2048
  policy:
    mode: "observe"
    autoEnableBlocking: false
```

Apply it:

```bash
kubectl apply -f agent.yaml
```

The operator will:
1. Create a ConfigMap with agent configuration
2. Create a Deployment running the agent
3. Update the Agent status

## Features

- ✅ Custom Resource Definition (CRD) for Agent
- ✅ Automatic Deployment management
- ✅ ConfigMap generation from Agent spec
- ✅ Status tracking
- ✅ Secret reference support
