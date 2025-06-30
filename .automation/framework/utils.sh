#!/usr/local/bin/bash
# Automation Framework Utilities
# Extended utility functions for automation scripts

# File and directory operations
backup_file() {
    local file="$1"
    local backup_dir="${2:-$AUTO_CACHE/backups}"
    
    if [ -f "$file" ]; then
        mkdir -p "$backup_dir"
        local timestamp=$(date +%Y%m%d_%H%M%S)
        local backup_name="$(basename "$file").backup.$timestamp"
        cp "$file" "$backup_dir/$backup_name"
        log "INFO" "Backed up $file to $backup_dir/$backup_name"
        echo "$backup_dir/$backup_name"
    else
        log "WARN" "File $file does not exist, skipping backup"
        return 1
    fi
}

# Template processing
process_template() {
    local template_file="$1"
    local output_file="$2"
    shift 2
    
    if [ ! -f "$template_file" ]; then
        log "ERROR" "Template file $template_file not found"
        return 1
    fi
    
    local content=$(cat "$template_file")
    
    # Replace variables passed as VAR=value
    for var_assignment in "$@"; do
        local var_name="${var_assignment%%=*}"
        local var_value="${var_assignment#*=}"
        content=$(echo "$content" | sed "s/{{$var_name}}/$var_value/g")
    done
    
    echo "$content" > "$output_file"
    log "INFO" "Processed template $template_file -> $output_file"
}

# Network utilities
download_file() {
    local url="$1"
    local output_file="$2"
    local max_retries="${3:-3}"
    
    for i in $(seq 1 $max_retries); do
        if curl -fsSL -o "$output_file" "$url"; then
            log "SUCCESS" "Downloaded $url to $output_file"
            return 0
        else
            log "WARN" "Download attempt $i failed for $url"
            sleep 2
        fi
    done
    
    log "ERROR" "Failed to download $url after $max_retries attempts"
    return 1
}

# JSON processing
get_json_value() {
    local json_file="$1"
    local json_path="$2"
    
    if [ ! -f "$json_file" ]; then
        log "ERROR" "JSON file $json_file not found"
        return 1
    fi
    
    jq -r "$json_path" "$json_file"
}

set_json_value() {
    local json_file="$1"
    local json_path="$2"
    local value="$3"
    
    local tmp_file=$(mktemp)
    jq "$json_path = \"$value\"" "$json_file" > "$tmp_file"
    mv "$tmp_file" "$json_file"
    log "DEBUG" "Set $json_path = $value in $json_file"
}

# Process management
run_with_timeout() {
    local timeout="$1"
    shift
    local command="$@"
    
    timeout "$timeout" bash -c "$command"
}

run_parallel() {
    local max_jobs="$1"
    shift
    local commands=("$@")
    
    local job_count=0
    for cmd in "${commands[@]}"; do
        if [ $job_count -ge $max_jobs ]; then
            wait -n  # Wait for any job to finish
            ((job_count--))
        fi
        
        eval "$cmd" &
        ((job_count++))
    done
    
    wait  # Wait for all remaining jobs
}

# System information
get_os_info() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macos"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "linux"
    else
        echo "unknown"
    fi
}

get_cpu_count() {
    if command_exists nproc; then
        nproc
    elif command_exists sysctl; then
        sysctl -n hw.ncpu
    else
        echo "1"
    fi
}

# Lock management for preventing concurrent runs
acquire_lock() {
    local lock_name="$1"
    local lock_file="$AUTO_CACHE/$lock_name.lock"
    local timeout="${2:-30}"
    
    local count=0
    while [ $count -lt $timeout ]; do
        if mkdir "$lock_file" 2>/dev/null; then
            echo $$ > "$lock_file/pid"
            log "DEBUG" "Acquired lock: $lock_name"
            return 0
        fi
        
        # Check if existing lock is stale
        if [ -f "$lock_file/pid" ]; then
            local existing_pid=$(cat "$lock_file/pid")
            if ! kill -0 "$existing_pid" 2>/dev/null; then
                log "WARN" "Removing stale lock: $lock_name"
                rm -rf "$lock_file"
                continue
            fi
        fi
        
        sleep 1
        ((count++))
    done
    
    log "ERROR" "Failed to acquire lock: $lock_name"
    return 1
}

release_lock() {
    local lock_name="$1"
    local lock_file="$AUTO_CACHE/$lock_name.lock"
    
    if [ -d "$lock_file" ]; then
        rm -rf "$lock_file"
        log "DEBUG" "Released lock: $lock_name"
    fi
}

# Cleanup function to release locks on exit
cleanup_locks() {
    for lock_file in "$AUTO_CACHE"/*.lock; do
        if [ -d "$lock_file" ] && [ -f "$lock_file/pid" ]; then
            local lock_pid=$(cat "$lock_file/pid")
            if [ "$lock_pid" = "$$" ]; then
                rm -rf "$lock_file"
                log "DEBUG" "Cleaned up lock: $(basename "$lock_file" .lock)"
            fi
        fi
    done
}

# Set cleanup trap
trap cleanup_locks EXIT

# Validation functions
validate_email() {
    local email="$1"
    [[ "$email" =~ ^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$ ]]
}

validate_url() {
    local url="$1"
    [[ "$url" =~ ^https?://[A-Za-z0-9.-]+.*$ ]]
}

validate_k8s_name() {
    local name="$1"
    [[ "$name" =~ ^[a-z0-9]([-a-z0-9]*[a-z0-9])?$ ]]
}

# Performance monitoring
start_timer() {
    export TIMER_START=$(date +%s)
}

end_timer() {
    local start_time="${TIMER_START:-$(date +%s)}"
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    local hours=$((duration / 3600))
    local minutes=$(((duration % 3600) / 60))
    local seconds=$((duration % 60))
    
    if [ $hours -gt 0 ]; then
        echo "${hours}h ${minutes}m ${seconds}s"
    elif [ $minutes -gt 0 ]; then
        echo "${minutes}m ${seconds}s"
    else
        echo "${seconds}s"
    fi
}