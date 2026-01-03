#!/bin/sh
# components/neovim/functions.sh - Neovim configuration management

# Default location for cloned config repos
NVIM_REPOS_DIR="${HOME}/Repos/personal"

# =============================================================================
# Setup Functions
# =============================================================================

# Setup Neovim with external config repo
nvim_setup() {
    local config_dir="${XDG_CONFIG_HOME:-$HOME/.config}/nvim"
    local repo_url=""
    local repo_name=""
    local repo_path=""

    echo "=== Neovim Configuration Setup ==="
    echo ""

    # Check if config already exists
    if [ -L "$config_dir" ]; then
        local current_target
        current_target=$(readlink "$config_dir")
        echo "Neovim config already linked to: $current_target"
        printf "Reconfigure? [y/N] "
        read -r response
        [ "$response" != "y" ] && [ "$response" != "Y" ] && return 0
        rm "$config_dir"
    elif [ -d "$config_dir" ]; then
        echo "Existing config directory found at $config_dir"
        printf "Backup and replace? [y/N] "
        read -r response
        if [ "$response" = "y" ] || [ "$response" = "Y" ]; then
            mv "$config_dir" "${config_dir}.backup.$(date +%Y%m%d%H%M%S)"
            echo "Backed up to ${config_dir}.backup.*"
        else
            return 0
        fi
    fi

    echo ""
    echo "Enter your Neovim config GitHub repo URL"
    echo "Examples:"
    echo "  https://github.com/username/nvim-config"
    echo "  git@github.com:username/kickstart.nvim.git"
    echo ""
    printf "Repo URL (or 'skip' to skip): "
    read -r repo_url

    if [ "$repo_url" = "skip" ] || [ -z "$repo_url" ]; then
        echo "Skipping Neovim config setup"
        return 0
    fi

    # Extract repo name from URL
    repo_name=$(basename "$repo_url" .git)
    repo_path="${NVIM_REPOS_DIR}/${repo_name}"

    # Ensure repos directory exists
    mkdir -p "$NVIM_REPOS_DIR"

    # Clone or update repo
    if [ -d "$repo_path" ]; then
        echo "Repo already exists at $repo_path"
        printf "Pull latest changes? [Y/n] "
        read -r response
        if [ "$response" != "n" ] && [ "$response" != "N" ]; then
            (cd "$repo_path" && git pull)
        fi
    else
        echo "Cloning $repo_url to $repo_path..."
        git clone "$repo_url" "$repo_path"
        if [ $? -ne 0 ]; then
            echo "Failed to clone repository"
            return 1
        fi
    fi

    # Create symlink
    ln -s "$repo_path" "$config_dir"
    echo ""
    echo "✓ Neovim config linked: $config_dir -> $repo_path"
    echo ""
    echo "Run 'nvim' to start Neovim and install plugins"
}

# Update Neovim config repo
nvim_update() {
    local config_dir="${XDG_CONFIG_HOME:-$HOME/.config}/nvim"

    if [ -L "$config_dir" ]; then
        local repo_path
        repo_path=$(readlink "$config_dir")
        echo "Updating Neovim config at $repo_path..."
        (cd "$repo_path" && git pull)
    elif [ -d "$config_dir" ]; then
        if [ -d "$config_dir/.git" ]; then
            echo "Updating Neovim config..."
            (cd "$config_dir" && git pull)
        else
            echo "Neovim config is not a git repository"
            return 1
        fi
    else
        echo "No Neovim config found. Run nvim_setup first."
        return 1
    fi
}

# Clean Neovim cache and state
nvim_clean() {
    local data_dir="${XDG_DATA_HOME:-$HOME/.local/share}/nvim"
    local cache_dir="${XDG_CACHE_HOME:-$HOME/.cache}/nvim"
    local state_dir="${XDG_STATE_HOME:-$HOME/.local/state}/nvim"

    echo "This will remove Neovim data, cache, and state directories:"
    echo "  - $data_dir"
    echo "  - $cache_dir"
    echo "  - $state_dir"
    echo ""
    printf "Continue? [y/N] "
    read -r response

    if [ "$response" = "y" ] || [ "$response" = "Y" ]; then
        rm -rf "$data_dir" "$cache_dir" "$state_dir"
        echo "Cleaned. Plugins will be reinstalled on next nvim launch."
    else
        echo "Cancelled."
    fi
}

# Check Neovim health
nvim_health() {
    local config_dir="${XDG_CONFIG_HOME:-$HOME/.config}/nvim"

    echo "=== Neovim Health Check ==="
    echo ""

    # Check nvim version
    if command -v nvim >/dev/null 2>&1; then
        echo "Neovim: $(nvim --version | head -1)"
    else
        echo "Neovim: NOT INSTALLED"
        return 1
    fi

    # Check config
    echo ""
    if [ -L "$config_dir" ]; then
        local target
        target=$(readlink "$config_dir")
        echo "Config: $config_dir -> $target"
        if [ -d "$target" ]; then
            echo "Status: OK"
        else
            echo "Status: BROKEN LINK (target doesn't exist)"
        fi
    elif [ -d "$config_dir" ]; then
        echo "Config: $config_dir (direct directory)"
        echo "Status: OK"
    else
        echo "Config: NOT FOUND"
        echo "Run 'nvim_setup' to configure"
    fi

    # Check for init file
    echo ""
    if [ -f "$config_dir/init.lua" ]; then
        echo "Init file: init.lua"
    elif [ -f "$config_dir/init.vim" ]; then
        echo "Init file: init.vim"
    else
        echo "Init file: NOT FOUND"
    fi

    # Show plugin manager if detectable
    echo ""
    if [ -d "$config_dir/lua/lazy" ] || grep -q "lazy" "$config_dir/init.lua" 2>/dev/null; then
        echo "Plugin manager: lazy.nvim (detected)"
    elif [ -d "${XDG_DATA_HOME:-$HOME/.local/share}/nvim/site/pack/packer" ]; then
        echo "Plugin manager: packer.nvim (detected)"
    fi
}
