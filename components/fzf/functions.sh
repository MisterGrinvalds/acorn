#!/bin/sh
# components/fzf/functions.sh - FZF-powered functions

# =============================================================================
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
