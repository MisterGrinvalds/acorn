---
description: View and analyze Kubernetes pod logs
argument-hint: <pod-name> [--follow] [--previous]
allowed-tools: Read, Bash
---

## Task

Help the user view and analyze Kubernetes pod logs.

## Basic Log Commands

### View Logs
```bash
# Current logs
kubectl logs <pod-name>
kl <pod-name>  # alias

# Follow logs in real-time
kubectl logs -f <pod-name>
klf <pod-name>  # dotfiles function
```

### Previous Container Logs
```bash
# After crash/restart
kubectl logs <pod-name> --previous

# Useful for CrashLoopBackOff debugging
```

### Multi-Container Pods
```bash
# List containers
kubectl get pod <pod> -o jsonpath='{.spec.containers[*].name}'

# Logs from specific container
kubectl logs <pod> -c <container>
klf <pod> <container>  # dotfiles function

# All containers
kubectl logs <pod> --all-containers
```

## Filtering and Searching

### Tail Logs
```bash
# Last N lines
kubectl logs <pod> --tail=100

# Last N lines, then follow
kubectl logs <pod> --tail=50 -f
```

### Time-Based
```bash
# Since time
kubectl logs <pod> --since=1h
kubectl logs <pod> --since=30m
kubectl logs <pod> --since=10s

# Since timestamp
kubectl logs <pod> --since-time="2024-01-01T10:00:00Z"
```

### Grep/Filter
```bash
# Pipe to grep
kubectl logs <pod> | grep ERROR
kubectl logs <pod> | grep -i "exception"

# With timestamps
kubectl logs <pod> --timestamps | grep "10:30"
```

## Multiple Pods

### By Label
```bash
# All pods with label
kubectl logs -l app=myapp

# Follow all pods
kubectl logs -l app=myapp -f

# All containers in labeled pods
kubectl logs -l app=myapp --all-containers
```

### By Deployment
```bash
# Logs from deployment's pods
kubectl logs deployment/myapp

# Follow deployment logs
kubectl logs -f deployment/myapp
```

## Log Aggregation Tools

### Stern (if installed)
```bash
# Follow multiple pods by pattern
stern "myapp-.*"

# With namespace
stern -n production "api-.*"
```

### kubectl plugins
```bash
# Using krew plugin manager
kubectl krew install tail
kubectl tail -l app=myapp
```

## Log Analysis Patterns

### Error Summary
```bash
# Count errors
kubectl logs <pod> | grep -c ERROR

# Unique errors
kubectl logs <pod> | grep ERROR | sort | uniq -c | sort -rn
```

### Request Tracing
```bash
# Find specific request
kubectl logs <pod> | grep "request-id-123"

# Time range
kubectl logs <pod> --since=5m | grep "request-id"
```

### Performance
```bash
# Slow requests (example pattern)
kubectl logs <pod> | grep "took" | awk '$NF > 1000'
```

## k9s Log Viewing

```bash
# Launch k9s
k9  # alias

# In k9s:
# - Navigate to Pods
# - Press 'l' to view logs
# - Press '0' for all containers
# - Press 'w' to wrap lines
# - '/' to search
```

## Best Practices

1. **Use --tail** to avoid overwhelming output
2. **Add timestamps** when debugging timing issues
3. **Use labels** to aggregate related pod logs
4. **Save logs** before pod deletion: `kubectl logs <pod> > pod.log`
5. **Check previous** logs after crashes

## Dotfiles Integration

- `klf <pod> [container]` - Follow logs
- `kl` - kubectl logs (alias)
- `k9` - k9s for visual log viewing
