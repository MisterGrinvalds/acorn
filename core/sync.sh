#!/bin/sh
# core/sync.sh - Dotfiles drift detection and synchronization
# Depends on: discovery.sh, xdg.sh, theme.sh
#
# Provides functions for:
#   - Git status checking
#   - Drift detection
#   - Auto-sync on shell startup
#   - Sync logging

# =============================================================================
# Configuration
# =============================================================================

SYNC_LOG="${XDG_STATE_HOME}/shell/sync.log"
SYNC_LAST_CHECK="${XDG_STATE_HOME}/shell/sync_last_check"
SYNC_AUTO_ENABLED="${DOTFILES_AUTO_SYNC:-1}"  # Enabled by default

# =============================================================================
# Logging
# =============================================================================

_sync_log() {
    local level="$1"
    local msg="$2"
    local timestamp

    timestamp=$(date '+%Y-%m-%d %H:%M:%S')

    # Log to file
    echo "[$timestamp] [$level] $msg" >> "$SYNC_LOG" 2>/dev/null

    # Also print if interactive
    case "$level" in
        info)    printf "${THEME_INFO}[sync]${THEME_RESET} %s\n" "$msg" ;;
        warn)    printf "${THEME_WARNING}[sync]${THEME_RESET} %s\n" "$msg" ;;
        error)   printf "${THEME_ERROR}[sync]${THEME_RESET} %s\n" "$msg" ;;
        success) printf "${THEME_SUCCESS}[sync]${THEME_RESET} %s\n" "$msg" ;;
    esac
}

# =============================================================================
# Git Operations
# =============================================================================

# Check if dotfiles repo is a git repo
_sync_is_git_repo() {
    git -C "$DOTFILES_ROOT" rev-parse --is-inside-work-tree >/dev/null 2>&1
}

# Get current branch
_sync_current_branch() {
    git -C "$DOTFILES_ROOT" rev-parse --abbrev-ref HEAD 2>/dev/null
}

# Get short status
_sync_git_status() {
    git -C "$DOTFILES_ROOT" status --porcelain 2>/dev/null
}

# Check if there are uncommitted changes
_sync_has_changes() {
    [ -n "$(_sync_git_status)" ]
}

# Check commits ahead/behind
_sync_commits_ahead() {
    git -C "$DOTFILES_ROOT" rev-list --count HEAD@{upstream}..HEAD 2>/dev/null || echo "0"
}

_sync_commits_behind() {
    git -C "$DOTFILES_ROOT" rev-list --count HEAD..HEAD@{upstream} 2>/dev/null || echo "0"
}

# Fetch latest from remote (quiet)
_sync_fetch() {
    git -C "$DOTFILES_ROOT" fetch --quiet 2>/dev/null
}

# =============================================================================
# Drift Detection
# =============================================================================

# Quick drift check (for shell startup)
dotfiles_check_drift() {
    local behind ahead modified untracked

    if ! _sync_is_git_repo; then
        _sync_log warn "Dotfiles directory is not a git repository"
        return 1
    fi

    # Fetch latest (background, quiet)
    _sync_fetch

    behind=$(_sync_commits_behind)
    ahead=$(_sync_commits_ahead)

    # Count modified and untracked files
    modified=$(git -C "$DOTFILES_ROOT" status --porcelain 2>/dev/null | grep -c '^ M\|^M ')
    untracked=$(git -C "$DOTFILES_ROOT" status --porcelain 2>/dev/null | grep -c '^??')

    # Build status message
    local has_drift=false
    local status_parts=""

    if [ "$behind" -gt 0 ]; then
        status_parts="${status_parts}${behind} behind, "
        has_drift=true
    fi

    if [ "$ahead" -gt 0 ]; then
        status_parts="${status_parts}${ahead} ahead, "
        has_drift=true
    fi

    if [ "$modified" -gt 0 ]; then
        status_parts="${status_parts}${modified} modified, "
        has_drift=true
    fi

    if [ "$untracked" -gt 0 ]; then
        status_parts="${status_parts}${untracked} untracked, "
        has_drift=true
    fi

    if [ "$has_drift" = "true" ]; then
        # Remove trailing comma and space
        status_parts=$(echo "$status_parts" | sed 's/, $//')
        printf "${THEME_WARNING}[dotfiles]${THEME_RESET} drift detected: %s\n" "$status_parts"
        return 1
    fi

    return 0
}

