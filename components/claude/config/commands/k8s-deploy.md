---
description: Deploy applications to Kubernetes
argument-hint: <manifest-file> [--dry-run]
allowed-tools: Read, Write, Bash
---

## Task

Help the user deploy applications to Kubernetes.

## Apply Manifests

### Basic Apply
```bash
# Apply single file
kubectl apply -f deployment.yaml
kaf deployment.yaml  # alias

# Apply directory
kubectl apply -f ./k8s/

# Apply from URL
kubectl apply -f https://example.com/manifest.yaml
```

### Dry Run
```bash
# Client-side dry run
kubectl apply -f deployment.yaml --dry-run=client

# Server-side dry run (validates against cluster)
kubectl apply -f deployment.yaml --dry-run=server

# Output what would be applied
kubectl apply -f deployment.yaml --dry-run=client -o yaml
```

## Deployment Strategies

### Rolling Update (Default)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
```

```bash
# Update image
kubectl set image deployment/myapp container=image:v2

# Watch rollout
kubectl rollout status deployment/myapp
```

### Recreate
```yaml
spec:
  strategy:
    type: Recreate
```

## Rollout Management

### Status
```bash
kubectl rollout status deployment/myapp
```

### History
```bash
kubectl rollout history deployment/myapp
kubectl rollout history deployment/myapp --revision=2
```

### Rollback
```bash
# Rollback to previous
kubectl rollout undo deployment/myapp

# Rollback to specific revision
kubectl rollout undo deployment/myapp --to-revision=2
```

### Pause/Resume
```bash
kubectl rollout pause deployment/myapp
# Make multiple changes
kubectl rollout resume deployment/myapp
```

## Scaling

```bash
# Manual scale
kubectl scale deployment/myapp --replicas=5

# Autoscale
kubectl autoscale deployment/myapp --min=2 --max=10 --cpu-percent=80
```

## Deployment Template

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  labels:
    app: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: myapp:1.0.0
        ports:
        - containerPort: 8080
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
        env:
        - name: ENV
          value: "production"
        - name: SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: myapp-secrets
              key: secret-key
```

## Service Exposure

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp
spec:
  selector:
    app: myapp
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP  # or LoadBalancer, NodePort
```

## Helm Deployment

```bash
# Install chart
helm install myapp ./chart
hin myapp ./chart  # alias

# Upgrade
helm upgrade myapp ./chart
hup myapp ./chart  # alias

# With values
helm install myapp ./chart -f values-prod.yaml

# Uninstall
helm uninstall myapp
hun myapp  # alias
```

## Verification

```bash
# Check deployment
kubectl get deployment myapp
kgd  # alias

# Check pods
kubectl get pods -l app=myapp
kpods myapp  # dotfiles function

# Check events
kubectl get events --sort-by='.lastTimestamp'
```

## Dotfiles Integration

- `kaf` - kubectl apply -f (alias)
- `kdf` - kubectl delete -f (alias)
- `kgd` - kubectl get deployments (alias)
- `kpods [filter]` - List pods with filter
- `hin` - helm install (alias)
- `hup` - helm upgrade (alias)
