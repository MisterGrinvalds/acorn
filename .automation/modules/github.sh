#!/bin/bash
# GitHub Automation Module

# Module configuration
readonly GITHUB_TEMPLATES_DIR="$AUTO_HOME/github/templates"
readonly GITHUB_WORKFLOWS_DIR="$AUTO_HOME/github/workflows"

# Initialize directories
mkdir -p "$GITHUB_TEMPLATES_DIR" "$GITHUB_WORKFLOWS_DIR"

# Help function for github module
github_help() {
    cat << EOF
GitHub Automation

USAGE:
    auto github <command> [options]

COMMANDS:
    repo <action>                Repository management (create, clone, fork)
    pr <action>                  Pull request operations (create, merge, list)
    issue <action>               Issue management (create, list, close)
    workflow <action>            GitHub Actions workflows
    release <action>             Release management (create, list, download)
    security                     Security scanning and setup
    templates                    Setup repository templates
    metrics                      Repository analytics
    sync                         Sync forks and branches

EXAMPLES:
    auto github repo create my-new-repo "Description"
    auto github pr create feature/new-feature "Add new feature"
    auto github issue create "Bug report" "Something is broken"
    auto github workflow setup-ci python
    auto github release create v1.0.0 "First release"
    auto github security setup
EOF
}

# Utility functions
require_gh() {
    require_command gh
}

get_repo_info() {
    if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
        local remote_url=$(git config --get remote.origin.url)
        if [[ "$remote_url" =~ github\.com[:/]([^/]+)/([^/.]+) ]]; then
            echo "${BASH_REMATCH[1]}/${BASH_REMATCH[2]}"
        fi
    fi
}

