#!/bin/bash
# Kubernetes Automation Module

# Module configuration
readonly K8S_MANIFESTS_DIR="$AUTO_HOME/k8s/manifests"
readonly K8S_CONFIGS_DIR="$AUTO_HOME/k8s/configs"
readonly K8S_CHARTS_DIR="$AUTO_HOME/k8s/charts"

# Initialize directories
mkdir -p "$K8S_MANIFESTS_DIR" "$K8S_CONFIGS_DIR" "$K8S_CHARTS_DIR"

# Help function for k8s module
k8s_help() {
    cat << EOF
Kubernetes Automation

USAGE:
    auto k8s <command> [options]

COMMANDS:
    cluster <action>              Cluster management (info, switch, list)
    deploy <app> [environment]    Deploy application
    scale <deployment> <replicas> Scale deployment
    logs <pod> [follow]           Get pod logs
    exec <pod> [command]          Execute command in pod
    port-forward <pod> <ports>    Port forward to pod
    status [namespace]            Get cluster status
    manifests <action>            Manage manifests (generate, validate, apply)
    helm <action>                 Helm operations (install, upgrade, uninstall)
    monitoring                    Setup monitoring stack
    backup <namespace>            Backup namespace resources
    cleanup                       Clean up unused resources

EXAMPLES:
    auto k8s cluster info
    auto k8s deploy my-app production
    auto k8s scale my-deployment 3
    auto k8s logs my-pod --follow
    auto k8s port-forward my-pod 8080:80
    auto k8s manifests generate my-app
    auto k8s helm install prometheus monitoring/prometheus
    auto k8s monitoring setup
EOF
}

# Utility functions
require_kubectl() {
    require_command kubectl
}

require_helm() {
    require_command helm
}

get_current_context() {
    kubectl config current-context 2>/dev/null || echo "No context set"
}

get_current_namespace() {
    kubectl config view --minify --output 'jsonpath={..namespace}' 2>/dev/null || echo "default"
}

# Cluster management
cluster_info() {
    require_kubectl
    
    log "INFO" "Kubernetes Cluster Information"
    echo "================================"
    echo "Context: $(get_current_context)"
    echo "Namespace: $(get_current_namespace)"
    echo ""
    
    kubectl cluster-info
    echo ""
    
    log "INFO" "Node Status:"
    kubectl get nodes -o wide
    echo ""
    
    log "INFO" "Namespace Resources:"
    kubectl get all --namespace="$(get_current_namespace)"
}

cluster_switch() {
    require_kubectl
    local context="$1"
    
    if [ -z "$context" ]; then
        log "INFO" "Available contexts:"
        kubectl config get-contexts
        return 0
    fi
    
    kubectl config use-context "$context"
    log "SUCCESS" "Switched to context: $context"
}

# Application deployment
deploy_app() {
    require_kubectl
    local app_name="$1"
    local environment="${2:-default}"
    local manifest_file="$K8S_MANIFESTS_DIR/$app_name-$environment.yaml"
    
    if [ ! -f "$manifest_file" ]; then
        log "WARN" "Manifest not found: $manifest_file"
        log "INFO" "Generating manifest from template..."
        generate_manifest "$app_name" "$environment"
    fi
    
    if [ -f "$manifest_file" ]; then
        log "INFO" "Deploying $app_name to $environment environment..."
        kubectl apply -f "$manifest_file"
        
        # Wait for deployment to be ready
        if grep -q "kind: Deployment" "$manifest_file"; then
            local deployment_name=$(grep -A 5 "kind: Deployment" "$manifest_file" | grep "name:" | head -1 | awk '{print $2}')
            kubectl rollout status deployment/"$deployment_name" --timeout=300s
            log "SUCCESS" "Deployment $deployment_name is ready"
        fi
    else
        log "ERROR" "Could not find or generate manifest for $app_name"
        exit 1
    fi
}

