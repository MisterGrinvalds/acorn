#!/bin/bash
# Automation Framework Core
# Provides common utilities for all automation scripts

set -eo pipefail

# Framework configuration
export AUTO_HOME="${DOTFILES}/.automation"
export AUTO_CONFIG="${AUTO_HOME}/config"
export AUTO_LOGS="${AUTO_HOME}/logs"
export AUTO_CACHE="${AUTO_HOME}/cache"

# Ensure required directories exist
mkdir -p "$AUTO_CONFIG" "$AUTO_LOGS" "$AUTO_CACHE"

# Color codes for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly PURPLE='\033[0;35m'
readonly CYAN='\033[0;36m'
readonly WHITE='\033[1;37m'
readonly NC='\033[0m' # No Color

# Logging functions
log() {
    local level="$1"
    shift
    local message="$*"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    case "$level" in
        "INFO")  echo -e "${GREEN}[INFO]${NC} $message" ;;
        "WARN")  echo -e "${YELLOW}[WARN]${NC} $message" ;;
        "ERROR") echo -e "${RED}[ERROR]${NC} $message" ;;
        "DEBUG") echo -e "${BLUE}[DEBUG]${NC} $message" ;;
        "SUCCESS") echo -e "${GREEN}[SUCCESS]${NC} $message" ;;
    esac
    
    # Also log to file
    echo "[$timestamp] [$level] $message" >> "$AUTO_LOGS/automation.log"
}

# Utility functions
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

require_command() {
    if ! command_exists "$1"; then
        log "ERROR" "Required command '$1' not found. Please install it first."
        exit 1
    fi
}

confirm() {
    local message="$1"
    local default="${2:-n}"
    
    if [ "$default" = "y" ]; then
        prompt="[Y/n]"
    else
        prompt="[y/N]"
    fi
    
    echo -n -e "${YELLOW}$message $prompt:${NC} "
    read -r response
    
    case "$response" in
        [yY]|[yY][eE][sS]) return 0 ;;
        [nN]|[nN][oO]) return 1 ;;
        "") [ "$default" = "y" ] && return 0 || return 1 ;;
        *) return 1 ;;
    esac
}

# Configuration management
load_config() {
    local config_file="$AUTO_CONFIG/$1.conf"
    if [ -f "$config_file" ]; then
        source "$config_file"
        log "DEBUG" "Loaded configuration: $config_file"
    else
        log "DEBUG" "Configuration file not found: $config_file"
    fi
}

save_config() {
    local config_name="$1"
    local config_file="$AUTO_CONFIG/$config_name.conf"
    shift
    
    # Save key=value pairs to config file
    for var in "$@"; do
        echo "$var" >> "$config_file"
    done
    log "DEBUG" "Saved configuration: $config_file"
}

# Progress indicators
show_progress() {
    local current="$1"
    local total="$2"
    local task="$3"
    
    local percent=$((current * 100 / total))
    local filled=$((percent / 2))
    local empty=$((50 - filled))
    
    printf "\r${CYAN}Progress:${NC} ["
    printf "%${filled}s" | tr ' ' '='
    printf "%${empty}s" | tr ' ' '-'
    printf "] %d%% - %s" "$percent" "$task"
}

# Error handling
handle_error() {
    local exit_code=$?
    local line_number=$1
    log "ERROR" "Script failed at line $line_number with exit code $exit_code"
    exit $exit_code
}

# Set error trap
trap 'handle_error $LINENO' ERR

# Check for required tools on framework load
check_framework_requirements() {
    local missing_tools=()
    
    # Core tools
    for tool in curl jq git; do
        if ! command_exists "$tool"; then
            missing_tools+=("$tool")
        fi
    done
    
    if [ ${#missing_tools[@]} -gt 0 ]; then
        log "WARN" "Missing recommended tools: ${missing_tools[*]}"
        log "INFO" "Some automation features may not work without these tools"
    fi
}

# Initialize framework
init_framework() {
    log "INFO" "Initializing automation framework..."
    check_framework_requirements
    
    # Create default configuration if it doesn't exist
    if [ ! -f "$AUTO_CONFIG/automation.conf" ]; then
        cat > "$AUTO_CONFIG/automation.conf" << EOF
# Automation Framework Configuration
AUTO_LOG_LEVEL=INFO
AUTO_PARALLEL_JOBS=4
AUTO_TIMEOUT=300
AUTO_RETRY_COUNT=3
EOF
        log "INFO" "Created default configuration"
    fi
    
    load_config "automation"
    log "SUCCESS" "Automation framework initialized"
}

# Run initialization when sourced
init_framework