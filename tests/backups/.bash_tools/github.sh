# Enhanced Git and GitHub Tools

# Enhanced git aliases (keeping existing ones, adding more)
alias ga='git add'
alias gaa='git add -A'
alias gap='git add -p'
alias gb='git branch'
alias gba='git branch -a'
alias gbd='git branch -d'
alias gbD='git branch -D'
alias gc='git commit'
alias gcm='git commit -m'
alias gca='git commit -a'
alias gcam='git commit -am'
alias gco='git checkout'
alias gcb='git checkout -b'
alias gd='git diff'
alias gdc='git diff --cached'
alias gf='git fetch'
alias gl='git log --oneline --graph --decorate'
alias gll='git log --oneline --graph --decorate -10'
alias gp='git push'
alias gpu='git push -u origin'
alias gpl='git pull'
alias gr='git reset'
alias grh='git reset --hard'
alias grs='git reset --soft'
alias gst='git stash'
alias gstp='git stash pop'
alias gstl='git stash list'

# GitHub CLI aliases (requires gh to be installed)
alias ghpr='gh pr create'
alias ghprs='gh pr status'
alias ghprv='gh pr view'
alias ghprc='gh pr checkout'
alias ghprm='gh pr merge'
alias ghrepo='gh repo view'
alias ghissue='gh issue create'
alias ghissues='gh issue list'

# Advanced Git workflows
gitcleanup() {
    echo "Cleaning up merged branches..."
    git branch --merged | grep -v "\*\|main\|master\|develop" | xargs -n 1 git branch -d
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

# Push new branch to origin
pushbranch() {
    local branch=$(git branch --show-current)
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
    local branch=$(git branch --show-current)
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

# Clone and cd into repository
gclone() {
    if [ -z "$1" ]; then
        echo "Usage: gclone <repo-url>"
        return 1
    fi
    git clone "$1"
    local repo_name=$(basename "$1" .git)
    cd "$repo_name" || return 1
}

# GitHub repository creation (requires gh)
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
    echo "$description" >> README.md
    git add README.md
    git commit -m "Initial commit"
    
    gh repo create "$repo_name" --public --description "$description" --source .
    git push -u origin main
}