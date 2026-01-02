#!/bin/sh
# components/ghostty/functions.sh - Ghostty helper functions

# =============================================================================
# Config Management
# =============================================================================

# Open Ghostty config in editor
ghostty_config() {
    local config="${GHOSTTY_CONFIG:-$HOME/.config/ghostty/config}"

    if [ ! -f "$config" ]; then
        echo "Ghostty config not found: $config"
        return 1
    fi

    ${EDITOR:-vim} "$config"
    echo "Config saved. Press Cmd+Shift+, (macOS) or Ctrl+Shift+, (Linux) to reload."
}

# Switch Ghostty theme
ghostty_theme() {
    local theme="$1"
    local config="${GHOSTTY_CONFIG:-$HOME/.config/ghostty/config}"

    if [ -z "$theme" ]; then
        echo "Usage: ghostty_theme <theme-name>"
        echo ""
        echo "Available themes (run 'ghostty +list-themes' for full list):"
        echo "  Catppuccin Mocha    (dark)"
        echo "  Catppuccin Macchiato (dark)"
        echo "  Catppuccin Frappe   (dark)"
        echo "  Catppuccin Latte    (light)"
        echo "  Dracula"
        echo "  Gruvbox Dark"
        echo "  Nord"
        echo "  One Dark"
        echo "  Solarized Dark"
        echo "  Tokyo Night"
        return 1
    fi

    if [ ! -f "$config" ]; then
        echo "Ghostty config not found: $config"
        return 1
    fi

    # Check if theme line exists
    if grep -q "^theme = " "$config"; then
        # Update existing theme line
        sed -i.bak "s/^theme = .*/theme = $theme/" "$config"
        rm -f "${config}.bak"
    else
        # Add theme line after first comment block
        sed -i.bak "1a\\
theme = $theme
" "$config"
        rm -f "${config}.bak"
    fi

    echo "Theme set to: $theme"
    echo "Press Cmd+Shift+, (macOS) or Ctrl+Shift+, (Linux) to reload."
}

# Change Ghostty font
ghostty_font() {
    local font="$1"
    local size="$2"
    local config="${GHOSTTY_CONFIG:-$HOME/.config/ghostty/config}"

    if [ -z "$font" ]; then
        echo "Usage: ghostty_font <font-family> [size]"
        echo "Example: ghostty_font 'JetBrains Mono' 14"
        return 1
    fi

    if [ ! -f "$config" ]; then
        echo "Ghostty config not found: $config"
        return 1
    fi

    # Update font-family
    if grep -q "^font-family = " "$config"; then
        sed -i.bak "s/^font-family = .*/font-family = \"$font\"/" "$config"
        rm -f "${config}.bak"
    fi

    # Update font-size if provided
    if [ -n "$size" ]; then
        if grep -q "^font-size = " "$config"; then
            sed -i.bak "s/^font-size = .*/font-size = $size/" "$config"
            rm -f "${config}.bak"
        fi
    fi

    echo "Font updated: $font${size:+ (size: $size)}"
    echo "Press Cmd+Shift+, (macOS) or Ctrl+Shift+, (Linux) to reload."
}

# =============================================================================
# Backup and Restore
# =============================================================================

# Backup current config
ghostty_backup() {
    local config="${GHOSTTY_CONFIG:-$HOME/.config/ghostty/config}"
    local backup_dir="${XDG_DATA_HOME:-$HOME/.local/share}/ghostty/backups"
    local timestamp
    timestamp=$(date +%Y%m%d_%H%M%S)

    if [ ! -f "$config" ]; then
        echo "Ghostty config not found: $config"
        return 1
    fi

    mkdir -p "$backup_dir"
    cp "$config" "$backup_dir/config.$timestamp"
    echo "Backup saved: $backup_dir/config.$timestamp"
}

# List config backups
ghostty_backups() {
    local backup_dir="${XDG_DATA_HOME:-$HOME/.local/share}/ghostty/backups"

    if [ ! -d "$backup_dir" ]; then
        echo "No backups found"
        return 1
    fi

    echo "Ghostty config backups:"
    ls -la "$backup_dir"
}

# Restore config from backup
ghostty_restore() {
    local backup="$1"
    local config="${GHOSTTY_CONFIG:-$HOME/.config/ghostty/config}"
    local backup_dir="${XDG_DATA_HOME:-$HOME/.local/share}/ghostty/backups"

    if [ -z "$backup" ]; then
        echo "Usage: ghostty_restore <backup-file>"
        ghostty_backups
        return 1
    fi

    local backup_file="$backup_dir/$backup"
    if [ ! -f "$backup_file" ]; then
        backup_file="$backup"
    fi

    if [ ! -f "$backup_file" ]; then
        echo "Backup not found: $backup"
        return 1
    fi

    # Backup current before restore
    ghostty_backup

    cp "$backup_file" "$config"
    echo "Config restored from: $backup_file"
    echo "Press Cmd+Shift+, (macOS) or Ctrl+Shift+, (Linux) to reload."
}

# =============================================================================
# Info
# =============================================================================

# Show Ghostty info
ghostty_info() {
    echo "Ghostty Terminal Information"
    echo "============================"

    if command -v ghostty >/dev/null 2>&1; then
        echo "Version: $(ghostty --version 2>/dev/null || echo 'installed')"
    else
        echo "Version: not installed"
    fi

    echo ""
    echo "Configuration:"
    echo "  Config: ${GHOSTTY_CONFIG:-$HOME/.config/ghostty/config}"

    if [ -f "${GHOSTTY_CONFIG:-$HOME/.config/ghostty/config}" ]; then
        local theme font size
        theme=$(grep "^theme = " "${GHOSTTY_CONFIG:-$HOME/.config/ghostty/config}" 2>/dev/null | cut -d= -f2 | tr -d ' ')
        font=$(grep "^font-family = " "${GHOSTTY_CONFIG:-$HOME/.config/ghostty/config}" 2>/dev/null | cut -d= -f2 | tr -d ' "')
        size=$(grep "^font-size = " "${GHOSTTY_CONFIG:-$HOME/.config/ghostty/config}" 2>/dev/null | cut -d= -f2 | tr -d ' ')

        echo "  Theme: ${theme:-manual palette}"
        echo "  Font: ${font:-default} ${size:-default}"
    fi
}
