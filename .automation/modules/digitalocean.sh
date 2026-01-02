#!/usr/bin/env bash
# DigitalOcean Cloud Automation Module

# Module configuration
readonly DO_TEMPLATES_DIR="$AUTO_HOME/cloud/digitalocean/templates"
readonly DO_CONFIGS_DIR="$AUTO_HOME/cloud/digitalocean/configs"
readonly DO_API_URL="https://api.digitalocean.com/v2"

# Initialize directories
mkdir -p "$DO_TEMPLATES_DIR" "$DO_CONFIGS_DIR"

# Help function for DigitalOcean module
digitalocean_help() {
    cat << EOF
DigitalOcean Cloud Automation

USAGE:
    auto digitalocean <command> [options]
    auto do <command> [options]          # Short alias

COMMANDS:
    auth                          Authentication setup
    droplets <action>             Droplet management
    kubernetes <action>           Kubernetes cluster operations
    volumes <action>              Block storage management
    spaces <action>               Spaces (S3-compatible) operations
    databases <action>            Managed database operations
    load-balancers <action>       Load balancer management
    firewalls <action>            Firewall management
    monitoring                    Monitoring and alerts
    backups <action>              Backup management
    domains <action>              DNS and domain management
    cost                          Billing and usage analysis

EXAMPLES:
    auto do auth setup
    auto do droplets create web-server s-1vcpu-1gb nyc3
    auto do kubernetes create my-cluster
    auto do volumes create data-volume 100 nyc3
    auto do spaces create my-bucket nyc3
    auto do databases create my-db mysql
    auto do load-balancers create web-lb
    auto do cost analyze
EOF
}

# Utility functions
require_doctl() {
    if ! command_exists doctl; then
        log "ERROR" "doctl CLI not found. Please install it first:"
        log "INFO" "  macOS: brew install doctl"
        log "INFO" "  Linux: snap install doctl"
        log "INFO" "  Manual: https://docs.digitalocean.com/reference/doctl/how-to/install/"
        exit 1
    fi
}

do_api_call() {
    local method="$1"
    local endpoint="$2"
    local data="$3"
    
    local token=$(doctl auth list --format Name --no-header 2>/dev/null | head -1)
    if [ -z "$token" ]; then
        log "ERROR" "No DigitalOcean authentication found. Run: auto do auth setup"
        exit 1
    fi
    
    local curl_args=(
        -H "Authorization: Bearer $(doctl auth list --format Token --no-header | head -1)"
        -H "Content-Type: application/json"
    )
    
    if [ "$method" = "POST" ] || [ "$method" = "PUT" ]; then
        curl_args+=(-X "$method" -d "$data")
    elif [ "$method" = "DELETE" ]; then
        curl_args+=(-X "$method")
    fi
    
    curl -s "${curl_args[@]}" "$DO_API_URL/$endpoint"
}

# Authentication setup
do_auth_setup() {
    require_doctl
    
    log "INFO" "Setting up DigitalOcean authentication..."
    
    # Check if already authenticated
    if doctl account get >/dev/null 2>&1; then
        log "INFO" "Already authenticated with DigitalOcean"
        doctl account get
        return 0
    fi
    
    log "INFO" "Please visit: https://cloud.digitalocean.com/account/api/tokens"
    log "INFO" "Create a new personal access token with read/write permissions"
    echo ""
    echo -n "Enter your DigitalOcean API token: "
    read -rs token
    echo ""
    
    # Validate token
    doctl auth init --access-token "$token"
    
    if doctl account get >/dev/null 2>&1; then
        log "SUCCESS" "DigitalOcean authentication successful"
        doctl account get
    else
        log "ERROR" "Authentication failed. Please check your token."
        exit 1
    fi
}

# Droplet operations
do_droplets_list() {
    require_doctl
    
    log "INFO" "Listing DigitalOcean droplets:"
    doctl compute droplet list --format "ID,Name,PublicIPv4,PrivateIPv4,Status,Size,Region"
}

