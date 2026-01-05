// Package shell provides shell integration generation and injection.
package shell

// RegisterAllComponents registers all known components with the manager.
func RegisterAllComponents(m *Manager) {
	m.RegisterComponent(GoComponent())
	m.RegisterComponent(VSCodeComponent())
	m.RegisterComponent(ToolsComponent())
	m.RegisterComponent(PythonComponent())
	m.RegisterComponent(TmuxComponent())
	m.RegisterComponent(ClaudeComponent())
	m.RegisterComponent(CloudFlareComponent())
	m.RegisterComponent(SecretsComponent())
	m.RegisterComponent(DatabaseComponent())
	m.RegisterComponent(FzfComponent())
}

// GoComponent returns the Go shell integration component.
func GoComponent() *Component {
	return &Component{
		Name:        "go",
		Description: "Go development environment and helpers",
		Env: `# Go environment setup
export GOPATH="${GOPATH:-$HOME/go}"
export GOBIN="$GOPATH/bin"
export GO111MODULE=on

# Add Go binary paths to PATH
case ":$PATH:" in
    *":$GOBIN:"*) ;;
    *) export PATH="$GOBIN:$PATH" ;;
esac

case ":$PATH:" in
    *":/usr/local/go/bin:"*) ;;
    *) export PATH="/usr/local/go/bin:$PATH" ;;
esac

# macOS Homebrew Go location
if [ -d "/opt/homebrew/opt/go/bin" ]; then
    case ":$PATH:" in
        *":/opt/homebrew/opt/go/bin:"*) ;;
        *) export PATH="/opt/homebrew/opt/go/bin:$PATH" ;;
    esac
fi
`,
		Aliases: `# Build and run
alias gob='go build'
alias gor='go run'
alias goi='go install'

# Testing
alias got='go test'
alias gotv='go test -v'
alias gotc='go test -cover'

# Module management
alias gom='go mod'
alias gomi='go mod init'
alias gomt='go mod tidy'
alias gomd='go mod download'

# Dependencies
alias gog='go get'
alias gou='go get -u'

# Code quality
alias gof='go fmt ./...'
alias gov='go vet ./...'

# Info
alias gover='go version'
alias goenv='go env'
`,
		Functions: `# Initialize new Go project (wrapper for acorn go new)
gonew() {
    if [ -z "$1" ]; then
        echo "Usage: gonew <module-name>"
        return 1
    fi
    acorn go new "$1" && cd "$1"
}

# Create new Cobra CLI project (wrapper for acorn go cobra new)
cobranew() {
    if [ -z "$1" ]; then
        echo "Usage: cobranew <app-name>"
        return 1
    fi
    acorn go cobra new "$1" && cd "$1"
}

# Add command to Cobra project (wrapper for acorn go cobra add)
cobradd() {
    if [ -z "$1" ]; then
        echo "Usage: cobradd <command-name>"
        return 1
    fi
    acorn go cobra add "$1"
}

# Run tests with optional pattern filter (wrapper for acorn go test)
gotest() {
    acorn go test "$@"
}

# Run tests with coverage report (wrapper for acorn go cover)
gotestcover() {
    acorn go cover
}

# Run benchmarks (wrapper for acorn go bench)
gobench() {
    acorn go bench "$@"
}

# Build for multiple platforms (wrapper for acorn go build-all)
gobuildall() {
    acorn go build-all "${1:-app}"
}

# Clean build artifacts (wrapper for acorn go clean)
goclean() {
    acorn go clean
}
`,
	}
}

// VSCodeComponent returns the VS Code shell integration component.
func VSCodeComponent() *Component {
	return &Component{
		Name:        "vscode",
		Description: "VS Code integration and project helpers",
		Aliases: `# Open current directory
alias c='code .'

# Open with options
alias cg='code --goto'
alias cn='code --new-window'
alias ca='code --add'
alias cr='code --reuse-window'

# Diff mode
alias cdiff='code --diff'

# Extensions
alias cext='code --list-extensions'
alias cexti='code --install-extension'
alias cextu='code --uninstall-extension'
`,
		Functions: `# Open VS Code workspace (wrapper for acorn vscode workspace)
workspace() {
    if [ -z "$1" ]; then
        echo "Usage: workspace <workspace-name>"
        return 1
    fi
    acorn vscode workspace "$1"
}

# List available workspaces (wrapper for acorn vscode workspaces)
workspaces() {
    acorn vscode workspaces
}

# Create new VS Code project with common settings (wrapper for acorn vscode project new)
newproject() {
    if [ -z "$1" ]; then
        echo "Usage: newproject <project-name> [language]"
        return 1
    fi
    acorn vscode project new "$@"
}

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
        1) acorn vscode ext install ;;
        2) acorn vscode config sync ;;
        3) acorn vscode ext install && acorn vscode config sync ;;
        *) echo "Invalid choice" ;;
    esac
}

# Install extensions from dotfiles list (wrapper for acorn vscode ext install)
vscode_install_extensions() {
    acorn vscode ext install "$@"
}

# Sync settings and keybindings from dotfiles (wrapper for acorn vscode config sync)
vscode_sync_config() {
    acorn vscode config sync "$@"
}

# Install essential VS Code extensions (wrapper for acorn vscode ext essentials)
vscode_install_essentials() {
    acorn vscode ext essentials
}

# Export extensions list (wrapper for acorn vscode ext export)
vscode_export_extensions() {
    acorn vscode ext export "${1:-vscode-extensions.txt}"
}

# Import extensions from list (wrapper for acorn vscode ext install)
vscode_import_extensions() {
    acorn vscode ext install "${1:-vscode-extensions.txt}"
}
`,
	}
}

