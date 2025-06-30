#!/bin/bash
# System Administration Automation Module

# Module configuration
readonly SYSTEM_BACKUPS_DIR="$AUTO_CACHE/system-backups"
readonly SYSTEM_CONFIGS_DIR="$AUTO_CONFIG/system"
readonly SYSTEM_SCRIPTS_DIR="$AUTO_HOME/system/scripts"

# Initialize directories
mkdir -p "$SYSTEM_BACKUPS_DIR" "$SYSTEM_CONFIGS_DIR" "$SYSTEM_SCRIPTS_DIR"

# Help function for system module
system_help() {
    cat << EOF
System Administration Automation

USAGE:
    auto system <command> [options]

COMMANDS:
    setup                        Initial system setup and configuration
    backup <path> [destination]  Backup files or directories
    restore <backup-file>        Restore from backup
    monitor                      System monitoring and health checks
    cleanup                      System cleanup and maintenance
    security                     Security hardening and checks
    network                      Network configuration and diagnostics
    services                     Service management
    users                        User account management
    packages                     Package management automation
    update                       Update system and automation tools

EXAMPLES:
    auto system setup
    auto system backup /important/data
    auto system monitor --cpu --memory --disk
    auto system cleanup --temp --logs --cache
    auto system security --scan --harden
    auto system packages --update --install docker
EOF
}

# System information gathering
get_system_info() {
    log "INFO" "System Information"
    echo "==================="
    
    # OS Information
    echo "OS: $(get_os_info)"
    echo "Hostname: $(hostname)"
    echo "Uptime: $(uptime | awk '{print $3,$4}' | sed 's/,//')"
    
    # Hardware information
    echo "CPU Cores: $(get_cpu_count)"
    
    if [[ "$(get_os_info)" == "macos" ]]; then
        echo "Memory: $(sysctl -n hw.memsize | awk '{print $1/1024/1024/1024 " GB"}')"
        echo "Disk Usage:"
        df -h | grep -E '^/dev/'
    else
        echo "Memory: $(free -h | awk 'NR==2{print $2}')"
        echo "Disk Usage:"
        df -h | grep -E '^/dev/'
    fi
    
    echo ""
}

# System setup and configuration
system_setup() {
    log "INFO" "Starting system setup and configuration..."
    
    # Update package manager
    update_packages
    
    # Install essential tools
    install_essential_tools
    
    # Configure shell environment
    configure_shell_environment
    
    # Setup development tools
    setup_development_environment
    
    # Configure security settings
    configure_security_basics
    
    log "SUCCESS" "System setup completed"
}

update_packages() {
    log "INFO" "Updating package manager..."
    
    case "$(get_os_info)" in
        "macos")
            if command_exists brew; then
                brew update && brew upgrade
            else
                log "WARN" "Homebrew not installed. Installing..."
                /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
            fi
            ;;
        "linux")
            if command_exists apt-get; then
                sudo apt-get update && sudo apt-get upgrade -y
            elif command_exists yum; then
                sudo yum update -y
            elif command_exists dnf; then
                sudo dnf update -y
            fi
            ;;
    esac
}

install_essential_tools() {
    log "INFO" "Installing essential tools..."
    
    local tools=("curl" "wget" "git" "jq" "htop" "tree" "zip" "unzip")
    
    case "$(get_os_info)" in
        "macos")
            for tool in "${tools[@]}"; do
                if ! command_exists "$tool"; then
                    brew install "$tool"
                fi
            done
            ;;
        "linux")
            if command_exists apt-get; then
                sudo apt-get install -y "${tools[@]}"
            elif command_exists yum; then
                sudo yum install -y "${tools[@]}"
            fi
            ;;
    esac
}

# Backup operations
backup_files() {
    local source_path="$1"
    local destination="${2:-$SYSTEM_BACKUPS_DIR}"
    local backup_name="backup_$(date +%Y%m%d_%H%M%S)"
    
    if [ ! -e "$source_path" ]; then
        log "ERROR" "Source path does not exist: $source_path"
        exit 1
    fi
    
    mkdir -p "$destination"
    
    log "INFO" "Creating backup of $source_path..."
    
    if [ -d "$source_path" ]; then
        tar -czf "$destination/$backup_name.tar.gz" -C "$(dirname "$source_path")" "$(basename "$source_path")"
    else
        cp "$source_path" "$destination/$backup_name.$(basename "$source_path")"
    fi
    
    log "SUCCESS" "Backup created: $destination/$backup_name"
    echo "$destination/$backup_name"
}

