#!/bin/sh
# components/git/functions.sh - Git helper functions

# Clone and cd into repository
gclone() {
    git clone "$1" && cd "$(basename "$1" .git)" || return
}

# Create and checkout new branch
gcob() {
    if [ -z "$1" ]; then
        echo "Usage: gcob <branch-name>"
        return 1
    fi
    git checkout -b "$1"
}

# Push with upstream tracking
gpush() {
    local branch
    branch=$(git rev-parse --abbrev-ref HEAD)
    git push -u origin "$branch"
}

# Pull with rebase
gpull() {
    git pull --rebase origin "$(git rev-parse --abbrev-ref HEAD)"
}

# Interactive add
gadd() {
    git add -p "$@"
}

# Show git status with branch info
ginfo() {
    echo "Branch: $(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'not a git repo')"
    echo "Remote: $(git remote -v 2>/dev/null | head -1 | awk '{print $2}' || echo 'none')"
    echo ""
    git status -s 2>/dev/null
}

# Undo last commit (keep changes)
gundo() {
    git reset --soft HEAD~1
}

# Amend last commit without editing message
gamend() {
    git add --all
    git commit --amend --no-edit
}

# Show files changed in a commit
gshow() {
    git show --stat "${1:-HEAD}"
}

# Git blame with line numbers
gblame() {
    if [ -z "$1" ]; then
        echo "Usage: gblame <file>"
        return 1
    fi
    git blame -n "$1"
}

# Find commits by message
gfind() {
    if [ -z "$1" ]; then
        echo "Usage: gfind <search-term>"
        return 1
    fi
    git log --oneline --all --grep="$1"
}

# Show contribution stats
gcontrib() {
    git shortlog -sn --all
}

# Clean merged branches
gcleanb() {
    git branch --merged | grep -v '\*\|main\|master' | xargs -n 1 git branch -d
}