# Full drift audit
dotfiles_audit() {
    echo ""
    echo "Dotfiles Drift Report"
    echo "====================="
    echo ""

    if ! _sync_is_git_repo; then
        echo "Error: Dotfiles directory is not a git repository"
        return 1
    fi

    # Repository info
    echo "Repository: $DOTFILES_ROOT"
    echo "Branch: $(_sync_current_branch)"
    echo ""

    # Fetch and check remote status
    _sync_fetch
    local behind=$(_sync_commits_behind)
    local ahead=$(_sync_commits_ahead)

    if [ "$behind" -gt 0 ] || [ "$ahead" -gt 0 ]; then
        printf "Remote Status: "
        [ "$behind" -gt 0 ] && printf "${THEME_WARNING}%d behind${THEME_RESET} " "$behind"
        [ "$ahead" -gt 0 ] && printf "${THEME_INFO}%d ahead${THEME_RESET}" "$ahead"
        echo ""
    else
        printf "Remote Status: ${THEME_SUCCESS}up to date${THEME_RESET}\n"
    fi
    echo ""

    # File changes
    local status
    status=$(_sync_git_status)

    if [ -n "$status" ]; then
        echo "File Changes:"
        echo "-------------"
        echo "$status" | while IFS= read -r line; do
            local prefix=$(echo "$line" | cut -c1-2)
            local file=$(echo "$line" | cut -c4-)
            case "$prefix" in
                "M "|" M") printf "  ${THEME_PEACH}M${THEME_RESET} %s\n" "$file" ;;
                "A "|" A") printf "  ${THEME_GREEN}A${THEME_RESET} %s\n" "$file" ;;
                "D "|" D") printf "  ${THEME_RED}D${THEME_RESET} %s\n" "$file" ;;
                "??")      printf "  ${THEME_MAUVE}?${THEME_RESET} %s\n" "$file" ;;
                *)         printf "  %s %s\n" "$prefix" "$file" ;;
            esac
        done
        echo ""
    else
        printf "File Changes: ${THEME_SUCCESS}none${THEME_RESET}\n\n"
    fi

    # XDG compliance check
    echo "XDG Compliance:"
    echo "---------------"

    local xdg_issues=0

    # Check for legacy dotfiles
    for legacy_file in ~/.bashrc ~/.zshrc ~/.bash_profile; do
        if [ -f "$legacy_file" ] && ! [ -L "$legacy_file" ]; then
            printf "  ${THEME_WARNING}!${THEME_RESET} %s exists (should be symlink to bootstrap)\n" "$legacy_file"
            xdg_issues=$((xdg_issues + 1))
        fi
    done

    if [ "$xdg_issues" -eq 0 ]; then
        printf "  ${THEME_SUCCESS}All checks passed${THEME_RESET}\n"
    fi
    echo ""
}

# =============================================================================
# Sync Operations
# =============================================================================

# Show detailed dotfiles status
dotfiles_status() {
    echo "Dotfiles Status"
    echo "==============="
    echo ""
    echo "Location: $DOTFILES_ROOT"
    echo "Shell: ${CURRENT_SHELL:-unknown}"
    echo "Platform: ${CURRENT_PLATFORM:-unknown}"
    echo ""

    # Check bootstrap files
    echo "Bootstrap files:"
    [ -f "$HOME/.bashrc" ] && echo "  ~/.bashrc: exists" || echo "  ~/.bashrc: missing"
    [ -f "$HOME/.zshrc" ] && echo "  ~/.zshrc: exists" || echo "  ~/.zshrc: missing"
    [ -f "$HOME/.bash_profile" ] && echo "  ~/.bash_profile: exists" || echo "  ~/.bash_profile: missing"
    echo ""

    # Check git status
    if _sync_is_git_repo; then
        echo "Git status:"
        git -C "$DOTFILES_ROOT" status --short --branch
    fi
}

# Pull latest changes
dotfiles_pull() {
    if ! _sync_is_git_repo; then
        echo "Error: Dotfiles directory is not a git repository"
        return 1
    fi

    _sync_log info "Pulling latest changes..."
    git -C "$DOTFILES_ROOT" pull --rebase

    if [ $? -eq 0 ]; then
        _sync_log success "Pull completed successfully"
    else
        _sync_log error "Pull failed"
        return 1
    fi
}

# Commit and push changes
dotfiles_push() {
    local message="${1:-Update dotfiles}"

    if ! _sync_is_git_repo; then
        echo "Error: Dotfiles directory is not a git repository"
        return 1
    fi

    if ! _sync_has_changes; then
        _sync_log info "No changes to commit"
        return 0
    fi

    _sync_log info "Committing changes..."
    git -C "$DOTFILES_ROOT" add -A
    git -C "$DOTFILES_ROOT" commit -m "$message"

    _sync_log info "Pushing to remote..."
    git -C "$DOTFILES_ROOT" push

    if [ $? -eq 0 ]; then
        _sync_log success "Push completed successfully"
    else
        _sync_log error "Push failed"
        return 1
    fi
}

# Full sync (pull then push)
dotfiles_sync() {
    dotfiles_pull || return 1

    if _sync_has_changes; then
        dotfiles_push "Sync local changes"
    fi
}

# =============================================================================
# Auto-Sync (Shell Startup)
# =============================================================================

# Run auto-sync check on shell startup
_sync_auto_check() {
    # Skip if disabled
    [ "$SYNC_AUTO_ENABLED" != "1" ] && return 0

    # Skip if not a git repo
    _sync_is_git_repo || return 0

    # Run quick drift check (non-blocking)
    dotfiles_check_drift
}

# Enable/disable auto-sync
dotfiles_auto_sync() {
    local action="${1:-status}"

    case "$action" in
        on|enable|1)
            export DOTFILES_AUTO_SYNC=1
            _sync_log success "Auto-sync enabled"
            ;;
        off|disable|0)
            export DOTFILES_AUTO_SYNC=0
            _sync_log info "Auto-sync disabled"
            ;;
        status|*)
            if [ "$DOTFILES_AUTO_SYNC" = "1" ]; then
                echo "Auto-sync: enabled"
            else
                echo "Auto-sync: disabled"
            fi
            ;;
    esac
}

# =============================================================================
# Bootstrap Management (from inject.sh)
# =============================================================================

