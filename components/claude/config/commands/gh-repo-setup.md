---
description: Set up a new GitHub repository
argument-hint: <repo-name> [--private]
allowed-tools: Read, Write, Bash
---

## Task

Help the user create and set up a new GitHub repository.

## Quick Setup

Using dotfiles function:
```bash
newrepo myproject "My awesome project"
# Creates local repo, initializes, and pushes to GitHub
```

## Manual Setup

### Create Repository

```bash
# Public repository
gh repo create myproject --public

# Private repository
gh repo create myproject --private

# With description
gh repo create myproject --public --description "My project description"

# Clone after creation
gh repo create myproject --public --clone

# From existing local repo
cd existing-project
gh repo create --source . --public
```

### Initialize Local Repository

```bash
mkdir myproject
cd myproject
git init
echo "# myproject" > README.md
git add README.md
git commit -m "Initial commit"
gh repo create myproject --source . --public
git push -u origin main
```

## Repository Configuration

### Add Description and Topics
```bash
gh repo edit --description "New description"
gh repo edit --add-topic golang,cli,tool
```

### Set Visibility
```bash
gh repo edit --visibility private
gh repo edit --visibility public
```

### Enable Features
```bash
gh repo edit --enable-issues
gh repo edit --enable-wiki
gh repo edit --enable-projects
```

### Branch Protection
```bash
# Via API (gh repo edit doesn't support all settings)
gh api repos/:owner/:repo/branches/main/protection \
  -X PUT \
  -F required_status_checks='{"strict":true,"contexts":[]}' \
  -F enforce_admins=false \
  -F required_pull_request_reviews='{"required_approving_review_count":1}'
```

## Essential Files

### README.md
```markdown
# Project Name

Brief description

## Installation

## Usage

## Contributing

## License
```

### .gitignore
```bash
# Generate from template
curl -s https://www.toptal.com/developers/gitignore/api/node > .gitignore
```

### LICENSE
```bash
# Add license during creation
gh repo create myproject --license mit

# Or add later
gh api licenses/mit --jq '.body' > LICENSE
```

## Repository Templates

### Create from Template
```bash
gh repo create myproject --template owner/template-repo
```

### Make Repository a Template
```bash
gh repo edit --enable-template
```

## Clone and Fork

### Clone
```bash
gh repo clone owner/repo
ghrepoc owner/repo  # alias

# Clone your own
gh repo clone myproject
```

### Fork
```bash
gh repo fork owner/repo
gh repo fork owner/repo --clone  # Fork and clone

forkclone owner/repo  # dotfiles function
```

## Repository Info

```bash
# View repository
gh repo view
ghrepo  # alias

# View in browser
gh repo view --web

# List your repos
gh repo list

# List org repos
gh repo list myorg
```

## Dotfiles Integration

- `newrepo <name> [desc]` - Create and push new repo
- `forkclone <owner/repo>` - Fork and clone
- `ghrepo` - gh repo view (alias)
- `ghrepoc` - gh repo clone (alias)
- `ghrepof` - gh repo fork (alias)
