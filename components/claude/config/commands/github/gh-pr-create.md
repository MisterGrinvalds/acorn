---
description: Create a GitHub pull request
argument-hint: [--draft] [--web]
allowed-tools: Read, Bash
---

## Task

Help the user create a pull request using GitHub CLI.

## Quick PR Creation

Using dotfiles function:
```bash
quickpr
# Pushes current branch and opens PR creation in browser
```

## Manual PR Creation

### Interactive
```bash
gh pr create
# Prompts for title, body, reviewers, etc.
```

### With Options
```bash
gh pr create \
  --title "Add user authentication" \
  --body "This PR adds login/logout functionality" \
  --reviewer username1,username2 \
  --assignee @me \
  --label enhancement
```

### Draft PR
```bash
gh pr create --draft
# or with dotfiles: quickpr (then convert to draft in web)
```

### From Template
```bash
gh pr create --template .github/PULL_REQUEST_TEMPLATE.md
```

## PR Creation Options

| Flag | Description |
|------|-------------|
| `--title` | PR title |
| `--body` | PR description |
| `--draft` | Create as draft |
| `--reviewer` | Request reviewers |
| `--assignee` | Assign to users |
| `--label` | Add labels |
| `--milestone` | Set milestone |
| `--project` | Add to project |
| `--base` | Base branch (default: main) |
| `--head` | Head branch (default: current) |
| `--web` | Open in browser after |

## PR Body Template

```markdown
## Summary
Brief description of changes

## Changes
- Added X
- Modified Y
- Removed Z

## Testing
- [ ] Unit tests pass
- [ ] Manual testing done

## Related Issues
Closes #123
```

## Workflow Examples

### Feature PR
```bash
# 1. Create branch
newbranch feature/add-search

# 2. Make changes and commit
git add .
git commit -m "Add search functionality"

# 3. Push and create PR
quickpr
# or: gh pr create --title "Add search functionality"
```

### Bug Fix PR
```bash
newbranch fix/login-error
# fix the bug
git add . && git commit -m "Fix login error (#456)"
gh pr create --title "Fix login error" --label bug
```

### Draft PR for Early Feedback
```bash
gh pr create --draft --title "WIP: Refactor auth module"
# Later, mark ready:
gh pr ready
```

## After Creating

```bash
# View your PR
gh pr view

# Check status
prstatus  # dotfiles function

# Watch CI checks
prchecks  # dotfiles function
gh pr checks

# Open in browser
gh pr view --web
```

## Dotfiles Integration

- `quickpr` - Push and create PR (opens web)
- `pushbranch` - Push current branch
- `newbranch <name>` - Create feature branch
- `ghpr` - gh pr create (alias)
- `ghprs` - gh pr status (alias)