// ToolsComponent returns the tools shell integration component.
func ToolsComponent() *Component {
	return &Component{
		Name:        "tools",
		Description: "System tools management and version checking",
		Aliases: `# Automation framework shortcuts
alias tools='tools_status'
alias tools-list='acorn tools list'
alias tools-check='acorn tools check'
alias tools-update='acorn tools update'
alias versions='quick_versions'
alias system-update='smart_update'
alias whichx='which_enhanced'
`,
		Functions: `# Show all tool versions and status (wrapper for acorn tools status)
tools_status() {
    acorn tools status "$@"
}

# Check if specific tools are installed (wrapper for acorn tools check)
check_tool() {
    acorn tools check "$@"
}

# Show outdated tools (wrapper for acorn tools outdated)
tools_outdated() {
    acorn tools outdated "$@"
}

# Update tools via brew/go/npm (wrapper for acorn tools update)
update_tools() {
    acorn tools update "$@"
}

# Quick check for common tools
tools_quick() {
    acorn tools status --quick
}

# Quick version checks for common tools
quick_versions() {
    echo "=== Quick Version Check ==="
    echo "System Tools:"
    command -v git >/dev/null && echo "  git: $(git --version 2>/dev/null | head -1)" || echo "  git: Not installed"
    command -v curl >/dev/null && echo "  curl: $(curl --version 2>/dev/null | head -1)" || echo "  curl: Not installed"
    command -v jq >/dev/null && echo "  jq: $(jq --version 2>/dev/null)" || echo "  jq: Not installed"
    echo "Languages:"
    command -v go >/dev/null && echo "  go: $(go version 2>/dev/null)" || echo "  go: Not installed"
    command -v node >/dev/null && echo "  node: $(node --version 2>/dev/null)" || echo "  node: Not installed"
    command -v python3 >/dev/null && echo "  python3: $(python3 --version 2>/dev/null)" || echo "  python3: Not installed"
    echo "Cloud Tools:"
    command -v aws >/dev/null && echo "  aws: $(aws --version 2>/dev/null | head -1)" || echo "  aws: Not installed"
    command -v kubectl >/dev/null && echo "  kubectl: $(kubectl version --client --short 2>/dev/null)" || echo "  kubectl: Not installed"
    echo "Development Tools:"
    command -v docker >/dev/null && echo "  docker: $(docker --version 2>/dev/null)" || echo "  docker: Not installed"
    command -v gh >/dev/null && echo "  gh: $(gh --version 2>/dev/null | head -1)" || echo "  gh: Not installed"
}

# Smart package manager detection and update
smart_update() {
    echo "Smart system update..."
    if command -v brew >/dev/null 2>&1; then
        echo "Updating Homebrew packages..."
        brew update && brew upgrade
    elif command -v apt-get >/dev/null 2>&1; then
        echo "Updating apt packages..."
        sudo apt-get update && sudo apt-get upgrade
    elif command -v dnf >/dev/null 2>&1; then
        echo "Updating dnf packages..."
        sudo dnf upgrade
    elif command -v pacman >/dev/null 2>&1; then
        echo "Updating pacman packages..."
        sudo pacman -Syu
    else
        echo "No supported package manager found"
        return 1
    fi
    echo "Smart update completed!"
}

# Enhanced which command with more info
which_enhanced() {
    local tool="$1"
    if [ -z "$tool" ]; then
        echo "Usage: which_enhanced <command>"
        return 1
    fi
    if command -v "$tool" >/dev/null 2>&1; then
        echo "$tool found:"
        echo "  Location: $(which "$tool")"
        if "$tool" --version >/dev/null 2>&1; then
            echo "  Version: $($tool --version 2>/dev/null | head -1)"
        elif "$tool" version >/dev/null 2>&1; then
            echo "  Version: $($tool version 2>/dev/null | head -1)"
        else
            echo "  Version: Unknown"
        fi
    else
        echo "$tool not found"
    fi
}
`,
	}
}

// PythonComponent returns the Python shell integration component.
func PythonComponent() *Component {
	return &Component{
		Name:        "python",
		Description: "Python development with UV package manager",
		Env: `# Python environment setup
export ENVS_LOCATION="${ENVS_LOCATION:-$HOME/.virtualenvs}"
export PYTHONSTARTUP="${XDG_CONFIG_HOME:-$HOME/.config}/python/pythonrc"
export IPYTHONDIR="${XDG_CONFIG_HOME:-$HOME/.config}/ipython"
export UV_CACHE_DIR="${XDG_CACHE_HOME:-$HOME/.cache}/uv"
`,
		Aliases: `# Python shortcuts
alias py='python3'
alias py2='python2'
alias pip='pip3'

# UV package manager (direct uv commands)
alias uvs='uv sync'
alias uva='uv add'
alias uvr='uv remove'
alias uvx='uv run'
alias uvi='uv init'

# Virtual environment shortcuts
alias mkv='mkvenv'
alias dv='dvenv'
`,
		Functions: `# Create Python virtual environment (wrapper for acorn python venv new + activate)
mkvenv() {
    local name="${1:-.venv}"
    acorn python venv new "$name"
    # Activate the environment after creation
    if [ -f "$name/bin/activate" ]; then
        . "$name/bin/activate"
        echo "Activated: $name"
    fi
}

# Activate a Python virtual environment
venv() {
    local name="${1:-.venv}"
    if [ -f "$name/bin/activate" ]; then
        . "$name/bin/activate"
    elif [ -n "$ENVS_LOCATION" ] && [ -f "$ENVS_LOCATION/$name/bin/activate" ]; then
        . "$ENVS_LOCATION/$name/bin/activate"
    else
        echo "Virtual environment not found: $name"
        return 1
    fi
}

# Deactivate current virtual environment
dvenv() {
    if [ -n "$VIRTUAL_ENV" ]; then
        deactivate
    else
        echo "No active virtual environment"
    fi
}

# List virtual environments (wrapper for acorn python venv list)
lsvenv() {
    acorn python venv list "$@"
}

# Initialize UV project (wrapper for acorn python init)
pyinit() {
    acorn python init "$@"
}

# Sync dependencies (wrapper for acorn python sync)
pysync() {
    acorn python sync
}

# Add packages (wrapper for acorn python add)
pyadd() {
    acorn python add "$@"
}

# Remove packages (wrapper for acorn python remove)
pyrm() {
    acorn python remove "$@"
}

# Run command in project environment (wrapper for acorn python run)
pyrun() {
    acorn python run "$@"
}

# Show Python environment (wrapper for acorn python env)
pyenv() {
    acorn python env "$@"
}

# FastAPI development environment (wrapper for acorn python fastapi)
fastapi_env() {
    local name="${1:-.venv}"
    acorn python fastapi "$name"
    # Activate the environment after setup
    if [ -f "$name/bin/activate" ]; then
        . "$name/bin/activate"
        echo "Activated: $name"
    fi
}

# Install IPython (wrapper for acorn python setup ipython)
setup_ipython() {
    acorn python setup ipython
}

# Install development tools (wrapper for acorn python setup devtools)
setup_devtools() {
    acorn python setup devtools
}
`,
	}
}