# Manifest generation
generate_manifest() {
    local app_name="$1"
    local environment="${2:-default}"
    local image="${3:-$app_name:latest}"
    local port="${4:-8080}"
    
    local manifest_file="$K8S_MANIFESTS_DIR/$app_name-$environment.yaml"
    
    cat > "$manifest_file" << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: $app_name
  labels:
    app: $app_name
    environment: $environment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: $app_name
  template:
    metadata:
      labels:
        app: $app_name
        environment: $environment
    spec:
      containers:
      - name: $app_name
        image: $image
        ports:
        - containerPort: $port
        env:
        - name: ENVIRONMENT
          value: "$environment"
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: $port
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: $port
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: $app_name-service
  labels:
    app: $app_name
spec:
  selector:
    app: $app_name
  ports:
  - protocol: TCP
    port: 80
    targetPort: $port
  type: ClusterIP
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: $app_name-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
spec:
  rules:
  - host: $app_name-$environment.local
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: $app_name-service
            port:
              number: 80
EOF
    
    log "SUCCESS" "Generated manifest: $manifest_file"
}

# Scaling operations
scale_deployment() {
    require_kubectl
    local deployment="$1"
    local replicas="$2"
    
    if [ -z "$deployment" ] || [ -z "$replicas" ]; then
        log "ERROR" "Usage: auto k8s scale <deployment> <replicas>"
        exit 1
    fi
    
    kubectl scale deployment "$deployment" --replicas="$replicas"
    kubectl rollout status deployment/"$deployment"
    log "SUCCESS" "Scaled $deployment to $replicas replicas"
}

# Log management
get_logs() {
    require_kubectl
    local pod="$1"
    local follow="${2:-false}"
    
    if [ -z "$pod" ]; then
        log "INFO" "Available pods:"
        kubectl get pods
        return 0
    fi
    
    if [ "$follow" = "--follow" ] || [ "$follow" = "-f" ]; then
        kubectl logs -f "$pod"
    else
        kubectl logs "$pod"
    fi
}

# Port forwarding
port_forward() {
    require_kubectl
    local pod="$1"
    local ports="$2"
    
    if [ -z "$pod" ] || [ -z "$ports" ]; then
        log "ERROR" "Usage: auto k8s port-forward <pod> <local-port:remote-port>"
        exit 1
    fi
    
    log "INFO" "Port forwarding $pod on $ports (Ctrl+C to stop)"
    kubectl port-forward "$pod" "$ports"
}

# Helm operations
helm_install() {
    require_helm
    local chart="$1"
    local release_name="$2"
    local namespace="${3:-default}"
    
    if [ -z "$chart" ] || [ -z "$release_name" ]; then
        log "ERROR" "Usage: auto k8s helm install <chart> <release-name> [namespace]"
        exit 1
    fi
    
    helm install "$release_name" "$chart" --namespace "$namespace" --create-namespace
    log "SUCCESS" "Helm chart $chart installed as $release_name"
}

helm_upgrade() {
    require_helm
    local release_name="$1"
    local chart="$2"
    local namespace="${3:-default}"
    
    helm upgrade "$release_name" "$chart" --namespace "$namespace"
    log "SUCCESS" "Helm release $release_name upgraded"
}

# Monitoring setup
setup_monitoring() {
    require_helm
    log "INFO" "Setting up monitoring stack (Prometheus + Grafana)"
    
    # Add Prometheus Helm repository
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo update
    
    # Create monitoring namespace
    kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -
    
    # Install Prometheus
    helm upgrade --install prometheus prometheus-community/kube-prometheus-stack \
        --namespace monitoring \
        --set grafana.adminPassword=admin123 \
        --set alertmanager.enabled=true
    
    log "SUCCESS" "Monitoring stack installed"
    log "INFO" "Access Grafana: kubectl port-forward -n monitoring svc/prometheus-grafana 3000:80"
    log "INFO" "Default credentials: admin/admin123"
}

