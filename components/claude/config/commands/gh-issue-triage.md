---
description: Triage and manage GitHub issues
argument-hint: [action: list|create|label|assign|close]
allowed-tools: Read, Bash
---

## Task

Help the user triage and manage GitHub issues using gh CLI.

## Actions

Based on `$ARGUMENTS`:

### list
List and filter issues:

```bash
# All open issues
gh issue list

# Filter by label
gh issue list --label bug
gh issue list --label "good first issue"

# Filter by assignee
gh issue list --assignee @me
gh issue list --assignee username

# Filter by state
gh issue list --state closed
gh issue list --state all

# Search
gh issue list --search "login error"

# Using alias
ghissues  # gh issue list
```

### create
Create new issue:

```bash
# Interactive
gh issue create

# With options
gh issue create \
  --title "Bug: Login fails" \
  --body "Steps to reproduce..." \
  --label bug \
  --assignee @me

# From template
gh issue create --template bug_report.md

# Using alias
ghissue  # gh issue create
```

### label
Manage issue labels:

```bash
# Add label
gh issue edit 123 --add-label bug

# Remove label
gh issue edit 123 --remove-label wontfix

# Multiple labels
gh issue edit 123 --add-label "bug,priority:high"

# List all labels
gh label list
```

### assign
Assign issues:

```bash
# Assign to self
gh issue edit 123 --add-assignee @me

# Assign to others
gh issue edit 123 --add-assignee username1,username2

# Remove assignee
gh issue edit 123 --remove-assignee username
```

### close
Close issues:

```bash
# Close issue
gh issue close 123

# Close with comment
gh issue close 123 --comment "Fixed in #456"

# Reopen
gh issue reopen 123
```

## Triage Workflow

### 1. Review New Issues
```bash
# List unlabeled issues
gh issue list --label ""

# Or search for untriaged
gh issue list --search "no:label"
```

### 2. Categorize
```bash
# Add appropriate labels
gh issue edit 123 --add-label bug
gh issue edit 456 --add-label enhancement
```

### 3. Prioritize
```bash
gh issue edit 123 --add-label "priority:high"
gh issue edit 456 --milestone "v2.0"
```

### 4. Assign
```bash
gh issue edit 123 --add-assignee developer1
```

### 5. Link to Project
```bash
gh issue edit 123 --add-project "Sprint 5"
```

## Common Label Categories

### Type
- `bug` - Something isn't working
- `enhancement` - New feature request
- `documentation` - Docs improvements
- `question` - Further information needed

### Priority
- `priority:critical` - Must fix immediately
- `priority:high` - Important
- `priority:medium` - Normal priority
- `priority:low` - Nice to have

### Status
- `needs-triage` - Needs review
- `in-progress` - Being worked on
- `blocked` - Waiting on something
- `wontfix` - Won't be addressed

### Community
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention needed

## Bulk Operations

```bash
# Close all issues with label
gh issue list --label wontfix --json number \
  | jq -r '.[].number' \
  | xargs -I {} gh issue close {}

# Add label to multiple issues
for i in 1 2 3; do gh issue edit $i --add-label reviewed; done
```

## Dotfiles Integration

- `ghissue` - gh issue create (alias)
- `ghissues` - gh issue list (alias)
- `ghissuev` - gh issue view (alias)
