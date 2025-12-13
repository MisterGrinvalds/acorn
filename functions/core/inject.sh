#!/bin/sh
# Dotfiles injection and management functions
# These functions install/remove the bootstrap configuration

# Get the dotfiles root directory
_dotfiles_root() {
    if [ -n "$DOTFILES_ROOT" ]; then
        echo "$DOTFILES_ROOT"
    else
        # Fallback: find this script's location
        echo "${XDG_CONFIG_HOME:-$HOME/.config}/dotfiles"
    fi
}

# Install bootstrap files to user's home
dotfiles_inject() {
    local dotfiles_root
    dotfiles_root=$(_dotfiles_root)

    echo "Installing dotfiles bootstrap..."
    echo "Dotfiles location: $dotfiles_root"

    # Create XDG directories
    mkdir -p "${XDG_CONFIG_HOME:-$HOME/.config}"
    mkdir -p "${XDG_DATA_HOME:-$HOME/.local/share}"
    mkdir -p "${XDG_CACHE_HOME:-$HOME/.cache}"
    mkdir -p "${XDG_STATE_HOME:-$HOME/.local/state}"

    # Create .bashrc bootstrap
    if [ ! -f "$HOME/.bashrc" ] || ! grep -q "dotfiles" "$HOME/.bashrc" 2>/dev/null; then
        echo "Creating ~/.bashrc bootstrap..."
        cat > "$HOME/.bashrc" << EOF
# Dotfiles bootstrap - sources configuration from repo
export DOTFILES_ROOT="$dotfiles_root"
[ -f "\$DOTFILES_ROOT/shell/init.sh" ] && . "\$DOTFILES_ROOT/shell/init.sh"
EOF
        echo "Created ~/.bashrc"
    else
        echo "~/.bashrc already configured"
    fi

    # Create .zshrc bootstrap
    if [ ! -f "$HOME/.zshrc" ] || ! grep -q "dotfiles" "$HOME/.zshrc" 2>/dev/null; then
        echo "Creating ~/.zshrc bootstrap..."
        cat > "$HOME/.zshrc" << EOF
# Dotfiles bootstrap - sources configuration from repo
export DOTFILES_ROOT="$dotfiles_root"
[ -f "\$DOTFILES_ROOT/shell/init.sh" ] && . "\$DOTFILES_ROOT/shell/init.sh"
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
    local dotfiles_root
    dotfiles_root=$(_dotfiles_root)

    echo "Updating dotfiles..."

    if [ -d "$dotfiles_root/.git" ]; then
        cd "$dotfiles_root" || return 1
        git pull --ff-only
        cd - > /dev/null || return 1
        echo ""
        echo "Dotfiles updated. Reloading..."
        dotfiles_reload
    else
        echo "Not a git repository: $dotfiles_root"
        return 1
    fi
}

# Reload shell configuration without restart
dotfiles_reload() {
    local dotfiles_root
    dotfiles_root=$(_dotfiles_root)

    echo "Reloading shell configuration..."

    if [ -f "$dotfiles_root/shell/init.sh" ]; then
        . "$dotfiles_root/shell/init.sh"
        echo "Configuration reloaded"
    else
        echo "init.sh not found at $dotfiles_root/shell/init.sh"
        return 1
    fi
}

# Show current dotfiles status
dotfiles_status() {
    local dotfiles_root
    dotfiles_root=$(_dotfiles_root)

    echo "Dotfiles Status"
    echo "==============="
    echo ""
    echo "Location: $dotfiles_root"
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
    if [ -d "$dotfiles_root/.git" ]; then
        echo "Git status:"
        cd "$dotfiles_root" || return 1
        git status --short --branch
        cd - > /dev/null || return 1
    fi
}

