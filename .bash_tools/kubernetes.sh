# Kubernetes Development Tools

# Kubectl aliases for common operations
alias k='kubectl'
alias kd='kubectl describe'
alias kg='kubectl get'
alias kl='kubectl logs'
alias kx='kubectl exec -it'
alias kdel='kubectl delete'
alias kaf='kubectl apply -f'
alias kdf='kubectl delete -f'

# Kubernetes context and namespace management
alias kctx='kubectl config current-context'
alias kns='kubectl config view --minify --output jsonpath={..namespace}'
alias kgctx='kubectl config get-contexts'
alias kgns='kubectl get namespaces'

# Helm aliases
alias h='helm'
alias hls='helm list'
alias hla='helm list -A'
alias hin='helm install'
alias hup='helm upgrade'
alias hun='helm uninstall'
alias hval='helm get values'

# k9s shortcut
alias k9='k9s'

# Quick pod operations
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
        echo "Usage: klf <pod-name>"
        return 1
    fi
    kubectl logs -f "$1"
}

# Port forward helper
kpf() {
    if [ -z "$2" ]; then
        echo "Usage: kpf <pod-name> <local-port:remote-port>"
        return 1
    fi
    kubectl port-forward "$1" "$2"
}

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

# Kubernetes environment variables for common tools
export KUBECONFIG="$HOME/.kube/config"