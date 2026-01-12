---
description: Debug a failing or problematic Kubernetes pod
argument-hint: <pod-name>
allowed-tools: Read, Bash
---

## Task

Help the user debug a failing or problematic Kubernetes pod.

## Debug Workflow

### Step 1: Get Pod Status
```bash
# Check pod status
kubectl get pod <pod-name> -o wide
kgp | grep <pod-name>  # using alias

# Check all pods in namespace
kpods  # dotfiles function
```

### Step 2: Describe Pod
```bash
kubectl describe pod <pod-name>
kd pod <pod-name>  # using alias

# Look for:
# - Events section (errors, warnings)
# - Conditions (Ready, ContainersReady)
# - Container statuses
# - Image pull status
```

### Step 3: Check Logs
```bash
# Current logs
kubectl logs <pod-name>
kl <pod-name>  # alias

# Previous container logs (if crashed)
kubectl logs <pod-name> --previous

# Follow logs
klf <pod-name>  # dotfiles function

# Specific container (multi-container pod)
kubectl logs <pod-name> -c <container>
klf <pod-name> <container>
```

### Step 4: Exec Into Pod
```bash
# If pod is running
kubectl exec -it <pod-name> -- /bin/sh
ksh <pod-name>  # dotfiles function

# Specific container
kubectl exec -it <pod-name> -c <container> -- /bin/sh
```

## Common Issues and Solutions

### ImagePullBackOff
```bash
# Check image name and tag
kubectl describe pod <pod> | grep -A5 "Image:"

# Check image pull secrets
kubectl get pod <pod> -o jsonpath='{.spec.imagePullSecrets}'

# Verify secret exists
kubectl get secrets
```

### CrashLoopBackOff
```bash
# Check previous logs
kubectl logs <pod> --previous

# Check exit code
kubectl describe pod <pod> | grep -A5 "Last State"

# Common causes:
# - Application error
# - Missing config/secrets
# - Resource limits too low
```

### Pending State
```bash
# Check events
kubectl describe pod <pod> | grep -A10 "Events"

# Common causes:
# - No nodes with resources
# - PVC not bound
# - Node selector mismatch

# Check node resources
kubectl describe nodes | grep -A5 "Allocated resources"
```

### OOMKilled
```bash
# Check memory limits
kubectl describe pod <pod> | grep -A5 "Limits"

# Solution: Increase memory limit
kubectl set resources deployment/<deploy> --limits=memory=512Mi
```

### Readiness Probe Failed
```bash
# Check probe configuration
kubectl describe pod <pod> | grep -A10 "Readiness"

# Test endpoint manually
kubectl exec <pod> -- wget -qO- localhost:8080/health
```

## Network Debugging

```bash
# Test DNS
kubectl exec <pod> -- nslookup kubernetes

# Test service connectivity
kubectl exec <pod> -- wget -qO- <service>:<port>

# Check service endpoints
kubectl get endpoints <service>
```

## Resource Check

```bash
# Pod resource usage
kubectl top pod <pod>

# Container resource usage
kubectl top pod <pod> --containers
```

## Quick Debug Commands

```bash
# All-in-one status
kubectl get pod <pod> -o yaml

# Events only
kubectl get events --field-selector involvedObject.name=<pod>

# Watch pod status
kubectl get pod <pod> -w
kwatch  # dotfiles function
```

## Dotfiles Integration

- `kpods [filter]` - List pods with filter
- `klf <pod> [container]` - Follow logs
- `ksh <pod>` - Exec into pod
- `kd` - kubectl describe (alias)
- `kl` - kubectl logs (alias)
