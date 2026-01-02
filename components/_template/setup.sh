#!/bin/sh
# components/_template/setup.sh - Installation and setup script
#
# This script is run on-demand via: make setup-<component>
# It should be idempotent (safe to run multiple times).

set -e

COMPONENT_NAME="template"
COMPONENT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Source core modules for XDG functions
. "${COMPONENT_DIR}/../../core/discovery.sh"
. "${COMPONENT_DIR}/../../core/xdg.sh"

echo "Setting up ${COMPONENT_NAME} component..."

# =============================================================================
# Check Prerequisites
# =============================================================================

check_prerequisites() {
    local missing=""

    # Example: Check for required tools
    # if ! command -v some-tool >/dev/null 2>&1; then
    #     missing="${missing} some-tool"
    # fi

    if [ -n "$missing" ]; then
        echo "Missing prerequisites:${missing}"
        echo "Install with: brew install${missing}"
        return 1
    fi

    return 0
}

# =============================================================================
# Install Dependencies
# =============================================================================

install_dependencies() {
    case "$CURRENT_PLATFORM" in
        darwin)
            # Homebrew installation
            # brew install package1 package2
            echo "No brew packages to install"
            ;;
        linux)
            # APT installation
            # sudo apt-get install -y package1 package2
            echo "No apt packages to install"
            ;;
    esac
}

# =============================================================================
# Configure Component
# =============================================================================

configure_component() {
    # Create XDG directories
    xdg_ensure_dirs "$COMPONENT_NAME"

    # Example: Create default config file
    # local config_file="$(xdg_config_dir $COMPONENT_NAME)/config.yaml"
    # if [ ! -f "$config_file" ]; then
    #     echo "Creating default config: $config_file"
    #     cat > "$config_file" << 'EOF'
    # # Default configuration
    # setting: value
    # EOF
    # fi

    echo "Configuration complete"
}

# =============================================================================
# Validate Setup
# =============================================================================

validate_setup() {
    local errors=0

    # Example: Check that tool is working
    # if ! some-tool --version >/dev/null 2>&1; then
    #     echo "Error: some-tool is not working"
    #     errors=$((errors + 1))
    # fi

    if [ $errors -eq 0 ]; then
        echo "Validation passed"
        return 0
    else
        echo "Validation failed with $errors error(s)"
        return 1
    fi
}

# =============================================================================
# Main
# =============================================================================

main() {
    check_prerequisites || exit 1
    install_dependencies
    configure_component
    validate_setup

    echo ""
    echo "${COMPONENT_NAME} setup complete!"
}

main "$@"
