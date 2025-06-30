#!/bin/bash
# Azure Cloud Automation Module

# Module configuration
readonly AZURE_TEMPLATES_DIR="$AUTO_HOME/cloud/azure/templates"
readonly AZURE_CONFIGS_DIR="$AUTO_HOME/cloud/azure/configs"
readonly AZURE_SCRIPTS_DIR="$AUTO_HOME/cloud/azure/scripts"

# Initialize directories
mkdir -p "$AZURE_TEMPLATES_DIR" "$AZURE_CONFIGS_DIR" "$AZURE_SCRIPTS_DIR"

# Help function for Azure module
azure_help() {
    cat << EOF
Azure Cloud Automation

USAGE:
    auto azure <command> [options]
    auto az <command> [options]           # Short alias

COMMANDS:
    auth                          Authentication and subscription management
    vm <action>                   Virtual machine management
    storage <action>              Storage account and blob operations
    functions <action>            Azure Functions management
    aks <action>                  Azure Kubernetes Service operations
    sql <action>                  Azure SQL database operations
    cosmosdb <action>             Cosmos DB operations
    appservice <action>           App Service and web app operations
    keyvault <action>             Key Vault management
    monitor                       Monitoring and diagnostics
    cost                          Cost analysis and optimization
    resource-groups <action>      Resource group management

EXAMPLES:
    auto az auth setup
    auto az vm create web-server Standard_B1s eastus
    auto az storage create mystorageaccount eastus
    auto az functions deploy my-function-app
    auto az aks create my-cluster eastus
    auto az sql create my-database
    auto az appservice deploy my-web-app
    auto az cost analyze --resource-group my-rg
EOF
}

# Utility functions
require_az_cli() {
    if ! command_exists az; then
        log "ERROR" "Azure CLI not found. Please install it first:"
        log "INFO" "  macOS: brew install azure-cli"
        log "INFO" "  Linux: curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash"
        log "INFO" "  Manual: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli"
        exit 1
    fi
}

get_azure_subscription() {
    az account show --query id --output tsv 2>/dev/null || echo ""
}

get_azure_location() {
    az account list-locations --query '[0].name' --output tsv 2>/dev/null || echo "eastus"
}

# Authentication and subscription management
azure_auth_setup() {
    require_az_cli
    
    log "INFO" "Setting up Azure authentication..."
    
    # Check if already logged in
    if az account show >/dev/null 2>&1; then
        log "INFO" "Already authenticated with Azure"
        az account show --output table
        return 0
    fi
    
    # Interactive login
    log "INFO" "Starting Azure authentication..."
    az login
    
    # Show available subscriptions
    log "INFO" "Available subscriptions:"
    az account list --output table
    
    # Set default subscription if multiple exist
    local sub_count=$(az account list --query 'length(@)' --output tsv)
    if [ "$sub_count" -gt 1 ]; then
        echo ""
        echo -n "Enter subscription ID to set as default (or press Enter for current): "
        read -r subscription_id
        
        if [ -n "$subscription_id" ]; then
            az account set --subscription "$subscription_id"
            log "SUCCESS" "Default subscription set to: $subscription_id"
        fi
    fi
    
    log "SUCCESS" "Azure authentication completed"
}

azure_auth_list_subscriptions() {
    require_az_cli
    
    log "INFO" "Azure subscriptions:"
    az account list --output table
}

# Resource group operations
azure_resource_groups_list() {
    require_az_cli
    
    log "INFO" "Listing resource groups:"
    az group list --output table
}

azure_resource_groups_create() {
    require_az_cli
    local rg_name="$1"
    local location="${2:-$(get_azure_location)}"
    
    if [ -z "$rg_name" ]; then
        log "ERROR" "Resource group name is required"
        exit 1
    fi
    
    log "INFO" "Creating resource group: $rg_name"
    az group create --name "$rg_name" --location "$location"
    log "SUCCESS" "Resource group created: $rg_name"
}

# Virtual machine operations
azure_vm_list() {
    require_az_cli
    local resource_group="$1"
    
    if [ -n "$resource_group" ]; then
        log "INFO" "Listing VMs in resource group: $resource_group"
        az vm list --resource-group "$resource_group" --output table
    else
        log "INFO" "Listing all VMs:"
        az vm list --output table
    fi
}