do_droplets_create() {
    require_doctl
    local name="$1"
    local size="${2:-s-1vcpu-1gb}"
    local region="${3:-nyc3}"
    local image="${4:-ubuntu-22-04-x64}"
    local ssh_keys="$5"
    
    if [ -z "$name" ]; then
        log "ERROR" "Droplet name is required"
        exit 1
    fi
    
    local create_args=(
        "$name"
        --size "$size"
        --image "$image"
        --region "$region"
        --enable-monitoring
        --enable-private-networking
    )
    
    # Add SSH keys if provided
    if [ -n "$ssh_keys" ]; then
        create_args+=(--ssh-keys "$ssh_keys")
    else
        # Try to get default SSH key
        local default_key=$(doctl compute ssh-key list --format ID --no-header | head -1)
        if [ -n "$default_key" ]; then
            create_args+=(--ssh-keys "$default_key")
        fi
    fi
    
    log "INFO" "Creating droplet: $name"
    
    local droplet_id=$(doctl compute droplet create "${create_args[@]}" --format ID --no-header)
    
    if [ -n "$droplet_id" ]; then
        log "SUCCESS" "Droplet created with ID: $droplet_id"
        
        # Wait for droplet to be active
        log "INFO" "Waiting for droplet to become active..."
        while true; do
            local status=$(doctl compute droplet get "$droplet_id" --format Status --no-header)
            if [ "$status" = "active" ]; then
                break
            fi
            sleep 5
        done
        
        # Show droplet info
        doctl compute droplet get "$droplet_id"
    else
        log "ERROR" "Failed to create droplet"
        exit 1
    fi
}

do_droplets_delete() {
    require_doctl
    local droplet_id="$1"
    
    if [ -z "$droplet_id" ]; then
        log "ERROR" "Droplet ID or name is required"
        exit 1
    fi
    
    # Show droplet info before deletion
    log "INFO" "Droplet to be deleted:"
    doctl compute droplet get "$droplet_id"
    
    if confirm "Delete this droplet?"; then
        doctl compute droplet delete "$droplet_id" --force
        log "SUCCESS" "Droplet deleted: $droplet_id"
    fi
}

do_droplets_resize() {
    require_doctl
    local droplet_id="$1"
    local new_size="$2"
    local disk="${3:-false}"
    
    if [ -z "$droplet_id" ] || [ -z "$new_size" ]; then
        log "ERROR" "Droplet ID and new size are required"
        exit 1
    fi
    
    local resize_args=("$droplet_id" --size "$new_size")
    [ "$disk" = "true" ] && resize_args+=(--resize-disk)
    
    log "INFO" "Resizing droplet $droplet_id to $new_size"
    doctl compute droplet-action resize "${resize_args[@]}"
    log "SUCCESS" "Resize operation initiated"
}

# Kubernetes operations
do_kubernetes_list() {
    require_doctl
    
    log "INFO" "Listing Kubernetes clusters:"
    doctl kubernetes cluster list
}

do_kubernetes_create() {
    require_doctl
    local cluster_name="$1"
    local region="${2:-nyc3}"
    local node_size="${3:-s-2vcpu-2gb}"
    local node_count="${4:-3}"
    
    if [ -z "$cluster_name" ]; then
        log "ERROR" "Cluster name is required"
        exit 1
    fi
    
    log "INFO" "Creating Kubernetes cluster: $cluster_name"
    
    doctl kubernetes cluster create "$cluster_name" \
        --region "$region" \
        --node-pool "name=worker-pool;size=$node_size;count=$node_count;auto-scale=true;min-nodes=1;max-nodes=5" \
        --maintenance-window "day=sunday;start=02:00;duration=4h00m0s" \
        --auto-upgrade=true
    
    log "SUCCESS" "Kubernetes cluster creation initiated: $cluster_name"
    log "INFO" "Cluster will take several minutes to become ready"
}

do_kubernetes_kubeconfig() {
    require_doctl
    local cluster_name="$1"
    
    if [ -z "$cluster_name" ]; then
        log "ERROR" "Cluster name is required"
        exit 1
    fi
    
    log "INFO" "Downloading kubeconfig for cluster: $cluster_name"
    doctl kubernetes cluster kubeconfig save "$cluster_name"
    log "SUCCESS" "Kubeconfig saved. You can now use kubectl"
}