# Repository management
repo_create() {
    require_gh
    local repo_name="$1"
    local description="${2:-A new repository}"
    local visibility="${3:-public}"
    
    if [ -z "$repo_name" ]; then
        log "ERROR" "Repository name is required"
        exit 1
    fi
    
    log "INFO" "Creating GitHub repository: $repo_name"
    
    # Create repository
    gh repo create "$repo_name" --description "$description" --"$visibility"
    
    # Clone if we're not in a git directory
    if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
        git clone "https://github.com/$(gh api user --jq .login)/$repo_name.git"
        cd "$repo_name"
    fi
    
    # Setup basic files
    if [ ! -f "README.md" ]; then
        cat > README.md << EOF
# $repo_name

$description

## Installation

\`\`\`bash
# Installation instructions here
\`\`\`

## Usage

\`\`\`bash
# Usage examples here
\`\`\`

## Contributing

1. Fork the repository
2. Create your feature branch (\`git checkout -b feature/amazing-feature\`)
3. Commit your changes (\`git commit -m 'Add some amazing feature'\`)
4. Push to the branch (\`git push origin feature/amazing-feature\`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
EOF
    fi
    
    if [ ! -f ".gitignore" ]; then
        cat > .gitignore << 'EOF'
# General
.DS_Store
Thumbs.db
*.log

# IDE
.vscode/
.idea/
*.swp
*.swo

# Environment
.env
.env.local
.env.*.local

# Dependencies
node_modules/
vendor/

# Build outputs
dist/
build/
bin/
target/

# Testing
coverage/
.nyc_output/
.pytest_cache/

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db
EOF
    fi
    
    # Initial commit if needed
    if ! git log --oneline -1 >/dev/null 2>&1; then
        git add .
        git commit -m "Initial commit"
        git push -u origin main
    fi
    
    log "SUCCESS" "Repository $repo_name created successfully"
}

repo_clone() {
    require_gh
    local repo="$1"
    local directory="${2:-}"
    
    if [ -z "$repo" ]; then
        log "ERROR" "Repository is required (format: user/repo or full URL)"
        exit 1
    fi
    
    if [ -n "$directory" ]; then
        gh repo clone "$repo" "$directory"
        cd "$directory"
    else
        gh repo clone "$repo"
        cd "$(basename "$repo")"
    fi
    
    log "SUCCESS" "Repository cloned successfully"
}

# Pull request management
pr_create() {
    require_gh
    local title="$1"
    local body="${2:-}"
    local draft="${3:-false}"
    
    if [ -z "$title" ]; then
        log "ERROR" "Pull request title is required"
        exit 1
    fi
    
    # Ensure we're on a feature branch
    local current_branch=$(git branch --show-current)
    if [ "$current_branch" = "main" ] || [ "$current_branch" = "master" ]; then
        log "ERROR" "Cannot create PR from main/master branch"
        exit 1
    fi
    
    # Push current branch
    git push -u origin "$current_branch"
    
    # Create PR
    local pr_args=(--title "$title")
    [ -n "$body" ] && pr_args+=(--body "$body")
    [ "$draft" = "true" ] && pr_args+=(--draft)
    
    gh pr create "${pr_args[@]}"
    log "SUCCESS" "Pull request created"
}

pr_merge() {
    require_gh
    local pr_number="$1"
    local merge_method="${2:-merge}"
    
    if [ -z "$pr_number" ]; then
        log "ERROR" "PR number is required"
        exit 1
    fi
    
    case "$merge_method" in
        "merge"|"squash"|"rebase")
            gh pr merge "$pr_number" --"$merge_method"
            log "SUCCESS" "PR #$pr_number merged using $merge_method"
            ;;
        *)
            log "ERROR" "Invalid merge method: $merge_method (use: merge, squash, rebase)"
            exit 1
            ;;
    esac
}

# Issue management
issue_create() {
    require_gh
    local title="$1"
    local body="${2:-}"
    local labels="${3:-}"
    
    if [ -z "$title" ]; then
        log "ERROR" "Issue title is required"
        exit 1
    fi
    
    local issue_args=(--title "$title")
    [ -n "$body" ] && issue_args+=(--body "$body")
    [ -n "$labels" ] && issue_args+=(--label "$labels")
    
    gh issue create "${issue_args[@]}"
    log "SUCCESS" "Issue created"
}

# Workflow management
setup_ci_workflow() {
    local language="${1:-auto-detect}"
    local workflow_dir=".github/workflows"
    
    mkdir -p "$workflow_dir"
    
    # Auto-detect language if not specified
    if [ "$language" = "auto-detect" ]; then
        if [ -f "package.json" ]; then
            language="node"
        elif [ -f "go.mod" ]; then
            language="go"
        elif [ -f "requirements.txt" ] || [ -f "pyproject.toml" ]; then
            language="python"
        else
            log "WARN" "Could not auto-detect language, defaulting to generic workflow"
            language="generic"
        fi
    fi
    
    case "$language" in
        "python")
            cat > "$workflow_dir/python-ci.yml" << 'EOF'
name: Python CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        python-version: [3.8, 3.9, '3.10', '3.11']

    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Python ${{ matrix.python-version }}
      uses: actions/setup-python@v4
      with:
        python-version: ${{ matrix.python-version }}
    
    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        pip install -r requirements.txt
        pip install pytest black isort flake8
    
    - name: Lint with flake8
      run: |
        flake8 . --count --select=E9,F63,F7,F82 --show-source --statistics
        flake8 . --count --exit-zero --max-complexity=10 --max-line-length=88 --statistics
    
    - name: Format check with black
      run: black --check .
    
    - name: Import sort check
      run: isort --check-only .
    
    - name: Test with pytest
      run: pytest --cov=. --cov-report=xml
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.xml
EOF
            ;;
        "go")
            cat > "$workflow_dir/go-ci.yml" << 'EOF'
name: Go CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Build
      run: go build -v ./...
    
    - name: Test
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out

  lint:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
EOF
            ;;
        "node")
            cat > "$workflow_dir/node-ci.yml" << 'EOF'
name: Node.js CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    strategy:
      matrix:
        node-version: [16.x, 18.x, 20.x]
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Use Node.js ${{ matrix.node-version }}
      uses: actions/setup-node@v4
      with:
        node-version: ${{ matrix.node-version }}
        cache: 'npm'
    
    - run: npm ci
    - run: npm run build --if-present
    - run: npm test
    - run: npm run lint --if-present
EOF
            ;;
    esac
    
    log "SUCCESS" "CI workflow created for $language"
}

# Security setup
setup_security() {
    require_gh
    local repo=$(get_repo_info)
    
    if [ -z "$repo" ]; then
        log "ERROR" "Not in a GitHub repository"
        exit 1
    fi
    
    log "INFO" "Setting up security features for $repo"
    
    # Enable vulnerability alerts
    gh api --method PUT "repos/$repo/vulnerability-alerts" || log "WARN" "Could not enable vulnerability alerts"
    
    # Enable automated security fixes
    gh api --method PUT "repos/$repo/automated-security-fixes" || log "WARN" "Could not enable automated security fixes"
    
    # Create security policy template
    mkdir -p .github
    if [ ! -f ".github/SECURITY.md" ]; then
        cat > .github/SECURITY.md << 'EOF'
# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

Please report security vulnerabilities by emailing [security@yourproject.com](mailto:security@yourproject.com).

**Please do not report security vulnerabilities through public GitHub issues.**

We will acknowledge your email within 48 hours and will send a more detailed response 
within 48 hours indicating the next steps in handling your report.
EOF
    fi
    
    # Create CodeQL workflow
    mkdir -p .github/workflows
    cat > .github/workflows/codeql.yml << 'EOF'
name: "CodeQL"

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  schedule:
    - cron: '42 17 * * 1'

jobs:
  analyze:
    name: Analyze
    runs-on: ubuntu-latest
    permissions:
      actions: read
      contents: read
      security-events: write

    strategy:
      fail-fast: false
      matrix:
        language: [ 'javascript', 'python', 'go' ]

    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Initialize CodeQL
      uses: github/codeql-action/init@v2
      with:
        languages: ${{ matrix.language }}

    - name: Autobuild
      uses: github/codeql-action/autobuild@v2

    - name: Perform CodeQL Analysis
      uses: github/codeql-action/analyze@v2
EOF
    
    log "SUCCESS" "Security features configured"
}

# Main github module function
github_main() {
    local command="${1:-help}"
    shift || true
    
    case "$command" in
        "help"|"-h"|"--help")
            github_help
            ;;
        "repo")
            local action="$1"
            case "$action" in
                "create") repo_create "$2" "$3" "$4" ;;
                "clone") repo_clone "$2" "$3" ;;
                "fork") gh repo fork "$2" ;;
                "list") gh repo list ;;
                *) log "ERROR" "Unknown repo action: $action" ;;
            esac
            ;;
        "pr")
            local action="$1"
            case "$action" in
                "create") pr_create "$2" "$3" "$4" ;;
                "merge") pr_merge "$2" "$3" ;;
                "list") gh pr list ;;
                "view") gh pr view "$2" ;;
                "checkout") gh pr checkout "$2" ;;
                *) log "ERROR" "Unknown pr action: $action" ;;
            esac
            ;;
        "issue")
            local action="$1"
            case "$action" in
                "create") issue_create "$2" "$3" "$4" ;;
                "list") gh issue list ;;
                "view") gh issue view "$2" ;;
                "close") gh issue close "$2" ;;
                *) log "ERROR" "Unknown issue action: $action" ;;
            esac
            ;;
        "workflow")
            local action="$1"
            case "$action" in
                "setup-ci") setup_ci_workflow "$2" ;;
                "list") gh workflow list ;;
                "run") gh workflow run "$2" ;;
                "view") gh run view "$2" ;;
                *) log "ERROR" "Unknown workflow action: $action" ;;
            esac
            ;;
        "release")
            local action="$1"
            case "$action" in
                "create") gh release create "$2" --title "$2" --notes "$3" ;;
                "list") gh release list ;;
                "view") gh release view "$2" ;;
                "download") gh release download "$2" ;;
                *) log "ERROR" "Unknown release action: $action" ;;
            esac
            ;;
        "security")
            setup_security
            ;;
        "metrics")
            require_gh
            local repo=$(get_repo_info)
            [ -n "$repo" ] && gh repo view "$repo" --json stargazerCount,forkCount,issues,pullRequests
            ;;
        *)
            log "ERROR" "Unknown command: $command"
            github_help
            exit 1
            ;;
    esac
}