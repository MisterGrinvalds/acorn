---
description: Interactive coaching session to learn Kubernetes
argument-hint: [skill-level: beginner|intermediate|advanced]
allowed-tools: Read, Glob, Grep, Bash
---

## Task

Guide the user through learning Kubernetes interactively.

## Approach

1. **Assess level** - Ask about Kubernetes experience
2. **Set goals** - Identify what they want to manage
3. **Progressive exercises** - Start simple, build complexity
4. **Real-time practice** - Have them run kubectl commands
5. **Reinforce** - Summarize and suggest next steps

## Skill Levels

### Beginner
- What is Kubernetes?
- Pods and containers
- Basic kubectl commands
- Viewing resources and logs
- Namespaces basics

### Intermediate
- Deployments and ReplicaSets
- Services and networking
- ConfigMaps and Secrets
- Context switching
- Using dotfiles functions
- Helm basics

### Advanced
- RBAC and security
- Custom resources
- Helm chart development
- Debugging complex issues
- Performance optimization
- k9s power usage

## Interactive Exercises

### Beginner Exercises
```bash
# Exercise 1: Check cluster info
kinfo  # dotfiles function

# Exercise 2: List pods
kgp  # kubectl get pods (alias)

# Exercise 3: List all namespaces
kgns  # kubectl get namespaces (alias)

# Exercise 4: Describe a pod
kd pod <pod-name>

# Exercise 5: View logs
kl <pod-name>
```

### Intermediate Exercises
```bash
# Exercise 6: Switch namespace
knsuse kube-system

# Exercise 7: Get all resources
kall  # dotfiles function

# Exercise 8: Follow logs
klf <pod-name>  # dotfiles function

# Exercise 9: Exec into pod
ksh <pod-name>  # dotfiles function

# Exercise 10: Port forward
kpf <pod-name> 8080:80
```

### Advanced Exercises
```bash
# Exercise 11: Use k9s
k9  # alias

# Exercise 12: Helm operations
hls  # helm list
hla  # helm list all namespaces

# Exercise 13: Debug failing pod
kubectl describe pod <failing-pod>
kubectl logs <pod> --previous

# Exercise 14: Clean up
kcleanpods  # Delete evicted pods
```

## Context

@components/kubernetes/functions.sh
@components/kubernetes/aliases.sh

## Coaching Style

- Always start with `kinfo` to confirm context
- Use dotfiles aliases to build muscle memory
- Emphasize namespace awareness
- Show k9s for visual exploration
- Build toward production debugging skills