# Install bootstrap files to user's home
dotfiles_inject() {
    echo "Installing dotfiles bootstrap..."
    echo "Dotfiles location: $DOTFILES_ROOT"

    # Create XDG directories
    mkdir -p "${XDG_CONFIG_HOME:-$HOME/.config}"
    mkdir -p "${XDG_DATA_HOME:-$HOME/.local/share}"
    mkdir -p "${XDG_CACHE_HOME:-$HOME/.cache}"
    mkdir -p "${XDG_STATE_HOME:-$HOME/.local/state}"

    # Create .bashrc bootstrap
    if [ ! -f "$HOME/.bashrc" ] || ! grep -q "DOTFILES_ROOT" "$HOME/.bashrc" 2>/dev/null; then
        echo "Creating ~/.bashrc bootstrap..."
        cat > "$HOME/.bashrc" << EOF
# Dotfiles bootstrap - component-based architecture
export DOTFILES_ROOT="$DOTFILES_ROOT"
[ -f "\$DOTFILES_ROOT/core/bootstrap.sh" ] && . "\$DOTFILES_ROOT/core/bootstrap.sh"
EOF
        echo "Created ~/.bashrc"
    else
        echo "~/.bashrc already configured"
    fi

    # Create .zshrc bootstrap
    if [ ! -f "$HOME/.zshrc" ] || ! grep -q "DOTFILES_ROOT" "$HOME/.zshrc" 2>/dev/null; then
        echo "Creating ~/.zshrc bootstrap..."
        cat > "$HOME/.zshrc" << EOF
# Dotfiles bootstrap - component-based architecture
export DOTFILES_ROOT="$DOTFILES_ROOT"
[ -f "\$DOTFILES_ROOT/core/bootstrap.sh" ] && . "\$DOTFILES_ROOT/core/bootstrap.sh"
EOF
        echo "Created ~/.zshrc"
    else
        echo "~/.zshrc already configured"
    fi

    # Create .bash_profile that sources .bashrc
    if [ ! -f "$HOME/.bash_profile" ] || ! grep -q "bashrc" "$HOME/.bash_profile" 2>/dev/null; then
        echo "Creating ~/.bash_profile..."
        cat > "$HOME/.bash_profile" << 'EOF'
# Source .bashrc for login shells
[ -f "$HOME/.bashrc" ] && . "$HOME/.bashrc"
EOF
        echo "Created ~/.bash_profile"
    else
        echo "~/.bash_profile already configured"
    fi

    echo ""
    echo "Bootstrap installation complete!"
    echo "Restart your shell or run: source ~/.bashrc"
}

# Remove all injected configuration
dotfiles_eject() {
    echo "Removing dotfiles bootstrap..."

    local files_removed=0

    # Remove .bashrc if it's our bootstrap
    if [ -f "$HOME/.bashrc" ] && grep -q "DOTFILES_ROOT" "$HOME/.bashrc" 2>/dev/null; then
        rm "$HOME/.bashrc"
        echo "Removed ~/.bashrc"
        files_removed=$((files_removed + 1))
    fi

    # Remove .zshrc if it's our bootstrap
    if [ -f "$HOME/.zshrc" ] && grep -q "DOTFILES_ROOT" "$HOME/.zshrc" 2>/dev/null; then
        rm "$HOME/.zshrc"
        echo "Removed ~/.zshrc"
        files_removed=$((files_removed + 1))
    fi

    # Remove .bash_profile if it's our bootstrap
    if [ -f "$HOME/.bash_profile" ] && grep -q "bashrc" "$HOME/.bash_profile" 2>/dev/null; then
        rm "$HOME/.bash_profile"
        echo "Removed ~/.bash_profile"
        files_removed=$((files_removed + 1))
    fi

    if [ $files_removed -eq 0 ]; then
        echo "No bootstrap files found to remove"
    else
        echo ""
        echo "Removed $files_removed bootstrap file(s)"
        echo "Your shell configuration has been reset"
    fi
}

# Update dotfiles from git and reload
dotfiles_update() {
    echo "Updating dotfiles..."

    if _sync_is_git_repo; then
        dotfiles_pull || return 1
        echo ""
        echo "Dotfiles updated. Reloading..."
        dotfiles_reload
    else
        echo "Not a git repository: $DOTFILES_ROOT"
        return 1
    fi
}

# Reload shell configuration without restart
dotfiles_reload() {
    echo "Reloading shell configuration..."

    if [ -f "$DOTFILES_ROOT/core/bootstrap.sh" ]; then
        . "$DOTFILES_ROOT/core/bootstrap.sh"
        echo "Configuration reloaded"
    else
        echo "bootstrap.sh not found at $DOTFILES_ROOT/core/bootstrap.sh"
        return 1
    fi
}

# =============================================================================
# Config Linking
# =============================================================================

# Symlink app configs to their expected locations
# =============================================================================
# Component-Driven Config Linking
# =============================================================================

# Internal: Expand path (handle ~ and environment variables)
_sync_expand_path() {
    local path="$1"
    # Expand ~ to $HOME
    path="${path/#\~/$HOME}"
    # Expand environment variables
    eval echo "$path"
}

