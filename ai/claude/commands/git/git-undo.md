---
description: Undo Git operations safely
argument-hint: [operation: commit|add|merge|rebase|push]
allowed-tools: Read, Bash
---

## Task

Help the user undo various Git operations safely.

## Undo Operations

Based on `$ARGUMENTS`:

### commit
Undo the last commit:

```bash
# Undo commit, keep changes staged
git reset --soft HEAD~1
gundo  # dotfiles function (same as above)

# Undo commit, keep changes unstaged
git reset HEAD~1
# or: git reset --mixed HEAD~1

# Undo commit, discard changes (DANGEROUS)
git reset --hard HEAD~1
grhh  # dotfiles alias
```

Undo multiple commits:
```bash
# Undo last 3 commits, keep changes
git reset --soft HEAD~3
```

### add
Unstage files:

```bash
# Unstage specific file
git reset HEAD <file>
grh <file>  # dotfiles alias

# Unstage all files
git reset HEAD
grh  # dotfiles alias

# Discard changes in working directory (DANGEROUS)
git checkout -- <file>
# or: git restore <file>
```

### merge
Undo a merge:

```bash
# During merge (before commit)
git merge --abort

# After merge commit (not pushed)
git reset --hard HEAD~1

# After merge is pushed (safe - creates revert commit)
git revert -m 1 <merge-commit>
```

### rebase
Undo a rebase:

```bash
# During rebase
git rebase --abort
gra  # dotfiles alias

# After rebase completed
# Find pre-rebase state
git reflog
# Look for "rebase (start)" or the commit before rebase
git reset --hard <pre-rebase-commit>
```

### push
Undo a push:

```bash
# If you're the only one who pulled:
git reset --hard <previous-commit>
git push --force-with-lease
gpf  # dotfiles alias (safer than --force)

# If others may have pulled (SAFER):
git revert <bad-commit>
git push

# Revert multiple commits
git revert <oldest>..<newest>
```

## Recovery with Reflog

The reflog tracks all HEAD movements:

```bash
# View reflog
git reflog

# Example output:
# abc123 HEAD@{0}: commit: Current commit
# def456 HEAD@{1}: rebase finished
# ghi789 HEAD@{2}: rebase (start)
# jkl012 HEAD@{3}: commit: Before rebase

# Restore to any previous state
git reset --hard HEAD@{3}
# or: git reset --hard jkl012
```

## Safety Guidelines

### Safe Operations (reversible)
- `git reset --soft` - Keeps all changes
- `git revert` - Creates undo commit
- `git stash` - Temporarily saves changes

### Dangerous Operations (may lose work)
- `git reset --hard` - Discards changes
- `git push --force` - Overwrites remote
- `git clean -fd` - Deletes untracked files

### Before Risky Operations
```bash
# Create backup branch
git branch backup-branch

# Or stash changes
git stash
```

## Common Scenarios

### "I committed to wrong branch"
```bash
# Save the commit hash
git log -1  # note the hash

# Undo commit on current branch
git reset --soft HEAD~1

# Switch to correct branch
git checkout correct-branch

# Apply changes there
git commit
```

### "I need to change last commit message"
```bash
git commit --amend
# or: gca (dotfiles alias)
```

### "I accidentally deleted a branch"
```bash
# Find the commit in reflog
git reflog

# Recreate branch
git checkout -b <branch-name> <commit-hash>
```

## Dotfiles Integration

- `gundo` - Undo last commit (soft reset)
- `grh` - git reset HEAD
- `grhh` - git reset HEAD --hard
- `grhs` - git reset HEAD --soft
- `gra` - git rebase --abort
- `gpf` - git push --force-with-lease
