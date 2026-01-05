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