// TmuxComponent returns the tmux shell integration component.
func TmuxComponent() *Component {
	return &Component{
		Name:        "tmux",
		Description: "Tmux session management with TPM and smug",
		Env: `# Tmux environment setup (XDG-compliant)
export TMUX_CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/tmux"
export TMUX_PLUGIN_DIR="${TMUX_CONFIG_DIR}/plugins"
export TMUX_TPM_DIR="${TMUX_PLUGIN_DIR}/tpm"
export TMUX_CONF="${TMUX_CONFIG_DIR}/tmux.conf"

# Smug session management with git sync
export SMUG_REPO="https://github.com/MisterGrinvalds/fmux.git"
export SMUG_REPO_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/smug-sessions"
export SMUG_CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/smug"
`,
		Aliases: `# Basic tmux shortcuts
alias tm='tmux'
alias tma='tmux attach-session'
alias tmat='tmux attach-session -t'
alias tmn='tmux new-session'
alias tmns='tmux new-session -s'
alias tml='tmux list-sessions'
alias tmk='tmux kill-session -t'
alias tmka='tmux kill-server'

# Attach to last session or create new
alias tmx='tmux attach-session 2>/dev/null || tmux new-session'

# Quick session access
alias tm0='tmux attach-session -t 0'
alias tm1='tmux attach-session -t 1'
alias tmdev='tmux attach-session -t dev 2>/dev/null || dev_session'
`,
		Functions: `# =============================================================================
# Development Sessions (must stay in shell for tmux attach)
# =============================================================================

# Create development session with multiple panes
dev_session() {
    local session_name="${1:-dev}"
    tmux new-session -d -s "$session_name"
    tmux split-window -h -t "$session_name"
    tmux split-window -v -t "$session_name:0.1"
    tmux send-keys -t "$session_name:0.0" 'echo "Editor pane"' Enter
    tmux send-keys -t "$session_name:0.1" 'echo "Main terminal"' Enter
    tmux send-keys -t "$session_name:0.2" 'echo "Logs/monitoring pane"' Enter
    tmux select-pane -t "$session_name:0.0"
    tmux attach-session -t "$session_name"
}

# Create Kubernetes development session
k8s_session() {
    local session_name="k8s"
    tmux new-session -d -s "$session_name"
    tmux new-window -t "$session_name" -n "kubectl"
    tmux new-window -t "$session_name" -n "logs"
    tmux new-window -t "$session_name" -n "k9s"
    tmux send-keys -t "$session_name:kubectl" 'echo "kubectl commands"' Enter
    tmux split-window -h -t "$session_name:logs"
    tmux send-keys -t "$session_name:logs.0" 'echo "App logs"' Enter
    tmux send-keys -t "$session_name:logs.1" 'echo "System logs"' Enter
    tmux send-keys -t "$session_name:k9s" 'k9s' Enter
    tmux select-window -t "$session_name:1"
    tmux attach-session -t "$session_name"
}

# Create project session (auto-detects project type)
project_session() {
    local project_path="${1:-.}"
    local session_name
    session_name=$(basename "$(realpath "$project_path")")
    cd "$project_path" || return 1
    tmux new-session -d -s "$session_name" -c "$project_path"
    if [ -f "go.mod" ]; then
        tmux send-keys -t "$session_name" 'echo "Go project detected"' Enter
        tmux new-window -t "$session_name" -n "test" -c "$project_path"
        tmux new-window -t "$session_name" -n "run" -c "$project_path"
    elif [ -f "requirements.txt" ] || [ -f "pyproject.toml" ]; then
        tmux send-keys -t "$session_name" 'echo "Python project detected"' Enter
        tmux new-window -t "$session_name" -n "test" -c "$project_path"
        tmux new-window -t "$session_name" -n "server" -c "$project_path"
    elif [ -f "package.json" ]; then
        tmux send-keys -t "$session_name" 'echo "Node.js project detected"' Enter
        tmux new-window -t "$session_name" -n "dev" -c "$project_path"
        tmux new-window -t "$session_name" -n "test" -c "$project_path"
    fi
    tmux new-window -t "$session_name" -n "editor" -c "$project_path"
    tmux send-keys -t "$session_name:editor" 'code .' Enter
    tmux select-window -t "$session_name:1"
    tmux attach-session -t "$session_name"
}

# =============================================================================
# Session Management with FZF (must stay in shell for fzf/tmux attach)
# =============================================================================

# Quick session switcher using fzf
tswitch() {
    if ! command -v fzf >/dev/null 2>&1; then
        echo "fzf is required for tswitch"
        return 1
    fi
    local session
    session=$(tmux list-sessions -F '#S' | fzf --prompt="Switch to session: ")
    if [ -n "$session" ]; then
        tmux attach-session -t "$session"
    fi
}

# Kill session with fzf selection
tkill() {
    if ! command -v fzf >/dev/null 2>&1; then
        echo "fzf is required for tkill"
        return 1
    fi
    local session
    session=$(tmux list-sessions -F '#S' | fzf --prompt="Kill session: ")
    if [ -n "$session" ]; then
        tmux kill-session -t "$session"
        echo "Killed session: $session"
    fi
}

# Attach to session or create new with name
tmux_attach() {
    local session="${1:-main}"
    if tmux has-session -t "$session" 2>/dev/null; then
        tmux attach-session -t "$session"
    else
        echo "Creating new session: $session"
        tmux new-session -s "$session"
    fi
}

# =============================================================================
# TPM Management (wrappers for acorn tmux tpm commands)
# =============================================================================

# Install TPM (wrapper for acorn tmux tpm install)
tmux_install_tpm() {
    acorn tmux tpm install
}

# Update TPM (wrapper for acorn tmux tpm update)
tmux_update_tpm() {
    acorn tmux tpm update
}

# Install all plugins (wrapper for acorn tmux tpm plugins-install)
tmux_install_plugins() {
    acorn tmux tpm plugins-install
}

# Update all plugins (wrapper for acorn tmux tpm plugins-update)
tmux_update_plugins() {
    acorn tmux tpm plugins-update
}

# =============================================================================
# Configuration Management
# =============================================================================

# Edit tmux config (must stay in shell for editor)
tmux_config() {
    local config="${TMUX_CONF:-$HOME/.config/tmux/tmux.conf}"
    if [ ! -f "$config" ]; then
        echo "Tmux config not found: $config"
        return 1
    fi
    ${EDITOR:-vim} "$config"
    echo "Config saved. Run 'tmux_reload' or prefix + r to reload."
}

# Reload tmux config (wrapper for acorn tmux config reload)
tmux_reload() {
    acorn tmux config reload
}

# Show tmux info (wrapper for acorn tmux info)
tmux_info() {
    acorn tmux info "$@"
}

# =============================================================================
# Smug Session Management (wrappers for acorn tmux smug commands)
# =============================================================================

# List available smug sessions (wrapper for acorn tmux smug list)
smug_list() {
    acorn tmux smug list "$@"
}

# Start smug session with fzf selection (must stay in shell for fzf/attach)
smug_start() {
    local session="$1"
    local smug_dir="${SMUG_CONFIG_DIR:-$HOME/.config/smug}"
    if ! command -v smug >/dev/null 2>&1; then
        echo "smug not installed. Run: acorn tmux smug install"
        return 1
    fi
    # If no session given, use fzf to select
    if [ -z "$session" ]; then
        if ! command -v fzf >/dev/null 2>&1; then
            echo "Usage: smug_start <session_name>"
            smug_list
            return 1
        fi
        session=$(ls "$smug_dir"/*.yml 2>/dev/null | xargs -I {} basename {} .yml | fzf --prompt="Start session: ")
        [ -z "$session" ] && return 0
    fi
    shift 2>/dev/null
    smug start "$session" "$@"
}

# Stop smug session with fzf selection (must stay in shell for fzf)
smug_stop() {
    local session="$1"
    if ! command -v smug >/dev/null 2>&1; then
        echo "smug not installed"
        return 1
    fi
    if [ -z "$session" ]; then
        if ! command -v fzf >/dev/null 2>&1; then
            echo "Usage: smug_stop <session_name>"
            return 1
        fi
        session=$(tmux list-sessions -F '#S' 2>/dev/null | fzf --prompt="Stop session: ")
        [ -z "$session" ] && return 0
    fi
    smug stop "$session"
}

# Create a new smug session config (wrapper for acorn tmux smug new)
smug_new() {
    acorn tmux smug new "$@"
}

# Edit existing smug config (must stay in shell for editor/fzf)
smug_edit() {
    local name="$1"
    local smug_dir="${SMUG_CONFIG_DIR:-$HOME/.config/smug}"
    if ! command -v smug >/dev/null 2>&1; then
        echo "smug not installed"
        return 1
    fi
    if [ -z "$name" ]; then
        if command -v fzf >/dev/null 2>&1; then
            name=$(ls "$smug_dir"/*.yml 2>/dev/null | xargs -I {} basename {} .yml | fzf --prompt="Edit config: ")
            [ -z "$name" ] && return 0
        else
            echo "Usage: smug_edit <session_name>"
            return 1
        fi
    fi
    local config_file="$smug_dir/$name.yml"
    if [ ! -f "$config_file" ]; then
        echo "Config not found: $config_file"
        return 1
    fi
    ${EDITOR:-vim} "$config_file"
}

# Install smug (wrapper for acorn tmux smug install)
smug_install() {
    acorn tmux smug install
}

# Link smug configs (wrapper for acorn tmux smug link)
smug_link_configs() {
    acorn tmux smug link
}

# =============================================================================
# Smug Git Sync (wrappers for acorn tmux smug repo commands)
# =============================================================================

# Initialize smug sessions repo (wrapper for acorn tmux smug repo-init)
smug_repo_init() {
    acorn tmux smug repo-init
}

# Show smug repo status (wrapper for acorn tmux smug status)
smug_status() {
    acorn tmux smug status "$@"
}

# Pull latest sessions (wrapper for acorn tmux smug pull)
smug_pull() {
    acorn tmux smug pull
}

# Commit and push changes (wrapper for acorn tmux smug push)
smug_push() {
    acorn tmux smug push "$@"
}

# Full sync (wrapper for acorn tmux smug sync)
smug_sync() {
    acorn tmux smug sync
}

# =============================================================================
# Window Alerts (must stay in shell for tmux commands)
# =============================================================================

# Set alert on current window (default: red)
tmux_alert() {
    if [ -z "$TMUX" ]; then
        echo "Error: Not in a tmux session"
        return 1
    fi
    local window_id
    window_id=$(tmux display-message -p '#{window_id}')
    tmux set-window-option -t "$window_id" @alert 1
    tmux set-window-option -t "$window_id" window-status-style "fg=#f38ba8,bold,bg=#45475a"
    echo "Alert set for window $window_id"
}

# Set alert with custom color
tmux_alert_color() {
    if [ -z "$TMUX" ]; then
        echo "Error: Not in a tmux session"
        return 1
    fi
    local color="${1:-#f38ba8}"
    local window_id
    window_id=$(tmux display-message -p '#{window_id}')
    tmux set-window-option -t "$window_id" @alert 1
    tmux set-window-option -t "$window_id" window-status-style "fg=$color,bold,bg=#45475a"
    echo "Alert set with color: $color"
}

# Priority alert levels using Catppuccin Mocha colors
tmux_alert_high() { tmux_alert_color "#f38ba8"; }   # red
tmux_alert_medium() { tmux_alert_color "#f9e2af"; } # yellow
tmux_alert_low() { tmux_alert_color "#94e2d5"; }    # teal
`,
	}
}