backup_system_configs() {
    log "INFO" "Backing up system configurations..."
    
    local config_backup_dir="$SYSTEM_BACKUPS_DIR/configs_$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$config_backup_dir"
    
    # Common configuration files to backup
    local config_files=(
        "$HOME/.bash_profile"
        "$HOME/.bashrc" 
        "$HOME/.zshrc"
        "$HOME/.gitconfig"
        "$HOME/.ssh/config"
        "/etc/hosts"
    )
    
    for config_file in "${config_files[@]}"; do
        if [ -f "$config_file" ]; then
            cp "$config_file" "$config_backup_dir/" 2>/dev/null || sudo cp "$config_file" "$config_backup_dir/" 2>/dev/null
            log "DEBUG" "Backed up: $config_file"
        fi
    done
    
    # Create archive
    tar -czf "$config_backup_dir.tar.gz" -C "$(dirname "$config_backup_dir")" "$(basename "$config_backup_dir")"
    rm -rf "$config_backup_dir"
    
    log "SUCCESS" "System configurations backed up to: $config_backup_dir.tar.gz"
}

# System monitoring
monitor_system() {
    local duration="${1:-60}"
    local interval="${2:-5}"
    
    log "INFO" "Monitoring system for $duration seconds (interval: ${interval}s)"
    
    local end_time=$(($(date +%s) + duration))
    
    while [ $(date +%s) -lt $end_time ]; do
        clear
        echo "System Monitor - $(date)"
        echo "================================="
        
        # CPU Usage
        if command_exists top; then
            echo "CPU Usage:"
            if [[ "$(get_os_info)" == "macos" ]]; then
                top -l 1 | grep "CPU usage" | awk '{print $3, $5}' | sed 's/%//g'
            else
                top -bn1 | grep "Cpu(s)" | awk '{print $2}' | sed 's/%us,//g'
            fi
        fi
        
        # Memory Usage
        echo ""
        echo "Memory Usage:"
        if [[ "$(get_os_info)" == "macos" ]]; then
            vm_stat | awk '/free/ {free=$3} /active/ {active=$3} END {print "Free: " free*4096/1024/1024 " MB, Active: " active*4096/1024/1024 " MB"}'
        else
            free -h | awk 'NR==2{printf "Used: %s/%s (%.2f%%)\n", $3,$2,$3*100/$2 }'
        fi
        
        # Disk Usage
        echo ""
        echo "Disk Usage:"
        df -h | awk 'NR>1 && /^\/dev/ {printf "%-20s %s/%s (%s)\n", $1, $3, $2, $5}'
        
        # Network
        echo ""
        echo "Network Connections:"
        netstat -an 2>/dev/null | grep ESTABLISHED | wc -l | awk '{print "Active connections: " $1}'
        
        sleep "$interval"
    done
}

# System cleanup
cleanup_system() {
    log "INFO" "Starting system cleanup..."
    
    # Clean temporary files
    cleanup_temp_files
    
    # Clean log files
    cleanup_old_logs
    
    # Clean package caches
    cleanup_package_cache
    
    # Clean browser caches (optional)
    if confirm "Clean browser caches?"; then
        cleanup_browser_cache
    fi
    
    log "SUCCESS" "System cleanup completed"
}