# Internal: Link configs for a single component
_sync_link_component_configs() {
    local component="$1"
    local comp_dir="${DOTFILES_ROOT}/components/${component}"
    local yaml_file="${comp_dir}/component.yaml"

    [ -f "$yaml_file" ] || return 0

    # Check if yq is available
    if ! command -v yq &>/dev/null; then
        return 0
    fi

    # Check if component has config section
    local file_count
    file_count=$(yq -r '.config.files | length // 0' "$yaml_file" 2>/dev/null)
    [ "$file_count" = "0" ] || [ -z "$file_count" ] && return 0

    # Check platform support
    local platforms
    platforms=$(yq -r '.platforms // []' "$yaml_file" 2>/dev/null)
    if [ "$platforms" != "[]" ] && [ "$platforms" != "null" ]; then
        if ! echo "$platforms" | grep -q "$CURRENT_PLATFORM"; then
            return 0
        fi
    fi

    # Create directories first
    local dir_count i
    dir_count=$(yq -r '.config.directories | length // 0' "$yaml_file" 2>/dev/null)
    for i in $(seq 0 $((dir_count - 1))); do
        local target perms
        target=$(yq -r ".config.directories[$i].target" "$yaml_file")
        perms=$(yq -r ".config.directories[$i].permissions // \"\"" "$yaml_file")
        target=$(_sync_expand_path "$target")
        mkdir -p "$target"
        [ -n "$perms" ] && [ "$perms" != "null" ] && chmod "$perms" "$target" 2>/dev/null
    done

    # Process each config file
    local linked_count=0
    for i in $(seq 0 $((file_count - 1))); do
        local source target method platform perms
        source=$(yq -r ".config.files[$i].source" "$yaml_file")
        target=$(yq -r ".config.files[$i].target" "$yaml_file")
        method=$(yq -r ".config.files[$i].method // \"symlink\"" "$yaml_file")
        platform=$(yq -r ".config.files[$i].platform // \"\"" "$yaml_file")
        perms=$(yq -r ".config.files[$i].permissions // \"\"" "$yaml_file")

        # Skip if platform-specific and not matching
        if [ -n "$platform" ] && [ "$platform" != "null" ] && [ "$platform" != "$CURRENT_PLATFORM" ]; then
            continue
        fi

        # Resolve paths
        local source_path="${comp_dir}/${source}"
        local target_path
        target_path=$(_sync_expand_path "$target")

        # Ensure source exists
        if [ ! -e "$source_path" ]; then
            continue
        fi

        # Ensure target directory exists
        mkdir -p "$(dirname "$target_path")"

        # Deploy based on method
        case "$method" in
            symlink)
                ln -sf "$source_path" "$target_path"
                ;;
            copy)
                cp "$source_path" "$target_path"
                ;;
        esac

        # Apply permissions if specified
        if [ -n "$perms" ] && [ "$perms" != "null" ]; then
            chmod "$perms" "$target_path" 2>/dev/null
        fi

        linked_count=$((linked_count + 1))
    done

    if [ "$linked_count" -gt 0 ]; then
        printf "  ${THEME_SUCCESS}%s${THEME_RESET}: %d file(s) linked\n" "$component" "$linked_count"
    fi
}

# Internal: Unlink configs for a single component
_sync_unlink_component_configs() {
    local component="$1"
    local comp_dir="${DOTFILES_ROOT}/components/${component}"
    local yaml_file="${comp_dir}/component.yaml"

    [ -f "$yaml_file" ] || return 0

    # Check if yq is available
    if ! command -v yq &>/dev/null; then
        return 0
    fi

    # Check if component has config section
    local file_count
    file_count=$(yq -r '.config.files | length // 0' "$yaml_file" 2>/dev/null)
    [ "$file_count" = "0" ] || [ -z "$file_count" ] && return 0

    # Process each config file
    local unlinked_count=0
    local i
    for i in $(seq 0 $((file_count - 1))); do
        local target method
        target=$(yq -r ".config.files[$i].target" "$yaml_file")
        method=$(yq -r ".config.files[$i].method // \"symlink\"" "$yaml_file")

        local target_path
        target_path=$(_sync_expand_path "$target")

        # Only remove symlinks (not copies, to be safe)
        if [ "$method" = "symlink" ] && [ -L "$target_path" ]; then
            rm "$target_path"
            unlinked_count=$((unlinked_count + 1))
        elif [ "$method" = "copy" ] && [ -f "$target_path" ]; then
            # For copies, we could optionally remove but it's safer not to
            # User may have modified the copy
            :
        fi
    done

    if [ "$unlinked_count" -gt 0 ]; then
        printf "  ${THEME_PEACH}%s${THEME_RESET}: %d file(s) unlinked\n" "$component" "$unlinked_count"
    fi
}

# Link app configurations from all components
dotfiles_link_configs() {
    _sync_log info "Linking component configurations..."

    # Check for yq
    if ! command -v yq &>/dev/null; then
        _sync_log error "yq is required for config linking. Install with: brew install yq"
        return 1
    fi

    echo ""

    local component
    for comp_dir in "${DOTFILES_ROOT}/components"/*/; do
        [ -d "$comp_dir" ] || continue
        component=$(basename "$comp_dir")
        [ "$component" = "_template" ] && continue

        _sync_link_component_configs "$component"
    done

    echo ""
    _sync_log success "Config linking complete!"
}

# Unlink app configs from all components
dotfiles_unlink_configs() {
    _sync_log info "Unlinking component configurations..."

    # Check for yq
    if ! command -v yq &>/dev/null; then
        _sync_log error "yq is required for config unlinking. Install with: brew install yq"
        return 1
    fi

    echo ""

    local component
    for comp_dir in "${DOTFILES_ROOT}/components"/*/; do
        [ -d "$comp_dir" ] || continue
        component=$(basename "$comp_dir")
        [ "$component" = "_template" ] && continue

        _sync_unlink_component_configs "$component"
    done

    echo ""
    _sync_log success "Config unlinking complete!"
}

# =============================================================================
# Component-Level Sync Operations
# =============================================================================