azure_vm_create() {
    require_az_cli
    local vm_name="$1"
    local vm_size="${2:-Standard_B1s}"
    local location="${3:-$(get_azure_location)}"
    local resource_group="${4:-$vm_name-rg}"
    local image="${5:-UbuntuLTS}"
    local admin_username="${6:-azureuser}"
    
    if [ -z "$vm_name" ]; then
        log "ERROR" "VM name is required"
        exit 1
    fi
    
    log "INFO" "Creating VM: $vm_name"
    
    # Create resource group if it doesn't exist
    if ! az group show --name "$resource_group" >/dev/null 2>&1; then
        log "INFO" "Creating resource group: $resource_group"
        az group create --name "$resource_group" --location "$location"
    fi
    
    # Create VM
    az vm create \
        --resource-group "$resource_group" \
        --name "$vm_name" \
        --image "$image" \
        --size "$vm_size" \
        --location "$location" \
        --admin-username "$admin_username" \
        --generate-ssh-keys \
        --output table
    
    log "SUCCESS" "VM created: $vm_name"
    
    # Get public IP
    local public_ip=$(az vm show -d --resource-group "$resource_group" --name "$vm_name" --query publicIps --output tsv)
    log "INFO" "Public IP: $public_ip"
    log "INFO" "SSH command: ssh $admin_username@$public_ip"
}

azure_vm_start() {
    require_az_cli
    local vm_name="$1"
    local resource_group="$2"
    
    if [ -z "$vm_name" ] || [ -z "$resource_group" ]; then
        log "ERROR" "VM name and resource group are required"
        exit 1
    fi
    
    log "INFO" "Starting VM: $vm_name"
    az vm start --resource-group "$resource_group" --name "$vm_name"
    log "SUCCESS" "VM started: $vm_name"
}

azure_vm_stop() {
    require_az_cli
    local vm_name="$1"
    local resource_group="$2"
    
    if [ -z "$vm_name" ] || [ -z "$resource_group" ]; then
        log "ERROR" "VM name and resource group are required"
        exit 1
    fi
    
    log "INFO" "Stopping VM: $vm_name"
    az vm deallocate --resource-group "$resource_group" --name "$vm_name"
    log "SUCCESS" "VM stopped: $vm_name"
}

# Storage operations
azure_storage_list() {
    require_az_cli
    
    log "INFO" "Listing storage accounts:"
    az storage account list --output table
}

azure_storage_create() {
    require_az_cli
    local account_name="$1"
    local location="${2:-$(get_azure_location)}"
    local resource_group="${3:-$account_name-rg}"
    local sku="${4:-Standard_LRS}"
    
    if [ -z "$account_name" ]; then
        log "ERROR" "Storage account name is required"
        exit 1
    fi
    
    log "INFO" "Creating storage account: $account_name"
    
    # Create resource group if it doesn't exist
    if ! az group show --name "$resource_group" >/dev/null 2>&1; then
        log "INFO" "Creating resource group: $resource_group"
        az group create --name "$resource_group" --location "$location"
    fi
    
    # Create storage account
    az storage account create \
        --name "$account_name" \
        --resource-group "$resource_group" \
        --location "$location" \
        --sku "$sku" \
        --encryption-services blob file
    
    log "SUCCESS" "Storage account created: $account_name"
}

azure_storage_blob_upload() {
    require_az_cli
    local file_path="$1"
    local container_name="$2"
    local account_name="$3"
    local blob_name="${4:-$(basename "$file_path")}"
    
    if [ -z "$file_path" ] || [ -z "$container_name" ] || [ -z "$account_name" ]; then
        log "ERROR" "File path, container name, and account name are required"
        exit 1
    fi
    
    if [ ! -f "$file_path" ]; then
        log "ERROR" "File not found: $file_path"
        exit 1
    fi
    
    log "INFO" "Uploading $file_path to blob storage"
    
    # Create container if it doesn't exist
    az storage container create \
        --name "$container_name" \
        --account-name "$account_name" \
        --auth-mode login >/dev/null 2>&1 || true
    
    # Upload blob
    az storage blob upload \
        --file "$file_path" \
        --container-name "$container_name" \
        --name "$blob_name" \
        --account-name "$account_name" \
        --auth-mode login
    
    log "SUCCESS" "File uploaded as blob: $blob_name"
}

# Azure Functions operations
azure_functions_list() {
    require_az_cli
    local resource_group="$1"
    
    if [ -n "$resource_group" ]; then
        log "INFO" "Listing function apps in resource group: $resource_group"
        az functionapp list --resource-group "$resource_group" --output table
    else
        log "INFO" "Listing all function apps:"
        az functionapp list --output table
    fi
}

