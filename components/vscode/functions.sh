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
# Extensions Management
# =============================================================================

# Install essential VS Code extensions
vscode_install_essentials() {
    echo "Installing essential VS Code extensions..."

    code --install-extension ms-python.python
    code --install-extension golang.go
    code --install-extension github.vscode-pull-request-github
    code --install-extension eamodio.gitlens
    code --install-extension ms-kubernetes-tools.vscode-kubernetes-tools
    code --install-extension ms-azuretools.vscode-docker

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