// ClaudeComponent returns the Claude Code shell integration component.
func ClaudeComponent() *Component {
	return &Component{
		Name:        "claude",
		Description: "Claude Code management and utilities",
		Env: `# Claude configuration paths
export CLAUDE_DIR="${HOME}/.claude"
export CLAUDE_CONFIG="${HOME}/.claude.json"
export CLAUDE_SETTINGS="${CLAUDE_DIR}/settings.json"
export CLAUDE_LOCAL="${CLAUDE_DIR}/settings.local.json"
export CLAUDE_STATS="${CLAUDE_DIR}/stats-cache.json"
export CLAUDE_PROJECTS="${CLAUDE_DIR}/projects"
`,
		Aliases: `# Quick CLI access (native claude commands)
alias cc='claude'
alias ccc='claude --continue'
alias ccr='claude --resume'
alias ccp='claude -p'

# Acorn wrappers
alias cc-stats='acorn claude stats'
alias cc-tokens='acorn claude stats tokens'
alias cc-info='acorn claude info'
alias cc-perms='acorn claude permissions'
alias cc-settings='acorn claude settings'
alias cc-help='acorn claude help'
`,
		Functions: `# Show Claude Code info (wrapper for acorn claude info)
claude_info() {
    acorn claude info "$@"
}

# View usage statistics (wrapper for acorn claude stats)
claude_stats() {
    acorn claude stats "$@"
}

# View token usage (wrapper for acorn claude stats tokens)
claude_tokens() {
    acorn claude stats tokens "$@"
}

# View daily usage (wrapper for acorn claude stats daily)
claude_daily() {
    acorn claude stats daily "$@"
}

# View permissions (wrapper for acorn claude permissions)
claude_permissions() {
    acorn claude permissions "$@"
}

# Add permission (wrapper for acorn claude permissions add)
claude_permission_add() {
    acorn claude permissions add "$@"
}

# Remove permission (wrapper for acorn claude permissions remove)
claude_permission_remove() {
    acorn claude permissions remove "$@"
}

# View settings (wrapper for acorn claude settings)
claude_settings() {
    acorn claude settings "$@"
}

# Edit settings (must stay in shell for $EDITOR)
claude_settings_edit() {
    local file="${1:-global}"
    local target

    case "$file" in
        global|g) target="${CLAUDE_SETTINGS}" ;;
        local|l)  target="${CLAUDE_LOCAL}" ;;
        config|c) target="${CLAUDE_CONFIG}" ;;
        *)
            echo "Usage: claude_settings_edit [global|local|config]"
            return 1
            ;;
    esac

    if [ -f "$target" ]; then
        ${EDITOR:-vim} "$target"
    else
        echo "File not found: $target"
        return 1
    fi
}

# List projects (wrapper for acorn claude projects)
claude_projects() {
    acorn claude projects "$@"
}

# List MCP servers (wrapper for acorn claude mcp)
claude_mcp() {
    acorn claude mcp "$@"
}

# Add MCP server (wrapper for acorn claude mcp add)
claude_mcp_add() {
    acorn claude mcp add "$@"
}

# List custom commands (wrapper for acorn claude commands)
claude_commands() {
    acorn claude commands "$@"
}

# Clear cache (wrapper for acorn claude clear)
claude_clear() {
    acorn claude clear "$@"
}

# Open Claude directory (must stay in shell for cd)
claude_dir() {
    local dir="${CLAUDE_DIR:-$HOME/.claude}"
    if [ -d "$dir" ]; then
        cd "$dir" || return 1
        echo "Changed to $dir"
        ls -la
    else
        echo "Claude directory not found: $dir"
        return 1
    fi
}

# Aggregate agents/commands (wrapper for acorn claude aggregate)
claude_aggregate() {
    acorn claude aggregate "$@"
}

# List aggregated items (wrapper for acorn claude aggregate list)
claude_list() {
    acorn claude aggregate list "$@"
}

# Show help (wrapper for acorn claude help)
claude_help() {
    acorn claude help "$@"
}
`,
	}
}

