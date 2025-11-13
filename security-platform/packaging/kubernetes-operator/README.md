# Kubernetes Operator

Kubernetes Operator for managing Security Platform Agents.

## Installation

```bash
# Install CRDs
kubectl apply -f config/crd/bases/

# Install operator
kubectl apply -f config/manager/manager.yaml
```

## Usage

1. Create a Secret with auth key:

```bash
kubectl create secret generic security-platform-auth \
  --from-literal=auth-key=your-auth-key
```

2. Create an Agent resource:

```bash
kubectl apply -f examples/agent.yaml
```

3. Check status:

```bash
kubectl get agents
kubectl describe agent my-agent
```

## Custom Resource Definition

The Agent CRD allows you to configure:
- Control plane URL
- Authentication (via Secret reference)
- Service name and environment
- Replicas and image
- Telemetry and policy settings