# Backup operations
backup_namespace() {
    require_kubectl
    local namespace="${1:-$(get_current_namespace)}"
    local backup_dir="$AUTO_CACHE/k8s-backups/$(date +%Y%m%d_%H%M%S)"
    
    mkdir -p "$backup_dir"
    
    log "INFO" "Backing up namespace: $namespace"
    
    # Backup all resources
    local resources=("deployments" "services" "configmaps" "secrets" "ingresses" "persistentvolumeclaims")
    
    for resource in "${resources[@]}"; do
        kubectl get "$resource" -n "$namespace" -o yaml > "$backup_dir/$resource.yaml" 2>/dev/null || true
    done
    
    # Create archive
    tar -czf "$backup_dir.tar.gz" -C "$(dirname "$backup_dir")" "$(basename "$backup_dir")"
    rm -rf "$backup_dir"
    
    log "SUCCESS" "Namespace backup saved: $backup_dir.tar.gz"
}

# Cleanup operations
cleanup_resources() {
    require_kubectl
    
    log "INFO" "Cleaning up unused Kubernetes resources..."
    
    # Remove completed jobs older than 1 hour
    kubectl get jobs --all-namespaces -o json | \
        jq -r '.items[] | select(.status.conditions[]?.type == "Complete") | select(.status.completionTime | fromdateiso8601 < (now - 3600)) | "\(.metadata.namespace) \(.metadata.name)"' | \
        while read -r namespace job; do
            kubectl delete job "$job" -n "$namespace"
            log "INFO" "Deleted completed job: $job in namespace $namespace"
        done
    
    # Remove evicted pods
    kubectl get pods --all-namespaces --field-selector=status.phase=Failed -o json | \
        jq -r '.items[] | select(.status.reason == "Evicted") | "\(.metadata.namespace) \(.metadata.name)"' | \
        while read -r namespace pod; do
            kubectl delete pod "$pod" -n "$namespace"
            log "INFO" "Deleted evicted pod: $pod in namespace $namespace"
        done
    
    log "SUCCESS" "Cleanup completed"
}

# Main k8s module function
k8s_main() {
    local command="${1:-help}"
    shift || true
    
    case "$command" in
        "help"|"-h"|"--help")
            k8s_help
            ;;
        "cluster")
            local action="${1:-info}"
            case "$action" in
                "info") cluster_info ;;
                "switch") cluster_switch "$2" ;;
                "list") kubectl config get-contexts ;;
                *) log "ERROR" "Unknown cluster action: $action" ;;
            esac
            ;;
        "deploy")
            deploy_app "$1" "$2"
            ;;
        "scale")
            scale_deployment "$1" "$2"
            ;;
        "logs")
            get_logs "$1" "$2"
            ;;
        "port-forward"|"pf")
            port_forward "$1" "$2"
            ;;
        "manifests")
            local action="$1"
            case "$action" in
                "generate") generate_manifest "$2" "$3" "$4" "$5" ;;
                "validate") kubectl apply --dry-run=client -f "$K8S_MANIFESTS_DIR/$2" ;;
                "apply") kubectl apply -f "$K8S_MANIFESTS_DIR/$2" ;;
                *) log "ERROR" "Unknown manifests action: $action" ;;
            esac
            ;;
        "helm")
            local action="$1"
            case "$action" in
                "install") helm_install "$2" "$3" "$4" ;;
                "upgrade") helm_upgrade "$2" "$3" "$4" ;;
                "list") helm list --all-namespaces ;;
                "uninstall") helm uninstall "$2" --namespace "${3:-default}" ;;
                *) log "ERROR" "Unknown helm action: $action" ;;
            esac
            ;;
        "monitoring")
            setup_monitoring
            ;;
        "backup")
            backup_namespace "$1"
            ;;
        "cleanup")
            cleanup_resources
            ;;
        "status")
            local namespace="${1:-$(get_current_namespace)}"
            kubectl get all -n "$namespace"
            ;;
        *)
            log "ERROR" "Unknown command: $command"
            k8s_help
            exit 1
            ;;
    esac
}