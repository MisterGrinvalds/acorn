#!/bin/sh
# VS Code Integration

# VS Code aliases
alias c='code .'
alias cg='code --goto'
alias cn='code --new-window'
alias ca='code --add'

# VS Code workspace helpers
workspace() {
    if [ -z "$1" ]; then
        echo "Usage: workspace <workspace-name>"
        return 1
    fi
    code "$HOME/.vscode/workspaces/$1.code-workspace"
}

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

# Quick open VS Code with specific file types (requires fzf)
cpy() {
    find . -name "*.py" | head -20 | fzf | xargs code
}
cgo() {
    find . -name "*.go" | head -20 | fzf | xargs code
}
cts() {
    find . -name "*.ts" -o -name "*.js" | head -20 | fzf | xargs code
}

# VS Code extensions management helpers
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
