# VS Code Integration

# VS Code aliases
alias code='code'
alias c='code .'                    # Open current directory in VS Code
alias cg='code --goto'              # Open file at specific line
alias cn='code --new-window'        # Open new VS Code window
alias ca='code --add'               # Add folder to current workspace

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
    "python.formatting.provider": "black",
    "python.sortImports.args": ["--profile", "black"],
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
        "source.organizeImports": true
    },
    "python.testing.pytestEnabled": true,
    "python.testing.unittestEnabled": false,
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
    "editor.codeActionsOnSave": {
        "source.organizeImports": true
    },
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
    "editor.codeActionsOnSave": {
        "source.fixAll.eslint": true,
        "source.organizeImports": true
    },
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
    "editor.rulers": [80, 120],
    "editor.wordWrap": "bounded",
    "editor.wordWrapColumn": 120
}
EOF
            ;;
    esac
    
    # Create launch.json for debugging
    cat > .vscode/launch.json << 'EOF'
{
    "version": "0.2.0",
    "configurations": []
}
EOF
    
    echo "VS Code project '$project_name' created with $language settings"
    code .
}

# Quick open VS Code with specific file types
cpy() { find . -name "*.py" | head -20 | fzf | xargs code; }
cgo() { find . -name "*.go" | head -20 | fzf | xargs code; }
cts() { find . -name "*.ts" -o -name "*.js" | head -20 | fzf | xargs code; }

# VS Code extensions management helpers
vscode_install_essentials() {
    echo "Installing essential VS Code extensions..."
    
    # General
    code --install-extension ms-vscode.vscode-json
    code --install-extension redhat.vscode-yaml
    code --install-extension ms-vscode.vscode-typescript-next
    
    # Python
    code --install-extension ms-python.python
    code --install-extension ms-python.pylint
    code --install-extension ms-python.black-formatter
    code --install-extension ms-python.isort
    
    # Go
    code --install-extension golang.go
    
    # Git & GitHub
    code --install-extension github.vscode-pull-request-github
    code --install-extension eamodio.gitlens
    
    # Kubernetes
    code --install-extension ms-kubernetes-tools.vscode-kubernetes-tools
    
    # Docker
    code --install-extension ms-azuretools.vscode-docker
    
    echo "Essential extensions installed!"
}