# Check health of a specific component or all components
# Usage: component_health [component_name]
component_health() {
    local component="${1:-}"
    local component_dir yaml_file

    if [ -n "$component" ]; then
        # Single component health check
        component_dir="${DOTFILES_ROOT}/components/${component}"
        yaml_file="${component_dir}/component.yaml"

        if [ ! -f "$yaml_file" ]; then
            printf "${THEME_ERROR}Component not found: %s${THEME_RESET}\n" "$component"
            return 1
        fi

        _component_health_check "$component"
    else
        # All components health check
        echo ""
        echo "Component Health Report"
        echo "======================="
        echo ""

        local total=0 healthy=0 warnings=0 errors=0

        for comp_dir in "${DOTFILES_ROOT}/components"/*/; do
            [ -d "$comp_dir" ] || continue
            local comp
            comp=$(basename "$comp_dir")
            [ "$comp" = "_template" ] && continue
            [ -f "${comp_dir}component.yaml" ] || continue

            total=$((total + 1))
            local status
            status=$(_component_health_check "$comp" quiet)
            case "$status" in
                healthy) healthy=$((healthy + 1)) ;;
                warning) warnings=$((warnings + 1)) ;;
                error)   errors=$((errors + 1)) ;;
            esac
        done

        echo ""
        echo "Summary: $total components"
        printf "  ${THEME_SUCCESS}Healthy: %d${THEME_RESET}\n" "$healthy"
        [ "$warnings" -gt 0 ] && printf "  ${THEME_WARNING}Warnings: %d${THEME_RESET}\n" "$warnings"
        [ "$errors" -gt 0 ] && printf "  ${THEME_ERROR}Errors: %d${THEME_RESET}\n" "$errors"
    fi
}

# Internal: Check health of a single component
_component_health_check() {
    local component="$1"
    local quiet="${2:-}"
    local component_dir="${DOTFILES_ROOT}/components/${component}"
    local yaml_file="${component_dir}/component.yaml"
    local status="healthy"
    local issues=""

    # Check YAML is valid
    if ! yq '.' "$yaml_file" >/dev/null 2>&1; then
        issues="${issues}  - Invalid component.yaml\n"
        status="error"
    fi

    # Check required tools
    local tools tool
    tools=$(yq -r '.requires.tools // [] | .[]' "$yaml_file" 2>/dev/null)
    for tool in $tools; do
        if ! command -v "$tool" >/dev/null 2>&1; then
            issues="${issues}  - Missing tool: ${tool}\n"
            [ "$status" != "error" ] && status="warning"
        fi
    done

    # Check shell files syntax
    for sh_file in "${component_dir}"/*.sh; do
        [ -f "$sh_file" ] || continue
        if ! bash -n "$sh_file" 2>/dev/null; then
            issues="${issues}  - Syntax error: $(basename "$sh_file")\n"
            status="error"
        fi
    done

    # Check component dependencies
    local deps dep
    deps=$(yq -r '.requires.components // [] | .[]' "$yaml_file" 2>/dev/null)
    for dep in $deps; do
        if [ ! -d "${DOTFILES_ROOT}/components/${dep}" ]; then
            issues="${issues}  - Missing dependency: ${dep}\n"
            [ "$status" != "error" ] && status="warning"
        fi
    done

    # Output
    if [ "$quiet" = "quiet" ]; then
        echo "$status"
    else
        local symbol color
        case "$status" in
            healthy) symbol="✓"; color="$THEME_SUCCESS" ;;
            warning) symbol="⚠"; color="$THEME_WARNING" ;;
            error)   symbol="✗"; color="$THEME_ERROR" ;;
        esac

        local desc
        desc=$(yq -r '.description // "No description"' "$yaml_file" 2>/dev/null)
        printf "${color}%s${THEME_RESET} %s - %s\n" "$symbol" "$component" "$desc"

        if [ -n "$issues" ]; then
            printf "$issues"
        fi
    fi

    [ "$status" = "healthy" ] && return 0 || return 1
}

# Show drift per component
# Usage: component_drift [component_name]
component_drift() {
    local component="${1:-}"

    if ! _sync_is_git_repo; then
        echo "Error: Not a git repository"
        return 1
    fi

    # Fetch to get accurate status
    _sync_fetch

    echo ""
    echo "Component Drift Report"
    echo "======================"
    echo ""

    local git_status
    git_status=$(_sync_git_status)

    if [ -z "$git_status" ]; then
        printf "${THEME_SUCCESS}No uncommitted changes${THEME_RESET}\n"
        return 0
    fi

    # Group changes by component
    local current_component=""
    local has_changes=false

    for comp_dir in "${DOTFILES_ROOT}/components"/*/; do
        [ -d "$comp_dir" ] || continue
        local comp
        comp=$(basename "$comp_dir")
        [ "$comp" = "_template" ] && continue

        # Skip if specific component requested and this isn't it
        [ -n "$component" ] && [ "$comp" != "$component" ] && continue

        # Find changes in this component
        local comp_changes
        comp_changes=$(echo "$git_status" | grep " components/${comp}/")

        if [ -n "$comp_changes" ]; then
            has_changes=true
            printf "${THEME_PEACH}%s${THEME_RESET}\n" "$comp"
            echo "$comp_changes" | while IFS= read -r line; do
                local prefix file
                prefix=$(echo "$line" | cut -c1-2)
                file=$(echo "$line" | cut -c4- | sed "s|components/${comp}/||")
                # Handle empty file (new directory)
                [ -z "$file" ] && file="(new component)"
                case "$prefix" in
                    "M "|" M") printf "  ${THEME_PEACH}M${THEME_RESET} %s\n" "$file" ;;
                    "A "|" A") printf "  ${THEME_GREEN}A${THEME_RESET} %s\n" "$file" ;;
                    "D "|" D") printf "  ${THEME_RED}D${THEME_RESET} %s\n" "$file" ;;
                    "??")      printf "  ${THEME_MAUVE}?${THEME_RESET} %s\n" "$file" ;;
                    *)         printf "  %s %s\n" "$prefix" "$file" ;;
                esac
            done
            echo ""
        fi
    done

    # Show core/ changes
    local core_changes
    core_changes=$(echo "$git_status" | grep " core/")
    if [ -n "$core_changes" ]; then
        has_changes=true
        printf "${THEME_BLUE}core${THEME_RESET}\n"
        echo "$core_changes" | while IFS= read -r line; do
            local prefix file
            prefix=$(echo "$line" | cut -c1-2)
            file=$(echo "$line" | cut -c4- | sed 's|core/||')
            case "$prefix" in
                "M "|" M") printf "  ${THEME_PEACH}M${THEME_RESET} %s\n" "$file" ;;
                "??")      printf "  ${THEME_MAUVE}?${THEME_RESET} %s\n" "$file" ;;
                *)         printf "  %s %s\n" "$prefix" "$file" ;;
            esac
        done
        echo ""
    fi

    # Show config/ changes
    local config_changes
    config_changes=$(echo "$git_status" | grep " config/")
    if [ -n "$config_changes" ]; then
        has_changes=true
        printf "${THEME_TEAL}config${THEME_RESET}\n"
        echo "$config_changes" | while IFS= read -r line; do
            local prefix file
            prefix=$(echo "$line" | cut -c1-2)
            file=$(echo "$line" | cut -c4- | sed 's|config/||')
            case "$prefix" in
                "M "|" M") printf "  ${THEME_PEACH}M${THEME_RESET} %s\n" "$file" ;;
                "??")      printf "  ${THEME_MAUVE}?${THEME_RESET} %s\n" "$file" ;;
                *)         printf "  %s %s\n" "$prefix" "$file" ;;
            esac
        done
        echo ""
    fi

    # Show other changes (not in components/, core/, or config/)
    local other_changes
    other_changes=$(echo "$git_status" | grep -v " components/" | grep -v " core/" | grep -v " config/")
    if [ -n "$other_changes" ]; then
        has_changes=true
        printf "${THEME_OVERLAY1}other${THEME_RESET}\n"
        echo "$other_changes" | while IFS= read -r line; do
            local prefix file
            prefix=$(echo "$line" | cut -c1-2)
            file=$(echo "$line" | cut -c4-)
            case "$prefix" in
                "M "|" M") printf "  ${THEME_PEACH}M${THEME_RESET} %s\n" "$file" ;;
                "A "|" A") printf "  ${THEME_GREEN}A${THEME_RESET} %s\n" "$file" ;;
                "D "|" D") printf "  ${THEME_RED}D${THEME_RESET} %s\n" "$file" ;;
                "??")      printf "  ${THEME_MAUVE}?${THEME_RESET} %s\n" "$file" ;;
                *)         printf "  %s %s\n" "$prefix" "$file" ;;
            esac
        done
        echo ""
    fi

    if [ "$has_changes" = false ]; then
        printf "${THEME_SUCCESS}No changes in components${THEME_RESET}\n"
    fi
}

