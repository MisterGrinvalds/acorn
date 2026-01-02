#!/bin/sh
# components/github/functions.sh - GitHub CLI workflow functions

# =============================================================================
# Branch & PR Workflow
# =============================================================================

# Clean up merged branches
gitcleanup() {
    echo "Cleaning up merged branches..."
    git branch --merged | grep -v "\*\|main\|master\|develop" | xargs -n 1 git branch -d 2>/dev/null
    echo "Pruning remote tracking branches..."
    git remote prune origin
}

# Quick commit with message
qcommit() {
    if [ -z "$1" ]; then
        echo "Usage: qcommit <commit-message>"
        return 1
    fi
    git add -A
    git commit -m "$1"
}

# Push current branch to origin
pushbranch() {
    local branch
    branch=$(git branch --show-current)
    git push -u origin "$branch"
}

# Create and checkout new branch
newbranch() {
    if [ -z "$1" ]; then
        echo "Usage: newbranch <branch-name>"
        return 1
    fi
    git checkout -b "$1"
}

# Quick PR creation workflow
quickpr() {
    local branch
    branch=$(git branch --show-current)

    if [ "$branch" = "main" ] || [ "$branch" = "master" ]; then
        echo "Cannot create PR from main/master branch"
        return 1
    fi

    echo "Pushing branch '$branch' to origin..."
    git push -u origin "$branch"

    if command -v gh >/dev/null 2>&1; then
        echo "Creating pull request..."
        gh pr create --web
    else
        echo "GitHub CLI (gh) not installed. Please create PR manually."
    fi
}

# =============================================================================
# Repository Management
# =============================================================================

# Create new GitHub repository (requires gh)
newrepo() {
    if [ -z "$1" ]; then
        echo "Usage: newrepo <repo-name> [description]"
        return 1
    fi

    if ! command -v gh >/dev/null 2>&1; then
        echo "GitHub CLI (gh) is required for this function"
        return 1
    fi

    local repo_name="$1"
    local description="${2:-A new repository}"

    mkdir "$repo_name"
    cd "$repo_name" || return 1
    git init
    echo "# $repo_name" > README.md
    echo "" >> README.md
    echo "$description" >> README.md
    git add README.md
    git commit -m "Initial commit"

    gh repo create "$repo_name" --public --description "$description" --source .
    git push -u origin main
}

# Fork and clone a repository
forkclone() {
    if [ -z "$1" ]; then
        echo "Usage: forkclone <owner/repo>"
        return 1
    fi

    if ! command -v gh >/dev/null 2>&1; then
        echo "GitHub CLI (gh) is required"
        return 1
    fi

    gh repo fork "$1" --clone
    local repo_name
    repo_name=$(basename "$1")
    cd "$repo_name" || return 1
}

# =============================================================================
# Status & Info
# =============================================================================

# Git status with helpful info
gstat() {
    echo "=== Git Status ==="
    git status -sb
    echo ""
    echo "=== Recent Commits ==="
    git log --oneline -5
    echo ""
    echo "=== Stashes ==="
    git stash list
}

# Show PR status for current branch
prstatus() {
    if ! command -v gh >/dev/null 2>&1; then
        echo "GitHub CLI (gh) not installed"
        return 1
    fi

    gh pr status
}

# View PR checks
prchecks() {
    if ! command -v gh >/dev/null 2>&1; then
        echo "GitHub CLI (gh) not installed"
        return 1
    fi

    gh pr checks
}

# =============================================================================
# Actions & Workflows
# =============================================================================

# Watch current workflow run
watchrun() {
    if ! command -v gh >/dev/null 2>&1; then
        echo "GitHub CLI (gh) not installed"
        return 1
    fi

    gh run watch
}

# Rerun failed workflow
rerun() {
    if ! command -v gh >/dev/null 2>&1; then
        echo "GitHub CLI (gh) not installed"
        return 1
    fi

    gh run rerun --failed
}