cleanup_temp_files() {
    log "INFO" "Cleaning temporary files..."
    
    case "$(get_os_info)" in
        "macos")
            # Clean user temp directory
            rm -rf "$HOME"/Library/Caches/Temporary\ Items/* 2>/dev/null
            # Clean system temp
            sudo rm -rf /tmp/* 2>/dev/null
            ;;
        "linux")
            # Clean temp directories
            sudo rm -rf /tmp/* 2>/dev/null
            sudo rm -rf /var/tmp/* 2>/dev/null
            ;;
    esac
    
    log "DEBUG" "Temporary files cleaned"
}

cleanup_old_logs() {
    log "INFO" "Cleaning old log files..."
    
    # Clean logs older than 30 days
    find /var/log -name "*.log" -mtime +30 -delete 2>/dev/null || true
    
    # Clean automation logs older than 7 days
    find "$AUTO_LOGS" -name "*.log" -mtime +7 -delete 2>/dev/null || true
    
    log "DEBUG" "Old log files cleaned"
}

cleanup_package_cache() {
    log "INFO" "Cleaning package caches..."
    
    case "$(get_os_info)" in
        "macos")
            if command_exists brew; then
                brew cleanup
            fi
            ;;
        "linux")
            if command_exists apt-get; then
                sudo apt-get clean
                sudo apt-get autoremove -y
            elif command_exists yum; then
                sudo yum clean all
            fi
            ;;
    esac
    
    log "DEBUG" "Package caches cleaned"
}

# Security operations
security_scan() {
    log "INFO" "Running security scan..."
    
    # Check for security updates
    check_security_updates
    
    # Scan for suspicious processes
    scan_processes
    
    # Check file permissions
    check_file_permissions
    
    # Check network connections
    check_network_connections
    
    log "SUCCESS" "Security scan completed"
}

check_security_updates() {
    log "INFO" "Checking for security updates..."
    
    case "$(get_os_info)" in
        "macos")
            softwareupdate -l 2>/dev/null | grep -i security || log "INFO" "No security updates available"
            ;;
        "linux")
            if command_exists apt-get; then
                apt list --upgradable 2>/dev/null | grep -i security || log "INFO" "No security updates available"
            fi
            ;;
    esac
}

# Service management
manage_services() {
    local action="$1"
    local service="$2"
    
    case "$action" in
        "list")
            case "$(get_os_info)" in
                "macos")
                    launchctl list | grep -v "com.apple"
                    ;;
                "linux")
                    if command_exists systemctl; then
                        systemctl list-units --type=service --state=running
                    else
                        service --status-all
                    fi
                    ;;
            esac
            ;;
        "start"|"stop"|"restart"|"status")
            if [ -z "$service" ]; then
                log "ERROR" "Service name required"
                exit 1
            fi
            
            case "$(get_os_info)" in
                "macos")
                    launchctl "$action" "$service"
                    ;;
                "linux")
                    if command_exists systemctl; then
                        sudo systemctl "$action" "$service"
                    else
                        sudo service "$service" "$action"
                    fi
                    ;;
            esac
            ;;
    esac
}

# Main system module function
system_main() {
    local command="${1:-help}"
    shift || true
    
    case "$command" in
        "help"|"-h"|"--help")
            system_help
            ;;
        "setup")
            system_setup
            ;;
        "info")
            get_system_info
            ;;
        "backup")
            if [ "$1" = "--configs" ]; then
                backup_system_configs
            else
                backup_files "$1" "$2"
            fi
            ;;
        "monitor")
            monitor_system "$1" "$2"
            ;;
        "cleanup")
            cleanup_system
            ;;
        "security")
            local action="${1:-scan}"
            case "$action" in
                "scan") security_scan ;;
                *) log "ERROR" "Unknown security action: $action" ;;
            esac
            ;;
        "services")
            manage_services "$1" "$2"
            ;;
        "update")
            update_packages
            ;;
        "packages")
            local action="${1:-list}"
            case "$action" in
                "update") update_packages ;;
                "install") 
                    shift
                    for package in "$@"; do
                        case "$(get_os_info)" in
                            "macos") brew install "$package" ;;
                            "linux") 
                                if command_exists apt-get; then
                                    sudo apt-get install -y "$package"
                                elif command_exists yum; then
                                    sudo yum install -y "$package"
                                fi
                                ;;
                        esac
                    done
                    ;;
                *) log "ERROR" "Unknown packages action: $action" ;;
            esac
            ;;
        *)
            log "ERROR" "Unknown command: $command"
            system_help
            exit 1
            ;;
    esac
}