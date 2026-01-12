---
description: Explain Kubernetes concepts, resources, and workflows
argument-hint: [topic]
allowed-tools: Read, Glob, Grep
---

## Task

Explain the requested topic about Kubernetes. If no specific topic provided, give an overview of Kubernetes and the dotfiles integration.

## Topics

If `$ARGUMENTS` is provided, explain that specific topic:

- **pods** - Pod lifecycle and management
- **deployments** - Deployment strategies and rollouts
- **services** - Service types and networking
- **namespaces** - Namespace isolation
- **configmaps** - Configuration management
- **secrets** - Secret management
- **ingress** - Ingress controllers and routing
- **volumes** - Persistent volumes and claims
- **rbac** - Role-based access control
- **helm** - Helm package manager
- **k9s** - k9s terminal UI
- **contexts** - Context and cluster management
- **probes** - Liveness and readiness probes
- **hpa** - Horizontal pod autoscaling

## Context

Reference these files for accurate information:
@components/kubernetes/component.yaml
@components/kubernetes/functions.sh
@components/kubernetes/aliases.sh

## Response Format

1. **Concept overview** - What it is and why it matters
2. **Key commands** - Essential kubectl commands
3. **Dotfiles integration** - Available functions and aliases
4. **Examples** - Practical usage examples
5. **Best practices** - Production recommendations