# Volume operations
do_volumes_list() {
    require_doctl
    
    log "INFO" "Listing block storage volumes:"
    doctl compute volume list
}

do_volumes_create() {
    require_doctl
    local name="$1"
    local size="$2"
    local region="${3:-nyc3}"
    local filesystem="${4:-ext4}"
    
    if [ -z "$name" ] || [ -z "$size" ]; then
        log "ERROR" "Volume name and size are required"
        exit 1
    fi
    
    log "INFO" "Creating volume: $name ($size GB)"
    
    doctl compute volume create "$name" \
        --size "$size" \
        --region "$region" \
        --fs-type "$filesystem"
    
    log "SUCCESS" "Volume created: $name"
}

do_volumes_attach() {
    require_doctl
    local volume_id="$1"
    local droplet_id="$2"
    
    if [ -z "$volume_id" ] || [ -z "$droplet_id" ]; then
        log "ERROR" "Volume ID and Droplet ID are required"
        exit 1
    fi
    
    log "INFO" "Attaching volume $volume_id to droplet $droplet_id"
    doctl compute volume-action attach "$volume_id" "$droplet_id"
    log "SUCCESS" "Volume attachment initiated"
}

# Spaces operations (S3-compatible object storage)
do_spaces_list() {
    require_doctl
    
    log "INFO" "Listing Spaces buckets:"
    doctl compute cdn list
}

do_spaces_create() {
    require_doctl
    local bucket_name="$1"
    local region="${2:-nyc3}"
    
    if [ -z "$bucket_name" ]; then
        log "ERROR" "Bucket name is required"
        exit 1
    fi
    
    # Note: doctl doesn't have direct Spaces bucket creation
    # This would typically be done via s3cmd or AWS CLI with DigitalOcean endpoints
    log "INFO" "Creating Spaces bucket: $bucket_name"
    log "INFO" "You can use s3cmd or AWS CLI with endpoint: https://$region.digitaloceanspaces.com"
    
    # Example s3cmd command (if available)
    if command_exists s3cmd; then
        s3cmd mb "s3://$bucket_name" --host="$region.digitaloceanspaces.com" --host-bucket="%(bucket)s.$region.digitaloceanspaces.com"
        log "SUCCESS" "Spaces bucket created: $bucket_name"
    else
        log "INFO" "Install s3cmd to create Spaces buckets automatically"
        log "INFO" "Manual command: s3cmd mb s3://$bucket_name --host=$region.digitaloceanspaces.com"
    fi
}

# Database operations
do_databases_list() {
    require_doctl
    
    log "INFO" "Listing managed databases:"
    doctl databases list
}

do_databases_create() {
    require_doctl
    local name="$1"
    local engine="${2:-mysql}"
    local size="${3:-db-s-1vcpu-1gb}"
    local region="${4:-nyc3}"
    local num_nodes="${5:-1}"
    
    if [ -z "$name" ]; then
        log "ERROR" "Database name is required"
        exit 1
    fi
    
    log "INFO" "Creating managed database: $name"
    
    doctl databases create "$name" \
        --engine "$engine" \
        --size "$size" \
        --region "$region" \
        --num-nodes "$num_nodes"
    
    log "SUCCESS" "Database creation initiated: $name"
    log "INFO" "Database will take several minutes to become ready"
}

# Load balancer operations
do_load_balancers_list() {
    require_doctl
    
    log "INFO" "Listing load balancers:"
    doctl compute load-balancer list
}

do_load_balancers_create() {
    require_doctl
    local name="$1"
    local region="${2:-nyc3}"
    local algorithm="${3:-round_robin}"
    local forwarding_rules="${4:-entry_protocol:http,entry_port:80,target_protocol:http,target_port:80}"
    
    if [ -z "$name" ]; then
        log "ERROR" "Load balancer name is required"
        exit 1
    fi
    
    log "INFO" "Creating load balancer: $name"
    
    doctl compute load-balancer create \
        --name "$name" \
        --region "$region" \
        --algorithm "$algorithm" \
        --forwarding-rules "$forwarding_rules"
    
    log "SUCCESS" "Load balancer created: $name"
}