# Symlink app configs to their expected locations
dotfiles_link_configs() {
    local dotfiles_root
    dotfiles_root=$(_dotfiles_root)

    echo "Linking app configurations..."

    # Git config
    if [ -f "$dotfiles_root/config/git/config" ]; then
        ln -sf "$dotfiles_root/config/git/config" "$HOME/.gitconfig"
        echo "Linked: ~/.gitconfig"
    fi

    # SSH config (be careful with permissions)
    if [ -f "$dotfiles_root/config/ssh/config" ]; then
        mkdir -p "$HOME/.ssh"
        chmod 700 "$HOME/.ssh"
        ln -sf "$dotfiles_root/config/ssh/config" "$HOME/.ssh/config"
        chmod 600 "$HOME/.ssh/config"
        echo "Linked: ~/.ssh/config"
    fi

    # Conda config
    if [ -f "$dotfiles_root/config/conda/condarc" ]; then
        ln -sf "$dotfiles_root/config/conda/condarc" "$HOME/.condarc"
        echo "Linked: ~/.condarc"
    fi

    # Karabiner config (macOS only)
    if [ "$CURRENT_PLATFORM" = "darwin" ] && [ -f "$dotfiles_root/config/karabiner/karabiner.json" ]; then
        mkdir -p "$HOME/.config/karabiner"
        ln -sf "$dotfiles_root/config/karabiner/karabiner.json" "$HOME/.config/karabiner/karabiner.json"
        echo "Linked: ~/.config/karabiner/karabiner.json"
    fi

    # Ghostty config
    if [ -f "$dotfiles_root/config/ghostty/config" ]; then
        mkdir -p "$HOME/.config/ghostty"
        ln -sf "$dotfiles_root/config/ghostty/config" "$HOME/.config/ghostty/config"
        echo "Linked: ~/.config/ghostty/config"
    fi

    # VS Code config
    if [ -d "$dotfiles_root/config/vscode" ]; then
        local vscode_dir
        if [ "$CURRENT_PLATFORM" = "darwin" ]; then
            vscode_dir="$HOME/Library/Application Support/Code/User"
        else
            vscode_dir="$HOME/.config/Code/User"
        fi
        mkdir -p "$vscode_dir"

        if [ -f "$dotfiles_root/config/vscode/settings.json" ]; then
            ln -sf "$dotfiles_root/config/vscode/settings.json" "$vscode_dir/settings.json"
            echo "Linked: VS Code settings.json"
        fi
        if [ -f "$dotfiles_root/config/vscode/keybindings.json" ]; then
            ln -sf "$dotfiles_root/config/vscode/keybindings.json" "$vscode_dir/keybindings.json"
            echo "Linked: VS Code keybindings.json"
        fi
    fi

    # IntelliJ config (reference only - requires manual setup due to version directories)
    if [ -d "$dotfiles_root/config/intellij" ]; then
        echo "IntelliJ configs available at: $dotfiles_root/config/intellij/"
        echo "  - Manual linking required due to version-specific paths"
    fi

    # Claude Code config
    if [ -f "$dotfiles_root/config/claude/settings.json" ]; then
        mkdir -p "$HOME/.claude"
        ln -sf "$dotfiles_root/config/claude/settings.json" "$HOME/.claude/settings.json"
        echo "Linked: ~/.claude/settings.json"
    fi

    echo ""
    echo "Config linking complete!"
}

# Unlink app configs
dotfiles_unlink_configs() {
    echo "Unlinking app configurations..."

    [ -L "$HOME/.gitconfig" ] && rm "$HOME/.gitconfig" && echo "Unlinked: ~/.gitconfig"
    [ -L "$HOME/.ssh/config" ] && rm "$HOME/.ssh/config" && echo "Unlinked: ~/.ssh/config"
    [ -L "$HOME/.condarc" ] && rm "$HOME/.condarc" && echo "Unlinked: ~/.condarc"
    [ -L "$HOME/.config/karabiner/karabiner.json" ] && rm "$HOME/.config/karabiner/karabiner.json" && echo "Unlinked: ~/.config/karabiner/karabiner.json"
    [ -L "$HOME/.config/ghostty/config" ] && rm "$HOME/.config/ghostty/config" && echo "Unlinked: ~/.config/ghostty/config"
    [ -L "$HOME/.claude/settings.json" ] && rm "$HOME/.claude/settings.json" && echo "Unlinked: ~/.claude/settings.json"

    # VS Code configs
    local vscode_dir
    if [ "$CURRENT_PLATFORM" = "darwin" ]; then
        vscode_dir="$HOME/Library/Application Support/Code/User"
    else
        vscode_dir="$HOME/.config/Code/User"
    fi
    [ -L "$vscode_dir/settings.json" ] && rm "$vscode_dir/settings.json" && echo "Unlinked: VS Code settings.json"
    [ -L "$vscode_dir/keybindings.json" ] && rm "$vscode_dir/keybindings.json" && echo "Unlinked: VS Code keybindings.json"

    echo ""
    echo "Config unlinking complete!"
}

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

EXAMPLES:
  dotfiles_inject         # First-time setup
  dotfiles_link_configs   # Link git, ssh configs
  dotfiles_update         # Update from git
  dotfiles_status         # Check current state
EOF
}

# Aliases
alias df-inject='dotfiles_inject'
alias df-eject='dotfiles_eject'
alias df-update='dotfiles_update'
alias df-reload='dotfiles_reload'
alias df-status='dotfiles_status'
alias df-link='dotfiles_link_configs'
alias df-unlink='dotfiles_unlink_configs'
alias df-help='dotfiles_help'
