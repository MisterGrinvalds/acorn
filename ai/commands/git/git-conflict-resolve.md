---
description: Help resolve Git merge or rebase conflicts
allowed-tools: Read, Bash, Glob, Grep
---

## Task

Help the user understand and resolve Git conflicts during merge or rebase operations.

## Process

1. **Identify conflict state** - Check current status
2. **List conflicted files** - Show what needs resolution
3. **Explain conflict markers** - Help understand the format
4. **Guide resolution** - Step through fixing each file
5. **Complete operation** - Finish merge/rebase

## Step 1: Check Status

```bash
# See current state
git status

# During merge conflict:
# "You have unmerged paths"
# "both modified: <file>"

# During rebase conflict:
# "rebase in progress"
# "fix conflicts and then run git rebase --continue"
```

## Step 2: Understand Conflict Markers

```
<<<<<<< HEAD (or <<<<<<< yours)
Your changes (current branch)
=======
Their changes (incoming branch)
>>>>>>> branch-name (or >>>>>>> theirs)
```

For 3-way diff:
```
<<<<<<< HEAD
Your version
||||||| merged common ancestor
Original version
=======
Their version
>>>>>>> feature-branch
```

## Step 3: Resolution Options

### Manual Resolution
1. Open conflicted file in editor
2. Decide which changes to keep
3. Remove conflict markers
4. Save file

### Using Git Commands
```bash
# Keep your version entirely
git checkout --ours <file>

# Keep their version entirely
git checkout --theirs <file>

# Interactive merge tool
git mergetool
```

### Using VS Code
```bash
# VS Code shows inline conflict resolution
code <conflicted-file>
# Click "Accept Current", "Accept Incoming", "Accept Both", or "Compare"
```

## Step 4: Complete Resolution

### For Merge
```bash
# After resolving all conflicts
git add <resolved-files>
git commit  # Will use merge commit message

# Or abort
git merge --abort
```

### For Rebase
```bash
# After resolving conflicts in current commit
git add <resolved-files>
git rebase --continue  # or: grc

# Skip this commit
git rebase --skip

# Abort entire rebase
git rebase --abort  # or: gra
```

## Common Conflict Scenarios

### Same Line Modified
Both branches changed the same line → manual decision needed

### File Deleted vs Modified
One branch deleted, other modified → decide if file should exist

### Rename Conflicts
Both branches renamed same file → choose final name

## Prevention Tips

1. **Pull frequently** - Keep branch up to date
2. **Small PRs** - Less chance of conflicts
3. **Communicate** - Know what others are working on
4. **Rebase before PR** - Resolve conflicts locally

## Helpful Commands

```bash
# Show conflicted files
git diff --name-only --diff-filter=U

# Show conflict details
git diff

# Show what you're merging in
git log HEAD..MERGE_HEAD --oneline

# Show base version (common ancestor)
git show :1:<file>  # Base
git show :2:<file>  # Ours
git show :3:<file>  # Theirs
```

## Dotfiles Integration

- `gd` - git diff (see conflicts)
- `grc` - git rebase --continue
- `gra` - git rebase --abort
- `gs` - git status (check state)
