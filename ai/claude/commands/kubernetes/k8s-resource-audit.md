---
description: Audit Kubernetes resources and cluster health
argument-hint: [namespace]
allowed-tools: Read, Bash
---

## Task

Help the user audit Kubernetes resources and cluster health.

## Quick Audit

Using dotfiles function:
```bash
kall  # Shows pods, services, deployments in current namespace
kall <namespace>  # Specific namespace
```

## Resource Inventory

### All Resources in Namespace
```bash
# Common resources
kubectl get all -n <namespace>

# Including configmaps, secrets, etc.
kubectl get all,cm,secret,ing -n <namespace>

# All resource types
kubectl api-resources --verbs=list -o name | xargs -n 1 kubectl get -n <namespace> 2>/dev/null
```

### Cluster-Wide Resources
```bash
# Nodes
kubectl get nodes -o wide
kgn  # alias

# Namespaces
kubectl get namespaces
kgns  # alias

# Persistent Volumes
kubectl get pv

# Cluster roles
kubectl get clusterroles
```

## Health Checks

### Node Health
```bash
# Node status
kubectl get nodes

# Node conditions
kubectl describe nodes | grep -A5 "Conditions:"

# Node resources
kubectl top nodes
```

### Pod Health
```bash
# Unhealthy pods
kubectl get pods --all-namespaces | grep -v Running | grep -v Completed

# Pods with restarts
kubectl get pods --all-namespaces -o wide | awk '$5 > 0'

# Evicted pods
kubectl get pods --all-namespaces -o json | jq -r '.items[] | select(.status.reason=="Evicted") | .metadata.name'

# Clean evicted pods
kcleanpods  # dotfiles function
```

### Resource Usage
```bash
# Pod resource usage
kubectl top pods -A

# High memory pods
kubectl top pods -A --sort-by=memory

# High CPU pods
kubectl top pods -A --sort-by=cpu
```

## Security Audit

### RBAC
```bash
# Service accounts
kubectl get serviceaccounts -A

# Roles and bindings
kubectl get roles,rolebindings -A
kubectl get clusterroles,clusterrolebindings

# Who can do what
kubectl auth can-i --list
kubectl auth can-i create pods -n <namespace>
```

### Secrets
```bash
# List secrets (not values)
kubectl get secrets -A

# Unused secrets (requires scripting)
# Compare secrets list with pod references
```

### Pod Security
```bash
# Pods running as root
kubectl get pods -A -o json | jq -r '.items[] | select(.spec.containers[].securityContext.runAsUser == 0 or .spec.securityContext.runAsUser == 0) | .metadata.name'

# Privileged pods
kubectl get pods -A -o json | jq -r '.items[] | select(.spec.containers[].securityContext.privileged == true) | .metadata.name'
```

## Resource Quotas

```bash
# View quotas
kubectl get resourcequota -A

# Describe quota usage
kubectl describe resourcequota -n <namespace>

# Limit ranges
kubectl get limitrange -A
```

## Networking Audit

### Services
```bash
# All services
kubectl get services -A

# Services without endpoints
kubectl get endpoints -A | awk '$2 == "<none>"'
```

### Ingress
```bash
# All ingress rules
kubectl get ingress -A

# Describe for details
kubectl describe ingress -A
```

### Network Policies
```bash
kubectl get networkpolicies -A
```

## Storage Audit

```bash
# Persistent Volume Claims
kubectl get pvc -A

# Unbound PVCs
kubectl get pvc -A | grep -v Bound

# Persistent Volumes
kubectl get pv
```

## Events

```bash
# Recent events
kubectl get events -A --sort-by='.lastTimestamp' | tail -20

# Warning events
kubectl get events -A --field-selector type=Warning

# Events for specific resource
kubectl get events --field-selector involvedObject.name=<pod-name>
```

## Report Generation

```bash
# Save resource summary
kubectl get all -A -o wide > cluster-audit.txt

# JSON for processing
kubectl get all -A -o json > cluster-audit.json
```

## Dotfiles Integration

- `kall [ns]` - Get pods, services, deployments
- `kcleanpods` - Delete evicted pods
- `kgp` - kubectl get pods (alias)
- `kgs` - kubectl get services (alias)
- `kgd` - kubectl get deployments (alias)
- `kgn` - kubectl get nodes (alias)