// CloudFlareComponent returns the CloudFlare shell integration component.
func CloudFlareComponent() *Component {
	return &Component{
		Name:        "cloudflare",
		Description: "CloudFlare CLI (wrangler) integration for Workers, Pages, R2, KV, and D1",
		Env: `# CloudFlare/Wrangler environment setup (XDG-compliant)
export WRANGLER_HOME="${XDG_CONFIG_HOME:-$HOME/.config}/wrangler"

# Ensure wrangler config directory exists
if [ ! -d "${WRANGLER_HOME}" ]; then
    mkdir -p "${WRANGLER_HOME}" 2>/dev/null
fi
`,
		Aliases: `# Core wrangler aliases
alias wr='wrangler'
alias wrd='wrangler dev'
alias wrp='wrangler pages'
alias wrr2='wrangler r2'
alias wrkv='wrangler kv'

# Workers management
alias wrlist='wrangler deployments list'
alias wrtail='wrangler tail'
alias wrpub='wrangler deploy'

# Pages management
alias wrplist='wrangler pages project list'
alias wrpdeploy='wrangler pages deploy'

# R2 storage
alias wrr2list='wrangler r2 bucket list'

# KV storage
alias wrkvlist='wrangler kv namespace list'

# D1 database
alias wrd1='wrangler d1'
alias wrd1list='wrangler d1 list'

# Secrets management
alias wrsecret='wrangler secret'
alias wrsecrets='wrangler secret list'

# Login/logout
alias wrlogin='wrangler login'
alias wrlogout='wrangler logout'
alias wrwhoami='wrangler whoami'
`,
		Functions: `# Check CloudFlare CLI status (wrapper for acorn cf status)
cf_status() {
    acorn cf status "$@"
}

# Show current account (wrapper for acorn cf whoami)
cf_whoami() {
    acorn cf whoami
}

# List Workers (wrapper for acorn cf workers)
cf_workers() {
    acorn cf workers
}

# List Pages projects (wrapper for acorn cf pages)
cf_pages() {
    acorn cf pages
}

# List R2 buckets (wrapper for acorn cf r2 list)
cf_r2_buckets() {
    acorn cf r2 list
}

# List KV namespaces (wrapper for acorn cf kv list)
cf_kv_namespaces() {
    acorn cf kv list
}

# List D1 databases (wrapper for acorn cf d1 list)
cf_d1_databases() {
    acorn cf d1 list
}

# Tail worker logs (wrapper for acorn cf logs)
cf_logs() {
    acorn cf logs "$@"
}

# Deploy current worker (wrapper for acorn cf deploy)
cf_deploy() {
    acorn cf deploy "$@"
}

# Initialize Worker project (wrapper for acorn cf init worker)
cf_worker_init() {
    acorn cf init worker "$@"
}

# Initialize Pages project (wrapper for acorn cf init pages)
cf_pages_init() {
    acorn cf init pages "$@"
}

# Create R2 bucket (wrapper for acorn cf r2 create)
cf_r2_create() {
    acorn cf r2 create "$@"
}

# Create KV namespace (wrapper for acorn cf kv create)
cf_kv_create() {
    acorn cf kv create "$@"
}

# Create D1 database (wrapper for acorn cf d1 create)
cf_d1_create() {
    acorn cf d1 create "$@"
}

# Put secret (wrapper for acorn cf secret-put)
cf_secret_put() {
    acorn cf secret-put "$@"
}

# List secrets (wrapper for acorn cf secrets)
cf_secrets() {
    acorn cf secrets
}

# Show overview (wrapper for acorn cf overview)
cf_overview() {
    acorn cf overview "$@"
}

# Login (wrapper for acorn cf login)
cf_login() {
    acorn cf login
}

# Logout (wrapper for acorn cf logout)
cf_logout() {
    acorn cf logout
}

# Show help
cf_help() {
    echo "CloudFlare Component Functions"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Status & Info:"
    echo "  cf_status           Check CLI status and auth"
    echo "  cf_whoami           Show current account"
    echo "  cf_overview         Overview of all resources"
    echo ""
    echo "List Resources:"
    echo "  cf_workers          List Workers"
    echo "  cf_pages            List Pages projects"
    echo "  cf_r2_buckets       List R2 buckets"
    echo "  cf_kv_namespaces    List KV namespaces"
    echo "  cf_d1_databases     List D1 databases"
    echo ""
    echo "Create Resources:"
    echo "  cf_worker_init      Create new Worker project"
    echo "  cf_pages_init       Create new Pages project"
    echo "  cf_r2_create        Create R2 bucket"
    echo "  cf_kv_create        Create KV namespace"
    echo "  cf_d1_create        Create D1 database"
    echo ""
    echo "Operations:"
    echo "  cf_deploy           Deploy current worker"
    echo "  cf_logs <worker>    Tail worker logs"
    echo "  cf_secret_put       Add worker secret"
    echo "  cf_secrets          List worker secrets"
    echo ""
    echo "Aliases:"
    echo "  wr, wrd, wrp, wrr2, wrkv, wrd1"
    echo "  wrlogin, wrlogout, wrwhoami"
}
`,
	}
}

