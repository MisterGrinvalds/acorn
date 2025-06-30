# tmux Helper Functions

# Create development session with multiple panes
dev_session() {
    local session_name="${1:-dev}"
    
    # Create new session
    tmux new-session -d -s "$session_name"
    
    # Split window horizontally (editor on top, terminal on bottom)
    tmux split-window -h -t "$session_name"
    
    # Split the right pane vertically (terminal and logs)
    tmux split-window -v -t "$session_name:0.1"
    
    # Optional: set up pane purposes
    tmux send-keys -t "$session_name:0.0" 'echo "Editor pane - use: code ."' Enter
    tmux send-keys -t "$session_name:0.1" 'echo "Main terminal"' Enter
    tmux send-keys -t "$session_name:0.2" 'echo "Logs/monitoring pane"' Enter
    
    # Select the first pane
    tmux select-pane -t "$session_name:0.0"
    
    # Attach to session
    tmux attach-session -t "$session_name"
}

# Create Kubernetes development session
k8s_session() {
    local session_name="k8s"
    
    tmux new-session -d -s "$session_name"
    
    # Create windows for different k8s tasks
    tmux new-window -t "$session_name" -n "kubectl"
    tmux new-window -t "$session_name" -n "logs"
    tmux new-window -t "$session_name" -n "k9s"
    
    # Set up kubectl window
    tmux send-keys -t "$session_name:kubectl" 'echo "kubectl commands"' Enter
    
    # Set up logs window with split panes
    tmux split-window -h -t "$session_name:logs"
    tmux send-keys -t "$session_name:logs.0" 'echo "App logs"' Enter
    tmux send-keys -t "$session_name:logs.1" 'echo "System logs"' Enter
    
    # Set up k9s window
    tmux send-keys -t "$session_name:k9s" 'k9s' Enter
    
    # Go to first window
    tmux select-window -t "$session_name:1"
    tmux attach-session -t "$session_name"
}

# Create project session (auto-detects project type)
project_session() {
    local project_path="${1:-.}"
    local session_name
    
    # Get project name from directory
    session_name=$(basename "$(realpath "$project_path")")
    
    cd "$project_path" || return 1
    
    # Create new session
    tmux new-session -d -s "$session_name" -c "$project_path"
    
    # Auto-detect project type and set up accordingly
    if [ -f "go.mod" ]; then
        # Go project
        tmux send-keys -t "$session_name" 'echo "Go project detected"' Enter
        tmux new-window -t "$session_name" -n "test" -c "$project_path"
        tmux send-keys -t "$session_name:test" 'echo "Run: go test ./..."' Enter
        tmux new-window -t "$session_name" -n "run" -c "$project_path"
        tmux send-keys -t "$session_name:run" 'echo "Run: go run main.go"' Enter
        
    elif [ -f "requirements.txt" ] || [ -f "pyproject.toml" ] || [ -f "Pipfile" ]; then
        # Python project
        tmux send-keys -t "$session_name" 'echo "Python project detected"' Enter
        tmux new-window -t "$session_name" -n "test" -c "$project_path"
        tmux send-keys -t "$session_name:test" 'echo "Run: pytest"' Enter
        tmux new-window -t "$session_name" -n "server" -c "$project_path"
        
        # Check for FastAPI
        if grep -q "fastapi\|uvicorn" requirements.txt pyproject.toml 2>/dev/null; then
            tmux send-keys -t "$session_name:server" 'echo "FastAPI detected - Run: uvdev"' Enter
        fi
        
    elif [ -f "package.json" ]; then
        # Node.js/TypeScript project
        tmux send-keys -t "$session_name" 'echo "Node.js project detected"' Enter
        tmux new-window -t "$session_name" -n "dev" -c "$project_path"
        tmux send-keys -t "$session_name:dev" 'echo "Run: npm run dev"' Enter
        tmux new-window -t "$session_name" -n "test" -c "$project_path"
        tmux send-keys -t "$session_name:test" 'echo "Run: npm test"' Enter
    fi
    
    # Always create an editor window
    tmux new-window -t "$session_name" -n "editor" -c "$project_path"
    tmux send-keys -t "$session_name:editor" 'code .' Enter
    
    # Go back to first window
    tmux select-window -t "$session_name:1"
    tmux attach-session -t "$session_name"
}

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