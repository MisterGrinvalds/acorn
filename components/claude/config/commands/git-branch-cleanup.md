---
description: Clean up old and merged Git branches
argument-hint: [scope: local|remote|all]
allowed-tools: Read, Bash
---

## Task

Help the user clean up old, merged, or stale Git branches.

## Cleanup Scopes

Based on `$ARGUMENTS`:

### local
Clean up local branches that have been merged:

```bash
# List merged branches (excluding main/master)
git branch --merged | grep -v '\*\|main\|master'

# Delete merged branches
gcleanb  # dotfiles function
# or manually:
git branch --merged | grep -v '\*\|main\|master' | xargs -n 1 git branch -d
```

### remote
Clean up remote tracking branches:

```bash
# Prune remote-tracking branches
git fetch --prune

# List remote branches
git branch -r

# Delete a remote branch
git push origin --delete <branch-name>
```

### all
Full cleanup of both local and remote:

```bash
# 1. Fetch and prune
git fetch --all --prune

# 2. Delete local merged branches
gcleanb

# 3. List remaining branches for review
git branch -a
```

## Safety Checks

Before deleting, always verify:

```bash
# Show what would be deleted (dry run)
git branch --merged | grep -v '\*\|main\|master'

# Check if branch has unmerged commits
git log main..<branch-name>

# Check branch age
git for-each-ref --sort=-committerdate refs/heads/ \
  --format='%(committerdate:short) %(refname:short)'
```

## Advanced Cleanup

### Find stale branches (older than N days)
```bash
# Branches not updated in 30 days
git for-each-ref --sort=-committerdate refs/heads/ \
  --format='%(committerdate:relative) %(refname:short)' | \
  grep -E 'months|weeks'
```

### Find branches with no remote
```bash
git branch -vv | grep ': gone]'
```

### Delete branches with no remote
```bash
git branch -vv | grep ': gone]' | awk '{print $1}' | xargs git branch -d
```

## Dotfiles Integration

- `gcleanb` - Clean merged local branches (excludes main/master)
- `gb` - List local branches
- `gba` - List all branches (including remote)
- `gbd <branch>` - Delete branch (safe)
- `gbD <branch>` - Force delete branch

## Caution

- `-d` only deletes if merged, `-D` forces deletion
- Cannot delete current branch (checkout main first)
- Remote deletion is permanent
- Always verify before bulk operations
