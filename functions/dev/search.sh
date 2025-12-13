#!/bin/sh
# FZF-powered search and navigation functions
# Requires: fzf, fd (optional: bat, git)

# =============================================================================
# File Search
# =============================================================================

# Interactive file finder with preview
fzf_files() {
    local file
    file=$(fzf --preview 'head -100 {}')
    [ -n "$file" ] && ${EDITOR:-vim} "$file"
}
alias ff='fzf_files'

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
alias fcd='fzf_cd'

# =============================================================================
# Git Integration
# =============================================================================

# Interactive git branch checkout
fzf_git_branch() {
    local branch
    branch=$(git branch --all | grep -v HEAD | sed 's/.* //' | sed 's#remotes/origin/##' | sort -u | fzf)
    [ -n "$branch" ] && git checkout "$branch"
}
alias fgb='fzf_git_branch'

# Interactive git log browser
fzf_git_log() {
    git log --oneline --color=always | fzf --ansi --preview 'git show --color=always {1}'
}
alias fgl='fzf_git_log'

# Interactive git stash browser
fzf_git_stash() {
    local stash
    stash=$(git stash list | fzf --preview 'git stash show -p {1}' | cut -d: -f1)
    [ -n "$stash" ] && git stash apply "$stash"
}
alias fgs='fzf_git_stash'

# Interactive git add
fzf_git_add() {
    local files
    files=$(git status -s | fzf -m --preview 'git diff --color=always {2}' | awk '{print $2}')
    [ -n "$files" ] && echo "$files" | xargs git add
}
alias fga='fzf_git_add'

# =============================================================================
# Process Management
# =============================================================================

# Interactive process killer
fzf_kill() {
    local pid
    pid=$(ps aux | sed 1d | fzf -m | awk '{print $2}')
    [ -n "$pid" ] && echo "$pid" | xargs kill -"${1:-9}"
}
alias fkill='fzf_kill'

# =============================================================================
# History Search
# =============================================================================

# Interactive history search and execute
fzf_history() {
    local cmd
    cmd=$(history | fzf --tac --no-sort | sed 's/^[ ]*[0-9]*[ ]*//')
    [ -n "$cmd" ] && eval "$cmd"
}
alias fh='fzf_history'

# =============================================================================
# Kubernetes Integration (if kubectl available)
# =============================================================================

if command -v kubectl >/dev/null 2>&1; then
    # Interactive pod selector
    fzf_k8s_pod() {
        local pod
        pod=$(kubectl get pods --all-namespaces -o wide | fzf | awk '{print $2}')
        echo "$pod"
    }

    # Interactive pod logs
    fzf_k8s_logs() {
        local selection
        selection=$(kubectl get pods --all-namespaces -o wide | fzf)
        if [ -n "$selection" ]; then
            local ns pod
            ns=$(echo "$selection" | awk '{print $1}')
            pod=$(echo "$selection" | awk '{print $2}')
            kubectl logs -n "$ns" "$pod" -f
        fi
    }
    alias fkl='fzf_k8s_logs'

    # Interactive pod exec
    fzf_k8s_exec() {
        local selection
        selection=$(kubectl get pods --all-namespaces -o wide | fzf)
        if [ -n "$selection" ]; then
            local ns pod
            ns=$(echo "$selection" | awk '{print $1}')
            pod=$(echo "$selection" | awk '{print $2}')
            kubectl exec -n "$ns" -it "$pod" -- "${1:-/bin/sh}"
        fi
    }
    alias fkx='fzf_k8s_exec'

    # Interactive namespace switcher
    fzf_k8s_ns() {
        local ns
        ns=$(kubectl get namespaces -o name | sed 's/namespace\///' | fzf)
        [ -n "$ns" ] && kubectl config set-context --current --namespace="$ns"
    }
    alias fkns='fzf_k8s_ns'

    # Interactive context switcher
    fzf_k8s_ctx() {
        local ctx
        ctx=$(kubectl config get-contexts -o name | fzf)
        [ -n "$ctx" ] && kubectl config use-context "$ctx"
    }
    alias fkctx='fzf_k8s_ctx'
fi

# =============================================================================
# Docker Integration (if docker available)
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

    # Interactive image selector
    fzf_docker_image() {
        docker images --format "table {{.Repository}}:{{.Tag}}\t{{.Size}}\t{{.CreatedSince}}" | fzf | awk '{print $1}'
    }
    alias fdi='fzf_docker_image'
fi

# =============================================================================
# Environment Variable Search
# =============================================================================

# Interactive environment variable browser
fzf_env() {
    local var
    var=$(env | fzf | cut -d= -f1)
    [ -n "$var" ] && echo "${var}=$(printenv "$var")"
}
alias fenv='fzf_env'
