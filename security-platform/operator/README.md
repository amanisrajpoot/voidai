# Security Platform Kubernetes Operator

Kubernetes operator for managing Security Platform agents.

## Quick Start

### Install Operator

```bash
# Install CRDs
kubectl apply -f config/crd/bases/

# Install operator
kubectl apply -f config/manager/manager.yaml
```

### Deploy Agent

```yaml
apiVersion: securityplatform.io/v1alpha1
kind: Agent
metadata:
  name: my-agent
spec:
  controlPlaneURL: "https://api.securityplatform.com"
  authKeySecretRef:
    name: security-platform-secret
    key: auth-key
  serviceName: "my-service"
  environment: "production"
  replicas: 2
  image: "securityplatform/agent:latest"
  resources:
    requests:
      cpu: "100m"
      memory: "128Mi"
    limits:
      cpu: "500m"
      memory: "512Mi"
```

### Create Secret

```bash
kubectl create secret generic security-platform-secret \
  --from-literal=auth-key=your-token-here
```

## Features

- ✅ Automatic agent deployment
- ✅ Secret management
- ✅ Resource management
- ✅ Scaling support
- ✅ Simple CRD-based configuration

## Documentation

See [docs/DEPLOYMENT.md](../../docs/DEPLOYMENT.md) for full documentation.
