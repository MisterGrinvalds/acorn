#!/bin/sh
# components/kubernetes/env.sh - Kubernetes environment variables

# Kubeconfig location
export KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}"

# Helm XDG paths
export HELM_CONFIG_HOME="${XDG_CONFIG_HOME:-$HOME/.config}/helm"
export HELM_DATA_HOME="${XDG_DATA_HOME:-$HOME/.local/share}/helm"
export HELM_CACHE_HOME="${XDG_CACHE_HOME:-$HOME/.cache}/helm"