# Monitoring and cost analysis
do_monitoring_setup() {
    require_doctl
    
    log "INFO" "DigitalOcean monitoring is automatically enabled for new resources"
    log "INFO" "You can view metrics in the DigitalOcean control panel"
    
    # List droplets with monitoring status
    log "INFO" "Droplets with monitoring:"
    doctl compute droplet list --format "Name,Monitoring"
}

do_cost_analyze() {
    require_doctl
    
    log "INFO" "Analyzing DigitalOcean usage and costs..."
    
    # Account info
    log "INFO" "Account Information:"
    doctl account get
    
    echo ""
    log "INFO" "Resource Summary:"
    
    # Count resources
    local droplet_count=$(doctl compute droplet list --format ID --no-header | wc -l)
    local volume_count=$(doctl compute volume list --format ID --no-header | wc -l)
    local k8s_count=$(doctl kubernetes cluster list --format ID --no-header | wc -l)
    local lb_count=$(doctl compute load-balancer list --format ID --no-header | wc -l)
    
    echo "  Droplets: $droplet_count"
    echo "  Volumes: $volume_count"
    echo "  Kubernetes Clusters: $k8s_count"
    echo "  Load Balancers: $lb_count"
    
    echo ""
    log "INFO" "For detailed billing information, visit:"
    log "INFO" "  https://cloud.digitalocean.com/account/billing"
}

# Main DigitalOcean module function
digitalocean_main() {
    local command="${1:-help}"
    shift || true
    
    case "$command" in
        "help"|"-h"|"--help")
            digitalocean_help
            ;;
        "auth")
            do_auth_setup
            ;;
        "droplets")
            local action="$1"
            case "$action" in
                "list") do_droplets_list ;;
                "create") do_droplets_create "$2" "$3" "$4" "$5" "$6" ;;
                "delete") do_droplets_delete "$2" ;;
                "resize") do_droplets_resize "$2" "$3" "$4" ;;
                *) log "ERROR" "Unknown droplets action: $action" ;;
            esac
            ;;
        "kubernetes"|"k8s")
            local action="$1"
            case "$action" in
                "list") do_kubernetes_list ;;
                "create") do_kubernetes_create "$2" "$3" "$4" "$5" ;;
                "kubeconfig") do_kubernetes_kubeconfig "$2" ;;
                *) log "ERROR" "Unknown kubernetes action: $action" ;;
            esac
            ;;
        "volumes")
            local action="$1"
            case "$action" in
                "list") do_volumes_list ;;
                "create") do_volumes_create "$2" "$3" "$4" "$5" ;;
                "attach") do_volumes_attach "$2" "$3" ;;
                *) log "ERROR" "Unknown volumes action: $action" ;;
            esac
            ;;
        "spaces")
            local action="$1"
            case "$action" in
                "list") do_spaces_list ;;
                "create") do_spaces_create "$2" "$3" ;;
                *) log "ERROR" "Unknown spaces action: $action" ;;
            esac
            ;;
        "databases")
            local action="$1"
            case "$action" in
                "list") do_databases_list ;;
                "create") do_databases_create "$2" "$3" "$4" "$5" "$6" ;;
                *) log "ERROR" "Unknown databases action: $action" ;;
            esac
            ;;
        "load-balancers"|"lb")
            local action="$1"
            case "$action" in
                "list") do_load_balancers_list ;;
                "create") do_load_balancers_create "$2" "$3" "$4" "$5" ;;
                *) log "ERROR" "Unknown load-balancers action: $action" ;;
            esac
            ;;
        "monitoring")
            do_monitoring_setup
            ;;
        "cost")
            do_cost_analyze
            ;;
        *)
            log "ERROR" "Unknown command: $command"
            digitalocean_help
            exit 1
            ;;
    esac
}

# Alias for shorter command
do_main() {
    digitalocean_main "$@"
}