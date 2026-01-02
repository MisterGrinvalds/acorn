#!/bin/sh
# components/kubernetes/aliases.sh - Kubernetes and Helm aliases

# =============================================================================
# kubectl aliases
# =============================================================================
alias k='kubectl'
alias kd='kubectl describe'
alias kg='kubectl get'
alias kl='kubectl logs'
alias kx='kubectl exec -it'
alias kdel='kubectl delete'
alias kaf='kubectl apply -f'
alias kdf='kubectl delete -f'

# Common resources
alias kgp='kubectl get pods'
alias kgs='kubectl get services'
alias kgd='kubectl get deployments'
alias kgn='kubectl get nodes'
alias kgcm='kubectl get configmaps'
alias kgsec='kubectl get secrets'

# Context and namespace
alias kctx='kubectl config current-context'
alias kns='kubectl config view --minify --output jsonpath={..namespace}'
alias kgctx='kubectl config get-contexts'
alias kgns='kubectl get namespaces'

# =============================================================================
# Helm aliases
# =============================================================================
alias h='helm'
alias hls='helm list'
alias hla='helm list -A'
alias hin='helm install'
alias hup='helm upgrade'
alias hun='helm uninstall'
alias hval='helm get values'
alias hs='helm search'
alias hsr='helm search repo'

# =============================================================================
# k9s shortcut
# =============================================================================
alias k9='k9s'