# Interactive component sync - review and decide per component
# Usage: component_sync [component_name]
component_sync() {
    local component="${1:-}"

    if ! _sync_is_git_repo; then
        echo "Error: Not a git repository"
        return 1
    fi

    local git_status
    git_status=$(_sync_git_status)

    if [ -z "$git_status" ]; then
        printf "${THEME_SUCCESS}No changes to sync${THEME_RESET}\n"
        return 0
    fi

    echo ""
    echo "Interactive Component Sync"
    echo "=========================="
    echo ""
    echo "For each component with changes, choose:"
    echo "  [c]ommit  - Stage and commit changes"
    echo "  [r]evert  - Discard changes"
    echo "  [s]kip    - Leave as-is"
    echo "  [d]iff    - Show diff"
    echo "  [q]uit    - Stop processing"
    echo ""

    # Process components with changes
    for comp_dir in "${DOTFILES_ROOT}/components"/*/; do
        [ -d "$comp_dir" ] || continue
        local comp
        comp=$(basename "$comp_dir")
        [ "$comp" = "_template" ] && continue

        # Skip if specific component requested
        [ -n "$component" ] && [ "$comp" != "$component" ] && continue

        local comp_changes
        comp_changes=$(echo "$git_status" | grep "components/${comp}/")

        if [ -n "$comp_changes" ]; then
            _sync_handle_component "$comp" "$comp_changes" || break
        fi
    done

    # Handle core/ changes
    if [ -z "$component" ]; then
        local core_changes
        core_changes=$(echo "$git_status" | grep "^.. core/")
        if [ -n "$core_changes" ]; then
            _sync_handle_section "core" "$core_changes" || return
        fi
    fi

    echo ""
    printf "${THEME_SUCCESS}Sync complete${THEME_RESET}\n"
}

# Internal: Handle sync for a component
_sync_handle_component() {
    local component="$1"
    local changes="$2"

    printf "\n${THEME_PEACH}Component: %s${THEME_RESET}\n" "$component"
    echo "Changes:"
    echo "$changes" | while IFS= read -r line; do
        printf "  %s\n" "$line"
    done

    printf "\n[c]ommit, [r]evert, [s]kip, [d]iff, [q]uit? "
    read -r action

    case "$action" in
        c|C|commit)
            local files
            files=$(echo "$changes" | awk '{print $2}')
            echo "$files" | xargs git -C "$DOTFILES_ROOT" add
            git -C "$DOTFILES_ROOT" commit -m "Update ${component} component"
            printf "${THEME_SUCCESS}Committed %s changes${THEME_RESET}\n" "$component"
            ;;
        r|R|revert)
            local files
            files=$(echo "$changes" | grep -v "^??" | awk '{print $2}')
            if [ -n "$files" ]; then
                echo "$files" | xargs git -C "$DOTFILES_ROOT" checkout --
                printf "${THEME_WARNING}Reverted %s changes${THEME_RESET}\n" "$component"
            fi
            # Handle untracked files
            local untracked
            untracked=$(echo "$changes" | grep "^??" | awk '{print $2}')
            if [ -n "$untracked" ]; then
                printf "Remove untracked files? [y/N] "
                read -r confirm
                if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
                    echo "$untracked" | xargs rm -rf
                    printf "${THEME_WARNING}Removed untracked files${THEME_RESET}\n"
                fi
            fi
            ;;
        d|D|diff)
            local files
            files=$(echo "$changes" | grep -v "^??" | awk '{print $2}')
            if [ -n "$files" ]; then
                echo "$files" | xargs git -C "$DOTFILES_ROOT" diff
            fi
            # Re-prompt after showing diff
            _sync_handle_component "$component" "$changes"
            ;;
        s|S|skip)
            printf "${THEME_INFO}Skipped %s${THEME_RESET}\n" "$component"
            ;;
        q|Q|quit)
            return 1
            ;;
        *)
            printf "${THEME_WARNING}Invalid choice, skipping${THEME_RESET}\n"
            ;;
    esac
    return 0
}

# Internal: Handle sync for a section (core, config, etc.)
_sync_handle_section() {
    local section="$1"
    local changes="$2"

    printf "\n${THEME_BLUE}Section: %s${THEME_RESET}\n" "$section"
    echo "Changes:"
    echo "$changes" | while IFS= read -r line; do
        printf "  %s\n" "$line"
    done

    printf "\n[c]ommit, [r]evert, [s]kip, [d]iff, [q]uit? "
    read -r action

    case "$action" in
        c|C|commit)
            local files
            files=$(echo "$changes" | awk '{print $2}')
            echo "$files" | xargs git -C "$DOTFILES_ROOT" add
            git -C "$DOTFILES_ROOT" commit -m "Update ${section}"
            printf "${THEME_SUCCESS}Committed %s changes${THEME_RESET}\n" "$section"
            ;;
        r|R|revert)
            local files
            files=$(echo "$changes" | grep -v "^??" | awk '{print $2}')
            if [ -n "$files" ]; then
                echo "$files" | xargs git -C "$DOTFILES_ROOT" checkout --
                printf "${THEME_WARNING}Reverted %s changes${THEME_RESET}\n" "$section"
            fi
            ;;
        d|D|diff)
            local files
            files=$(echo "$changes" | grep -v "^??" | awk '{print $2}')
            if [ -n "$files" ]; then
                echo "$files" | xargs git -C "$DOTFILES_ROOT" diff
            fi
            _sync_handle_section "$section" "$changes"
            ;;
        s|S|skip)
            printf "${THEME_INFO}Skipped %s${THEME_RESET}\n" "$section"
            ;;
        q|Q|quit)
            return 1
            ;;
    esac
    return 0
}

# Validate component configurations
# Usage: component_validate [component_name]
component_validate() {
    local component="${1:-}"
    local errors=0

    echo ""
    echo "Component Validation"
    echo "===================="
    echo ""

    for comp_dir in "${DOTFILES_ROOT}/components"/*/; do
        [ -d "$comp_dir" ] || continue
        local comp
        comp=$(basename "$comp_dir")
        [ "$comp" = "_template" ] && continue

        # Skip if specific component requested
        [ -n "$component" ] && [ "$comp" != "$component" ] && continue

        local yaml_file="${comp_dir}component.yaml"
        local comp_errors=0

        # Check YAML exists
        if [ ! -f "$yaml_file" ]; then
            printf "${THEME_ERROR}✗${THEME_RESET} %s: missing component.yaml\n" "$comp"
            errors=$((errors + 1))
            continue
        fi

        # Validate YAML syntax
        if ! yq '.' "$yaml_file" >/dev/null 2>&1; then
            printf "${THEME_ERROR}✗${THEME_RESET} %s: invalid YAML syntax\n" "$comp"
            errors=$((errors + 1))
            continue
        fi

        # Check required fields
        local name version description category
        name=$(yq -r '.name // ""' "$yaml_file")
        version=$(yq -r '.version // ""' "$yaml_file")
        description=$(yq -r '.description // ""' "$yaml_file")
        category=$(yq -r '.category // ""' "$yaml_file")

        local field_errors=""
        [ -z "$name" ] && field_errors="${field_errors}name, " && comp_errors=$((comp_errors + 1))
        [ -z "$version" ] && field_errors="${field_errors}version, " && comp_errors=$((comp_errors + 1))
        [ -z "$description" ] && field_errors="${field_errors}description, " && comp_errors=$((comp_errors + 1))
        [ -z "$category" ] && field_errors="${field_errors}category, " && comp_errors=$((comp_errors + 1))

        # Validate shell files
        for sh_file in "${comp_dir}"/*.sh; do
            [ -f "$sh_file" ] || continue
            if ! bash -n "$sh_file" 2>/dev/null; then
                field_errors="${field_errors}$(basename "$sh_file") syntax, "
                comp_errors=$((comp_errors + 1))
            fi
        done

        # Validate config files exist
        local config_count
        config_count=$(yq -r '.config.files | length // 0' "$yaml_file" 2>/dev/null)
        if [ "$config_count" != "0" ] && [ -n "$config_count" ]; then
            local i
            for i in $(seq 0 $((config_count - 1))); do
                local source method
                source=$(yq -r ".config.files[$i].source" "$yaml_file")
                method=$(yq -r ".config.files[$i].method // \"symlink\"" "$yaml_file")

                # Check source file exists
                if [ ! -e "${comp_dir}${source}" ]; then
                    field_errors="${field_errors}config:${source} missing, "
                    comp_errors=$((comp_errors + 1))
                fi

                # Validate method
                if [ "$method" != "symlink" ] && [ "$method" != "copy" ]; then
                    field_errors="${field_errors}config:invalid method '${method}', "
                    comp_errors=$((comp_errors + 1))
                fi
            done
        fi

        if [ $comp_errors -eq 0 ]; then
            printf "${THEME_SUCCESS}✓${THEME_RESET} %s (%s)\n" "$comp" "$version"
        else
            field_errors=$(echo "$field_errors" | sed 's/, $//')
            printf "${THEME_ERROR}✗${THEME_RESET} %s: missing/invalid: %s\n" "$comp" "$field_errors"
            errors=$((errors + comp_errors))
        fi
    done

    echo ""
    if [ $errors -eq 0 ]; then
        printf "${THEME_SUCCESS}All components valid${THEME_RESET}\n"
    else
        printf "${THEME_ERROR}%d validation error(s)${THEME_RESET}\n" "$errors"
    fi

    return $errors
}

# Quick overview of all components and their status
component_overview() {
    echo ""
    echo "Component Overview"
    echo "=================="
    echo ""

    # Group by category
    for category in core dev cloud ai database editor; do
        local has_category=false

        for comp_dir in "${DOTFILES_ROOT}/components"/*/; do
            [ -d "$comp_dir" ] || continue
            local comp
            comp=$(basename "$comp_dir")
            [ "$comp" = "_template" ] && continue

            local yaml_file="${comp_dir}component.yaml"
            [ -f "$yaml_file" ] || continue

            local comp_category
            comp_category=$(yq -r '.category // "unknown"' "$yaml_file" 2>/dev/null)

            if [ "$comp_category" = "$category" ]; then
                if [ "$has_category" = false ]; then
                    printf "${THEME_BLUE}[%s]${THEME_RESET}\n" "$category"
                    has_category=true
                fi

                # Check if tools are installed
                local tools_ok=true
                local tools
                tools=$(yq -r '.requires.tools // [] | .[]' "$yaml_file" 2>/dev/null)
                for tool in $tools; do
                    command -v "$tool" >/dev/null 2>&1 || tools_ok=false
                done

                local desc
                desc=$(yq -r '.description // ""' "$yaml_file" 2>/dev/null | head -c 45)

                if [ "$tools_ok" = true ]; then
                    printf "  ${THEME_SUCCESS}✓${THEME_RESET} %-15s %s\n" "$comp" "$desc"
                else
                    printf "  ${THEME_WARNING}○${THEME_RESET} %-15s %s\n" "$comp" "$desc"
                fi
            fi
        done

        [ "$has_category" = true ] && echo ""
    done
}

