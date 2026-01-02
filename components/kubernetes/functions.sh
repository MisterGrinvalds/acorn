#!/bin/sh
# components/kubernetes/functions.sh - Kubernetes helper functions

# =============================================================================
# Pod Operations
# =============================================================================

# List pods with optional filter
kpods() {
    if [ -z "$1" ]; then
        kubectl get pods
    else
        kubectl get pods | grep "$1"
    fi
}

# Get pod logs with follow
klf() {
    if [ -z "$1" ]; then
        echo "Usage: klf <pod-name> [container]"
        return 1
    fi
    if [ -n "$2" ]; then
        kubectl logs -f "$1" -c "$2"
    else
        kubectl logs -f "$1"
    fi
}

# Port forward helper
kpf() {
    if [ -z "$2" ]; then
        echo "Usage: kpf <pod-name> <local-port:remote-port>"
        return 1
    fi
    kubectl port-forward "$1" "$2"
}

# Exec into pod
ksh() {
    if [ -z "$1" ]; then
        echo "Usage: ksh <pod-name> [command]"
        return 1
    fi
    local cmd="${2:-/bin/sh}"
    kubectl exec -it "$1" -- "$cmd"
}

# =============================================================================
# Context and Namespace Management
# =============================================================================

# Quick context switching
kuse() {
    if [ -z "$1" ]; then
        kubectl config get-contexts
        return 0
    fi
    kubectl config use-context "$1"
}

# Namespace switching
knsuse() {
    if [ -z "$1" ]; then
        kubectl get namespaces
        return 0
    fi
    kubectl config set-context --current --namespace="$1"
}

# Show current context and namespace
kinfo() {
    echo "Context:   $(kubectl config current-context)"
    echo "Namespace: $(kubectl config view --minify --output 'jsonpath={..namespace}' 2>/dev/null || echo 'default')"
    echo "Server:    $(kubectl config view --minify --output 'jsonpath={.clusters[0].cluster.server}')"
}

# =============================================================================
# Resource Helpers
# =============================================================================

# Get all resources in namespace
kall() {
    local ns="${1:-$(kubectl config view --minify --output 'jsonpath={..namespace}' 2>/dev/null)}"
    ns="${ns:-default}"
    echo "=== Pods ==="
    kubectl get pods -n "$ns"
    echo ""
    echo "=== Services ==="
    kubectl get services -n "$ns"
    echo ""
    echo "=== Deployments ==="
    kubectl get deployments -n "$ns"
}

# Watch pods
kwatch() {
    kubectl get pods -w "$@"
}

# Delete all evicted pods
kcleanpods() {
    kubectl get pods --all-namespaces -o json | \
        jq -r '.items[] | select(.status.reason=="Evicted") | "\(.metadata.namespace) \(.metadata.name)"' | \
        while read -r ns name; do
            kubectl delete pod -n "$ns" "$name"
        done
}
