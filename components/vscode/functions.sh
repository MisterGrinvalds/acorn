#!/bin/sh
# components/vscode/functions.sh - VS Code project helpers

# =============================================================================
# Workspace Management
# =============================================================================

# Open VS Code workspace
workspace() {
    if [ -z "$1" ]; then
        echo "Usage: workspace <workspace-name>"
        return 1
    fi
    code "$HOME/.vscode/workspaces/$1.code-workspace"
}

# List available workspaces
workspaces() {
    local ws_dir="$HOME/.vscode/workspaces"
    if [ -d "$ws_dir" ]; then
        ls -1 "$ws_dir" | sed 's/.code-workspace$//'
    else
        echo "No workspaces directory found at $ws_dir"
    fi
}

# =============================================================================
# Project Initialization
# =============================================================================

# Create new VS Code project with common settings
newproject() {
    if [ -z "$1" ]; then
        echo "Usage: newproject <project-name> [language]"
        return 1
    fi

    local project_name="$1"
    local language="${2:-general}"

    mkdir -p "$project_name"
    cd "$project_name" || return 1

    # Create .vscode directory with common settings
    mkdir -p .vscode

    # Create settings.json based on language
    case "$language" in
        "python"|"py")
            cat > .vscode/settings.json << 'EOF'
{
    "python.defaultInterpreterPath": "./.venv/bin/python",
    "python.terminal.activateEnvironment": true,
    "editor.formatOnSave": true,
    "python.testing.pytestEnabled": true,
    "python.testing.pytestArgs": ["."]
}
EOF
            ;;
        "go"|"golang")
            cat > .vscode/settings.json << 'EOF'
{
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "editor.formatOnSave": true,
    "go.testFlags": ["-v"],
    "go.testTimeout": "30s"
}
EOF
            ;;
        "typescript"|"ts"|"node")
            cat > .vscode/settings.json << 'EOF'
{
    "typescript.preferences.importModuleSpecifier": "relative",
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "esbenp.prettier-vscode",
    "typescript.updateImportsOnFileMove.enabled": "always"
}
EOF
            ;;
        *)
            cat > .vscode/settings.json << 'EOF'
{
    "editor.formatOnSave": true,
    "files.trimTrailingWhitespace": true,
    "files.insertFinalNewline": true,
    "editor.rulers": [80, 120]
}
EOF
            ;;
    esac

    echo "VS Code project '$project_name' created with $language settings"
    code .
}

# =============================================================================
# Quick Open by Language (requires fzf)
# =============================================================================

# Open Python files with fzf
cpy() {
    if ! command -v fzf >/dev/null 2>&1; then
        echo "fzf required"
        return 1
    fi
    find . -name "*.py" -type f 2>/dev/null | head -50 | fzf --preview 'head -50 {}' | xargs -r code
}

# Open Go files with fzf
cgo() {
    if ! command -v fzf >/dev/null 2>&1; then
        echo "fzf required"
        return 1
    fi
    find . -name "*.go" -type f 2>/dev/null | head -50 | fzf --preview 'head -50 {}' | xargs -r code
}

# Open TypeScript/JavaScript files with fzf
cts() {
    if ! command -v fzf >/dev/null 2>&1; then
        echo "fzf required"
        return 1
    fi
    find . \( -name "*.ts" -o -name "*.tsx" -o -name "*.js" -o -name "*.jsx" \) -type f 2>/dev/null | head -50 | fzf --preview 'head -50 {}' | xargs -r code
}

# =============================================================================
# Setup (Install vs Config)
# =============================================================================

# Full VS Code setup (extensions + config)
vscode_setup() {
    echo "VS Code Setup"
    echo "============="
    echo "1) Install extensions only"
    echo "2) Sync config only"
    echo "3) Both (full setup)"
    echo ""
    printf "Choice [1-3]: "
    read -r choice

    case "$choice" in
        1) vscode_install_extensions ;;
        2) vscode_sync_config ;;
        3) vscode_install_extensions && vscode_sync_config ;;
        *) echo "Invalid choice" ;;
    esac
}

# Install extensions from dotfiles list
vscode_install_extensions() {
    local ext_file="${DOTFILES_ROOT:-$HOME/.config/dotfiles}/config/vscode/extensions.txt"

    if [ ! -f "$ext_file" ]; then
        echo "No extensions list found at $ext_file"
        echo "Run vscode_export_extensions to create one"
        return 1
    fi

    echo "Installing VS Code extensions..."
    while IFS= read -r ext; do
        [ -z "$ext" ] && continue
        echo "  → $ext"
        code --install-extension "$ext" --force 2>/dev/null
    done < "$ext_file"
    echo "Extensions installed!"
}

# Sync settings and keybindings from dotfiles
vscode_sync_config() {
    local src_dir="${DOTFILES_ROOT:-$HOME/.config/dotfiles}/config/vscode"
    local dest_dir="$HOME/Library/Application Support/Code/User"

    # Linux path
    if [ "$(uname)" != "Darwin" ]; then
        dest_dir="$HOME/.config/Code/User"
    fi

    if [ ! -d "$dest_dir" ]; then
        echo "VS Code user directory not found: $dest_dir"
        echo "Is VS Code installed?"
        return 1
    fi

    echo "Syncing VS Code config..."

    # Backup existing
    if [ -f "$dest_dir/settings.json" ]; then
        cp "$dest_dir/settings.json" "$dest_dir/settings.json.backup"
        echo "  → Backed up existing settings.json"
    fi
    if [ -f "$dest_dir/keybindings.json" ]; then
        cp "$dest_dir/keybindings.json" "$dest_dir/keybindings.json.backup"
        echo "  → Backed up existing keybindings.json"
    fi

    # Copy new config
    cp "$src_dir/settings.json" "$dest_dir/settings.json"
    echo "  → Synced settings.json"

    if [ -f "$src_dir/keybindings.json" ]; then
        cp "$src_dir/keybindings.json" "$dest_dir/keybindings.json"
        echo "  → Synced keybindings.json"
    fi

    echo "VS Code config synced!"
}

# =============================================================================
# Extensions Management
# =============================================================================

# Install essential VS Code extensions (quick setup)
vscode_install_essentials() {
    echo "Installing essential VS Code extensions..."

    code --install-extension ms-python.python
    code --install-extension golang.go
    code --install-extension github.vscode-pull-request-github
    code --install-extension eamodio.gitlens
    code --install-extension ms-kubernetes-tools.vscode-kubernetes-tools
    code --install-extension ms-azuretools.vscode-docker
    code --install-extension catppuccin.catppuccin-vsc
    code --install-extension catppuccin.catppuccin-vsc-icons

    echo "Essential extensions installed!"
}

# Export extensions list
vscode_export_extensions() {
    code --list-extensions > vscode-extensions.txt
    echo "Extensions exported to vscode-extensions.txt"
}

# Import extensions from list
vscode_import_extensions() {
    local file="${1:-vscode-extensions.txt}"
    if [ ! -f "$file" ]; then
        echo "File not found: $file"
        return 1
    fi

    while IFS= read -r ext; do
        code --install-extension "$ext"
    done < "$file"
}
