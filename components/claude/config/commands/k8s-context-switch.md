---
description: Switch Kubernetes contexts and namespaces
argument-hint: [context-name] [--namespace]
allowed-tools: Read, Bash
---

## Task

Help the user switch between Kubernetes contexts and namespaces.

## Context Management

### View Current Context
```bash
# Current context
kubectl config current-context
kctx  # alias

# Full info
kinfo  # dotfiles function
# Shows: context, namespace, server
```

### List Contexts
```bash
kubectl config get-contexts
kgctx  # alias

# Using dotfiles function (shows list)
kuse
```

### Switch Context
```bash
# Using dotfiles function
kuse <context-name>

# Direct kubectl
kubectl config use-context <context-name>
```

## Namespace Management

### View Current Namespace
```bash
kubectl config view --minify -o jsonpath='{..namespace}'
kns  # alias
```

### List Namespaces
```bash
kubectl get namespaces
kgns  # alias

# Using dotfiles function (shows list)
knsuse
```

### Switch Namespace
```bash
# Using dotfiles function
knsuse <namespace>

# Direct kubectl
kubectl config set-context --current --namespace=<namespace>
```

## Common Workflows

### Switch Environment
```bash
# Example: dev -> staging -> prod
kuse dev-cluster
knsuse development

kuse staging-cluster
knsuse staging

kuse prod-cluster
knsuse production
```

### Quick Context Check Before Commands
```bash
# Always verify before destructive operations
kinfo
# Then proceed with command
kubectl delete pod <pod>
```

## Context Configuration

### View Config
```bash
kubectl config view
kubectl config view --minify  # Current context only
```

### Add Context (from kubeconfig)
```bash
# Merge kubeconfigs
export KUBECONFIG=~/.kube/config:~/.kube/new-config
kubectl config view --flatten > ~/.kube/merged-config
```

### Rename Context
```bash
kubectl config rename-context old-name new-name
```

### Delete Context
```bash
kubectl config delete-context <context-name>
```

## FZF Integration

If fzf component is loaded:
```bash
# Interactive context selection
fzf_k8s_ns  # or: fkns (alias)

# Select namespace with fzf
kubectl get namespaces -o name | sed 's/namespace\///' | fzf
```

## Kubeconfig Management

### Multiple Kubeconfigs
```bash
# Set for session
export KUBECONFIG=/path/to/config

# Multiple files
export KUBECONFIG=~/.kube/config:~/.kube/work-config

# Per-command
kubectl --kubeconfig=/path/to/config get pods
```

### Default Location
```
~/.kube/config
```

## Safety Tips

1. **Always check context** before destructive operations
2. **Use different prompts** for different clusters (e.g., colored PS1)
3. **Alias dangerous commands** to require confirmation
4. **Set default namespace** to avoid mistakes

### Shell Prompt Integration
```bash
# Add to PS1 (shows context:namespace)
export PS1="[\$(kubectl config current-context):\$(kubectl config view --minify -o jsonpath='{..namespace}')] $PS1"
```

## Dotfiles Integration

- `kuse [context]` - Switch context (or list)
- `knsuse [namespace]` - Switch namespace (or list)
- `kinfo` - Show current context, namespace, server
- `kctx` - Current context (alias)
- `kns` - Current namespace (alias)
- `kgctx` - List contexts (alias)
- `kgns` - List namespaces (alias)