// SecretsComponent returns the secrets shell integration component.
func SecretsComponent() *Component {
	return &Component{
		Name:        "secrets",
		Description: "Secrets management and credential loading",
		Env: `# Secrets directory (XDG-compliant)
export SECRETS_DIR="${SECRETS_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/secrets}"

# Auto-load secrets if enabled
if [ "$AUTO_LOAD_SECRETS" = "true" ]; then
    if [ -f "$SECRETS_DIR/.env" ]; then
        set -a
        . "$SECRETS_DIR/.env"
        set +a
    fi
fi
`,
		Aliases: `# Secrets management aliases
alias secrets-load='load_secrets'
alias secrets-status='acorn secrets status'
alias secrets-validate='acorn secrets validate'
alias secrets-list='acorn secrets list'
alias check-aws='acorn secrets check aws'
alias check-azure='acorn secrets check azure'
alias check-github='acorn secrets check github'
alias check-do='acorn secrets check digitalocean'
`,
		Functions: `# Load secrets into current shell environment
# Note: This must be a shell function to affect the current shell
load_secrets() {
    local secrets_file="${SECRETS_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/secrets}/.env"

    if [ -f "$secrets_file" ]; then
        if [ -r "$secrets_file" ]; then
            set -a
            . "$secrets_file"
            set +a
            echo "Secrets loaded into environment"
        else
            echo "Cannot read secrets file (check permissions)"
            return 1
        fi
    else
        echo "No secrets file found at: $secrets_file"
        echo "Run: acorn secrets init"
        return 1
    fi
}

# Show secrets status (wrapper for acorn secrets status)
secrets_status() {
    acorn secrets status "$@"
}

# Validate secrets (wrapper for acorn secrets validate)
validate_secrets() {
    acorn secrets validate "$@"
}

# List secrets (wrapper for acorn secrets list)
list_secrets() {
    acorn secrets list "$@"
}

# Check AWS credentials (wrapper for acorn secrets check)
check_aws_key() {
    acorn secrets check aws
}

# Check Azure credentials (wrapper for acorn secrets check)
check_azure_key() {
    acorn secrets check azure
}

# Check GitHub token (wrapper for acorn secrets check)
check_github_key() {
    acorn secrets check github
}

# Check DigitalOcean token (wrapper for acorn secrets check)
check_digitalocean_key() {
    acorn secrets check digitalocean
}

# Check OpenAI API key (wrapper for acorn secrets check)
check_openai_key() {
    acorn secrets check openai
}

# Check Anthropic API key (wrapper for acorn secrets check)
check_anthropic_key() {
    acorn secrets check anthropic
}

# Check all credentials (wrapper for acorn secrets check)
check_all_keys() {
    acorn secrets check
}

# Initialize secrets file (wrapper for acorn secrets init)
secrets_init() {
    acorn secrets init
}

# Show secrets path (wrapper for acorn secrets path)
secrets_path() {
    acorn secrets path
}
`,
	}
}

