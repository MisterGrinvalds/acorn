---
description: Interactive coaching session to learn GitHub CLI workflows
argument-hint: [skill-level: beginner|intermediate|advanced]
allowed-tools: Read, Glob, Grep, Bash
---

## Task

Guide the user through learning GitHub CLI workflows interactively.

## Approach

1. **Assess level** - Ask about GitHub/git experience
2. **Set goals** - Identify what workflows they need
3. **Progressive exercises** - Start simple, build complexity
4. **Real-time practice** - Have them run gh commands
5. **Reinforce** - Summarize and suggest next steps

## Skill Levels

### Beginner
- Installing and authenticating gh
- Creating and viewing PRs
- Basic issue management
- Cloning and forking repos
- Using gh browse

### Intermediate
- PR review workflow
- Issue triage and labels
- GitHub Actions basics
- Using dotfiles functions
- Branch management

### Advanced
- gh api for custom queries
- Complex workflow automation
- Release management
- Repository templates
- GitHub Apps integration

## Interactive Exercises

### Beginner Exercises
```bash
# Exercise 1: Check authentication
gh auth status

# Exercise 2: View repository
ghrepo  # or: gh repo view

# Exercise 3: List open PRs
ghprl  # or: gh pr list

# Exercise 4: Create an issue
ghissue  # or: gh issue create

# Exercise 5: Open in browser
ghweb  # or: gh browse
```

### Intermediate Exercises
```bash
# Exercise 6: Create feature branch and PR
newbranch feature/test
# make changes
quickpr  # dotfiles function

# Exercise 7: Check PR status
prstatus
prchecks

# Exercise 8: Review a PR
gh pr checkout 123
gh pr review --approve

# Exercise 9: Clean up
gitcleanup  # dotfiles function
```

### Advanced Exercises
```bash
# Exercise 10: Query with gh api
gh api repos/:owner/:repo/issues

# Exercise 11: Create a release
gh release create v1.0.0 --generate-notes

# Exercise 12: Manage Actions
gh workflow list
gh run watch
```

## Context

@components/github/functions.sh
@components/github/aliases.sh

## Coaching Style

- Start with authentication
- Use aliases for common operations
- Show web fallbacks when helpful
- Build toward full PR workflow
- Emphasize automation benefits
