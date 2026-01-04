---
description: Guide to Git rebasing - when and how to use it
argument-hint: [type: basic|interactive|onto]
allowed-tools: Read, Bash
---

## Task

Guide the user through Git rebasing scenarios and best practices.

## Rebase Types

Based on `$ARGUMENTS`:

### basic
Simple rebase to update branch with latest main:

```bash
# Update your feature branch with latest main
git checkout feature/my-feature
git fetch origin
git rebase origin/main

# Or using aliases
gco feature/my-feature
gf
gr origin/main
```

If conflicts occur:
```bash
# Fix conflicts in files, then:
git add <fixed-files>
git rebase --continue  # or: grc

# To abort:
git rebase --abort  # or: gra
```

### interactive
Interactive rebase for history cleanup:

```bash
# Rebase last N commits
git rebase -i HEAD~3  # or: gri HEAD~3

# Rebase from a specific commit
git rebase -i <commit-hash>
```

Interactive commands:
| Command | Action |
|---------|--------|
| `pick` | Keep commit as-is |
| `reword` | Change commit message |
| `edit` | Stop to amend commit |
| `squash` | Combine with previous commit |
| `fixup` | Combine, discard message |
| `drop` | Remove commit |

Common patterns:
```bash
# Squash all into one commit
pick abc123 First commit
squash def456 WIP
squash ghi789 More WIP

# Reorder commits
pick ghi789 Third (now first)
pick abc123 First (now second)
pick def456 Second (now third)
```

### onto
Rebase onto a different base:

```bash
# Move commits to a new base
git rebase --onto <new-base> <old-base> <branch>

# Example: Move feature from develop to main
git rebase --onto main develop feature/my-feature
```

## When to Rebase vs Merge

### Use Rebase When
- Updating feature branch with main (before PR)
- Cleaning up local commit history
- You want linear history
- Commits haven't been pushed/shared

### Use Merge When
- Integrating feature into main
- Preserving branch history matters
- Commits are already public
- Multiple people working on branch

## Golden Rules

1. **Never rebase public/shared branches**
2. **Always use `--force-with-lease` not `--force`**
3. **Create backup branch before complex rebases**
4. **Know your escape: `git rebase --abort`**

## Recovery

If rebase goes wrong:
```bash
# Abort current rebase
git rebase --abort

# Find pre-rebase state in reflog
git reflog
git reset --hard <pre-rebase-commit>
```

## Dotfiles Aliases

- `gr` - git rebase
- `gri` - git rebase -i (interactive)
- `grc` - git rebase --continue
- `gra` - git rebase --abort
- `gpf` - git push --force-with-lease