azure_functions_create() {
    require_az_cli
    local app_name="$1"
    local resource_group="${2:-$app_name-rg}"
    local location="${3:-$(get_azure_location)}"
    local runtime="${4:-python}"
    local storage_account="${5:-${app_name}storage}"
    
    if [ -z "$app_name" ]; then
        log "ERROR" "Function app name is required"
        exit 1
    fi
    
    log "INFO" "Creating function app: $app_name"
    
    # Create resource group if it doesn't exist
    if ! az group show --name "$resource_group" >/dev/null 2>&1; then
        az group create --name "$resource_group" --location "$location"
    fi
    
    # Create storage account for function app
    az storage account create \
        --name "$storage_account" \
        --resource-group "$resource_group" \
        --location "$location" \
        --sku Standard_LRS >/dev/null 2>&1 || true
    
    # Create function app
    az functionapp create \
        --resource-group "$resource_group" \
        --name "$app_name" \
        --storage-account "$storage_account" \
        --runtime "$runtime" \
        --consumption-plan-location "$location" \
        --functions-version 4
    
    log "SUCCESS" "Function app created: $app_name"
    log "INFO" "URL: https://$app_name.azurewebsites.net"
}

azure_functions_deploy() {
    require_az_cli
    local app_name="$1"
    local source_path="${2:-.}"
    
    if [ -z "$app_name" ]; then
        log "ERROR" "Function app name is required"
        exit 1
    fi
    
    if [ ! -d "$source_path" ]; then
        log "ERROR" "Source path not found: $source_path"
        exit 1
    fi
    
    log "INFO" "Deploying function app: $app_name"
    
    # Deploy using zip deployment
    if [ -f "$source_path/requirements.txt" ] || [ -f "$source_path/package.json" ]; then
        az functionapp deployment source config-zip \
            --resource-group "$(az functionapp show --name "$app_name" --query resourceGroup --output tsv)" \
            --name "$app_name" \
            --src "$source_path"
    else
        log "WARN" "No requirements.txt or package.json found. Deploying as-is."
        az functionapp deployment source config-zip \
            --resource-group "$(az functionapp show --name "$app_name" --query resourceGroup --output tsv)" \
            --name "$app_name" \
            --src "$source_path"
    fi
    
    log "SUCCESS" "Function app deployed: $app_name"
}

# Azure Kubernetes Service operations
azure_aks_list() {
    require_az_cli
    
    log "INFO" "Listing AKS clusters:"
    az aks list --output table
}

azure_aks_create() {
    require_az_cli
    local cluster_name="$1"
    local resource_group="${2:-$cluster_name-rg}"
    local location="${3:-$(get_azure_location)}"
    local node_count="${4:-3}"
    local node_size="${5:-Standard_DS2_v2}"
    
    if [ -z "$cluster_name" ]; then
        log "ERROR" "Cluster name is required"
        exit 1
    fi
    
    log "INFO" "Creating AKS cluster: $cluster_name"
    
    # Create resource group if it doesn't exist
    if ! az group show --name "$resource_group" >/dev/null 2>&1; then
        az group create --name "$resource_group" --location "$location"
    fi
    
    # Create AKS cluster
    az aks create \
        --resource-group "$resource_group" \
        --name "$cluster_name" \
        --location "$location" \
        --node-count "$node_count" \
        --node-vm-size "$node_size" \
        --enable-addons monitoring \
        --generate-ssh-keys \
        --enable-managed-identity
    
    log "SUCCESS" "AKS cluster created: $cluster_name"
}

azure_aks_get_credentials() {
    require_az_cli
    local cluster_name="$1"
    local resource_group="$2"
    
    if [ -z "$cluster_name" ] || [ -z "$resource_group" ]; then
        log "ERROR" "Cluster name and resource group are required"
        exit 1
    fi
    
    log "INFO" "Getting AKS credentials for: $cluster_name"
    az aks get-credentials --resource-group "$resource_group" --name "$cluster_name"
    log "SUCCESS" "Credentials configured. You can now use kubectl"
}

# Azure SQL operations
azure_sql_list() {
    require_az_cli
    
    log "INFO" "Listing SQL servers:"
    az sql server list --output table
}

