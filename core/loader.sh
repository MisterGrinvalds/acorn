#!/bin/sh
# core/loader.sh - Component discovery and loading
# Depends on: discovery.sh, xdg.sh, theme.sh
#
# Scans components/ directory for component.yaml files,
# resolves dependencies, and loads components in order.
#
# Environment variables:
#   DOTFILES_DISABLE_<COMPONENT> - Set to 1 to disable a component
#   DOTFILES_COMPONENTS          - Comma-separated list to load only specific components

# Enable word splitting in zsh (bash splits by default)
# shellcheck disable=SC3040
[ -n "$ZSH_VERSION" ] && setopt shwordsplit 2>/dev/null

# =============================================================================
# Configuration
# =============================================================================

COMPONENTS_DIR="${DOTFILES_ROOT}/components"
COMPONENT_REGISTRY="${XDG_DATA_HOME}/shell/component-registry.json"

# =============================================================================
# Utility Functions
# =============================================================================

# Log a message with color
_loader_log() {
    local level="$1"
    local msg="$2"

    case "$level" in
        info)    printf "${THEME_INFO}[loader]${THEME_RESET} %s\n" "$msg" ;;
        warn)    printf "${THEME_WARNING}[loader]${THEME_RESET} %s\n" "$msg" ;;
        error)   printf "${THEME_ERROR}[loader]${THEME_RESET} %s\n" "$msg" ;;
        success) printf "${THEME_SUCCESS}[loader]${THEME_RESET} %s\n" "$msg" ;;
    esac
}

# Show warning for missing tool (only once per component)
_loader_warn_missing() {
    local component="$1"
    local tool="$2"

    if ! xdg_warning_shown "$component"; then
        _loader_log warn "Component '$component' skipped: '$tool' not installed"
        xdg_warning_mark "$component"
    fi
}

# Check if a command exists
_loader_has_command() {
    command -v "$1" >/dev/null 2>&1
}

# Read a YAML field using yq
_loader_yaml_get() {
    local file="$1"
    local field="$2"
    yq -r "$field // empty" "$file" 2>/dev/null
}

# Read a YAML array as space-separated values
_loader_yaml_array() {
    local file="$1"
    local field="$2"
    yq -r "$field // [] | .[]" "$file" 2>/dev/null | tr '\n' ' '
}

# =============================================================================
# Component Discovery
# =============================================================================