// DatabaseComponent returns the database shell integration component.
func DatabaseComponent() *Component {
	return &Component{
		Name:        "database",
		Description: "Database tools and service management",
		Aliases: `# PostgreSQL
alias pg='pgcli'
alias psqlc='psql'

# MySQL
if command -v mycli >/dev/null 2>&1; then
    alias my='mycli'
fi

# MongoDB
if command -v mongosh >/dev/null 2>&1; then
    alias mongo='mongosh'
    alias msh='mongosh'
fi

# Redis
if command -v iredis >/dev/null 2>&1; then
    alias rd='iredis'
elif command -v redis-cli >/dev/null 2>&1; then
    alias rd='redis-cli'
fi

# SQLite
alias sq='sqlite3'
alias sqr='sqlite3 -readonly'
alias sqh='sqlite3 -header -column'

# Neo4j
if command -v cypher-shell >/dev/null 2>&1; then
    alias neo='cypher-shell'
fi

# macOS Brew service management
if [ "$CURRENT_PLATFORM" = "darwin" ]; then
    # PostgreSQL
    alias pgstart='brew services start postgresql@14'
    alias pgstop='brew services stop postgresql@14'
    alias pgrestart='brew services restart postgresql@14'
    alias pgstatus='brew services info postgresql@14'
    # MySQL
    alias mystart='brew services start mysql'
    alias mystop='brew services stop mysql'
    alias myrestart='brew services restart mysql'
    alias mystatus='brew services info mysql'
    # MongoDB
    alias mongostart='brew services start mongodb-community'
    alias mongostop='brew services stop mongodb-community'
    alias mongorestart='brew services restart mongodb-community'
    alias mongostatus='brew services info mongodb-community'
    # Redis
    alias rdstart='brew services start redis'
    alias rdstop='brew services stop redis'
    alias rdrestart='brew services restart redis'
    alias rdstatus='brew services info redis'
    # Neo4j
    alias neostart='brew services start neo4j'
    alias neostop='brew services stop neo4j'
    alias neorestart='brew services restart neo4j'
    alias neostatus='brew services info neo4j'
    # Kafka
    alias kafkastart='brew services start kafka'
    alias kafkastop='brew services stop kafka'
    alias kafkarestart='brew services restart kafka'
    alias kafkastatus='brew services info kafka'
    alias zkstart='brew services start zookeeper'
    alias zkstop='brew services stop zookeeper'
fi

# Kafka tools
if command -v kafka-console-producer >/dev/null 2>&1; then
    alias kprod='kafka-console-producer'
    alias kcons='kafka-console-consumer'
    alias ktop='kafka-topics'
    alias kgroups='kafka-consumer-groups'
fi
`,
		Functions: `# PostgreSQL local connection
pglocal() {
    pgcli -h localhost -U "${1:-postgres}" "${2:-postgres}"
}

# MySQL local connection
mylocal() {
    mycli -h localhost -u "${1:-root}" "${2:-}"
}

# MongoDB local connection
mongolocal() {
    mongosh "mongodb://localhost:27017/${1:-test}"
}

# Redis local connection
rdlocal() {
    if command -v iredis >/dev/null 2>&1; then
        iredis -h localhost -p "${1:-6379}"
    else
        redis-cli -h localhost -p "${1:-6379}"
    fi
}

# Neo4j local connection
neolocal() {
    cypher-shell -u "${1:-neo4j}" -p "${2:-neo4j}" -a "bolt://localhost:7687"
}

# Database status (wrapper for acorn db status)
db_status() {
    acorn db status "$@"
}

# Start all databases (wrapper for acorn db start-all)
db_start_all() {
    acorn db start-all
}

# Stop all databases (wrapper for acorn db stop-all)
db_stop_all() {
    acorn db stop-all
}
`,
	}
}