azure_sql_create() {
    require_az_cli
    local server_name="$1"
    local database_name="${2:-$server_name-db}"
    local resource_group="${3:-$server_name-rg}"
    local location="${4:-$(get_azure_location)}"
    local admin_user="${5:-sqladmin}"
    local admin_password="$6"
    
    if [ -z "$server_name" ]; then
        log "ERROR" "SQL server name is required"
        exit 1
    fi
    
    if [ -z "$admin_password" ]; then
        echo -n "Enter SQL admin password: "
        read -rs admin_password
        echo ""
    fi
    
    log "INFO" "Creating SQL server: $server_name"
    
    # Create resource group if it doesn't exist
    if ! az group show --name "$resource_group" >/dev/null 2>&1; then
        az group create --name "$resource_group" --location "$location"
    fi
    
    # Create SQL server
    az sql server create \
        --name "$server_name" \
        --resource-group "$resource_group" \
        --location "$location" \
        --admin-user "$admin_user" \
        --admin-password "$admin_password"
    
    # Create database
    az sql db create \
        --resource-group "$resource_group" \
        --server "$server_name" \
        --name "$database_name" \
        --service-objective Basic
    
    log "SUCCESS" "SQL server and database created"
    log "INFO" "Server: $server_name.database.windows.net"
    log "INFO" "Database: $database_name"
    log "INFO" "Admin user: $admin_user"
}

# Cost analysis
azure_cost_analyze() {
    require_az_cli
    local resource_group="$1"
    local time_period="${2:-Month}"
    
    log "INFO" "Analyzing Azure costs..."
    
    if [ -n "$resource_group" ]; then
        log "INFO" "Cost analysis for resource group: $resource_group"
        # Note: Azure CLI cost management commands are limited
        # Most cost analysis is done through the portal
        az group show --name "$resource_group" --output table
        az resource list --resource-group "$resource_group" --output table
    else
        log "INFO" "Overall resource usage:"
        az resource list --output table | head -20
    fi
    
    log "INFO" "For detailed cost analysis, visit:"
    log "INFO" "  https://portal.azure.com/#blade/Microsoft_Azure_CostManagement/Menu/overview"
}

# Main Azure module function
azure_main() {
    local command="${1:-help}"
    shift || true
    
    case "$command" in
        "help"|"-h"|"--help")
            azure_help
            ;;
        "auth")
            local action="${1:-setup}"
            case "$action" in
                "setup") azure_auth_setup ;;
                "list-subscriptions") azure_auth_list_subscriptions ;;
                *) log "ERROR" "Unknown auth action: $action" ;;
            esac
            ;;
        "resource-groups"|"rg")
            local action="$1"
            case "$action" in
                "list") azure_resource_groups_list ;;
                "create") azure_resource_groups_create "$2" "$3" ;;
                *) log "ERROR" "Unknown resource-groups action: $action" ;;
            esac
            ;;
        "vm")
            local action="$1"
            case "$action" in
                "list") azure_vm_list "$2" ;;
                "create") azure_vm_create "$2" "$3" "$4" "$5" "$6" "$7" ;;
                "start") azure_vm_start "$2" "$3" ;;
                "stop") azure_vm_stop "$2" "$3" ;;
                *) log "ERROR" "Unknown vm action: $action" ;;
            esac
            ;;
        "storage")
            local action="$1"
            case "$action" in
                "list") azure_storage_list ;;
                "create") azure_storage_create "$2" "$3" "$4" "$5" ;;
                "upload") azure_storage_blob_upload "$2" "$3" "$4" "$5" ;;
                *) log "ERROR" "Unknown storage action: $action" ;;
            esac
            ;;
        "functions")
            local action="$1"
            case "$action" in
                "list") azure_functions_list "$2" ;;
                "create") azure_functions_create "$2" "$3" "$4" "$5" "$6" ;;
                "deploy") azure_functions_deploy "$2" "$3" ;;
                *) log "ERROR" "Unknown functions action: $action" ;;
            esac
            ;;
        "aks")
            local action="$1"
            case "$action" in
                "list") azure_aks_list ;;
                "create") azure_aks_create "$2" "$3" "$4" "$5" "$6" ;;
                "get-credentials") azure_aks_get_credentials "$2" "$3" ;;
                *) log "ERROR" "Unknown aks action: $action" ;;
            esac
            ;;
        "sql")
            local action="$1"
            case "$action" in
                "list") azure_sql_list ;;
                "create") azure_sql_create "$2" "$3" "$4" "$5" "$6" "$7" ;;
                *) log "ERROR" "Unknown sql action: $action" ;;
            esac
            ;;
        "cost")
            azure_cost_analyze "$1" "$2"
            ;;
        *)
            log "ERROR" "Unknown command: $command"
            azure_help
            exit 1
            ;;
    esac
}

# Alias for shorter command
az_main() {
    azure_main "$@"
}