# Get list of all component directories
_loader_discover_components() {
    local comp_dir comp_name
    for comp_dir in "${COMPONENTS_DIR}"/*/; do
        [ -d "$comp_dir" ] || continue
        comp_name=$(basename "$comp_dir")
        # Skip directories starting with _ (disabled/template)
        case "$comp_name" in
            _*) continue ;;
        esac
        # Must have component.yaml
        [ -f "${comp_dir}component.yaml" ] || continue
        echo "$comp_name"
    done
}

# Check if a component is explicitly disabled
_loader_is_disabled() {
    local component="$1"
    local upper_name

    # Convert to uppercase for env var check
    upper_name=$(echo "$component" | tr '[:lower:]' '[:upper:]' | tr '-' '_')

    # Check DOTFILES_DISABLE_<COMPONENT>
    eval "[ \"\${DOTFILES_DISABLE_${upper_name}:-0}\" = \"1\" ]"
}

# Check if component is in the whitelist (if DOTFILES_COMPONENTS is set)
_loader_in_whitelist() {
    local component="$1"

    # If no whitelist, all components are allowed
    [ -z "$DOTFILES_COMPONENTS" ] && return 0

    # Check if component is in comma-separated list
    echo ",$DOTFILES_COMPONENTS," | grep -q ",${component},"
}

# Check if component's platform is supported
_loader_platform_supported() {
    local component="$1"
    local yaml_file="${COMPONENTS_DIR}/${component}/component.yaml"
    local platforms

    platforms=$(_loader_yaml_array "$yaml_file" '.platforms')

    # Empty means all platforms supported
    [ -z "$platforms" ] && return 0

    # Check if current platform is in list
    echo "$platforms" | grep -q "$CURRENT_PLATFORM"
}

# Check if component's shell is supported
_loader_shell_supported() {
    local component="$1"
    local yaml_file="${COMPONENTS_DIR}/${component}/component.yaml"
    local shells

    shells=$(_loader_yaml_array "$yaml_file" '.shells')

    # Empty means all shells supported
    [ -z "$shells" ] && return 0

    # Check if current shell is in list
    echo "$shells" | grep -q "$CURRENT_SHELL"
}

# Check if required tools are installed
_loader_check_tools() {
    local component="$1"
    local yaml_file="${COMPONENTS_DIR}/${component}/component.yaml"
    local tools tool

    tools=$(_loader_yaml_array "$yaml_file" '.requires.tools')

    for tool in $tools; do
        if ! _loader_has_command "$tool"; then
            _loader_warn_missing "$component" "$tool"
            return 1
        fi
    done

    return 0
}

# =============================================================================
# Dependency Resolution
# =============================================================================

# Get component dependencies
_loader_get_deps() {
    local component="$1"
    local yaml_file="${COMPONENTS_DIR}/${component}/component.yaml"
    _loader_yaml_array "$yaml_file" '.requires.components'
}

# Topological sort for dependency resolution
# Uses a simple DFS-based approach
_loader_resolve_order() {
    local components="$1"
    local resolved=""
    local visiting=""
    local component

    _visit() {
        local comp="$1"
        local dep

        # Already resolved
        echo "$resolved" | grep -q " $comp " && return 0

        # Circular dependency check
        if echo "$visiting" | grep -q " $comp "; then
            _loader_log error "Circular dependency detected: $comp"
            return 1
        fi

        visiting="$visiting $comp "

        # Visit dependencies first
        for dep in $(_loader_get_deps "$comp"); do
            # Only visit if dependency is in our component list
            if echo " $components " | grep -q " $dep "; then
                _visit "$dep" || return 1
            fi
        done

        resolved="$resolved $comp "
        visiting=$(echo "$visiting" | sed "s/ $comp //")
    }

    for component in $components; do
        _visit "$component" || return 1
    done

    echo "$resolved"
}

# =============================================================================
# Component Loading
# =============================================================================

# Source a file if it exists
_loader_source_if_exists() {
    local file="$1"
    if [ -f "$file" ]; then
        . "$file"
        return 0
    fi
    return 1
}

# Load a single component
_loader_load_component() {
    local component="$1"
    local comp_dir="${COMPONENTS_DIR}/${component}"
    local yaml_file="${comp_dir}/component.yaml"

    # Check prerequisites
    _loader_is_disabled "$component" && return 0
    _loader_in_whitelist "$component" || return 0
    _loader_platform_supported "$component" || return 0
    _loader_shell_supported "$component" || return 0
    _loader_check_tools "$component" || return 0

    # Ensure XDG directories for this component
    local xdg_config xdg_data xdg_cache xdg_state
    xdg_config=$(_loader_yaml_get "$yaml_file" '.xdg.config')
    xdg_data=$(_loader_yaml_get "$yaml_file" '.xdg.data')
    xdg_cache=$(_loader_yaml_get "$yaml_file" '.xdg.cache')
    xdg_state=$(_loader_yaml_get "$yaml_file" '.xdg.state')

    [ -n "$xdg_config" ] && mkdir -p "${XDG_CONFIG_HOME}/${xdg_config}" 2>/dev/null
    [ -n "$xdg_data" ] && mkdir -p "${XDG_DATA_HOME}/${xdg_data}" 2>/dev/null
    [ -n "$xdg_cache" ] && mkdir -p "${XDG_CACHE_HOME}/${xdg_cache}" 2>/dev/null
    [ -n "$xdg_state" ] && mkdir -p "${XDG_STATE_HOME}/${xdg_state}" 2>/dev/null

    # Load env.sh (always, for PATH setup etc.)
    _loader_source_if_exists "${comp_dir}/env.sh"

    # Load interactive-only files
    if [ "$IS_INTERACTIVE" = "true" ]; then
        _loader_source_if_exists "${comp_dir}/aliases.sh"
        _loader_source_if_exists "${comp_dir}/functions.sh"
        _loader_source_if_exists "${comp_dir}/completions.sh"
    fi

    return 0
}

# =============================================================================
# Main Loading Sequence
# =============================================================================

loader_run() {
    local components ordered component
    local start_time end_time

    # Check if yq is available
    if ! _loader_has_command yq; then
        _loader_log error "yq is required for component loading. Install with: brew install yq"
        return 1
    fi

    # Discover components
    components=$(_loader_discover_components)

    if [ -z "$components" ]; then
        # No components found, that's okay for initial setup
        return 0
    fi

    # Resolve load order
    ordered=$(_loader_resolve_order "$components")

    if [ $? -ne 0 ]; then
        _loader_log error "Failed to resolve component dependencies"
        return 1
    fi

    # Load components in order
    for component in $ordered; do
        _loader_load_component "$component"
    done

    return 0
}

# Export functions that components might need
export -f xdg_config_dir xdg_data_dir xdg_cache_dir xdg_state_dir xdg_ensure_dirs 2>/dev/null
