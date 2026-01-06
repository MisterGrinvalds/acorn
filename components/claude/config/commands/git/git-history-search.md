---
description: Search Git history for commits, changes, and code
argument-hint: <search-term>
allowed-tools: Read, Bash
---

## Task

Help the user search through Git history to find commits, changes, or code.

## Search Types

### Search Commit Messages
```bash
# Find commits with message containing term
git log --oneline --all --grep="$ARGUMENTS"

# Using dotfiles function
gfind "$ARGUMENTS"

# Case insensitive
git log --oneline --all --grep="$ARGUMENTS" -i
```

### Search Code Changes (Pickaxe)
```bash
# Find commits that added or removed this string
git log -S "$ARGUMENTS" --oneline

# With diff context
git log -S "$ARGUMENTS" -p

# Regex search
git log -G "$ARGUMENTS" --oneline
```

### Search File Content at Revision
```bash
# Search in specific commit
git grep "$ARGUMENTS" <commit>

# Search in all branches
git grep "$ARGUMENTS" $(git rev-list --all)

# Search with context
git grep -n -C 3 "$ARGUMENTS"
```

### Search by Author
```bash
# Commits by author
git log --author="name" --oneline

# Commits by author in date range
git log --author="name" --since="2024-01-01" --until="2024-06-01"
```

### Search by Date
```bash
# Commits since date
git log --since="2024-01-01" --oneline

# Commits in last 2 weeks
git log --since="2 weeks ago" --oneline

# Commits between dates
git log --since="2024-01-01" --until="2024-01-31"
```

### Search by File
```bash
# History of specific file
git log --oneline -- <file>

# With diffs
git log -p -- <file>

# Who changed this file
git log --format="%an" -- <file> | sort | uniq -c | sort -rn
```

## Advanced Searches

### Find When Line Was Added
```bash
# Blame shows who added each line
git blame <file>
gblame <file>  # dotfiles function with line numbers

# Blame specific line range
git blame -L 10,20 <file>
```

### Find When Code Was Deleted
```bash
# Search for deleted code
git log -S "deleted_function" --diff-filter=D

# Show the commit that deleted it
git log -S "deleted_function" -p
```

### Find Merge Commits
```bash
# All merge commits
git log --merges --oneline

# Merges to main
git log --merges --oneline main
```

### Find First Commit of File
```bash
git log --diff-filter=A --summary -- <file>
```

## Output Formatting

```bash
# Custom format
git log --pretty=format:"%h %ad %an: %s" --date=short

# One line with date
git log --oneline --format="%h %ad %s" --date=short

# Graph view
git log --oneline --graph --all
gla  # dotfiles alias
```

## FZF Integration

If fzf component is loaded:
```bash
# Interactive log browser with preview
fzf_git_log  # or: fgl
```

## Dotfiles Functions

- `gfind <term>` - Search commit messages
- `glog` - Pretty graph log
- `gl` - Last 20 commits oneline
- `gla` - All branches graph
- `gblame <file>` - Blame with line numbers
- `gshow [commit]` - Show commit details