# =============================================================================
# Help
# =============================================================================

# Help function
dotfiles_help() {
    cat << 'EOF'
Dotfiles Management Commands
============================

INSTALLATION:
  dotfiles_inject         Install bootstrap files (~/.bashrc, ~/.zshrc)
  dotfiles_eject          Remove all bootstrap files
  dotfiles_link_configs   Symlink app configs (git, ssh, etc.)
  dotfiles_unlink_configs Remove app config symlinks

MANAGEMENT:
  dotfiles_update         Git pull and reload configuration
  dotfiles_reload         Reload without restart
  dotfiles_status         Show current status
  dotfiles_audit          Full drift report

SYNC:
  dotfiles_pull           Pull latest changes from remote
  dotfiles_push           Commit and push local changes
  dotfiles_sync           Pull then push (full sync)
  dotfiles_check_drift    Quick drift check
  dotfiles_auto_sync      Enable/disable auto-sync on startup

COMPONENT MANAGEMENT:
  component_health        Check health of components (tools, syntax)
  component_drift         Show changes grouped by component
  component_sync          Interactive per-component sync
  component_validate      Validate component configurations
  component_overview      Quick overview of all components

EXAMPLES:
  dotfiles_inject         # First-time setup
  dotfiles_link_configs   # Link git, ssh configs
  dotfiles_update         # Update from git
  dotfiles_status         # Check current state
  dotfiles_audit          # Full drift report
  component_health        # Check all component health
  component_drift python  # Show drift in python component
  component_sync          # Interactive sync per component
EOF
}

# Convenience aliases
alias df-inject='dotfiles_inject'
alias df-eject='dotfiles_eject'
alias df-update='dotfiles_update'
alias df-reload='dotfiles_reload'
alias df-status='dotfiles_status'
alias df-audit='dotfiles_audit'
alias df-link='dotfiles_link_configs'
alias df-unlink='dotfiles_unlink_configs'
alias df-help='dotfiles_help'

# Component management aliases
alias comp-health='component_health'
alias comp-drift='component_drift'
alias comp-sync='component_sync'
alias comp-validate='component_validate'
alias comp-overview='component_overview'

# =============================================================================
# Initialize
# =============================================================================

# Create sync log directory
mkdir -p "$(dirname "$SYNC_LOG")" 2>/dev/null