// FzfComponent returns the FZF shell integration component.
func FzfComponent() *Component {
	return &Component{
		Name:        "fzf",
		Description: "Fuzzy finder with Catppuccin theme and shell integrations",
		Env: `# FZF Version Detection
if command -v fzf >/dev/null 2>&1; then
    FZF_VERSION=$(fzf --version | cut -d' ' -f1)
    export FZF_VERSION
fi

# FZF Default Commands (use fd for faster search)
case "$CURRENT_PLATFORM" in
    darwin)
        export FZF_DEFAULT_COMMAND='fd --type f --hidden --follow --exclude .git'
        export FZF_ALT_C_COMMAND='fd --type d --hidden --follow --exclude .git'
        ;;
    linux)
        if command -v fdfind >/dev/null 2>&1; then
            export FZF_DEFAULT_COMMAND='fdfind --type f --hidden --follow --exclude .git'
            export FZF_ALT_C_COMMAND='fdfind --type d --hidden --follow --exclude .git'
        elif command -v fd >/dev/null 2>&1; then
            export FZF_DEFAULT_COMMAND='fd --type f --hidden --follow --exclude .git'
            export FZF_ALT_C_COMMAND='fd --type d --hidden --follow --exclude .git'
        fi
        ;;
esac

export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"

# FZF Catppuccin Mocha Theme
export FZF_DEFAULT_OPTS="
  --extended
  --height 40%
  --layout=reverse
  --border
  --color=bg+:#313244,bg:#1e1e2e,spinner:#f5e0dc,hl:#f38ba8
  --color=fg:#cdd6f4,header:#f38ba8,info:#cba6f7,pointer:#f5e0dc
  --color=marker:#f5e0dc,fg+:#cdd6f4,prompt:#cba6f7,hl+:#f38ba8
  --bind='ctrl-/:toggle-preview'
  --preview-window='right:50%:hidden'
"

# Preview options
export FZF_CTRL_T_OPTS="--preview '[ -d {} ] && ls -la {} || head -100 {}'"
export FZF_ALT_C_OPTS="--preview 'ls -la {}'"

# FZF Location Detection
if [ -z "$FZF_LOCATION" ]; then
    if [ -d "/opt/homebrew/opt/fzf" ]; then
        FZF_LOCATION="/opt/homebrew/opt/fzf"
    elif [ -d "/home/linuxbrew/.linuxbrew/opt/fzf" ]; then
        FZF_LOCATION="/home/linuxbrew/.linuxbrew/opt/fzf"
    elif [ -d "/usr/share/fzf" ]; then
        FZF_LOCATION="/usr/share/fzf"
    elif [ -d "$HOME/.fzf" ]; then
        FZF_LOCATION="$HOME/.fzf"
    fi
fi
export FZF_LOCATION
`,
		Aliases: `# FZF shortcuts
alias ff='fzf_files'
alias fcd='fzf_cd'
alias fgb='fzf_git_branch'
alias fgl='fzf_git_log'
alias fgs='fzf_git_stash'
alias fga='fzf_git_add'
alias fkill='fzf_kill'
alias fh='fzf_history'
alias fenv='fzf_env'
`,
		Functions: `# =============================================================================
# File Search
# =============================================================================

# Interactive file finder with preview
fzf_files() {
    local file
    file=$(fzf --preview 'head -100 {}')
    [ -n "$file" ] && ${EDITOR:-vim} "$file"
}

# Find and edit file
fe() {
    local file
    file=$(fzf --query="$1" --select-1 --exit-0)
    [ -n "$file" ] && ${EDITOR:-vim} "$file"
}

# =============================================================================
# Directory Navigation
# =============================================================================

# Interactive cd with preview
fzf_cd() {
    local dir
    case "$CURRENT_PLATFORM" in
        darwin) dir=$(fd --type d --hidden --follow --exclude .git 2>/dev/null | fzf --preview 'ls -la {}') ;;
        linux)  dir=$(fdfind --type d --hidden --follow --exclude .git 2>/dev/null | fzf --preview 'ls -la {}') ;;
        *)      dir=$(find . -type d 2>/dev/null | fzf --preview 'ls -la {}') ;;
    esac
    [ -n "$dir" ] && cd "$dir" || return
}

# =============================================================================
# Git Integration
# =============================================================================

# Interactive git branch checkout
fzf_git_branch() {
    local branch
    branch=$(git branch --all | grep -v HEAD | sed 's/.* //' | sed 's#remotes/origin/##' | sort -u | fzf)
    [ -n "$branch" ] && git checkout "$branch"
}

# Interactive git log browser
fzf_git_log() {
    git log --oneline --color=always | fzf --ansi --preview 'git show --color=always {1}'
}

# Interactive git stash browser
fzf_git_stash() {
    local stash
    stash=$(git stash list | fzf --preview 'git stash show -p {1}' | cut -d: -f1)
    [ -n "$stash" ] && git stash apply "$stash"
}

# Interactive git add
fzf_git_add() {
    local files
    files=$(git status -s | fzf -m --preview 'git diff --color=always {2}' | awk '{print $2}')
    [ -n "$files" ] && echo "$files" | xargs git add
}

# =============================================================================
# Process Management
# =============================================================================

# Interactive process killer
fzf_kill() {
    local pid
    pid=$(ps aux | sed 1d | fzf -m | awk '{print $2}')
    [ -n "$pid" ] && echo "$pid" | xargs kill -"${1:-9}"
}

# =============================================================================
# History Search
# =============================================================================

# Interactive history search and execute
fzf_history() {
    local cmd
    cmd=$(history | fzf --tac --no-sort | sed 's/^[ ]*[0-9]*[ ]*//')
    [ -n "$cmd" ] && eval "$cmd"
}

# =============================================================================
# Environment
# =============================================================================

# Interactive environment variable browser
fzf_env() {
    local var
    var=$(env | fzf | cut -d= -f1)
    [ -n "$var" ] && echo "${var}=$(printenv "$var")"
}

# =============================================================================
# Kubernetes (if kubectl available)
# =============================================================================

if command -v kubectl >/dev/null 2>&1; then
    # Interactive pod selector
    fzf_k8s_pod() {
        kubectl get pods --all-namespaces -o wide | fzf | awk '{print $2}'
    }

    # Interactive pod logs
    fzf_k8s_logs() {
        local selection ns pod
        selection=$(kubectl get pods --all-namespaces -o wide | fzf)
        if [ -n "$selection" ]; then
            ns=$(echo "$selection" | awk '{print $1}')
            pod=$(echo "$selection" | awk '{print $2}')
            kubectl logs -n "$ns" "$pod" -f
        fi
    }
    alias fkl='fzf_k8s_logs'

    # Interactive namespace switcher
    fzf_k8s_ns() {
        local ns
        ns=$(kubectl get namespaces -o name | sed 's/namespace\///' | fzf)
        [ -n "$ns" ] && kubectl config set-context --current --namespace="$ns"
    }
    alias fkns='fzf_k8s_ns'
fi

# =============================================================================
# Docker (if docker available)
# =============================================================================

if command -v docker >/dev/null 2>&1; then
    # Interactive container selector
    fzf_docker_container() {
        docker ps -a --format "table {{.ID}}\t{{.Names}}\t{{.Status}}\t{{.Image}}" | fzf | awk '{print $1}'
    }

    # Interactive container logs
    fzf_docker_logs() {
        local container
        container=$(fzf_docker_container)
        [ -n "$container" ] && docker logs -f "$container"
    }
    alias fdl='fzf_docker_logs'

    # Interactive container exec
    fzf_docker_exec() {
        local container
        container=$(fzf_docker_container)
        [ -n "$container" ] && docker exec -it "$container" "${1:-/bin/sh}"
    }
    alias fdx='fzf_docker_exec'
fi
`,
	}
}
